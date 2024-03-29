package main

import (
	"flag"
	"github.com/henlo-fiesta/hashcode-2021/model"
	"log"
	"os"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	files := []string{"a.txt", "b.txt", "c.txt", "d.txt", "e.txt", "f.txt"}
	for _, filename := range files {
		simulation, err := model.LoadSimulation(filename)
		if err != nil {
			log.Fatal(err)
		}
		optimize(simulation)

		f, err := os.Create("output/" + filename)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		if err := simulation.WriteSchedule(f); err != nil {
			log.Fatal(err)
		}
		f.Sync()
	}
}
