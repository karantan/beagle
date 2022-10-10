package spotter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_parseRawProcesses(t *testing.T) {
	stdout := `      4 123 php-fpm: pool foo_com
		2 234 php-fpm: pool bar_com
	`
	want := []Process{
		{123, "foo_com", time.Duration(4)},
		{234, "bar_com", time.Duration(2)},
	}
	got := parseRawProcesses(stdout)
	assert.Equal(t, want, got)

	var want2 []Process
	got2 := parseRawProcesses("")
	assert.Equal(t, want2, got2)
}
