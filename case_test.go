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

	simulation := model.NewSimulation("a.txt")
	optimize(simulation)
	simulation = model.NewSimulation("b.txt")
	optimize(simulation)
	simulation = model.NewSimulation("c.txt")
	optimize(simulation)
	simulation = model.NewSimulation("d.txt")
	optimize(simulation)
	simulation = model.NewSimulation("e.txt")
	optimize(simulation)
	simulation = model.NewSimulation("f.txt")
	optimize(simulation)
}
