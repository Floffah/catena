package environment

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Mode               string        `env:"ENVIRONMENT" envDefault:"development"`
	Port               string        `env:"PORT" envDefault:"8080"`
	DatabaseURL        string        `env:"DATABASE_URL,required"`
	ClerkSecretKey     string        `env:"CLERK_SECRET_KEY,required"`
	CORSAllowedOrigins []string      `env:"CORS_ALLOWED_ORIGINS" envDefault:"http://localhost:3000,http://127.0.0.1:3000"`
	CatenaGitRoot      string        `env:"CATENA_GIT_ROOT" envDefault:"/var/lib/catena/git"`
	GitTreeMaxEntries  int           `env:"CATENA_GIT_TREE_MAX_ENTRIES" envDefault:"50000"`
	GitTreeMaxBytes    int64         `env:"CATENA_GIT_TREE_MAX_BYTES" envDefault:"20971520"`
	GitTreeTimeout     time.Duration `env:"CATENA_GIT_TREE_TIMEOUT" envDefault:"10s"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
