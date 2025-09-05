package config

import (
	"fmt"
	"log"
	"reflect"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func New[T any](path string) (*T, error) {
	// Load .env file if available
	err := godotenv.Load(path)
	if err != nil {
		log.Println("File .env not found, reading from environment")
	}

	// Configure viper
	viper.AutomaticEnv()

	// Create an instance of the type to inspect its fields
	var cfg T
	t := reflect.TypeOf(cfg)

	// Handle both struct types and pointers to struct types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Ensure we're working with a struct
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("config type must be a struct or pointer to struct")
	}

	// Bind each field with a mapstructure tag to its environment variable
	for i := range t.NumField() {
		field := t.Field(i)
		envVar := field.Tag.Get("mapstructure")
		if envVar != "" {
			err = viper.BindEnv(envVar, envVar)
			if err != nil {
				return nil, fmt.Errorf("failed to bind environment variable %s: %w", envVar, err)
			}
		}
	}

	// Unmarshal the configuration
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	return &cfg, nil
}
