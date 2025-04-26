package config_test

import (
	"fmt"
	"log"
	"os"

	"github.com/andreiavrammsd/config"
)

type Configuration struct {
	Username string `env:"USERNAME"`
	Tag      string `env:"TAG" default:"none"`
}

func ExampleLoader_Env() {
	if err := os.Setenv("USERNAME", "msd"); err != nil {
		log.Fatal(err)
	}

	cfg := Configuration{}
	if err := config.Load(&cfg).Env(); err != nil {
		log.Fatalf("cannot load config: %s", err)
	}

	fmt.Println(cfg.Username)
	fmt.Println(cfg.Tag)

	// Output:
	// msd
	// none
}

func ExampleLoader_Bytes() {
	input := []byte(`USERNAME=msd # username`)

	cfg := Configuration{}
	if err := config.Load(&cfg).Bytes(input); err != nil {
		log.Fatalf("cannot load config: %s", err)
	}

	fmt.Println(cfg.Username)

	// Output:
	// msd
}
