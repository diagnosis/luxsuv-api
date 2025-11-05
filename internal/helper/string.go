package helper

import (
	"os"
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
