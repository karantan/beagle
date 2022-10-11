package isolater

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_addToCgroup(t *testing.T) {
	defer os.Truncate("../fixtures/grp/cgroup.procs", 0)
	want := "2\n"
	addToCgroup(2, "../fixtures/grp/cgroup.procs")
	got, _ := os.ReadFile("../fixtures/grp/cgroup.procs")
	assert.Equal(t, want, string(got))
}
