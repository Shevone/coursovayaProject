package accounts

import (
	"context"
	"fitnes-gateway/internal/models"
	"fmt"
	pb "github.com/Shevone/proto-fitnes/gen/go/accounts"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	grpcrlog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"os"
	"time"
)

var (
	roleUser    = "User"
	roleTrainer = "Trainer"
	roleAdmin   = "Admin"
	roleNew     = "New"
)

type AccountsService struct {
	api pb.AuthClient
	log *slog.Logger
}

// NewAccountClient конструктор клиента grpc
func NewAccountClient(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*AccountsService, error) {
	const op = "grpc.New"

	// Конфигурируем то, в каких случаях делаем retry
	retryOpts := []grpcretry.CallOption{
		// Коды при которых делаем retry
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		// Максимальное количество попыток новых
		grpcretry.WithMax(uint(0)),
		// таймаут ретраев
		grpcretry.WithPerRetryTimeout(timeout),
	}
	// Логирование запросов
	logOpt := []grpcrlog.Option{
		// На какие события мы реагируем
		// По сути тело запроса и ответа
		grpcrlog.WithLogOnEvents(grpcrlog.PayloadReceived, grpcrlog.PayloadSent),
	}

	// Подключение к сервису
	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// Оборачиваем 2 interceptors в 1
		// Похоже на middleware
		grpc.WithChainUnaryInterceptor(
			grpcrlog.UnaryClientInterceptor(InterceptorLogger(log), logOpt...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	// Засовываем наше подключение в сервис
	return &AccountsService{
		api: pb.NewAuthClient(cc),
		log: log,
	}, nil
}

// InterceptorLogger адаптирует slog logger под logger interceptor - a
// Обертка над логгером, чтоб им мог пользоватлься интерсептор

func InterceptorLogger(l *slog.Logger) grpcrlog.Logger {

	// Возвращает функцию, которая будет вызваться внутри интерсептора
	// А в свою очередь в эту функцию мы засунули наш логгер и метод Log

	return grpcrlog.LoggerFunc(
		func(ctx context.Context, level grpcrlog.Level, msg string, fields ...any) {
			l.Log(ctx, slog.Level(level), msg, fields...)
		},
	)
}

// CreateUser метод регистрации пользователя
func (aS *AccountsService) CreateUser(ctx *gin.Context, user models.User) (int64, error) {
	const op = "grpc.Account.CreateUser"

	registerRequest := &pb.RegisterRequest{
		Email:       user.Email,
		Password:    user.Password,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
		Surname:     user.Surname,
		Role:        pb.Role_New,
	}

	resp, err := aS.api.Register(ctx, registerRequest)

	if err != nil {
		return -1, fmt.Errorf("%s : %w", op, err)
	}
	return resp.UserId, nil
}

// LoginUser метод авторизации пользователя
func (aS *AccountsService) LoginUser(ctx *gin.Context, loginInput models.LoginInput) (string, error) {
	const op = "grpc.Account.LoginUser"

	loginRequest := &pb.LoginRequest{Email: loginInput.Email, Password: loginInput.Password}

	loginResp, err := aS.api.Login(ctx, loginRequest)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return loginResp.Token, nil
}

// ParseToken метод парсинга токена
func (aS *AccountsService) ParseToken(accessToken string) (*models.TokenClaims, error) {
	const op = "grpc.Account.ParseToken"
	// Парсим и проверяем токен
	token, err := jwt.ParseWithClaims(
		accessToken,
		&models.TokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%s: %s", op, "invalid signing method")
			}
			return []byte(os.Getenv("APP_SECRET")), nil
		})

	// Если произошла ошибка
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	claims, ok := token.Claims.(*models.TokenClaims)
	if !ok {
		return nil, fmt.Errorf("%s: %s", op, "token claims are not of type *tokenClaims")
	}

	return claims, nil
}

