package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const EnvKey string = "ENV"

var configs = make(map[string]interface{})

func init() {

	configFileName := "config"
	env, profileSet := os.LookupEnv(EnvKey)
	if profileSet {
		slog.Info("Setting active profile to " + env)
		configFileName = configFileName + "-" + env
	}

	if len(env) == 0 {
		slog.Info("Setting ENV to NONE")
	}

	viper.SetConfigName(configFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("Error reading config file", "error", err)
		panic(-1)
	}

	for _, k := range viper.AllKeys() {
		value := viper.GetString(k)
		if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
			val := getEnvOrPanic(strings.TrimSuffix(strings.TrimPrefix(value, "${"), "}"))
			viper.Set(k, val)
		}
	}
	err := viper.Unmarshal(&configs)
	if err != nil {
		slog.Error("Error creating configs", "error", err)
		panic(-1)
	}

	slog.Debug("Configuration loaded", "file", viper.ConfigFileUsed())
	slog.Info("Configuration loaded successfully")

}

func getEnvOrPanic(env string) string {
	res, set := os.LookupEnv(env)
	if !set {
		slog.Error("Config load failed")
		panic("Mandatory env variable not found:" + env)
	}
	return res
}

func GetConfig[T int | bool | string](key string) (*T, error) {
	return getRecursive[T](key, configs)
}

func getRecursive[T int | bool | string](key string, scope map[string]interface{}) (*T, error) {
	// Support nested keys using dot notation
	if before, after, found := strings.Cut(key, "."); found {
		nested, ok := scope[before].(map[string]interface{})
		if !ok {
			configError := errors.New("Config type assertion failed for key:" + before)
			slog.Error("Error fetching config", "error", configError)
			return nil, configError
		}
		return getRecursive[T](after, nested)
	}
	val, exists := scope[key]
	if !exists {
		configError := errors.New("Config key not found:" + key)
		slog.Error("Error fetching config", "error", configError)
		return nil, configError
	}
	typedVal, ok := val.(T)
	if !ok {
		typeString := fmt.Sprintf("%T", *new(T))
		configError := errors.New("Config type assertion failed for key:" + key +
			" type: " + typeString +
			" expected: " + fmt.Sprintf("%T", val))
		slog.Error("Error fetching config", "error", configError)
		return nil, configError
	}
	return &typedVal, nil
}
