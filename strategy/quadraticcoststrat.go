package strategy

import (
	"github.com/henlo-fiesta/hashcode-2021/model"
	"sort"
)

func QuadraticCostStrat(s *model.Simulation) {
	for i := range s.Intersections {
		inter := &s.Intersections[i]
		if inter.ActiveIn == 0 {
			continue
		}
		r := inter.TargetCycleLength
		n := make(model.Streets, 0)
		rules := make(map[*model.Street]*model.ScheduleEntry)
		maxBandwidth := 0
		for _, str := range inter.In {
			if str.Bandwidth > 0 && r > 0 {
				n = append(n, str)
				if str.Bandwidth > maxBandwidth {
					maxBandwidth = str.Bandwidth
				}
				sched := &model.ScheduleEntry{
					Street:   str,
					Duration: 1,
				}
				inter.Schedule = append(inter.Schedule, sched)
				rules[str] = sched
				r--
				inter.CycleTime++
			}
		}
		sbb:=StreetsByBandwidth{n}
		sort.Sort(sort.Reverse(sbb))
		for cost := float64(maxBandwidth); r > 0; cost *=0.8{
			for _, str := range n {
				if float64(str.Bandwidth) >= cost && r > 0 {
					r--
					rules[str].Duration++
					inter.CycleTime++
				}
			}
		}
	}
	//os.Exit(0)
}
