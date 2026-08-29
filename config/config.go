package config

type Config struct {
	DB        DBConfig
	JwtSecret string
	Auth      Auth
}

type DBConfig struct {
	Host     string `koanf:"host"`
	Port     string `koanf:"port"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	Database string `koanf:"database"`
}

type Auth struct {
	TokenExpiresMinute uint `koanf:"token_expires_minute"`
	TokenRefreshMinute uint `koanf:"token_refresh_minute"`
}
