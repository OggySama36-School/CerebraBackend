package config

import "os"

type Config struct {
	AllowedOrigins string
	Port           string
	PyURL          string
}

var App = &Config{}

func LoadConfig() {
	App.AllowedOrigins = os.Getenv("ALLOWED_ORIGINS")
	App.Port = os.Getenv("PORT")
	App.PyURL = os.Getenv("PYURL")
}
