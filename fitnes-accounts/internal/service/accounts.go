package service

import (
	"context"
	"errors"
	"fitnes-account/internal/lib"
	"fitnes-account/internal/models"
	"fitnes-account/internal/repository"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"os"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AccountService struct {
	logger      *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		newUser models.User,
	) (uid int64, err error)
	EditUser(
		ctx context.Context,
		editedProfile models.User,
	) error
	EditUserRole(
		ctx context.Context,
		userId int64,
		newRole string,
	) (string, error)
}
type UserProvider interface {
	User(
		ctx context.Context,
		email string,
	) (models.User, error)
	GetUserDataById(
		ctx context.Context,
		userid int64,
	) (models.User, error)
	GetUsers(
		ctx context.Context,
		page int64,
		limit int64,
	) ([]models.User, error)
}

// NewAccountService - конструктор сервисного слоя
func NewAccountService(
	logger *slog.Logger, usrSaver UserSaver, usrProvider UserProvider, tokenTTl time.Duration,
) *AccountService {
	return &AccountService{
		logger:      logger,
		usrSaver:    usrSaver,
		usrProvider: usrProvider,
		tokenTTL:    tokenTTl,
	}
}

// ==========================
// Методы сервисного слоя

// Login - метод авторизации
func (a *AccountService) Login(ctx context.Context, email string, password string) (token string, err error) {
	const op = "Accounts.Login"

	log := a.logger.With(
		slog.String("op", op),
		slog.String("username", email),
	)

	log.Info("attempting to login user")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			a.logger.Warn("user not found", err)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.logger.Error("failed to get user", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.logger.Info("invalid credentials", err)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	token, err = lib.NewToken(user, os.Getenv("APP_SECRET"), a.tokenTTL)
	if err != nil {
		a.logger.Error("failed to generate token", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *AccountService) RegisterNewUser(
	ctx context.Context, email string, password string, name string,
	surname string, patronymic string, role string, phoneNumber string,
) (userId int64, err error) {
	const op = "Accounts.RegisterNewUser"

	log := a.logger.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", err)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(
		ctx, models.User{
			Email: email, PassHash: passHash, Name: name,
			Surname: surname, Patronymic: patronymic, Role: role, PhoneNumber: phoneNumber})
	if err != nil {
		log.Error("failed to save user", err)

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *AccountService) EditUserProfile(
	ctx context.Context, userId int64, name string, surname string, patronymic string,
) error {
	const op = "Accounts.EditUserProfile"

	log := a.logger.With(
		slog.String("op", op),
		slog.Int64("userId", userId),
	)

	log.Info("edit suer profile")

	err := a.usrSaver.EditUser(ctx, models.User{ID: userId, Name: name, Surname: surname, Patronymic: patronymic})
	if err != nil {
		log.Error("failed to edit user Profile")
		return err
	}
	return nil
}

func (a *AccountService) GetUserData(ctx context.Context, userId int64) (user models.User, err error) {
	const op = "Accounts.GetUserData"

	log := a.logger.With(
		slog.String("op", op),
		slog.Int64("userId", userId),
	)

	log.Info("attempt to retrieve user data")

	user, err = a.usrProvider.GetUserDataById(ctx, userId)
	if err != nil {
		log.Error("Failed to get user data")
		return models.User{}, err
	}
	return user, nil
}

func (a *AccountService) GetUsers(ctx context.Context, page int64, limit int64) ([]models.User, error) {
	const op = "Accounts.GetUsers"

	log := a.logger.With(
		slog.String("op", op),
		slog.Int64("page", page),
		slog.Int64("limit", limit),
	)

	log.Info("attempt to get users")
	users, err := a.usrProvider.GetUsers(ctx, page, limit)

	if err != nil {
		log.Error("Failed to get users")
		return []models.User{}, err
	}
	return users, nil
}

func (a *AccountService) UpdateUserRole(ctx context.Context, userId int64, newRole string) (string, error) {
	const op = "Accounts.UpdateUserRole"

	log := a.logger.With(
		slog.String("op", op),
		slog.Int64("userId", userId),
		slog.String("New role", newRole),
	)

	log.Info("attempting to change role to user")

	message, err := a.usrSaver.EditUserRole(ctx, userId, newRole)
	if err != nil {
		log.Error("failed to change role")
		return "err", err
	}
	return message, nil
}
