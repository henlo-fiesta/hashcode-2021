package main

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/henlo-fiesta/hashcode-2021/model"
)

func atoi(raw string) int {
	num, err := strconv.Atoi(raw)
	if err != nil {
		log.Fatal(err)
	}
	return num
}

func main() {

	files := []string{"a.txt", "b.txt", "c.txt", "d.txt", "e.txt", "f.txt"}
	// files := []string{"a.txt", "b.txt", "c.txt", "e.txt", "f.txt"}
	// files := []string{"d.txt"}
	for _, filename := range files {
		// filename := "f.txt"
		// Read simulation input
		// data, err := ioutil.ReadFile("input/a.txt")
		// if err != nil {
		// 	panic("")
		// }

		file, err := os.Open("input/" + filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Scan()

		// First line
		line := scanner.Text()
		s := strings.Fields(line)
		fmt.Printf("%+v\n", s)

		// Time
		time := atoi(s[0])

		numIntersections, err := strconv.Atoi(s[1])
		if err != nil {
			log.Fatal(err)
		}

		numStreets, err := strconv.Atoi(s[2])
		if err != nil {
			log.Fatal(err)
		}

		numCars, err := strconv.Atoi(s[3])
		if err != nil {
			log.Fatal(err)
		}

		scorePerCar := atoi(s[4])

		simulation := &model.Simulation{
			Duration:      time,
			Intersections: make([]*model.Intersection, numIntersections),
			Bonus:         scorePerCar,
			Streets:       make(map[string]*model.Street),
			Cars:          make([]*model.Car, numCars),
		}
		for i := range simulation.Intersections {
			simulation.Intersections[i] = &model.Intersection{}
		}

		// streets := []model.Street{}
		// streets := map[string]model.Street{}
		for i := 0; i < numStreets; i++ {
			if scanner.Scan() {
				line := scanner.Text()
				streetString := strings.Fields(line)

				// Street Input: [start] [end] [name] [length]
				// streets = append(streets, model.Street{
				// 	Name:   streetString[2],
				// 	Start:  simulation.Intersections[atoi(streetString[0])],
				// 	End:    simulation.Intersections[atoi(streetString[1])],
				// 	Length: atoi(streetString[3]),
				// })
				street := model.Street{
					Name:   streetString[2],
					Start:  simulation.Intersections[atoi(streetString[0])],
					End:    simulation.Intersections[atoi(streetString[1])],
					Length: atoi(streetString[3]),
					Queue:  list.New(),
				}
				simulation.Streets[streetString[2]] = &street
				street.Start.Out = append(street.Start.Out, &street)
				street.End.In = append(street.End.In, &street)
			}
		}

		for i := 0; i < numCars; i++ {
			if scanner.Scan() {
				line := scanner.Text()
				carString := strings.Fields(line)

				pathLength := atoi(carString[0])

				car := &model.Car{
					Path: []*model.Street{},
				}

				// Assign the list in reverse the list since we wanted the path order
				for path := 1; path <= pathLength; path++ {
					car.Path = append(car.Path, simulation.Streets[carString[path]])
				}

				simulation.Cars[i] = car
			}
		}

		// Load car into the end of street
		for _, car := range simulation.Cars {
			car.Path[0].Queue.PushBack(car)
			car.Position = car.Path[0].Length + 1
		}

		cloneSim := simulation.Clone()
		dumStrat(simulation)

		for simulation.T <= simulation.Duration {
			// fmt.Printf("%+v\n\n", simulation.Cars)
			// fmt.Printf("%+v\n\n", simulation)
			simulation.Step()
		}
		fmt.Printf("%+v\n", simulation.Score)
		/*for _, str := range simulation.Streets {
			fmt.Printf("%20s cong=%04d\n", str.Name, str.Congestion)
		}*/
		for i := 0; i < 3; i++ {
			simulation := cloneSim
			cloneSim = simulation.Clone()
			simulation.Reset()
			dogStrats(simulation)
			for simulation.T <= simulation.Duration {
				simulation.Step()
			}
			fmt.Printf("%+v\n", simulation.Score)
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		// fmt.Println(simulation.Output())
		f, err := os.Create("output/" + filename)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		_, err = f.WriteString(simulation.Output())
		if err != nil {
			log.Fatal(err)
		}
		f.Sync()
	}
}
