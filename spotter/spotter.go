package spotter

import (
	"beagle/logger"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var log = logger.New("spotter")

type Process struct {
	PID      int
	Owner    string
	Duration time.Duration
}

// FindOldProcesses finds long running processes (running longer than `seconds`).
// Under the hood it runs this command:
// ps -eo etimes,pid,args | grep "<filter>" | awk '{if ($1 > <seconds>) print $0 }'
func FindOldProcesses(processFilter string, seconds int) []Process {
	cmdPS := exec.Command("ps", "-eo", "etimes,pid,args")
	cmdGrep := exec.Command("grep", processFilter)
	cmdAwk := exec.Command("awk", fmt.Sprintf("{if ($1 > %d) print $0 }", seconds))

	var out bytes.Buffer

	cmdGrep.Stdin, _ = cmdPS.StdoutPipe()
	cmdAwk.Stdin, _ = cmdGrep.StdoutPipe()
	cmdAwk.Stdout = &out

	// Start grep and awk processes, but wait for input
	cmdGrep.Start()
	cmdAwk.Start()

	// Start the ps process and send the output to grep
	cmdPS.Run()

	// process the ps ouput (pipe it to grep and awk)
	cmdGrep.Wait()
	cmdAwk.Wait()

	return parseRawProcesses(out.String())
}

// parseRawProcesses processes line
// `      4 1037007 php-fpm: pool foo_com
//        2 1037246 php-fpm: pool bar_com
// ` into Process struct
func parseRawProcesses(lines string) (procs []Process) {
	trimmed := strings.TrimSpace(lines)
	if trimmed == "" {
		return
	}
	pidsRaw := strings.Split(trimmed, "\n")
	for _, line := range pidsRaw {
		chunks := strings.Split(strings.TrimSpace(line), " ")

		duration, _ := strconv.Atoi(chunks[0])
		pid, _ := strconv.Atoi(chunks[1])
		procs = append(procs, Process{
			Duration: time.Duration(duration) * time.Second,
			PID:      pid,
			Owner:    chunks[len(chunks)-1],
		})
	}
	return
}
