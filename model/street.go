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
	PathCount  int
	Congestion int
}

func (s *Street) String() string {
	return fmt.Sprintf("%s(%d)", s.Name, s.Length)
}
