package usecase

import (
	"auth-service/internal/domain"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	userRepo     domain.UserRepository
	tokenService domain.TokenService
}

func NewAuthUseCase(userRepo domain.UserRepository, tokenService domain.TokenService) domain.AuthUseCase {
	return &authUseCase{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (a *authUseCase) Register(username, email, password string) (*domain.User, string, error) {
	// check username exist
	if _, err := a.userRepo.FindByUsername(username); err == nil {
		return nil, "", errors.New("username already exist")
	}

	// check email exist
	if _, err := a.userRepo.FindByEmail(email); err == nil {
		return nil, "", errors.New("email already exist")
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &domain.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := a.userRepo.Create(user); err != nil {
		return nil, "", err
	}

	// generate token
	token, err := a.tokenService.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil

}

func (a *authUseCase) Login(username, password string) (*domain.User, string, error) {
	user, err := a.userRepo.FindByUsername(username)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// generate token
	token, err := a.tokenService.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (a *authUseCase) ValidateToken(token string) (*domain.User, error) {
	// check if token is blacklisted
	if a.tokenService.IsTokenBlacklisted(token) {
		return nil, errors.New("token is blacklisted")
	}

	// validate token
	userID, err := a.tokenService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// get user by id
	user, err := a.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (a *authUseCase) Logout(token string) error {
	// add token to blacklist
	return a.tokenService.BlacklistToken(token)
}
