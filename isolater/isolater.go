package isolater

import (
	"beagle/logger"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

var log = logger.New("cgroup-mover")

const (
	CGROUP_PATH  = "/sys/fs/cgroup"
	CGROUP_PROCS = "cgroup.procs"
)

// Isolate moves all processes owned by user `user` to a `cgroup` cgroup.
func Isolate(user, cgroup string) {
	pids := findUserProcesses(user)

	if err := addToCgroup(pids, path.Join(CGROUP_PATH, cgroup, CGROUP_PROCS)); err != nil {
		log.Errorf("Error trying to add pids to cgroup (%s)", cgroup)
	} else {
		for _, p := range pids {
			log.Infof("%d -> %s", p, cgroup)
		}
	}
}

func findUserProcesses(user string) (childPids []int) {
	cmd := exec.Command("pgrep", "--uid", user)
	out, err := cmd.Output()
	if err != nil {
		log.Warnf("Problems finding processes for user %s", user)
		log.Error(err)
	}
	pidsRaw := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, p := range pidsRaw {
		i, _ := strconv.Atoi(p)
		childPids = append(childPids, i)
	}
	return
}

func addToCgroup(pids []int, cgroupProcsFile string) error {
	f, err := os.OpenFile(cgroupProcsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Error(err)
		return err
	}
	defer f.Close()

	for _, pid := range pids {
		if _, err := f.WriteString(fmt.Sprintf("%d\n", pid)); err != nil {
			log.Errorw("Couldn't write pid to the groupc.procs file", "pid", pid, "err", err.Error())
			return err
		}
	}
	return nil
}

func pidExists(pid int, pids []int) bool {
	for _, p := range pids {
		if p == pid {
			return true
		}
	}
	return false
}
