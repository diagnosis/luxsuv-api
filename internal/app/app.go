package app

import (
	"os"
	"time"

	"github.com/diagnosis/luxsuv-api-v2/internal/api"
	"github.com/diagnosis/luxsuv-api-v2/internal/logger"
	"github.com/diagnosis/luxsuv-api-v2/internal/secure"
	"github.com/diagnosis/luxsuv-api-v2/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	DB                         *pgxpool.Pool
	ServerHealthCheckerHandler *api.ServerHealthCheckerHandler
	UserHandler                *api.UserHandler
}

func NewApplication(pool *pgxpool.Pool) (*Application, error) {
	logger.Info(nil, "initializing application")

	serverHealthCheckerHandler := api.NewServerHealthCheckerHandler()

	userStore := store.NewPostgresUserStore(pool)
	refreshStore := store.NewPostgresRefreshTokenStore(pool)

	accessSecret := []byte(os.Getenv("JWT_ACCESS_SECRET"))
	refreshSecret := []byte(os.Getenv("JWT_REFRESH_SECRET"))
	issuer := os.Getenv("JWT_ISSUER")
	audience := os.Getenv("JWT_AUDIENCE")

	if len(accessSecret) < 32 || len(refreshSecret) < 32 {
		logger.Error(nil, "JWT secrets are invalid", "access_len", len(accessSecret), "refresh_len", len(refreshSecret))
		return nil, secure.ErrSecretsInvalid
	}

	signer, err := secure.NewSigner(
		issuer,
		audience,
		accessSecret,
		refreshSecret,
		15*time.Minute,
		7*24*time.Hour,
	)
	if err != nil {
		logger.Error(nil, "failed to create JWT signer", "error", err)
		return nil, err
	}

	logger.Info(nil, "JWT signer initialized", "issuer", issuer, "audience", audience)

	userHandler := api.NewUserHandler()
	userHandler.UserStore = userStore
	userHandler.RefreshStore = refreshStore
	userHandler.Signer = *signer

	logger.Info(nil, "application initialized successfully")

	return &Application{
		DB:                         pool,
		ServerHealthCheckerHandler: serverHealthCheckerHandler,
		UserHandler:                userHandler,
	}, nil
}
