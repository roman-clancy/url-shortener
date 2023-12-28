package config

import (
    "github.com/ilyakaznacheev/cleanenv"
    "log"
    "os"
    "time"
)

type Config struct {
    Env         string `yaml:"env"`
    StoragePath string `yaml:"storage_path" env-required:"true"`
    HttpServer  `yaml:"http_server"`
}

type HttpServer struct {
    Address     string        `yaml:"address"`
    IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
    Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
}

func MustLoad() *Config {
    var cfgPath, exists = os.LookupEnv("URL_SHORTENER_CFG_PATH")
    if !exists {
        log.Fatalf("Env variable %s doesn't set", "URL_SHORTENER_CFG_PATH")
    }

    if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
        log.Fatalf("File doesn't exists: %s", cfgPath)
    }
    var cfg Config

    if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
        log.Fatalf("Error while reading config file: %s", err)
    }

    return &cfg
}
