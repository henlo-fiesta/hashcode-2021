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

func printStats(sim *model.Simulation) {
	fmt.Printf("Simulation - t:%d #s:%d #i:%d #c:%d\n",
		sim.Duration,
		len(sim.Streets),
		len(sim.Intersections),
		len(sim.Cars))
	/* streets are useless, refer to intersections instead
	fmt.Printf("top 10 streets by bandwidth:\n")
	ss := StreetsByBandwidth{sim.GetStreets()}
	sort.Sort(ss)
	for i := ss.Len() - 1; i >= 0 && i > ss.Len()-11; i-- {
		s := ss.Streets[i]
		fmt.Printf(" - %s(%d): b:%d\n", s.Name, s.Length, s.Bandwidth)
	}*/
	fmt.Printf("top 10 intersections by bandwidth:\n")
	is := IntersectionsByBandwidth{sim.GetIntersections()}
	sort.Sort(is)
	for i := is.Len() - 1; i >= 0 && i > is.Len()-11; i-- {
		inter := is.Intersections[i]
		if inter.Bandwidth == 0 {
			break
		}
		hopefulCycles := float64(sim.Duration) / float64(inter.Bandwidth)
		actualCycles := int(math.Floor(hopefulCycles / 2))
		if actualCycles < 1 {
			fmt.Printf("...warning... cycle overflow, fixing to 1\n")
			actualCycles = 1
		}
		Z := inter.StdDev / inter.Mean
		if actualCycles > 2 {
			actualCycles = int(math.Floor(hopefulCycles / 3))
			if actualCycles < 2 {
				actualCycles = 2
			}
		}
		if actualCycles > 3 {
			divisor := Z
			if divisor < 4 {
				divisor = 4
			}
			actualCycles = int(math.Floor(hopefulCycles / divisor))
			if actualCycles < 3 {
				actualCycles = 3
			}
		}
		targetCycleLength := sim.Duration / actualCycles
		fmt.Printf(" - %d: b:%d #hc:%.2f #ac:%d(l=%d) avg:%.2f stddev:%.2f(%.2fZ)\n",
			inter.Id,
			inter.Bandwidth,
			hopefulCycles,
			actualCycles,
			targetCycleLength,
			inter.Mean,
			inter.StdDev,
			Z)
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
