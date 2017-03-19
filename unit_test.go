package main

import (
	"testing"
)

func TestReadline(t *testing.T) {
	r := readLine("./example/.dbignore")
	s := []string{"node_modules", ".idea"}

	for i := range s {
		if r[i] != s[i] {
			t.Errorf("Expected %s, got %s", s[i], r[i])
		}
	}
}

func TestReadDirNames(t *testing.T) {
	r := readDirNames("./example")
	s := []string{".dbignore"}
	for i := range r {
		if r[i].Name() != s[i] {
			t.Errorf("Expected %s, got %s", s[i], r[i].Name())
		}
	}
}

func TestExec(t *testing.T) {
	s := "5 euro?!"
	r := execute("echo", s)

	if s != r[:len(r)-1] {
		t.Errorf("Expected %s, got %s", s, r)
	}
}