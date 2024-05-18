package config

import (
	"bufio"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"os"
	"strings"
)

var (
	envPrefix = "lessons"
	envPath   = ".env"
)

// MustConfig returns AppConfig if all is well
func MustConfig(configStruct interface{}) interface{} {
	loadEnv(envPath)
	err := envconfig.Process(envPrefix, configStruct)
	if err != nil {
		panic(err)
	}
	return configStruct
}

// to make sure there's a colon at the port.
func processPort(port string) string {
	if port[0] == ':' {
		return port
	}
	return ":" + port
}

func loadEnv(dotEnvPath string) error {
	// Open the .env file
	file, err := os.Open(dotEnvPath)
	if err != nil {
		return fmt.Errorf("error opening .env file: %v", err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Split the line by '=' character
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			// Skip lines without a variable value
			continue
		}
		// Clean up the value from quotes, if any
		value := strings.Trim(parts[1], "\"")
		// Set the environment variable
		os.Setenv(parts[0], value)
	}

	// Check for any errors while scanning the file
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning .env file: %v", err)
	}

	return nil
}
