package main

import (
	"fmt"
	"github.com/henlo-fiesta/hashcode-2021/model"
	"github.com/henlo-fiesta/hashcode-2021/strategy"
)

func optimize(simulation *model.Simulation) int {
	cloneSim := simulation.Clone()
	strategy.CalcStats(simulation)
	//strategy.PrintStats(simulation)
	strategy.QuadraticCostStrat(simulation)
	//strategy.DumStrat(simulation)
	prevScore := simulation.Run()
	best := prevScore
	simulation.SaveBest()

	fmt.Printf("score=%d\n", simulation.Score)
	/*for _, str := range simulation.StreetsIndex {
		fmt.Printf("%20s cong=%04d\n", str.Name, str.Congestion)
	}*/
	consecStagnate := 0
	for i := 0; i < simulation.Duration; i++ {
		simulation := cloneSim
		cloneSim = simulation.Clone()
		simulation.Reset()
		//strategy.DogStrats(simulation, i)
		//strategy.CongestedFirstStrat(simulation)
		strategy.MatchPeakStrat(simulation)
		score := simulation.Run()
		if score > best {
			best = score
			simulation.SaveBest()
		}
		growth := float64(score)/float64(prevScore)*100 - 100
		if growth <= 0.2 || score-prevScore < 1000 {
			consecStagnate++
		} else {
			consecStagnate = 0
		}
		fmt.Printf("score=%d (%.2f%%)\n", simulation.Score, growth)
		if consecStagnate > 4 {
			break
		}
		prevScore = score
	}
	fmt.Printf("best score=%d\n\n", best)
	return best
}
