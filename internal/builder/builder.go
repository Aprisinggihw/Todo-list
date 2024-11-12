package builder

import (
	"todo-list/configs"
	"todo-list/internal/http/handler"
	"todo-list/internal/http/router"
	"todo-list/internal/repository"
	"todo-list/internal/service"
	"todo-list/pkg/cache"
	"todo-list/pkg/route"
	"todo-list/pkg/token"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func BuildPublicRoutes(cfg *configs.Config, db *gorm.DB, rdb *redis.Client) []route.Route {
	cacheable := cache.NewCacheable(rdb)
	userRepository := repository.NewUserRepository(db)
	tokenUseCase := token.NewTokenUseCase(cfg.JWT.SecretKey)
	userService := service.NewUserService(userRepository, tokenUseCase, cacheable)
	userHandler := handler.NewUserHandler(userService)
	return router.PublicRoutes(userHandler)
}

func BuildPrivateRoutes(cfg *configs.Config, db *gorm.DB, rdb *redis.Client) []route.Route {
	cacheable := cache.NewCacheable(rdb)
	userRepository := repository.NewUserRepository(db)
	tokenUseCase := token.NewTokenUseCase(cfg.JWT.SecretKey)
	userService := service.NewUserService(userRepository, tokenUseCase, cacheable)
	userHandler := handler.NewUserHandler(userService)
	todoRepository := repository.NewTodoRepository(db)
	todoService := service.NewTodoService(todoRepository, tokenUseCase, cacheable)
	todoHandler := handler.NewTodoHandler(todoService)
	return router.PrivateRoutes(userHandler,*todoHandler)
}
