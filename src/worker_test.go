package main

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ch := make(chan uint64)
	w, e := NewWorker(21, 31, ch)
	if e != nil {
		t.Error(e)
	}
	fmt.Println(w)

}

func TestGenerateId(t *testing.T) {
	ch := make(chan uint64)

	w, _ := NewWorker(21, 31, ch)

	var id uint64
	for i := 0; i < 10; i++ {
		id = w.GenerateId()
		time.Sleep(2 * time.Millisecond)
		fmt.Println(id)
	}
}

func BenchmarkGenerateId(b *testing.B) {
	ch := make(chan uint64)

	w, _ := NewWorker(21, 31, ch)
	for i := 0; i < b.N; i++ {
		w.GenerateId()
	}
}

func BenchmarkYieldId(b *testing.B) {
	ch := make(chan uint64)

	w, _ := NewWorker(21, 31, ch)
	go func() {
		for {
			id := <-w.Out
			if id == 0 {
				return
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		w.YieldId()
	}
}
