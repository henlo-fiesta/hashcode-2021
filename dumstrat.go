package main

import "github.com/henlo-fiesta/hashcode-2021/model"

func dumStrat(sim *model.Simulation) {
	for _, car := range sim.Cars {
		for _, str := range car.Path {
			str.PathCount++
		}
	}
	for _, inter := range sim.Intersections {
		sum := 0
		for _, str := range inter.In {
			sum += str.PathCount
		}
		if sum == 0 || len(inter.In) == 0 {
			continue
		}
		avg := float64(sum) / float64(len(inter.In))
		for _, str := range inter.In {
			if str.PathCount == 0 {
				continue
			}
			entry := model.ScheduleEntry{
				Street:   str,
				Duration: 1,
			}
			if float64(str.PathCount) < avg {
				entry.Duration = 0
			}
			inter.Schedule = append(inter.Schedule, &entry)
			inter.CycleTime += entry.Duration
		}
		// inter.CycleTime = len(inter.In)
	}
}
