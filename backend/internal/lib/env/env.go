package env

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// GetString получает строковое значение из окружения или использует fallback.
func GetString(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		log.Printf("Environment variable %s set to: %s", key, value)
		return value
	}
	log.Printf("Environment variable %s not set, using fallback value: %s", key, fallback)
	return fallback
}

// GetInt получает целочисленное значение из окружения или использует fallback.
func GetInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Invalid integer for %s: %v, using default value %d", key, err, fallback)
			return fallback
		}
		log.Printf("Environment variable %s set to: %d", key, intValue)
		return intValue
	}
	log.Printf("Environment variable %s not set, using fallback value: %d", key, fallback)
	return fallback
}

// GetDuration получает значение типа time.Duration из окружения или использует fallback.
func GetDuration(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		duration, err := time.ParseDuration(value)
		if err != nil {
			log.Printf("Invalid duration for %s: %v, using default value %v", key, err, fallback)
			return fallback
		}
		log.Printf("Environment variable %s set to: %v", key, duration)
		return duration
	}
	log.Printf("Environment variable %s not set, using fallback value: %v", key, fallback)
	return fallback
}

// GetList получает список строк из окружения или использует fallback.
func GetList(key string, separator string, fallback []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		items := strings.Split(value, separator)
		for i := range items {
			items[i] = strings.TrimSpace(items[i])
		}
		log.Printf("Environment variable %s set to: %v", key, items)
		return items
	}
	log.Printf("Environment variable %s not set, using fallback value: %v", key, fallback)
	return fallback
}
