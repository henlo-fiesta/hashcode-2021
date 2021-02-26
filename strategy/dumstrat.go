package strategy

import (
	"github.com/henlo-fiesta/hashcode-2021/model"
	"math"
)

func DumStrat(sim *model.Simulation) {
	for i := range sim.Cars {
		path := sim.Cars[i].Path
		for i, str := range path {
			if i < len(path)-1 {
				// meaning car crosses thru the intersection
				str.Bandwidth++
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
			if str.Bandwidth > 0 {
				inter.ActiveIn++
			}
		}
		if inter.ActiveIn > 0 {
			// calc stats
			inter.Mean = float64(inter.Bandwidth) / float64(inter.ActiveIn)
			for _, s := range inter.In {
				if s.Bandwidth > 0 {
					inter.Variance += math.Pow(float64(s.Bandwidth)-inter.Mean, 2)
				}
			}
			inter.StdDev = math.Sqrt(inter.Variance)
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
	printStats(sim)
}
