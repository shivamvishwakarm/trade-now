package instrument

type Instrument struct {
	Id       int
	Name     string
	Symbol   string
	Exchange string
	Segment  string
	LotSize  int
	TickSize int
	Status   string
}
