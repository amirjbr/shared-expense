package config

type Config struct {
	DB        DBConfig `koanf:"db"`
	JwtSecret string   `koanf:"jwt_secret"`
	Auth      Auth     `koanf:"auth"`
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
