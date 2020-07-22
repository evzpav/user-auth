package auth

import (
	"context"
	"fmt"

	"net/smtp"

	uuid "github.com/satori/go.uuid"
	"gitlab.com/evzpav/user-auth/internal/domain"
	"gitlab.com/evzpav/user-auth/pkg/errors"
	"gitlab.com/evzpav/user-auth/pkg/log"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userService     domain.UserService
	emailFrom       string
	emailPassword   string
	googleSigninCli domain.GoogleSigner
	platformURL     string
	log             log.Logger
}

func NewService(userService domain.UserService, emailFrom, emailPassword string, googleSigninCli domain.GoogleSigner, platformURL string, log log.Logger) *service {
	return &service{
		userService:     userService,
		emailFrom:       emailFrom,
		emailPassword:   emailPassword,
		googleSigninCli: googleSigninCli,
		platformURL:     platformURL,
		log:             log,
	}
}

func (s *service) GetGoogleSigninLink(state string) string {
	return s.googleSigninCli.GetLoginURL(state)
}

func (s *service) GetGoogleProfile(code string) (*domain.GoogleUser, error) {
	return s.googleSigninCli.GetProfile(code)
}

func (s *service) hashPassword(password string) (string, error) {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPW), nil
}

func (s *service) hashMatchesPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (s *service) GenerateToken() string {
	sID := uuid.NewV4()
	return sID.String()
}

func (s *service) Authenticate(ctx context.Context, authUser *domain.AuthUser) (*domain.User, error) {
	user, err := s.userService.FindByEmail(ctx, authUser.Email)
	if err != nil {
		authUser.Errors["Credentials"] = "invalid credentials"
		return nil, errors.NewNotAuthorized(domain.ErrInvalidCredentials)
	}

	if user == nil {
		authUser.Errors["Credentials"] = "invalid credentials"
		return nil, errors.NewNotAuthorized(domain.ErrInvalidCredentials)
	}

	if !s.hashMatchesPassword(user.Password, authUser.Password) {
		authUser.Errors["Credentials"] = "invalid credentials"
		return nil, errors.NewNotAuthorized(domain.ErrInvalidCredentials)
	}

	token := s.GenerateToken()
	authUser.Token = token
	user.Token = token

	if err := s.userService.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil

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

	hashedPassword, err := s.hashPassword(authUser.Password)
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

func (s *service) SignupWithGoogle(ctx context.Context, authUser *domain.AuthUser) error {
	existingUser, err := s.userService.FindByEmail(ctx, authUser.Email)
	if err != nil {
		return err
	}

	if existingUser != nil {
		return nil
	}

	hashedPassword, err := s.hashPassword(authUser.Password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    authUser.Email,
		Password: hashedPassword,
	}

	if authUser.Token != "" {
		user.Token = authUser.Token
	}

	if authUser.Name != "" {
		user.Name = authUser.Name
	}

	if authUser.GoogleID != "" {
		user.GoogleID = authUser.GoogleID
	}

	return s.userService.Create(ctx, user)
}

func (s *service) AuthenticateToken(ctx context.Context, token string) (*domain.User, error) {
	user, err := s.userService.FindByToken(ctx, token)
	if err != nil {
		return nil, errors.NewNotAuthorized(domain.ErrInvalidCredentials)
	}
	if user == nil {
		return nil, errors.NewNotAuthorized(domain.ErrInvalidCredentials)
	}

	user.Token = s.GenerateToken()

	return user, s.userService.Update(ctx, user)
}

func (s *service) generateResetPasswordLink(token string) string {
	return fmt.Sprintf("%s/password/new?token=%s", s.platformURL, token)
}

func (s *service) sendEmail(ctx context.Context, message []byte, toEmail string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", s.emailFrom, s.emailPassword, smtpHost)

	to := []string{toEmail}

	s.log.Debug().Sendf("Sending email with message: %s", message)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, s.emailFrom, to, message)
}

func (s *service) SetNewPassword(ctx context.Context, user *domain.User, password string) error {
	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.RecoveryToken = ""

	return s.userService.Update(ctx, user)

}

func (s *service) SetUserRecoveryToken(ctx context.Context, email string) (string, error) {

	user, err := s.userService.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", fmt.Errorf("invalid user")
	}

	user.RecoveryToken = s.GenerateToken()

	if err := s.userService.Update(ctx, user); err != nil {
		return "", err
	}

	return user.RecoveryToken, nil

}

func (s *service) SendResetPasswordLink(ctx context.Context, authUser *domain.AuthUser) {
	link := s.generateResetPasswordLink(authUser.RecoveryToken)

	msg := []byte("To: " + authUser.Email + "\r\n" +
		"Subject: Recover password - user-auth\r\n" +
		"\r\n" +
		"Reset password link. Copy it and paste it in the browser: \n" + link + "\r\n")

	if err := s.sendEmail(ctx, msg, authUser.Email); err != nil {
		s.log.Error().Err(err).Sendf("failed to send email")
		return
	}
	s.log.Info().Sendf("sent reset password link to %s", authUser.Email)
}
