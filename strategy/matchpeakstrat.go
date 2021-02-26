package strategy

import (
	"github.com/henlo-fiesta/hashcode-2021/model"
	"sort"
)

type RulesByPeak struct{ ruleset }

func (rbp RulesByPeak) Len() int { return len(rbp.ruleset) }
func (rbp RulesByPeak) Swap(i, j int) {
	rbp.ruleset[i], rbp.ruleset[j] = rbp.ruleset[j], rbp.ruleset[i]
}
func (rbp RulesByPeak) Less(i, j int) bool {
	return rbp.ruleset[i].Street.PeaksAt < rbp.ruleset[j].Street.PeaksAt
}

func MatchPeakStrat(s *model.Simulation) {
	for i := range s.Intersections {
		inter := &s.Intersections[i]
		if inter.CycleTime < 1 {
			continue
		}
		n := make([]*model.Street, len(inter.Schedule))
		for i, entry := range inter.Schedule {
			n[i] = entry.Street
			for t := 1; t < inter.CycleTime; t++ {
				entry.Street.PeaksAt += float64(t * entry.Street.CongestionAt[t])
			}
			entry.Street.PeaksAt /= float64(inter.CycleTime)
		}
		sort.Sort(RulesByPeak{inter.Schedule})
	}
}
