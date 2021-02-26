package model

type Intersection struct {
	Id int

	In  []*Street
	Out []*Street

	// In + Time
	Schedule  []*ScheduleEntry
	CycleTime int

	best       []*ScheduleEntry
	rulesIndex []*Street

	Bandwidth int
	ActiveIn  int
	Mean      float64
	Variance  float64
	StdDev    float64
}

type Intersections []*Intersection

func (is Intersections) Len() int      { return len(is) }
func (is Intersections) Swap(i, j int) { is[i], is[j] = is[j], is[i] }

func (inter *Intersection) CompileRules() {
	inter.rulesIndex = make([]*Street, inter.CycleTime)
	dt := 0
	for _, rule := range inter.Schedule {
		for i := 0; i < rule.Duration; i++ {
			inter.rulesIndex[dt+i] = rule.Street
		}
		dt += rule.Duration
	}
}

func (inter *Intersection) ApplyRules(t int) {
	if inter.CycleTime < 1 {
		return
	}
	goStreet := inter.rulesIndex[t%inter.CycleTime]
	for _, str := range inter.In {
		str.Go = str == goStreet
	}
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

func (inter *Intersection) ConciseSchedule() []*ScheduleEntry {
	concise := make([]*ScheduleEntry, 0)
	for _, entry := range inter.Schedule {
		if entry.Duration > 0 {
			concise = append(concise, entry)
		}
	}
	return concise
}
