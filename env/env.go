package env

import (
	"fmt"
	"os"
)

const (
	EnvMotorAPTAddr = "MOTOR_APT_ADDR"
)

func GetMotorAptAddr() (string, error) {
	return notEmptyEnvValue(EnvMotorAPTAddr)
}

func noEnvError(env string) error {
	return fmt.Errorf("la variable de entorno \"%v\" no fue establecida", env)
}

func emptyEnvError(env string) error {
	return fmt.Errorf("la variable de entorno \"%v\" está vacía", env)
}

func notEmptyEnvValue(env string) (string, error) {
	val, ok := os.LookupEnv(env)
	if !ok {
		return val, noEnvError(env)
	}
	if len(val) == 0 {
		return val, emptyEnvError(env)
	}
	return val, nil
}

func envValue(env string) (string, error) {
	path, ok := os.LookupEnv(env)
	if !ok {
		return path, noEnvError(env)
	}

	return path, nil
}
