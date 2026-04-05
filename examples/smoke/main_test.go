package main

import "testing"

func TestCandidateConfigPaths(t *testing.T) {
	paths := candidateConfigPaths()
	if len(paths) == 0 {
		t.Fatal("expected at least one candidate path")
	}
}
