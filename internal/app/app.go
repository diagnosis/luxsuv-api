package app

import (
	"github.com/diagnosis/luxsuv-api-v2/internal/api"
)

type Application struct {
	ServerHealthCheckerHandler *api.ServerHealthCheckerHandler
}

func NewApplication() *Application {
	serverHealthCheckerHandler := api.NewServerHealthCheckerHandler()
	return &Application{serverHealthCheckerHandler}
}
