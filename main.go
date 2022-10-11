package main

import (
	"beagle/logger"
	"beagle/spotter"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/slack-go/slack"
)

var log = logger.New("beagle")

func main() {
	var interval time.Duration
	var processFilter, slackChan string

	flag.DurationVar(&interval, "interval", 3*time.Minute, "Frequency of checking for running processes. It's also the max allowed running time for reporting it as a long-running process")
	flag.StringVar(&processFilter, "filter", "php-fpm: pool", "Process filter (i.e. '| grep <filter>')")
	flag.StringVar(&slackChan, "slack", "", "Slack Channel name where the notification will be posted")
	flag.Parse()

	slackHook := os.Getenv("SLACK_NOTIFICATION")

	log.Info("Checking for processes every ", interval)
	for range time.Tick(interval) {
		procs := spotter.FindOldProcesses(processFilter, int(interval.Seconds()))
		for _, p := range procs {
			log.Info(p)
			if slackHook != "" {
				msg := &slack.WebhookMessage{
					Text:    fmt.Sprintf("🐶 I spotted PHP-FPM pool for user %s has been running for more than %s", p.Owner, p.Duration.Round(time.Second)),
					Channel: slackChan,
				}
				if err := slack.PostWebhook(slackHook, msg); err != nil {
					log.Error(err)
				}
			}
		}
	}
}
