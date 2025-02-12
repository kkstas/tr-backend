package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/kkstas/tr-backend/internal/config"
)

type cfg struct {
	port   string
	dbName string
}

func getConfigs(getenv func(string) string) (*cfg, *config.Config, error) {
	errs := []string{}

	port := getenv("PORT")
	if port == "" {
		errs = append(errs, "PORT (number) is not defined")
	} else if _, err := strconv.Atoi(getenv("PORT")); err != nil {
		errs = append(errs, "PORT is not a valid number")
	}

	jwtSecretKey := getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		errs = append(errs, "JWT_SECRET_KEY (string) is not provided")
	}

	var enableRegister bool
	if getenv("ENABLE_REGISTER") == "true" {
		enableRegister = true
	}

	dbName := getenv("DB_NAME")
	if dbName == "" {
		errs = append(errs, "DB_NAME (string) is not defined")
	}

	if len(errs) != 0 {
		return nil, nil, errors.New(strings.Join(append([]string{"ERROR: Environment variables failed validation"}, errs...), "\n - "))
	}

	return &cfg{
			port:   port,
			dbName: dbName,
		},
		&config.Config{
			JWTSecretKey:   []byte(jwtSecretKey),
			EnableRegister: enableRegister,
		},
		nil
}
