package repository

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	MySQL struct {
		User     string
		Password string
		Host     string
		Db       string
	}
	FireStore struct {
		CredintialPath string
	}
}

var config *Config

func init() {
	config = &Config{}
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	// viper.AddConfigPath("/etc/appname")  // path to look for the config file in
	// viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(".")                // optionally look for config in the working directory
	viper.AddConfigPath("./repository")     // optionally look for config in the working directory
	viper.AddConfigPath("../../repository") // optionally look for config in the working directory

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("fail to decode config: ", err.Error())
	}
}
