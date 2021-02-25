package model

import "fmt"

type Street struct {
	Name   string
	Start  *Intersection
	End    *Intersection
	Length int
}

func (s *Street) String() string {
	return fmt.Sprintf("%s(%d)", s.Name, s.Length)
}
