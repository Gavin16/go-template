package conf

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// load config from json
func init() {
	log.SetPrefix("[go-template]")
	log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading common config file: %v ", err)
	}

	env := viper.GetString("ENV")
	if env == "" {
		env = "local"
	}

	viper.SetConfigName(fmt.Sprintf("config-%s", env))

	if err := viper.MergeInConfig(); err != nil {
		log.Fatalf("Error reading environment specific config file: %v", err)
	}

	log.Printf(">>> ... json config loaded successfully ...")
}
