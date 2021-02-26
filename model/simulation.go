package model

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Simulation struct {
	Duration      int
	Intersections []*Intersection
	Bonus         int
	Cars          []*Car
	Streets       map[string]*Street
	ss            []*Street
	T             int
	Score         int
}

func atoi(raw string) int {
	num, err := strconv.Atoi(raw)
	if err != nil {
		log.Fatal(err)
	}
	return num
}
func (s *Simulation) SaveBest() {
	for _, i := range s.Intersections {
		i.SaveBest()
	}
}
func NewSimulation(filename string) *Simulation {
	// filename := "f.txt"
	// Read simulation input
	// data, err := ioutil.ReadFile("input/a.txt")
	// if err != nil {
	// 	panic("")
	// }

	file, err := os.Open("input/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	// First line
	line := scanner.Text()
	s := strings.Fields(line)
	fmt.Printf("%+v\n", s)

	// Time
	time := atoi(s[0])

	numIntersections, err := strconv.Atoi(s[1])
	if err != nil {
		log.Fatal(err)
	}

	numStreets, err := strconv.Atoi(s[2])
	if err != nil {
		log.Fatal(err)
	}

	numCars, err := strconv.Atoi(s[3])
	if err != nil {
		log.Fatal(err)
	}

	scorePerCar := atoi(s[4])

	simulation := &Simulation{
		Duration:      time,
		Intersections: make([]*Intersection, numIntersections),
		Bonus:         scorePerCar,
		Streets:       make(map[string]*Street),
		ss:            make([]*Street, numStreets),
		Cars:          make([]*Car, numCars),
	}
	for i := range simulation.Intersections {
		simulation.Intersections[i] = &Intersection{}
	}

	// streets := []model.Street{}
	// streets := map[string]model.Street{}
	for i := 0; i < numStreets; i++ {
		if scanner.Scan() {
			line := scanner.Text()
			streetString := strings.Fields(line)

			// Street Input: [start] [end] [name] [length]
			// streets = append(streets, model.Street{
			// 	Name:   streetString[2],
			// 	Start:  simulation.Intersections[atoi(streetString[0])],
			// 	End:    simulation.Intersections[atoi(streetString[1])],
			// 	Length: atoi(streetString[3]),
			// })
			street := Street{
				Name:   streetString[2],
				Start:  simulation.Intersections[atoi(streetString[0])],
				End:    simulation.Intersections[atoi(streetString[1])],
				Length: atoi(streetString[3]),
				Queue:  list.New(),
			}
			simulation.Streets[streetString[2]] = &street
			simulation.ss[i] = &street
			street.Start.Out = append(street.Start.Out, &street)
			street.End.In = append(street.End.In, &street)
		}
	}

	for i := 0; i < numCars; i++ {
		if scanner.Scan() {
			line := scanner.Text()
			carString := strings.Fields(line)

			pathLength := atoi(carString[0])

			car := &Car{
				Path: []*Street{},
			}

			// Assign the list in reverse the list since we wanted the path order
			for path := 1; path <= pathLength; path++ {
				car.Path = append(car.Path, simulation.Streets[carString[path]])
			}

			simulation.Cars[i] = car
		}
	}

	// Load car into the end of street
	for _, car := range simulation.Cars {
		car.Path[0].Queue.PushBack(car)
		car.Position = car.Path[0].Length + 1
	}
	return simulation
}

func (sim *Simulation) Clone() *Simulation {
	var newSim = *sim
	newSim.Cars = make([]*Car, len(sim.Cars))
	for i := range newSim.Cars {
		var newCar = *sim.Cars[i]
		newSim.Cars[i] = &newCar
		newSim.Cars[i].Path = make([]*Street, len(sim.Cars[i].Path))
		copy(newSim.Cars[i].Path, sim.Cars[i].Path)
	}
	return &newSim
}

func (sim *Simulation) Reset() {
	for _, str := range sim.Streets {
		str.Queue.Init()
	}
	for _, cars := range sim.Cars {
		cars.Path[0].Queue.PushBack(cars)
	}
}

func (sim *Simulation) Step() {
	// set streets traffic light state
	for _, inter := range sim.Intersections {
		if inter.CycleTime < 1 {
			continue
		}
		at := sim.T % inter.CycleTime
		dt := 0
		for _, sched := range inter.Schedule {
			// fmt.Printf("%s dt=%d at=%d dur=%d\n", sched.Street.Name, dt, at, sched.Duration)
			sched.Street.Go = dt <= at && at < dt+sched.Duration
			dt += sched.Duration
		}
	}
	// dequeue Streets
	for _, str := range sim.ss {
		str.Congestion += str.Queue.Len()
		if str.Go {
			// fmt.Println(str.Name + " street is green")
			front := str.Queue.Front()
			if front == nil {
				continue
			}
			car := str.Queue.Remove(front).(*Car)
			car.Position = -1
			car.Path = car.Path[1:]
		}
	}
	// step cars
	for _, car := range sim.Cars {
		if car.done {
			continue
		}
		str := car.Path[0]
		// case queued
		if car.Position == str.Length+1 {
			continue
		}
		car.Position++
		// case fin
		if len(car.Path) == 1 && car.Position == str.Length {
			car.done = true
			sim.Score += sim.Bonus + (sim.Duration - sim.T)
			continue
		}
		// case end street
		if car.Position == str.Length {
			str.Queue.PushBack(car)
			car.Position++ // when car == str.Length +1 -> on queue
		}
	}

	sim.T++
}

func (sim *Simulation) Run() int {
	// make sched concise
	for _, i := range sim.Intersections {
		i.fullSched = make([]*ScheduleEntry, len(i.Schedule))
		copy(i.fullSched, i.Schedule)
		i.Schedule = make([]*ScheduleEntry, 0)
		for _, s := range i.fullSched {
			if s.Duration > 0 {
				i.Schedule = append(i.Schedule, s)
			}
		}
	}
	for sim.T <= sim.Duration {
		// fmt.Printf("%+v\n\n", simulation.Cars)
		// fmt.Printf("%+v\n\n", simulation)
		sim.Step()
	}
	for _, i := range sim.Intersections {
		i.Schedule = make([]*ScheduleEntry, len(i.fullSched))
		copy(i.Schedule, i.fullSched)
	}
	return sim.Score
}

func (sim *Simulation) Output() string {
	out := ""
	numInt := 0
	for id, inter := range sim.Intersections {
		if len(inter.Schedule) == 0 {
			continue
		}
		out += fmt.Sprintf("%d\n", id)
		out += fmt.Sprintf("%d\n", len(inter.Schedule))
		for _, entry := range inter.Schedule {
			out += fmt.Sprintf("%s %d\n", entry.Street.Name, entry.Duration)
		}
		numInt++
	}
	out = fmt.Sprintf("%d\n", numInt) + out
	return out
}
