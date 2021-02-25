package model

import "fmt"

type Car struct {
	Path     []*Street
	Position int
	done     bool
}

// We wanted to calculate the min time for a car to spend on the system
// We can calculate the time that a particular car would spend on the system by (its path + time waited on interection + number of cars on the intersection in front of us)

func (c *Car) MinRemaining() int {
	if c.done {
		return 0
	}
	p := c.Position
	t := 0
	for _, str := range c.Path {
		rem := str.Length - p
		p = 0
		if rem < 0 {
			continue
		}
		t += rem
	}
	return t
}

func (c *Car) String() string {
	s := "{"
	for i, street := range c.Path {
		if i == 0 {
			s += fmt.Sprintf("%s(%d/%d),", street.Name, c.Position, street.Length)
		} else {
			s += fmt.Sprintf("%s(%d),", street.Name, street.Length)
		}
	}
	s += "}"
	return s
}
