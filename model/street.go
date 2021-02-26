package model

import (
	"container/list"
	"fmt"
)

type Street struct {
	Name       string
	Start      *Intersection
	End        *Intersection
	Length     int
	Go         bool
	Queue      *list.List
	Bandwidth  int
	Congestion int
}

func (s *Street) String() string {
	return fmt.Sprintf("%s(%d)", s.Name, s.Length)
}

func (s *Street) InformIntersections() {
	s.Start.Out = append(s.Start.Out, s)
	s.End.In = append(s.End.In, s)
}
