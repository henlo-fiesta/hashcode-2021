package model

type Intersection struct {
	In  []*Street
	Out []*Street

	// In + Time
	Schedule  []*ScheduleEntry
	CycleTime int

	fullSched []*ScheduleEntry
	best      []*ScheduleEntry
}

type ScheduleEntry struct {
	Street   *Street
	Duration int
}

func (inter *Intersection) SaveBest() {
	inter.best = make([]*ScheduleEntry, len(inter.Schedule))
	for i := range inter.Schedule {
		var sched = *inter.Schedule[i]
		inter.best[i] = &sched
	}
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
