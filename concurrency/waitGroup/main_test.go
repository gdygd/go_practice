package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_printSomething(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w
	var wg sync.WaitGroup
	wg.Add(1)

	go printSomething("abcd", &wg)
	wg.Wait()

	_ = w.Close()
	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "abcd") {
		t.Errorf("Expected to find epsilon, but it is not there")
	}
}

func Test_UpdateMessage(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go updateMessage("abcd", &wg)
	wg.Wait()

	if msg != "abcd" {
		t.Errorf("Not match..")
	}

}
