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
FROM 123rate.rateperdistance
;

# rateperhour no number, maybe I should not include all hours?
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
FROM 123rate.rateperhour
;


#rateperprofile no number
SELECT  
ifnull(MOVESScenarioID, -1) AS MOVESScenarioID,
ifnull(MOVESRunID, -1) AS MOVESRunID,
ifnull(temperatureProfileID, -1) AS temperatureProfileID,
ifnull(yearID, -1) AS yearID,
ifnull(dayID, -1) AS dayID,
ifnull(hourID, -1) AS hourID,
ifnull(pollutantID, -1) AS pollutantID,
ifnull(processID, -1) AS processID,
ifnull(sourceTypeID, -1) AS sourceTypeID,
ifnull(regClassID, -1) AS regClassID,
ifnull(SCC, -1) AS SCC,
ifnull(fuelTypeID, -1) AS fuelTypeID,
ifnull(modelYearID, -1) AS modelYearID,
ifnull(temperature, -1) AS temperature,
ifnull(relHumidity, -1) AS relHumidity,
ifnull(ratePerVehicle, -1) AS ratePerVehicle
FROM 123rate.rateperprofile
;

#rateperstart
SELECT  
ifnull(MOVESScenarioID, -1) AS MOVESScenarioID,
ifnull(MOVESRunID, -1) AS MOVESRunID,
ifnull(yearID, -1) AS yearID,
ifnull(monthID, -1) AS monthID,
ifnull(dayID, -1) AS dayID,
ifnull(hourID, -1) AS hourID,
ifnull(zoneID, -1) AS zoneID,
ifnull(sourceTypeID, -1) AS sourceTypeID,
ifnull(regClassID, -1) AS regClassID,
ifnull(SCC, -1) AS SCC,
ifnull(fuelTypeID, -1) AS fuelTypeID,
ifnull(modelYearID, -1) AS modelYearID,
ifnull(pollutantID, -1) AS pollutantID,
ifnull(processID, -1) AS processID,
ifnull(temperature, -1) AS temperature,
ifnull(relHumidity, -1) AS relHumidity,
ifnull(ratePerStart, -1) AS ratePerStart
FROM 123rate.rateperstart
;

#ratepervehicle
SELECT  
ifnull(MOVESScenarioID, -1) AS MOVESScenarioID,
ifnull(MOVESRunID, -1) AS MOVESRunID,
ifnull(yearID, -1) AS yearID,
ifnull(monthID, -1) AS monthID,
ifnull(dayID, -1) AS dayID,
ifnull(hourID, -1) AS hourID,
ifnull(zoneID, -1) AS zoneID,
ifnull(pollutantID, -1) AS pollutantID,
ifnull(processID, -1) AS processID,
ifnull(sourceTypeID, -1) AS sourceTypeID,
ifnull(regClassID, -1) AS regClassID,
ifnull(SCC, -1) AS SCC,
ifnull(fuelTypeID, -1) AS fuelTypeID,
ifnull(modelYearID, -1) AS modelYearID,
ifnull(temperature, -1) AS temperature,
ifnull(relHumidity, -1) AS relHumidity,
ifnull(ratePerVehicle, -1) AS ratePerVehicle
FROM 123rate.ratepervehicle
;
#startpervehicle
SELECT  
ifnull(MOVESScenarioID, -1) AS MOVESScenarioID,
ifnull(MOVESRunID, -1) AS MOVESRunID,
ifnull(yearID, -1) AS yearID,
ifnull(monthID, -1) AS monthID,
ifnull(dayID, -1) AS dayID,
ifnull(hourID, -1) AS hourID,
ifnull(zoneID, -1) AS zoneID,
ifnull(sourceTypeID, -1) AS sourceTypeID,
ifnull(regClassID, -1) AS regClassID,
ifnull(SCC, -1) AS SCC,
ifnull(fuelTypeID, -1) AS fuelTypeID,
ifnull(modelYearID, -1) AS modelYearID,
ifnull(startsPerVehicle, -1) AS startsPerVehicle
FROM 123rate.startspervehicle
;