package main

import (
	"fmt"
	"os"
	"warchest/src/config"
)

const WarchestConfigEnv = "WARCHEST_CONFIG"

func main() {

	configPath := os.Getenv(WarchestConfigEnv)
	warchestConfig, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Failed loading config due to: %s", err)
	}

}
