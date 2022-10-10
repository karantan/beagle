package main

import (
	"beagle/spotter"
	"flag"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ProcessFilter string
	CGroup        string
	SlackHook     string
}

func init() {
	var cfgFile string
	flag.StringVar(&cfgFile, "config file", "config.yaml", "Beagle configuration file.")
	viper.SetConfigFile(cfgFile)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func main() {
	var tickerDuration time.Duration
	flag.DurationVar(&tickerDuration, "tickerDuration", 5*time.Second, "How frequently should I check for long-running processes?")
	flag.Parse()

	processFilter := viper.GetString("process-filter")
	maxTime := viper.GetInt("max-time")

	for range time.Tick(tickerDuration) {
		spotter.FindOldProcesses(processFilter, maxTime)
	}
}
