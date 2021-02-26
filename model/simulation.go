package model

import (
	"fmt"
	"io"
	"runtime"
)

type Simulation struct {
	Duration            int
	Intersections       []Intersection
	intersectionBatches [][]*Intersection
	Bonus               int
	Cars                []Car
	StreetsIndex        map[string]*Street
	Streets             []Street
	T                   int
	Score               int
	jobs                chan func()
	jobResults          chan bool
}

func (simulation *Simulation) GetStreets() Streets {
	// TODO cache?
	s := make(Streets, len(simulation.Streets))
	for i := range simulation.Streets {
		s[i] = &simulation.Streets[i]
	}
	return s
}
func (simulation *Simulation) GetIntersections() Intersections {
	is := make(Intersections, len(simulation.Intersections))
	for i := range simulation.Intersections {
		is[i] = &simulation.Intersections[i]
	}
	return is
}

func (simulation *Simulation) SaveBest() {
	for i := range simulation.Intersections {
		simulation.Intersections[i].SaveBest()
	}
}

func (simulation *Simulation) BuildIndex() {
	simulation.StreetsIndex = make(map[string]*Street)
	for i := range simulation.Streets {
		simulation.StreetsIndex[simulation.Streets[i].Name] = &simulation.Streets[i]
		simulation.Streets[i].InformIntersections()
	}
}

func (simulation *Simulation) Clone() *Simulation {
	var newSim = *simulation
	newSim.Cars = make([]Car, len(simulation.Cars))
	copy(newSim.Cars, simulation.Cars)
	for i := range newSim.Cars {
		newSim.Cars[i].Path = make([]*Street, len(simulation.Cars[i].Path))
		copy(newSim.Cars[i].Path, simulation.Cars[i].Path)
	}
	return &newSim
}

func workers(id int, jobs <-chan func(), results chan<- bool) {
	for j := range jobs {
		j()
		results <- true
	}
}

func (simulation *Simulation) StartWorkers(n int) {
	if n < 1 {
		n = runtime.NumCPU()
	}

	batchSize := len(simulation.Intersections) / n / 8
	if batchSize < 4 {
		batchSize = 4
	}
	for i := 0; i < len(simulation.Intersections); {
		j := i
		batch := make([]*Intersection, 0, batchSize)
		for ; j < i+batchSize && j < len(simulation.Intersections); j++ {
			batch = append(batch, &simulation.Intersections[j])
		}
		simulation.intersectionBatches = append(simulation.intersectionBatches, batch)
		i = j
	}

	simulation.jobs = make(chan func(), n)
	simulation.jobResults = make(chan bool, n)
	for i := 0; i < n; i++ {
		go workers(i+2, simulation.jobs, simulation.jobResults)
	}
}

func (simulation *Simulation) StopWorkers() {
	close(simulation.jobs)
	close(simulation.jobResults)
}

func (simulation *Simulation) Reset() {
	for i := range simulation.Streets {
		str := &simulation.Streets[i]
		str.Queue.Init()
		str.Congestion = 0
	}
	for i := range simulation.Cars {
		car := &simulation.Cars[i]
		car.Path[0].Queue.PushBack(car)
		car.Position = car.Path[0].Length + 1
	}
}

func (simulation *Simulation) Step() {
	// set streets traffic light state

	go func() {
		for i := range simulation.intersectionBatches {
			batch := simulation.intersectionBatches[i]
			simulation.jobs <- func() {
				for _, inter := range batch {
					inter.ApplyRules(simulation.T)
				}
			}
		}
	}()
	for range simulation.intersectionBatches {
		<-simulation.jobResults
	}

	// dequeue StreetsIndex
	for i := range simulation.Streets {
		str := &simulation.Streets[i]
		cong := str.Queue.Len()
		if cong==0{
			continue
		}
		str.Congestion += cong
		str.CongestionAt[simulation.T%str.End.CycleTime] += cong
		if str.Go {
			front := str.Queue.Front()
			if front == nil {
				continue
			}
			car := str.Queue.Remove(front).(*Car)
			// advance car to next street segment
			car.Path = car.Path[1:]
			car.Position = -1
		}
	}

	// step cars
	for i := range simulation.Cars {
		car := &simulation.Cars[i]
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
			simulation.Score += simulation.Bonus + (simulation.Duration - simulation.T)
			continue
		}

		// case end street
		if car.Position == str.Length {
			str.Queue.PushBack(car)
			car.Position++ // when car == str.Length +1 -> on queue
		}
	}

	simulation.T++
}

func (simulation *Simulation) Run() int {
	for i := range simulation.Intersections {
		inter := &simulation.Intersections[i]
		inter.CompileRules()
		for _, str := range inter.In {
			str.CongestionAt = make([]int, inter.CycleTime)
		}
	}
	for simulation.T <= simulation.Duration {
		simulation.Step()
	}
	return simulation.Score
}

func (simulation *Simulation) WriteSchedule(w io.Writer) error {
	m := make(map[int][]*ScheduleEntry)
	for id := range simulation.Intersections {
		inter := &simulation.Intersections[id]
		entries := inter.ConciseSchedule()
		if len(entries) == 0 {
			continue
		}
		m[id] = entries
	}

	if _, err := fmt.Fprintf(w, "%d\n", len(m)); err != nil {
		return err
	}
	for id, entries := range m {
		if _, err := fmt.Fprintf(w, "%d\n", id); err != nil {
			return nil
		}
		if _, err := fmt.Fprintf(w, "%d\n", len(entries)); err != nil {
			return nil
		}
		for _, entry := range entries {
			if _, err := fmt.Fprintf(w, "%s %d\n", entry.Street.Name, entry.Duration); err != nil {
				return nil
			}
		}
	}
	return nil
}
