package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

var Viper *viper.Viper

func InitViper() {

	// Initialize Viper
	v := viper.New()
	Viper = v
	LoadConfig()

}

func LoadConfig() {

	// Directly read the environment variable using os.Getenv
	env := os.Getenv("APP_ENV")
	if env == "" {
		log.Println("APP_ENV not set, defaulting to 'development'")
		env = "development"
	}

	if err := checkFileExists(env); err != nil {
		log.Println(err)
		log.Printf("Loaded From Enivronment Variables. APP_ENV: %s \n", env)
		Viper.AutomaticEnv()
	} else {
		log.Printf("File .env.%s exists\n", env)
		log.Printf("Loaded .env.%s file\n", env)

		// Set the configuration file based on the environment
		Viper.SetConfigFile(fmt.Sprintf(".env.%s", env))
		Viper.SetConfigType("dotenv")
		Viper.AddConfigPath("./")

		_ = Viper.ReadInConfig()
	}
}

func checkFileExists(env string) error {
	fileName := fmt.Sprintf(".env.%s", env)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist", fileName)
	}
	return nil
}
