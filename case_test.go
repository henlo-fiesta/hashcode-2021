package main

import (
	"flag"
	"github.com/henlo-fiesta/hashcode-2021/model"
	"log"
	"os"
	"runtime/pprof"
	"testing"
)

func TestCases(t *testing.T) {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	cases := []string{"a", "b", "c", "d", "e", "f"}
	for _, filename := range cases {
		simulation, err := model.LoadSimulation(filename)
		if err != nil {
			log.Fatal(err)
		}
		simulation.StartWorkers(0)
		optimize(simulation)
		simulation.StopWorkers()
	}
}
