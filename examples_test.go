package config_test

import (
	"fmt"
	"log"
	"os"

	"github.com/andreiavrammsd/config"
)

type Config struct {
	Username string `env:"USERNAME"`
}

func Example() {
	if err := os.Setenv("USERNAME", "msd"); err != nil {
		log.Fatal(err)
	}

	cfg := Config{}
	if err := config.Load(&cfg).Env(); err != nil {
		log.Fatalf("cannot load config: %s", err)
	}

	fmt.Println(cfg.Username)

	// Output:
	// msd
}
