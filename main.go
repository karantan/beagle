package main

import (
	"beagle/isolater"
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
	var interval, maxRunning time.Duration
	var processFilter, slackChan, cgroup string

	flag.DurationVar(&interval, "interval", 5*time.Minute, "Frequency of checking for running processes")
	flag.DurationVar(&maxRunning, "max-running", 20*time.Minute, "Max allowed running time for reporting it as a long-running process")
	flag.StringVar(&processFilter, "filter", "php-fpm: pool", "Process filter (i.e. '| grep <filter>')")
	flag.StringVar(&slackChan, "slack", "", "Slack Channel name where the notification will be posted")
	flag.StringVar(&cgroup, "cgroup", "", "Cgroup TO which long running processes will be moved.")
	flag.Parse()

	slackHook := os.Getenv("SLACK_NOTIFICATION")
	hostname, _ := os.Hostname()

	ticker := time.NewTicker(interval)
	log.Info("Checking for processes every ", interval)
	for ; true; <-ticker.C {
		procs := spotter.FindOldProcesses(processFilter, int(maxRunning.Seconds()))
		for _, p := range procs {
			log.Info(p)

			body := fmt.Sprintf("🐶 [%s] I spotted PHP-FPM pool for user %s has been running for more than %s", hostname, p.Owner, p.Duration.Round(time.Second))
			if cgroup != "" {
				err := isolater.Isolate(p.PID, cgroup)
				if err != nil {
					body = fmt.Sprintf("🐶 [%s] I failed trying to isolate %d process (%s) because it has been running for more than %s", hostname, p.PID, p.Owner, p.Duration.Round(time.Second))
				} else {
					body = fmt.Sprintf("🐶 [%s] I isolated %d process (%s) because it has been running for more than %s", hostname, p.PID, p.Owner, p.Duration.Round(time.Second))
				}
			}
			if slackHook != "" {
				msg := &slack.WebhookMessage{Text: body, Channel: slackChan}
				if err := slack.PostWebhook(slackHook, msg); err != nil {
					log.Error(err)
				}
			}

		}
	}
}
