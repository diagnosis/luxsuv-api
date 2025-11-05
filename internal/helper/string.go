package helper

import (
	"os"

	"github.com/google/uuid"
)

func getEnvDef(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func DeferOrString(p *string, def string) string {
	if p != nil && *p != "" {
		return *p
	}
	return def
}

func GenerateID() string {
	return uuid.NewString()
}
