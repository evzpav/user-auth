package auth

import (
	"context"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/evzpav/user-auth/internal/domain"
	"gitlab.com/evzpav/user-auth/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userService domain.UserService
}

func NewService(userService domain.UserService) *service {
	return &service{
		userService: userService,
	}
}

func (s *service) hash(password string) (string, error) {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPW), nil
}

func (s *service) hashMatchesPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (s *service) generateToken() string {
	sID := uuid.NewV4()
	return sID.String()
}

func (s *service) Authenticate(ctx context.Context, authUser *domain.AuthUser) error {
	u, err := s.userService.FindByEmail(ctx, authUser.Email)
	if err != nil {
		authUser.Errors["Credentials"] = "invalid credentials"
		return errors.NewNotAuthorized(domain.ErrInvalidCredentials)
	}

	if u == nil {
		authUser.Errors["Credentials"] = "invalid credentials"
		return errors.NewNotAuthorized(domain.ErrInvalidCredentials)
	}

	if !s.hashMatchesPassword(u.Password, authUser.Password) {
		authUser.Errors["Credentials"] = "invalid credentials"
		return errors.NewNotAuthorized(domain.ErrInvalidCredentials)
	}

	token := s.generateToken()
	authUser.Token = token

	// _, err = s.userService.Update(ctx, u)
	// if err != nil {
	// 	return nil, err
	// }

	return nil
}

func (s *service) Signup(ctx context.Context, authUser *domain.AuthUser) error {
	existingUser, err := s.userService.FindByEmail(ctx, authUser.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		authUser.Errors["Credentials"] = "email already being used"
		return fmt.Errorf("email already being used")
	}

	hashedPassword, err := s.hash(authUser.Password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    authUser.Email,
		Password: hashedPassword,
		Token:    authUser.Token,
	}

	return s.userService.Create(ctx, user)
}

func (s *service) Me(ctx context.Context) (*domain.User, error) {
	return nil, nil
}
