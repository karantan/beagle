package isolater

import (
	"beagle/logger"
	"fmt"
	"os"
	"path"
)

var log = logger.New("cgroup-mover")

const (
	CGROUP_PATH  = "/sys/fs/cgroup"
	CGROUP_PROCS = "cgroup.procs"
)

// Isolate moves process `pid` to a `cgroup` cgroup.
func Isolate(pid int, cgroup string) (err error) {
	if _, err = os.Stat(path.Join(CGROUP_PATH, cgroup)); os.IsNotExist(err) {
		log.Errorf("Cgroup %s doesn't exist", path.Join(CGROUP_PATH, cgroup))
		return
	}
	if err = addToCgroup(pid, path.Join(CGROUP_PATH, cgroup, CGROUP_PROCS)); err != nil {
		log.Errorf("Error trying to add pid to cgroup (%s)", cgroup)
	} else {
		log.Infof("%d -> %s", pid, cgroup)
	}
	return
}

func addToCgroup(pid int, cgroupProcsFile string) error {
	f, err := os.OpenFile(cgroupProcsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Error(err)
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("%d\n", pid)); err != nil {
		log.Errorw("Couldn't write pid to the groupc.procs file", "pid", pid, "err", err.Error())
		return err
	}
	return nil
}
