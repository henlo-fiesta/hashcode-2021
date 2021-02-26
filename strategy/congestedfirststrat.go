package strategy

import (
	"github.com/henlo-fiesta/hashcode-2021/model"
	"sort"
)

type ruleset []*model.ScheduleEntry
type sbc struct{ ruleset }

func (a sbc) Len() int      { return len(a.ruleset) }
func (a sbc) Swap(i, j int) { a.ruleset[i], a.ruleset[j] = a.ruleset[j], a.ruleset[i] }
func (a sbc) Less(i, j int) bool {
	return a.ruleset[i].Street.Congestion < a.ruleset[j].Street.Congestion
}

func CongestedFirstStrat(s *model.Simulation) {
	for i := range s.Intersections {
		inter := &s.Intersections[i]
		sort.Sort(sort.Reverse(sbc{inter.Schedule}))
	}
}
