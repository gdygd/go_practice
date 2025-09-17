package main

import (
	"testing"
	"time"
)

func Test_dine(t *testing.T) {
	eatTime = 0 * time.Second
	thinkTime = 0 * time.Second
	sleepTime = 0 * time.Second

	for i := 0; i < 10; i++ {
		phlist = []string{}
		dine()
		if len(phlist) != 5 {
			t.Errorf("length err.. %d", len(phlist))
		}
	}
}

func Test_dineWithVaryingDelays(t *testing.T) {
	var theTests = []struct {
		name  string
		delay time.Duration
	}{
		{"zero delay", time.Second * 0},
		{"quarter delay", time.Millisecond * 250},
		{"half second delay", time.Millisecond * 500},
	}

	for _, e := range theTests {
		phlist = []string{}
		eatTime = e.delay
		thinkTime = e.delay
		sleepTime = e.delay

		dine()
		if len(phlist) != 5 {
			t.Errorf("length err.. %d", len(phlist))
		}
	}
}
