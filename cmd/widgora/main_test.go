package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name   string
		args   []string
		code   int
		stdout string
		stderr string
	}{
		{"no command", nil, 2, "", "Usage: widgora"},
		{"version", []string{"version"}, 0, "widgora dev (commit none, built unknown)", ""},
		{"help", []string{"help"}, 0, "Usage: widgora", ""},
		{"unknown command", []string{"frobnicate"}, 2, "", `unknown command "frobnicate"`},
		{"serve help", []string{"serve", "-h"}, 0, "", "-listen"},
		{"serve bad log level", []string{"serve", "--log-level", "loud"}, 1, "", "invalid --log-level"},
		{"serve bad log format", []string{"serve", "--log-format", "xml"}, 1, "", "invalid --log-format"},
		{"serve extra argument", []string{"serve", "now"}, 1, "", `unexpected argument "now"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run(tt.args, &stdout, &stderr); code != tt.code {
				t.Fatalf("exit code = %d, want %d (stderr: %s)", code, tt.code, stderr.String())
			}
			if !strings.Contains(stdout.String(), tt.stdout) {
				t.Errorf("stdout = %q, want it to contain %q", stdout.String(), tt.stdout)
			}
			if !strings.Contains(stderr.String(), tt.stderr) {
				t.Errorf("stderr = %q, want it to contain %q", stderr.String(), tt.stderr)
			}
		})
	}
}
