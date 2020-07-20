package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/smtp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/evzpav/user-auth/internal/domain"
	"gitlab.com/evzpav/user-auth/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userService   domain.UserService
	emailFrom     string
	emailPassword string
	googleKey     string
	googleSecret  string
}

func NewService(userService domain.UserService, emailFrom, emailPassword, googleKey, googleSecret string) *service {
	return &service{
		userService:   userService,
		emailFrom:     emailFrom,
		emailPassword: emailPassword,
		googleKey:     googleKey,
		googleSecret:  googleSecret,
	}
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
	u.Token = token

	// u.Token, err = s.GenerateJWTToken(u)
	// if err != nil {
	// 	return errors.NewNotAuthorized(domain.ErrInvalidCredentials)
	// }

	return s.userService.Update(ctx, u)

}

func (s *service) Authenticate2(ctx context.Context, authUser *domain.AuthUser) (*domain.User, error) {
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

	token := s.generateToken()
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

func (s *service) Me(ctx context.Context) (*domain.User, error) {
	return nil, nil
}

func (s *service) AuthenticateToken(ctx context.Context, token string) (*domain.User, error) {
	user, err := s.userService.FindByToken(ctx, token)
	if err != nil {
		return nil, errors.NewNotAuthorized(domain.ErrInvalidCredentials)
	}
	if user == nil {
		return nil, errors.NewNotAuthorized(domain.ErrInvalidCredentials)
	}

	// claims, ok := s.validateToken(token)
	// if !ok {
	// 	return fmt.Errorf("invalid token")
	// }

	// id := claims["id"].(int)

	// user, err := s.userService.FindByID(ctx, id)
	// if err != nil {
	// 	return errors.NewNotAuthorized(domain.ErrInvalidCredentials)
	// }

	// if user == nil {
	// 	return errors.NewNotAuthorized(domain.ErrInvalidCredentials)
	// }

	// jwtToken, err := s.GenerateJWTToken(user)
	// if err != nil {
	// 	return errors.NewNotAuthorized(domain.ErrInvalidCredentials)
	// }

	user.Token = s.generateToken()

	return user, s.userService.Update(ctx, user)
}

// func (s *service) AuthenticateToken(ctx context.Context, token string) (*domain.User, error) {

// 	claims, ok := s.validateToken(token)
// 	if !ok {
// 		return nil, fmt.Errorf("invalid token")
// 	}

// 	id := claims["id"].(int)

// 	user, err := s.userService.FindByID(ctx, id)
// 	if err != nil {
// 		return nil, errors.NewNotAuthorized(domain.ErrInvalidCredentials)
// 	}

// 	if user == nil {
// 		return nil, errors.NewNotAuthorized(domain.ErrInvalidCredentials)
// 	}

// 	jwtToken, err := s.GenerateJWTToken(user)
// 	if err != nil {
// 		return nil, errors.NewNotAuthorized(domain.ErrInvalidCredentials)
// 	}

// 	user.Token = jwtToken

// 	if err := s.userService.Update(ctx, user); err != nil {
// 		return nil, err
// 	}

// 	return user, nil
// }

func (s *service) GenerateJWTToken(u *domain.User) (string, error) {
	expire := time.Now().Add(time.Hour * 1)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  u.ID,
		"e":   u.Email,
		"a":   u.Address,
		"p":   u.Phone,
		"exp": expire.Unix(),
	})

	tokenString, err := token.SignedString("key") //TODOD ADD VAR

	return tokenString, err
}

func (s *service) ParseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, fmt.Errorf("wrong signing method")
		}
		return "key", nil
	})

}

func (s *service) validateToken(token string) (map[string]interface{}, bool) {
	jwtToken, err := s.ParseToken(token)
	if err != nil || !jwtToken.Valid {
		return nil, false
	}
	return jwtToken.Claims.(jwt.MapClaims), true
}

func (s *service) SendEmail(ctx context.Context, message, toEmail string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", s.emailFrom, s.emailPassword, smtpHost)

	to := []string{toEmail}

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, s.emailFrom, to, []byte(message))
}

func (s *service) GoogleAuthentication(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore) {
	gothic.Store = store
	goth.UseProviders(
		google.New(s.googleKey, s.googleSecret, "http://localhost:5001/login/google/callback"),
	)

	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err == nil {
		fmt.Printf("gothUSER: %+v\n", gothUser)
		return
	}

	gothic.BeginAuthHandler(w, r)

}
