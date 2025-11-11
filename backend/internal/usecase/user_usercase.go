package usecase

import (
	"errors"

	"github.com/prabalesh/loco/backend/internal/domain"
	"github.com/prabalesh/loco/backend/internal/infrastructure/auth"
	"github.com/prabalesh/loco/backend/internal/usecase/interfaces"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	repo       interfaces.UserRepository
	jwtService *auth.JWTService
}

func NewUserUsecase(repo interfaces.UserRepository, jwtService *auth.JWTService) *UserUsecase {
	return &UserUsecase{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (u *UserUsecase) Register(req *domain.RegisterRequest) error {
	// Check if email already exists
	existing, err := u.repo.GetByEmail(req.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("email already registered")
	}

	// Check if username already exists
	existingUser, err := u.repo.GetByUsername(req.Username)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("username already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         "user",
	}

	return u.repo.Create(user)
}

func (u *UserUsecase) Login(req *domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := u.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := u.jwtService.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (u *UserUsecase) GetUserByID(id int) (*domain.User, error) {
	return u.repo.GetByID(id)
}
