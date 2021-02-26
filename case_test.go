package main

import (
	"flag"
	"fmt"
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
	sum:=0
	for _, filename := range cases {
		simulation, err := model.LoadSimulation(filename)
		if err != nil {
			log.Fatal(err)
		}
		simulation.StartWorkers(0)
		fmt.Printf("SIM %s\n", filename)
		sum+=optimize(simulation)
		simulation.StopWorkers()
		fmt.Println()
	}
	fmt.Printf("total score=%d\n", sum)
}
