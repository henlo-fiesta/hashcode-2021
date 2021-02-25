package model

type Intersection struct {
	In  []*Street
	Out []*Street
}

func (i *Intersection) String() string {
	str := "{in:"
	for _, s := range i.In {
		str += s.Name + ","
	}
	str += " out:"
	for _, s := range i.Out {
		str += s.Name + ","
	}
	return str + "}"
}
