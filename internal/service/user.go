package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"
	"todo-list/internal/entity"
	"todo-list/internal/repository"
	"todo-list/pkg/cache"
	"todo-list/pkg/token"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FindAll(ctx context.Context) ([]entity.User, error)
	Register(ctx context.Context, req *entity.UserReg) error
	Login(ctx context.Context, username, password string) (string, error)
}

type userService struct {
	userRepository repository.UserRepository
	tokenUseCase   token.TokenUseCase
	cacheable      cache.Cacheable
}

func NewUserService(
	userRepository repository.UserRepository,
	tokenUseCase token.TokenUseCase,
	cacheable cache.Cacheable,
) UserService {
	return &userService{userRepository, tokenUseCase, cacheable}
}

func (s *userService) FindAll(ctx context.Context) (result []entity.User, err error) {
	keyFindAll := "todo-list:users:find-all"
	data := s.cacheable.Get(keyFindAll)
	if data == "" {
		result, err = s.userRepository.FindAll(ctx)
		if err != nil {
			return nil, err
		}

		marshalledData, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}

		err = s.cacheable.Set(keyFindAll, marshalledData, 5*time.Minute)
		if err != nil {
			return nil, err
		}
	} else {
		err = json.Unmarshal([]byte(data), &result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// Logika registrasi user
func (s *userService) Register(ctx context.Context, req *entity.UserReg) error {
	// Periksa apakah username sudah ada
	_, err := s.userRepository.FindByUsername(ctx, req.Username)
	if err == nil {
		return errors.New("username already exists")
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}
	req.Password = string(hashedPassword)
	
	return s.userRepository.CreateUser(ctx, req)
}

func (s *userService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepository.FindByUsername(ctx, username)
	if err != nil {
		log.Println(err.Error())
		return "", errors.New("username or password invalid")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("username or password invalid")
	}

	expiredTime := time.Now().Local().Add(time.Minute * 5)

	claims := token.JwtCustomClaims{
		UserID:   uint(user.ID),
		Username: user.Username,
		Role:     user.Role,
		FullName: user.FullName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "todo-list",
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}

	token, err := s.tokenUseCase.GenerateAccessToken(claims)
	if err != nil {
		return "", errors.New("ada kesalahan di server")
	}
	return token, nil
}
