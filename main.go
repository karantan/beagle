package main

import (
	"beagle/isolater"
	"beagle/logger"
	"beagle/spotter"
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/slack-go/slack"
)

var log = logger.New("beagle")

func main() {
	var interval time.Duration
	var processFilter, slackChan, cgroup string

	flag.DurationVar(&interval, "interval", 3*time.Minute, "Frequency of checking for running processes. It's also the max allowed running time for reporting it as a long-running process")
	flag.StringVar(&processFilter, "filter", "php-fpm: pool", "Process filter (i.e. '| grep <filter>')")
	flag.StringVar(&slackChan, "slack", "", "Slack Channel name where the notification will be posted")
	flag.StringVar(&cgroup, "cgroup", "", "Cgroup TO which long running processes will be moved.")
	flag.Parse()

	if cgroup != "" {
		if _, err := os.Stat(path.Join(isolater.CGROUP_PATH, cgroup)); os.IsNotExist(err) {
			log.Panicf("Cgroup %s doesn't exist", path.Join(isolater.CGROUP_PATH, cgroup))
		}
	}

	slackHook := os.Getenv("SLACK_NOTIFICATION")

	ticker := time.NewTicker(interval)
	log.Info("Checking for processes every ", interval)
	for ; true; <-ticker.C {
		procs := spotter.FindOldProcesses(processFilter, int(interval.Seconds()))
		for _, p := range procs {
			log.Info(p)
			msg := &slack.WebhookMessage{
				Text:    fmt.Sprintf("🐶 I isolated %d process (%s) because it has been running for more than %s", p.PID, p.Owner, p.Duration.Round(time.Second)),
				Channel: slackChan,
			}

			if cgroup != "" {
				isolater.Isolate(p.PID, cgroup)
				msg = &slack.WebhookMessage{
					Text:    fmt.Sprintf("🐶 I spotted PHP-FPM pool for user %s has been running for more than %s", p.Owner, p.Duration.Round(time.Second)),
					Channel: slackChan,
				}
			}
			if slackHook != "" {
				if err := slack.PostWebhook(slackHook, msg); err != nil {
					log.Error(err)
				}
			}

		}
	}
}
