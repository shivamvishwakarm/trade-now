package database

import (
	"context"
	"database/sql"

	"github.com/shivamvishwakarm/trade-now/internal/instrument"
)

type InstrumentRepository struct {
	db *sql.DB
}

func NewInstrumentRepository(db *sql.DB) *InstrumentRepository {
	return &InstrumentRepository{
		db: db,
	}
}

func (i *InstrumentRepository) GetBySymbol(ctx context.Context, symbol string) (*instrument.Instrument, error) {

	return &instrument.Instrument{
		Name:     "not",
		Symbol:   "d",
		Id:       1,
		Segment:  "equity",
		LotSize:  1,
		TickSize: 1,
	}, nil
}
