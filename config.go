package main

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Configuration struct {
	URL string
	Token string
	From time.Time
	To time.Time
	Device int
	GPXFile string
	GPX ConfigGPX
}

type ConfigGPX struct {
	Title string
	Author string
}

func ParseConfig() {

	viper.SetConfigName("my_config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	pflag.String("from", "", "Start of date / time to export from")
	pflag.String("to", "", "End of date / time to export to")
	pflag.Int("device", 0, "Traccar Device ID")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// // see https://github.com/spf13/viper/issues/496#issuecomment-1152943336
	// tfrom := viper.GetTime("from")
	// viper.Set("from", tfrom)
	// tto := viper.GetTime("to")
	// viper.Set("to", tto)

	// var config Configuration
	// err = viper.Unmarshal(&config)
	// if err != nil {
	// 	fmt.Printf("Unable to decode into struct, %v", err)
	// }
	// viper.BindPFlag("config", )
	// pflag.Int("from", nil, "from")

}