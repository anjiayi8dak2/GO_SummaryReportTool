SELECT  ifnull(MOVESRunID, -1) AS MOVESRunID,
ifnull(MOVESRunID, -1) AS MOVESRunID,
ifnull(iterationID, -1) AS iterationID,
ifnull(yearID, -1) AS yearID,
ifnull(monthID, -1) AS monthID,
ifnull(dayID, -1) AS dayID,
ifnull(stateID, -1) AS stateID,
ifnull(countyID, -1) AS countyID,
ifnull(zoneID, -1) AS zoneID,
ifnull(linkID, -1) AS linkID,
ifnull(pollutantID, -1) AS pollutantID,
ifnull(processID, -1) AS processID,
ifnull(sourceTypeID, -1) AS sourceTypeID,
ifnull(regClassID, -1) AS regClassID,
ifnull(fuelTypeID, -1) AS fuelTypeID,
ifnull(fuelSubTypeID, -1) AS fuelSubTypeID,
ifnull(modelYearID, -1) AS modelYearID,
ifnull(roadTypeID, -1) AS roadTypeID,
ifnull(SCC, -1) AS SCC,
ifnull(engTechID, -1) AS engTechID,
ifnull(sectorID, -1) AS sectorID,
ifnull(hpID, -1) AS hpID,
ifnull(emissionQuant, 0) AS emissionQuant
FROM 123test.movesoutput;

# rateperdistance
SELECT  
ifnull(MOVESScenarioID, -1) AS MOVESScenarioID,
ifnull(MOVESRunID, -1) AS MOVESRunID,
ifnull(yearID, -1) AS yearID,
ifnull(monthID, -1) AS monthID,
ifnull(dayID, -1) AS dayID,
ifnull(hourID, -1) AS hourID,
ifnull(linkID, -1) AS linkID,
ifnull(pollutantID, -1) AS pollutantID,
ifnull(processID, -1) AS processID,
ifnull(sourceTypeID, -1) AS sourceTypeID,
ifnull(regClassID, -1) AS regClassID,
ifnull(SCC, -1) AS SCC,
ifnull(fuelTypeID, -1) AS fuelTypeID,
ifnull(modelYearID, -1) AS modelYearID,
ifnull(roadTypeID, -1) AS roadTypeID,
ifnull(avgSpeedBinID, -1) AS avgSpeedBinID,
ifnull(temperature, -1) AS temperature,
ifnull(relHumidity, -1) AS relHumidity,
ifnull(ratePerDistance, -1) AS ratePerDistance
FROM 123test.rateperdistance
;

# rateperhour
SELECT  
ifnull(MOVESScenarioID, -1) AS MOVESScenarioID,
ifnull(MOVESRunID, -1) AS MOVESRunID,
ifnull(yearID, -1) AS yearID,
ifnull(monthID, -1) AS monthID,
ifnull(dayID, -1) AS dayID,
ifnull(hourID, -1) AS hourID,
ifnull(linkID, -1) AS linkID,
ifnull(pollutantID, -1) AS pollutantID,
ifnull(processID, -1) AS processID,
ifnull(sourceTypeID, -1) AS sourceTypeID,
ifnull(regClassID, -1) AS regClassID,
ifnull(SCC, -1) AS SCC,
ifnull(fuelTypeID, -1) AS fuelTypeID,
ifnull(modelYearID, -1) AS modelYearID,
ifnull(roadTypeID, -1) AS roadTypeID,
ifnull(temperature, -1) AS temperature,
ifnull(relHumidity, -1) AS relHumidity,
ifnull(ratePerHour, -1) AS ratePerHour
FROM 123test.rateperhour
;


#rateperprofile
SELECT  
ifnull(666, -1) AS MOVESScenarioID,
ifnull(66, -1) AS MOVESRunID,
ifnull(66, -1) AS temperatureProfileID,
ifnull(66, -1) AS yearID,
ifnull(66, -1) AS dayID,
ifnull(66, -1) AS hourID,
ifnull(66, -1) AS pollutantID,
ifnull(66, -1) AS processID,
ifnull(66, -1) AS sourceTypeID,
ifnull(66, -1) AS regClassID,
ifnull(66, -1) AS SCC,
ifnull(66, -1) AS fuelTypeID,
ifnull(66, -1) AS modelYearID,
ifnull(66, -1) AS temperature,
ifnull(66, -1) AS relHumidity,
ifnull(66, -1) AS ratePerVehicle
FROM 123test.rateperprofile
;

#rateperstart
SELECT  
ifnull(66, -1) AS MOVESScenarioID,
ifnull(66, -1) AS MOVESRunID,
ifnull(66, -1) AS yearID,
ifnull(66, -1) AS monthID,
ifnull(66, -1) AS dayID,
ifnull(66, -1) AS hourID,
ifnull(66, -1) AS zoneID,
ifnull(66, -1) AS sourceTypeID,
ifnull(66, -1) AS regClassID,
ifnull(66, -1) AS SCC,
ifnull(66, -1) AS fuelTypeID,
ifnull(66, -1) AS modelYearID,
ifnull(66, -1) AS pollutantID,
ifnull(66, -1) AS processID,
ifnull(66, -1) AS temperature,
ifnull(66, -1) AS relHumidity,
ifnull(66, -1) AS ratePerStart
FROM 123test.rateperstart
;

#ratepervehicle
SELECT  
ifnull(66, -1) AS MOVESScenarioID,
ifnull(66, -1) AS MOVESRunID,
ifnull(66, -1) AS yearID,
ifnull(66, -1) AS monthID,
ifnull(66, -1) AS dayID,
ifnull(66, -1) AS hourID,
ifnull(66, -1) AS zoneID,
ifnull(66, -1) AS pollutantID,
ifnull(66, -1) AS processID,
ifnull(66, -1) AS sourceTypeID,
ifnull(66, -1) AS regClassID,
ifnull(66, -1) AS SCC,
ifnull(66, -1) AS fuelTypeID,
ifnull(66, -1) AS modelYearID,
ifnull(66, -1) AS temperature,
ifnull(66, -1) AS relHumidity,
ifnull(66, -1) AS ratePerVehicle
FROM 123test.ratepervehicle
;
#startpervehicle
SELECT  
ifnull(66, -1) AS MOVESScenarioID,
ifnull(66, -1) AS MOVESRunID,
ifnull(66, -1) AS yearID,
ifnull(66, -1) AS monthID,
ifnull(66, -1) AS dayID,
ifnull(66, -1) AS hourID,
ifnull(66, -1) AS zoneID,
ifnull(66, -1) AS sourceTypeID,
ifnull(66, -1) AS regClassID,
ifnull(66, -1) AS SCC,
ifnull(66, -1) AS fuelTypeID,
ifnull(66, -1) AS modelYearID,
ifnull(66, -1) AS startsPerVehicle
FROM 123test.startpervehicle
;