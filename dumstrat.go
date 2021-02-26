package main

import "github.com/henlo-fiesta/hashcode-2021/model"

func dumStrat(sim *model.Simulation) {
	for _, car := range sim.Cars {
		for _, str := range car.Path {
			str.PathCount++
		}
	}
	for _, inter := range sim.Intersections {
		max := inter.In[0]
		for _, str := range inter.In {
			if str.PathCount > max.PathCount {
				max = str
			}
		}
		for _, str := range inter.In {
			if str.PathCount == 0 {
				continue
			}
			entry := model.ScheduleEntry{
				Street:   str,
				Duration: 0,
			}
			if str.PathCount == max.PathCount {
				entry.Duration = 1
			}
			inter.Schedule = append(inter.Schedule, &entry)
			inter.CycleTime += entry.Duration
		}
	}
}
