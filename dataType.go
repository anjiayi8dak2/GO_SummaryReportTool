package main

// Movesoutput all the field in movesoutput table
type Movesoutput struct {
	MOVESRunID    int     `db:"MOVESRunID"`
	iterationID   int     `db:"iterationID"`
	yearID        int     `db:"yearID"`
	monthID       int     `db:"monthID"`
	dayID         int     `db:"dayID"`
	hourID        int     `db:"hourID"`
	stateID       int     `db:"stateID"`
	countyID      int     `db:"countyID"`
	zoneID        int     `db:"zoneID"`
	linkID        int     `db:"linkID"`
	pollutantID   int     `db:"pollutantID"`
	processID     int     `db:"processID"`
	sourceTypeID  int     `db:"sourceTypeID"`
	regClassID    int     `db:"regClassID"`
	fuelTypeID    int     `db:"fuelTypeID"`
	fuelSubTypeID int     `db:"fuelSubTypeID"`
	modelYearID   int     `db:"modelYearID"`
	roadTypeID    int     `db:"roadTypeID"`
	SCC           int     `db:"SCC"`
	engTechID     int     `db:"engTechID"`
	sectorID      int     `db:"sectorID"`
	hpID          int     `db:"hpID"`
	emissionQuant float64 `db:"emissionQuant"`
}

type dummy struct {
	MOVESRunID    int     `db:"MOVESRunID"`
	iterationID   int     `db:"iterationID"`
	yearID        int     `db:"yearID"`
	monthID       int     `db:"monthID"`
	dayID         int     `db:"dayID"`
	stateID       int     `db:"stateID"`
	countyID      int     `db:"countyID"`
	pollutantID   int     `db:"pollutantID"`
	processID     int     `db:"processID"`
	modelYearID   int     `db:"modelYearID"`
	SCC           int     `db:"SCC"`
	engTechID     int     `db:"engTechID"`
	sectorID      int     `db:"sectorID"`
	hpID          int     `db:"hpID"`
	emissionQuant float64 `db:"emissionQuant"`
}
