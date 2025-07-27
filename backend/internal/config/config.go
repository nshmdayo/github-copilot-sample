package config

import (
	"os"
	"strconv"

	"github.com/spf13/viper"
)

// Config holds all configuration values
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	CORS     CORSConfig     `mapstructure:"cors"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret          string `mapstructure:"secret"`
	ExpirationHours int    `mapstructure:"expiration_hours"`
	RefreshHours    int    `mapstructure:"refresh_hours"`
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

// LoadConfig loads configuration from environment variables and config file
func LoadConfig() (*Config, error) {
	config := &Config{}

	// Set default values
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.username", "postgres")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.dbname", "todoapp")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expiration_hours", 24)
	viper.SetDefault("jwt.refresh_hours", 168)
	viper.SetDefault("cors.allow_origins", []string{"*"})
	viper.SetDefault("cors.allow_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allow_headers", []string{"Origin", "Content-Type", "Accept", "Authorization"})
	viper.SetDefault("cors.allow_credentials", true)

	// Read from environment variables
	viper.AutomaticEnv()

	// Override with environment variables if present
	if port := os.Getenv("PORT"); port != "" {
		viper.Set("server.port", port)
	}
	if host := os.Getenv("HOST"); host != "" {
		viper.Set("server.host", host)
	}
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		viper.Set("server.mode", mode)
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		viper.Set("database.host", dbHost)
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		if port, err := strconv.Atoi(dbPort); err == nil {
			viper.Set("database.port", port)
		}
	}
	if dbUser := os.Getenv("DB_USERNAME"); dbUser != "" {
		viper.Set("database.username", dbUser)
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		viper.Set("database.password", dbPassword)
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		viper.Set("database.dbname", dbName)
	}
	if sslMode := os.Getenv("DB_SSLMODE"); sslMode != "" {
		viper.Set("database.sslmode", sslMode)
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		viper.Set("jwt.secret", jwtSecret)
	}
	if jwtExp := os.Getenv("JWT_EXPIRATION_HOURS"); jwtExp != "" {
		if hours, err := strconv.Atoi(jwtExp); err == nil {
			viper.Set("jwt.expiration_hours", hours)
		}
	}

	// Unmarshal to struct
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

// GetDSN returns database connection string
func (c *DatabaseConfig) GetDSN() string {
	return "host=" + c.Host +
		" port=" + strconv.Itoa(c.Port) +
		" user=" + c.Username +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode
}
