package config

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	ServerPort string
}

func LoadConfig() *Config {
	return &Config{
		DBHost:     "localhost",
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "spagetiX",
		DBPort:     "5432",
		ServerPort: "8080",
	}
}
