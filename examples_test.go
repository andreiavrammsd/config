package config_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/andreiavrammsd/config"
)

type Configuration struct {
	Username string `env:"USERNAME"`
	Tag      string `env:"TAG"      default:"none"`
	Timeout  int
}

func ExampleConfig_FromFile() {
	configuration := Configuration{}

	if err := config.New().FromFile(&configuration, "testdata/.env", "testdata/.example"); err != nil {
		log.Fatalf("cannot parse config: %s", err)
	}

	fmt.Println(configuration.Username)
	fmt.Println(configuration.Tag)
	fmt.Println(configuration.Timeout)

	// Output:
	// msd
	// none
	// 2000000000
}

func ExampleConfig_FromEnv() {
	if err := os.Setenv("USERNAME", "msd"); err != nil {
		log.Fatal(err)
	}

	configuration := Configuration{}

	if err := config.New().FromEnv(&configuration); err != nil {
		log.Fatalf("cannot parse config: %s", err)
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
		log.Fatalf("cannot parse config: %s", err)
	}

	fmt.Println(configuration.Username)

	// Output:
	// msd
}

func ExampleConfig_FromJSON() {
	configuration := Configuration{}
	input := json.RawMessage(`{"USERNAME": "msd"}`)

	c := config.New()

	if err := c.FromJSON(&configuration, input); err != nil {
		log.Fatalf("cannot parse config: %s", err)
	}

	fmt.Println(configuration.Username)

	// Output:
	// msd
}
