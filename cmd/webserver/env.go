package main

import (
	"errors"
	"os"
	"strings"
)

func validateEnvs() error {
	ers := []string{}

	if os.Getenv("JWT_SECRET_KEY") == "" {
		ers = append(ers, "JWT_SECRET_KEY is not provided")
	}

	if len(ers) != 0 {
		return errors.New(strings.Join(append([]string{"ERROR: Environment variables failed validation"}, ers...), "\n - "))
	}

	return nil
}
