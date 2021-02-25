package model

import "fmt"

type Simulation struct {
	Duration      int
	Intersections []*Intersection
	Bonus         int
	Cars          []*Car
	Streets       map[string]*Street
	T             int
	Score         int
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
		dT := 0
		for _, sched := range inter.Schedule {
			// fmt.Printf("%s dt=%d at=%d dur=%d\n", sched.Street.Name, dT, at, sched.Duration)
			sched.Street.Go = dT <= at && at < dT+sched.Duration
			dT += sched.Duration
		}
	}
	// dequeue Streets
	for _, str := range sim.Streets {
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
