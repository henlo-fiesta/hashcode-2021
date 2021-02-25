package model

type Car struct {
	Path []*Street
}

func (c *Car) String() string {
	s := "{"
	for _, street := range c.Path {
		s += street.Name + ","
	}
	s += "}"
	return s
}
