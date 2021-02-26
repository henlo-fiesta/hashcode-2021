package main

import (
	"sort"

	"github.com/henlo-fiesta/hashcode-2021/model"
)

type DogStreets []*model.ScheduleEntry

func (a DogStreets) Len() int           { return len(a) }
func (a DogStreets) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DogStreets) Less(i, j int) bool { return a[i].Street.Congestion > a[j].Street.Congestion }

func dogStrats(sim *model.Simulation, iter int) {
	for _, inter := range sim.Intersections {
		if len(inter.Schedule) == 0 {
			continue
		}
		congRank := make(DogStreets, len(inter.Schedule))
		last := len(congRank) - 1
		copy(congRank, inter.Schedule)
		sort.Sort(congRank)
		/*sum := 0
		for _, str := range congRank {
			sum += str.Street.Congestion
		}
		avg := float64(sum) / float64(len(congRank))
		if sum == 0 {
			continue
		}

		congRank = congRank[0 : len(congRank)/4]

		for _, sched := range congRank {
			if float64(sched.Street.Congestion)/avg > 2 {
				sched.Duration++
				inter.CycleTime++
			}
		}*/
		for i:=0;i<len(congRank);i++ {
			if float64(congRank[last].Street.Congestion)/float64(congRank[i].Street.Congestion) < 0.7 &&
				congRank[i].Duration<8{
				congRank[i].Duration++
				inter.CycleTime++
				continue
			}
		}
		if iter>5 && congRank[last].Duration > 1 {
			congRank[last].Duration--
			inter.CycleTime--
		}
	}
}
