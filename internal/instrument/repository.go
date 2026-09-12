package instrument

import "context"

type InstrumentRepository interface {
	GetBySymbol(ctx context.Context, symbol string) (*Instrument, error)
}
