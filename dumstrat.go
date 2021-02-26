package main

import "github.com/henlo-fiesta/hashcode-2021/model"

func dumStrat(sim *model.Simulation) {
	for i := range sim.Cars {
		path := sim.Cars[i].Path
		for i, str := range path {
			str.Bandwidth++
			if i < len(path)-1 {
				// meaning car crosses thru the intersection
				str.End.Bandwidth++
			}
		}
	}
	for i := range sim.Intersections {
		inter := &sim.Intersections[i]
		max := inter.In[0]
		for _, str := range inter.In {
			if str.Bandwidth > max.Bandwidth {
				max = str
			}
		}
		for _, str := range inter.In {
			if str.Bandwidth == 0 {
				continue
			}
			entry := model.ScheduleEntry{
				Street:   str,
				Duration: 0,
			}
			if str.Bandwidth == max.Bandwidth {
				entry.Duration = 1
			}
			inter.Schedule = append(inter.Schedule, &entry)
			inter.CycleTime += entry.Duration
		}
	}
}
