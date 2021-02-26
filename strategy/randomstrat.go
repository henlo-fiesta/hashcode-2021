package strategy

import (
	"github.com/henlo-fiesta/hashcode-2021/model"
	"math/rand"
)

func RandomStrat(s *model.Simulation) {
	for i := range s.Intersections {
		inter := &s.Intersections[i]
		rand.Shuffle(len(inter.Schedule), func(i, j int) {
			inter.Schedule[i], inter.Schedule[j] = inter.Schedule[j], inter.Schedule[i]
		})
	}
}
