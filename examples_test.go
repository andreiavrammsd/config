package config_test

import (
	"fmt"
	"log"
	"os"

	"github.com/andreiavrammsd/config"
)

type Configuration struct {
	Username string `env:"USERNAME"`
	Tag      string `env:"TAG"      default:"none"`
}

func ExampleConfig_FromEnv() {
	if err := os.Setenv("USERNAME", "msd"); err != nil {
		log.Fatal(err)
	}

	configuration := Configuration{}

	if err := config.New().FromEnv(&configuration); err != nil {
		log.Fatalf("cannot load config: %s", err)
	}

	fmt.Println(configuration.Username)
	fmt.Println(configuration.Tag)

	// Output:
	// msd
	// none
}

func ExampleConfig_FromBytes() {
	configuration := Configuration{}
	input := []byte(`USERNAME=msd # username`)

	c := config.New()

	if err := c.FromBytes(&configuration, input); err != nil {
		log.Fatalf("cannot load config: %s", err)
	}

	fmt.Println(configuration.Username)

	// Output:
	// msd
}