// GetUserProfile получить 1 профиль по id пользователя
func (aS *AccountsService) GetUserProfile(ctx *gin.Context, userId int64) (*models.User, error) {
	const op = "grpc.userProfile"

	loginRequest := &pb.GetDataRequest{UserId: userId}

	loginResp, err := aS.api.GetUserData(ctx, loginRequest)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	user := models.User{
		Id:          userId,
		Email:       loginResp.User.Email,
		Name:        loginResp.User.Name,
		Surname:     loginResp.User.Surname,
		Patronymic:  loginResp.User.Patronymic,
		PhoneNumber: loginResp.User.PhoneNumber,
	}
	switch loginResp.User.Role {
	case pb.Role_User:
		user.Role = roleUser
	case pb.Role_Trainer:
		user.Role = roleTrainer
	case pb.Role_Admin:
		user.Role = roleAdmin
	case pb.Role_New:
		user.Role = roleNew
	}

	return &user, nil
}

// EditUserProfile метод вызывающий редактирование профиля
func (aS *AccountsService) EditUserProfile(ctx *gin.Context, profileRequest models.EditProfileRequest) (string, error) {
	const op = "grpc.editProfile"
	userId, _ := profileRequest.GetId()
	updaterId := ctx.Value("userId").(int64)
	editReq := &pb.EditRequest{
		UserId:     userId,
		Name:       profileRequest.Name,
		Surname:    profileRequest.Surname,
		Patronymic: profileRequest.Patronymic,
		UpdaterId:  updaterId,
	}
	editResponse, err := aS.api.EditProfile(ctx, editReq)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return editResponse.Message, nil

}

// GetUsers Метод получения пользователей
func (aS *AccountsService) GetUsers(ctx *gin.Context, request models.PaginateRequest) ([]models.User, error) {
	const op = "grpc.getUsers"

	getUsersResponse, err := aS.api.GetUsers(ctx,
		&pb.GetUsersRequest{Page: request.Page, Limit: request.Limit})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	usersList := make([]models.User, 0, cap(getUsersResponse.Users))
	for _, userResp := range getUsersResponse.Users {
		user := models.User{
			Id:          userResp.Id,
			Email:       userResp.Email,
			Name:        userResp.Name,
			Surname:     userResp.Surname,
			Patronymic:  userResp.Patronymic,
			PhoneNumber: userResp.PhoneNumber,
		}
		switch userResp.Role {
		case pb.Role_User:
			user.Role = roleUser
		case pb.Role_Trainer:
			user.Role = roleTrainer
		case pb.Role_Admin:
			user.Role = roleAdmin
		case pb.Role_New:
			user.Role = roleNew

		}
		usersList = append(usersList, user)
	}

	return usersList, nil
}

func (aS *AccountsService) EditUserRole(ctx *gin.Context, request models.EditUserRoleRequest) (string, error) {
	const op = "grpc.editUserRole"
	updaterId := ctx.Value("userId").(int64)
	updateReq := &pb.UpdateUserRoleReq{UserId: request.UserId, UpdaterId: updaterId}

	switch request.NewRole {
	case 0:
		updateReq.Role = pb.Role_User
	case 1:
		updateReq.Role = pb.Role_Trainer
	case 2:
		updateReq.Role = pb.Role_Admin
	case 3:
		updateReq.Role = pb.Role_New

	default:
		return "", fmt.Errorf("%s: %w", op, "User role has invalid argument")
	}
	updateResp, err := aS.api.UpdateUserRole(ctx, updateReq)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return updateResp.Message, nil
}
func (aS *AccountsService) EditUserPassword(ctx *gin.Context, request models.UpdatePasswordRequest) (string, error) {
	const op = "grpc.editUserPassword"

	updaterId := ctx.Value("userId").(int64)
	editReq := &pb.UpdateUserPasswordReq{UserId: request.UserId, Password: request.Password, UpdaterId: updaterId}
	updateResp, err := aS.api.UpdateUserPassword(ctx, editReq)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return updateResp.Message, nil
}
