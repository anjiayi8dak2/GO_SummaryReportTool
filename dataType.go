package main

type Movesactivityoutput struct {
	MOVESRunID     int
	iterationID    int
	yearID         int
	monthID        int
	dayID          int
	hourID         int
	stateID        int
	countyID       int
	zoneID         int
	linkID         int
	sourceTypeID   int
	regClassID     int
	fuelTypeID     int
	fuelSubTypeID  int
	modelYearID    int
	roadTypeID     int
	SCC            int
	engTechID      int
	sectorID       int
	hpID           int
	activityTypeID int
	activity       float64
}

// Movesoutput all the field in movesoutput table
type Movesoutput struct {
	MOVESRunID    int
	iterationID   int
	yearID        int
	monthID       int
	dayID         int
	hourID        int
	stateID       int
	countyID      int
	zoneID        int
	linkID        int
	pollutantID   int
	processID     int
	sourceTypeID  int
	regClassID    int
	fuelTypeID    int
	fuelSubTypeID int
	modelYearID   int
	roadTypeID    int
	SCC           int
	engTechID     int
	sectorID      int
	hpID          int
	emissionQuant float64
}

type Rateperdistance struct {
	MOVESScenarioID string
	MOVESRunID      int
	yearID          int
	monthID         int
	dayID           int
	hourID          int
	linkID          int
	pollutantID     int
	processID       int
	sourceTypeID    int
	regClassID      int
	SCC             int
	fuelTypeID      int
	modelYearID     int
	roadTypeID      int
	avgSpeedBinID   int
	temperature     float64
	relHumidity     float64
	ratePerDistance float64
}

type Rateperhour struct {
	MOVESScenarioID string
	MOVESRunID      int
	yearID          int
	monthID         int
	dayID           int
	hourID          int
	linkID          int
	pollutantID     int
	processID       int
	sourceTypeID    int
	regClassID      int
	SCC             int
	fuelTypeID      int
	modelYearID     int
	roadTypeID      int
	temperature     float64
	relHumidity     float64
	ratePerHour     float64
}

type Rateperprofile struct {
	MOVESScenarioID      string
	MOVESRunID           int
	temperatureProfileID int
	yearID               int
	dayID                int
	hourID               int
	pollutantID          int
	processID            int
	sourceTypeID         int
	regClassID           int
	SCC                  int
	fuelTypeID           int
	modelYearID          int
	temperature          float64
	relHumidity          float64
	ratePerVehicle       float64
}

type Rateperstart struct {
	MOVESScenarioID string
	MOVESRunID      int
	yearID          int
	monthID         int
	dayID           int
	hourID          int
	zoneID          int
	sourceTypeID    int
	regClassID      int
	SCC             int
	fuelTypeID      int
	modelYearID     int
	pollutantID     int
	processID       int
	temperature     float64
	relHumidity     float64
	ratePerStart    float64
}

type Ratepervehicle struct {
	MOVESScenarioID string
	MOVESRunID      int
	yearID          int
	monthID         int
	dayID           int
	hourID          int
	zoneID          int
	pollutantID     int
	processID       int
	sourceTypeID    int
	regClassID      int
	SCC             int
	fuelTypeID      int
	modelYearID     int
	temperature     float64
	relHumidity     float64
	ratePerVehicle  float64
}

type Startspervehicle struct {
	MOVESScenarioID  string
	MOVESRunID       int
	yearID           int
	monthID          int
	dayID            int
	hourID           int
	zoneID           int
	sourceTypeID     int
	regClassID       int
	SCC              int
	fuelTypeID       int
	modelYearID      int
	startsPerVehicle float64
}

// reshape query result from long to wide, in go we use map to represent contingency table
// prepare for map
type wideTableShapeStruct struct {
	field1, field2 string
}
