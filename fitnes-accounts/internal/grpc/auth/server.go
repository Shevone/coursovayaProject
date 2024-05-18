package auth

import (
	"context"
	"fitnes-account/internal/models"
	"github.com/Shevone/proto-fitnes/gen/go/accounts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	roleUser    = "User"
	roleTrainer = "Trainer"
	roleAdmin   = "Admin"
)

type serverApi struct {
	accountsFitnesv1.UnimplementedAuthServer
	accountsService Accounts
}

// Accounts - Интерфейс нижнего слоя
type Accounts interface {
	Login(ctx context.Context, email string, password string) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string, password string,
		name string, surname string, patronymic string,
		role string, phoneNumber string) (userId int64, err error)
	EditUserProfile(
		ctx context.Context,
		userId int64,
		name string, surname string, patronymic string) error
	GetUserData(ctx context.Context, userId int64) (
		user models.User, err error)
	GetUsers(ctx context.Context, page int64, limit int64) (
		[]models.User, error)
	UpdateUserRole(ctx context.Context, userId int64, newRole string) (
		string, error)
}

// Register - метод регистрирующий наш обработчик(Accounts) на созданный grpc сервер
func Register(gRPCServer *grpc.Server, auth Accounts) {
	accountsFitnesv1.RegisterAuthServer(gRPCServer, &serverApi{accountsService: auth})
}

// ====================================
// Далее описаны методы, которые должен реализовавыать сервер по задумке protobuf

// Login - ручка для авторизации
func (s *serverApi) Login(ctx context.Context, in *accountsFitnesv1.LoginRequest) (*accountsFitnesv1.LoginResponse, error) {
	// Валидируем посутпившие значения
	if in.Email == "" {

		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if in.Password == "" {

		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	token, err := s.accountsService.Login(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {

		return nil, status.Error(codes.Internal, "failed to login")
	}
	return &accountsFitnesv1.LoginResponse{Token: token}, nil
}

// Register - ручка для регистрации
func (s *serverApi) Register(ctx context.Context, in *accountsFitnesv1.RegisterRequest) (*accountsFitnesv1.RegisterResponse, error) {
	if in.Email == "" {

		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if in.Password == "" {

		return nil, status.Error(codes.InvalidArgument, "password is required")
	}
	if in.Name == "" || in.Surname == "" {

		return nil, status.Error(codes.InvalidArgument, "all personal data is required")
	}
	if in.PhoneNumber == "" {

		return nil, status.Error(codes.InvalidArgument, "phone number is required")
	}
	userId, err := s.accountsService.RegisterNewUser(
		ctx, in.GetEmail(),
		in.GetPassword(),
		in.GetName(),
		in.GetSurname(),
		in.GetPatronymic(),
		roleUser,
		in.GetPhoneNumber())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to register")
	}
	return &accountsFitnesv1.RegisterResponse{UserId: userId}, nil
}

// EditProfile - для изменения данных
func (s *serverApi) EditProfile(ctx context.Context, in *accountsFitnesv1.EditRequest) (*accountsFitnesv1.EditResponse, error) {
	if in.Name == "" || in.Surname == "" {

		return nil, status.Error(codes.InvalidArgument, "all personal data is required")
	}
	err := s.accountsService.EditUserProfile(ctx, in.GetUserId(), in.GetName(), in.GetSurname(), in.GetPatronymic())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to edit user profile")
	}

	return &accountsFitnesv1.EditResponse{Message: "Профиль отредактирован"}, nil
}

// GetUserData - получить данные пользователя
func (s *serverApi) GetUserData(ctx context.Context, in *accountsFitnesv1.GetDataRequest) (*accountsFitnesv1.GetDataResponse, error) {
	if in.UserId < 0 {

		return nil, status.Error(codes.InvalidArgument, "user id must be greater than 0")
	}
	user, err := s.accountsService.GetUserData(ctx, in.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user profile data")
	}
	userModelResponse := &accountsFitnesv1.UserModel{
		Id:          user.ID,
		Email:       user.Email,
		Name:        user.Name,
		Surname:     user.Surname,
		Patronymic:  user.Patronymic,
		PhoneNumber: user.PhoneNumber,
	}
	switch user.Role {
	case roleUser:
		userModelResponse.Role = accountsFitnesv1.Role_User
	case roleAdmin:
		userModelResponse.Role = accountsFitnesv1.Role_Admin
	case roleTrainer:
		userModelResponse.Role = accountsFitnesv1.Role_Trainer
	}
	return &accountsFitnesv1.GetDataResponse{User: userModelResponse}, nil
}

// GetUsers - получить список всех пользователей с пагинацией
// метод предназначен для пользователей с правами admin, для того чтобы просматривать всех и менять роли
func (s *serverApi) GetUsers(ctx context.Context, in *accountsFitnesv1.GetUsersRequest) (*accountsFitnesv1.GetUsersResponse, error) {
	if in.Page < 0 {
		return nil, status.Error(codes.InvalidArgument, "page must be greater than 0")
	}
	if in.Limit < 1 {
		return nil, status.Error(codes.InvalidArgument, "counts el on page must be greater than 1")
	}

	users, err := s.accountsService.GetUsers(ctx, in.GetPage(), in.GetLimit())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get users")
	}

	getUserResp := &accountsFitnesv1.GetUsersResponse{}
	getUserResp.Users = make([]*accountsFitnesv1.UserModel, 0, cap(users))
	for _, user := range users {
		userModelResponse := &accountsFitnesv1.UserModel{
			Id:          user.ID,
			Email:       user.Email,
			Name:        user.Name,
			Surname:     user.Surname,
			Patronymic:  user.Patronymic,
			PhoneNumber: user.PhoneNumber,
		}
		switch user.Role {
		case roleUser:
			userModelResponse.Role = accountsFitnesv1.Role_User
		case roleAdmin:
			userModelResponse.Role = accountsFitnesv1.Role_Admin
		case roleTrainer:
			userModelResponse.Role = accountsFitnesv1.Role_Trainer
		}
		getUserResp.Users = append(getUserResp.Users, userModelResponse)
	}

	return getUserResp, nil
}

// UpdateUserRole - изменить роль определенного пользователя по его id
// Метод предназначен для пользователей с правами admin для выдачи ролей после регистрации пользователя
func (s *serverApi) UpdateUserRole(ctx context.Context, in *accountsFitnesv1.UpdateUserRoleReq) (*accountsFitnesv1.UpdateUserRoleResp, error) {
	if in.UserId < 0 {
		return nil, status.Error(codes.InvalidArgument, "user id must be greater than 0")
	}

	message, err := s.accountsService.UpdateUserRole(ctx, in.UserId, in.Role.String())

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to change user role")
	}
	return &accountsFitnesv1.UpdateUserRoleResp{Message: message}, nil
}
