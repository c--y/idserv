package main

import (
	"fmt"
	"testing"
)

var (
	config = &Config{Host: "localhost", Port: "8080", DataCenterId: 10, NumWorkers: 5}
)

func TestNewServer(t *testing.T) {
	s, _ := NewServer(config)
	fmt.Println(s)
	if s == nil {
		t.Error("NewServer() error")
	}
}
