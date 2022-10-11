package main

import (
	"beagle/logger"
	"beagle/spotter"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

var log = logger.New("beagle")

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
	flag.DurationVar(&tickerDuration, "tickerDuration", 5*time.Minute, "How frequently should I check for long-running processes?")
	flag.Parse()

	processFilter := viper.GetString("process-filter")
	maxTime := viper.GetInt("max-time")
	slackChan := viper.GetString("slack-channel")
	slackHook := os.Getenv("SLACK_NOTIFICATION")

	for range time.Tick(tickerDuration) {
		procs := spotter.FindOldProcesses(processFilter, maxTime)
		for _, p := range procs {
			log.Info(p)
			if slackHook != "" {
				msg := &slack.WebhookMessage{
					Text:    fmt.Sprintf("🐶: I spotted PHP-FPM pool for user %s has been running for more than %s", p.Owner, p.Duration.Round(time.Second)),
					Channel: slackChan,
				}
				if err := slack.PostWebhook(slackHook, msg); err != nil {
					log.Error(err)
				}
			}
		}
	}
}
