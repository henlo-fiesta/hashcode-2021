package strategy

import (
	"fmt"
	"github.com/henlo-fiesta/hashcode-2021/model"
	"math"
	"sort"
)

type StreetsByBandwidth struct {
	model.Streets
}

func (s StreetsByBandwidth) Less(i, j int) bool {
	return s.Streets[i].Bandwidth < s.Streets[j].Bandwidth
}

type IntersectionsByBandwidth struct {
	model.Intersections
}

func (is IntersectionsByBandwidth) Less(i, j int) bool {
	return is.Intersections[i].Bandwidth < is.Intersections[j].Bandwidth
}
func CalcStats(sim *model.Simulation) {
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
		for _, str := range inter.In {
			if str.Bandwidth > 0 {
				inter.ActiveIn++
			}
		}
		if inter.ActiveIn > 0 {
			inter.Mean = float64(inter.Bandwidth) / float64(inter.ActiveIn)
			for _, s := range inter.In {
				if s.Bandwidth > 0 {
					inter.Variance += math.Pow(float64(s.Bandwidth)-inter.Mean, 2)
				}
			}
			inter.StdDev = math.Sqrt(inter.Variance)
			inter.Z = inter.StdDev / inter.Mean
		}
	}
	for i := range sim.Intersections {
		inter := &sim.Intersections[i]
		if inter.Bandwidth == 0 {
			continue
		}
		inter.HopefulCycles = float64(sim.Duration) / float64(inter.Bandwidth)
		inter.ActualCycles = int(inter.HopefulCycles)
		if inter.ActualCycles < 1 {
			fmt.Printf("...warning... cycle overflow, fixing to 1\n")
			inter.ActualCycles = 1
		}
		/*shrink := 1.01
		if inter.ActualCycles > 7 {
			inter.ActualCycles = int(math.Floor(inter.HopefulCycles / shrink))
			if inter.ActualCycles < 7 {
				inter.ActualCycles = 7
			}
		}
		if inter.ActualCycles > 13 {
			divisor := inter.Z/4
			if divisor < shrink{
				divisor = shrink
			}
			inter.ActualCycles = int(math.Floor(inter.HopefulCycles / divisor))
			if inter.ActualCycles < 13 {
				inter.ActualCycles = 13
			}
		}*/
		inter.TargetCycleLength = sim.Duration / inter.ActualCycles
	}
}

func PrintStats(sim *model.Simulation) {
	fmt.Printf("Simulation - ruleset:%d #s:%d #i:%d #c:%d\n",
		sim.Duration,
		len(sim.Streets),
		len(sim.Intersections),
		len(sim.Cars))
	fmt.Printf("top 10 intersections by bandwidth:\n")
	is := IntersectionsByBandwidth{sim.GetIntersections()}
	sort.Sort(is)
	for i := is.Len() - 1; i >= 0 && i > is.Len()-11; i-- {
		inter := is.Intersections[i]
		if inter.Bandwidth == 0 {
			break
		}

		fmt.Printf(" - %d: b:%d #hc:%.2f #ac:%d(l=%d) avg:%.2f stddev:%.2f(%.2fZ)\n",
			inter.Id,
			inter.Bandwidth,
			inter.HopefulCycles,
			inter.ActualCycles,
			inter.TargetCycleLength,
			inter.Mean,
			inter.StdDev,
			inter.Z)
	}

	fmt.Printf("top 6 streets of top 3 intersections:\n")
	for i := is.Len() - 1; i >= 0 && i > is.Len()-4; i-- {
		inter := is.Intersections[i]
		fmt.Printf(" - INTER %d: b:%d mean:%.2f\n", inter.Id, inter.Bandwidth, inter.Mean)
		sbi := StreetsByBandwidth{inter.In}
		sum := float64(0)
		for _, s := range sbi.Streets {
			sum += float64(s.Bandwidth)
		}
		sort.Sort(sbi)
		for j := sbi.Len() - 1; j >= 0 && j > sbi.Len()-7; j-- {
			str := sbi.Streets[j]
			pc := float64(str.Bandwidth) * 100 / sum
			fmt.Printf("   - %s(%d): b:%d(%.2f%%) %+.2fZ\n",
				str.Name, str.Length, str.Bandwidth,
				pc, float64(str.Bandwidth)/inter.Mean-1)
		}
	}
}
