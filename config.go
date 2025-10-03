package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadFromJson() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	fmt.Println("membaca config file")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
