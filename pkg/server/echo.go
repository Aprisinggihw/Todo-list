package server

import (
	"todo-list/configs"
	"todo-list/pkg/response"
	"todo-list/pkg/route"
	"todo-list/pkg/token"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Server struct {
	*echo.Echo
}

func NewServer(cfg *configs.Config,
	publicRoutes, privateRoutes []route.Route) *Server {
	e := echo.New()
	e.HideBanner = true

	v1 := e.Group("/api/v1")

	if len(publicRoutes) > 0 {
		for _, route := range publicRoutes {
			v1.Add(route.Method, route.Path, route.Handler)
		}
	}

	if len(privateRoutes) > 0 {
		for _, route := range privateRoutes {
			v1.Add(route.Method, route.Path, route.Handler, JWTMiddleware(cfg.JWT.SecretKey), RBACMiddleware(route.Roles))
		}
	}
	return &Server{e}
}

func JWTMiddleware(secretKey string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(token.JwtCustomClaims)
		},
		SigningKey: []byte(secretKey),
		ErrorHandler: func(ctx echo.Context, err error) error {
			return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, "anda harus login untuk megakses resource ini."))
		},
	})
}

func RBACMiddleware(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			user := ctx.Get("user").(*jwt.Token)
			claims := user.Claims.(*token.JwtCustomClaims)
			
			if claims.Role != "admin"{
				// Simpan data user_id ke context
				ctx.Set("user_id", claims.UserID)
			}

			allowed := false
			for _, role := range roles {
				if role == claims.Role {
					allowed = true
					break
				}
			}
			
			if !allowed {
				return ctx.JSON(http.StatusForbidden, response.ErrorResponse(http.StatusForbidden, "anda tidak diizinkan untuk mengakses resource ini."))
			}
			
			return next(ctx)
		}
	}
}
