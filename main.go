package main

import (
	"beagle/isolater"
	"beagle/logger"
	"beagle/spotter"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/slack-go/slack"
)

var log = logger.New("beagle")
var store = StatsRegistry{isolatedProcesses: make(map[int]bool)}

type StatsRegistry struct {
	sync.Mutex
	isolatedProcesses map[int]bool
}

func main() {
	var interval, maxRunning, flush time.Duration
	var processFilter, slackChan, cgroup string

	flag.DurationVar(&interval, "interval", 5*time.Minute, "Frequency of checking for running processes")
	flag.DurationVar(&maxRunning, "max-running", 20*time.Minute, "Max allowed running time for reporting it as a long-running process")
	flag.DurationVar(&flush, "memory-flush", 30*time.Minute, "Maximum frequency of slack notifications. By default a message will be sent to slack every 30min even if --max-running or --interval are set to lower")
	flag.StringVar(&processFilter, "filter", "php-fpm: pool", "Process filter (i.e. '| grep <filter>')")
	flag.StringVar(&slackChan, "slack", "", "Slack Channel name where the notification will be posted")
	flag.StringVar(&cgroup, "cgroup", "", "Cgroup TO which long running processes will be moved.")
	flag.Parse()

	slackHook := os.Getenv("SLACK_NOTIFICATION")
	hostname, _ := os.Hostname()

	go FlushMemory(flush)

	ticker := time.NewTicker(interval)
	log.Info("Checking for processes every ", interval)
	for ; true; <-ticker.C {
		procs := spotter.FindOldProcesses(processFilter, int(maxRunning.Seconds()))
		for _, p := range procs {
			_, alreadySeen := store.isolatedProcesses[p.PID]
			if alreadySeen {
				log.Infow("Skipping", "Process", p)
				continue
			}
			log.Info(p)

			body := fmt.Sprintf("🐶 [%s] I spotted PHP-FPM pool for user %s has been running for more than %s", hostname, p.Owner, p.Duration.Round(time.Second))
			if cgroup != "" {
				err := isolater.Isolate(p.PID, cgroup)
				if err != nil {
					body = fmt.Sprintf("🐶 [%s] I failed trying to isolate %d process (%s) because it has been running for more than %s", hostname, p.PID, p.Owner, p.Duration.Round(time.Second))
				} else {
					body = fmt.Sprintf("🐶 [%s] I isolated %d process (%s) because it has been running for more than %s", hostname, p.PID, p.Owner, p.Duration.Round(time.Second))
					store.isolatedProcesses[p.PID] = true
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

// FlushMemory flushes temporary stored isolated processes in memory map every hour
func FlushMemory(t time.Duration) {
	for range time.Tick(t) {
		store.Lock()
		log.Info("Memory flushed")
		store.isolatedProcesses = make(map[int]bool)
		store.Unlock()
	}
}
