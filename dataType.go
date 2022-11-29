package main

import "database/sql"

// movesoutput represents an item from movesoutput table
type Movesoutput struct {
	MOVESRunID         int
	iterationID        sql.NullInt16
	yearID             sql.NullInt16
	monthID            sql.NullInt16
	dayID              sql.NullInt16
	hourID             sql.NullInt16
	stateID            sql.NullInt16
	countyID           sql.NullInt16
	zoneID             sql.NullInt16
	linkID             sql.NullInt16
	pollutantID        sql.NullInt16
	processID          sql.NullInt16
	sourceTypeID       sql.NullInt16
	regClassID         sql.NullInt16
	fuelTypeID         sql.NullInt16
	fuelSubTypeID      sql.NullInt16
	modelYearID        sql.NullInt16
	roadTypeID         sql.NullInt16
	SCC                sql.NullString
	engTechID          sql.NullInt16
	sectorID           sql.NullInt16
	hpID               sql.NullInt16
	emissionQuant      sql.NullFloat64
	emissionQuantMean  sql.NullFloat64
	emissionQuantSigma sql.NullFloat64
}

//type Movesoutput struct {
//	MOVESRunID         int
//	iterationID        *int
//	yearID             *int
//	monthID            *int
//	dayID              *int
//	hourID             *int
//	stateID            *int
//	countyID           *int
//	zoneID             *int
//	linkID             *int
//	pollutantID        *int
//	processID          *int
//	sourceTypeID       *int
//	regClassID         *int
//	fuelTypeID         *int
//	fuelSubTypeID      *int
//	modelYearID        *int
//	roadTypeID         *int
//	SCC                *int
//	engTechID          *int
//	sectorID           *int
//	hpID               *int
//	emissionQuant      *float64
//	emissionQuantMean  *float64
//	emissionQuantSigma *float64
//}
