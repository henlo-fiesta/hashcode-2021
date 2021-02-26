package model

import (
	"bufio"
	"container/list"
	"log"
	"os"
	"strconv"
	"strings"
)

func atoi(raw string) int {
	num, err := strconv.Atoi(raw)
	if err != nil {
		log.Fatal(err)
	}
	return num
}

func LoadSimulation(filename string) (*Simulation, error) {
	f, err := os.Open("input/" + filename + ".txt")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Print(err)
		}
	}()
	scanner := bufio.NewScanner(f)
	scanner.Scan()

	// First line
	line := scanner.Text()
	fields := strings.Fields(line)
	time := atoi(fields[0])
	numIntersections := atoi(fields[1])
	numStreets := atoi(fields[2])
	numCars := atoi(fields[3])
	scorePerCar := atoi(fields[4])

	simulation := &Simulation{
		Duration:      time,
		Intersections: make([]Intersection, numIntersections),
		Bonus:         scorePerCar,
		Streets:       make([]Street, numStreets),
		Cars:          make([]Car, numCars),
	}
	for i := 0; i < numIntersections; i++ {
		simulation.Intersections[i].Id = i
	}

	for i := 0; i < numStreets; i++ {
		if scanner.Scan() {
			line = scanner.Text()
			fields = strings.Fields(line)

			simulation.Streets[i] = Street{
				Name:   fields[2],
				Start:  &simulation.Intersections[atoi(fields[0])],
				End:    &simulation.Intersections[atoi(fields[1])],
				Length: atoi(fields[3]),
				Queue:  list.New(),
			}
		} else {
			return nil, scanner.Err()
		}
	}

	simulation.BuildIndex()

	for i := 0; i < numCars; i++ {
		if scanner.Scan() {
			line = scanner.Text()
			fields = strings.Fields(line)

			pathLength := atoi(fields[0])
			path := make([]*Street, pathLength)
			simulation.Cars[i] = Car{Path: path}
			for pathId := 0; pathId < pathLength; pathId++ {
				path[pathId] = simulation.StreetsIndex[fields[pathId+1]]
			}
		} else {
			return nil, scanner.Err()
		}
	}
	return simulation, nil
}
