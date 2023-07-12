package main

import (
	"encoding/csv"
	"fmt"
	"github.com/kardianos/osext"
	"os"
)

// get the absolute path from current executable
func getAbsolutePath() string {
	folderPath, pathPrr := osext.ExecutableFolder()
	if pathPrr != nil {
		fmt.Println("getAbsolutePath error", pathPrr)
	}
	fmt.Println("getAbsolutePath print dir ", folderPath)
	return folderPath
}

func clearMap(m map[wideTableShapeStruct]float64) {
	for key := range m {
		delete(m, key)
	}
}

// copy map (destination, source)
func mapCopy[M1, M2 ~map[K]V, K comparable, V any](dst M1, src M2) {
	for k, v := range src {
		dst[k] = v
	}
}

// remove duplicated string in a slice, and return updated one
func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// remove duplicated integer in a slice, and return updated one
func removeDuplicateInt(intSlice []int) []int {
	allKeys := make(map[int]bool)
	list := []int{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// remove duplicated generic type in a slice, and return updated one
func removeDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// count how many characters that the longest element in the []string has, then x 10 pixel and return the resolution
func getColWidths(data [][]string) []float32 {
	res := make([]float32, 0)
	for _, row := range data {
		for i, col := range row {
			cur := float32(len(fmt.Sprint(col)) * 10)
			if len(res) <= i {
				res = append(res, cur)
			} else {
				if res[i] <= cur {
					res[i] = cur
				}
			}
		}
	}
	return res
}

// take [][]string and download csv file in the same folder that main.go is located
// TODO make download navigation, something that open window file explorer and download with saved path
func csvExport(data [][]string) error {
	file, err := os.Create("SummaryReportTool_Result.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		if err = writer.Write(value); err != nil {
			return err // let's return errors if necessary, rather than having a one-size-fits-all error handler
		}
	}
	return nil
}

// RemoveElementFromSlice take []string, and bad_string, scan the whole slice then delete the first bad_string
func RemoveElementFromSlice(l []string, item string) []string {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

// input [][]string, output decoded [][]string IDs in English
// decoder can be found in decoder.go, it predefined the key[value] combinations
func decodeDataTable(Matrix [][]string) {
	//outer loop for each column
	for col := 0; col < len(Matrix[0]); col++ {
		//inner loop for each row
		switch Matrix[0][col] {
		case "dayID": // if the header of column is dayID, then we use the dayID map for decode.
			for row := 1; row < len(Matrix); row++ { //skip the first element that must be header
				val, ok := decoder_dayOfAnyWeek[Matrix[row][col]]
				// If the key exists
				if ok {
					Matrix[row][col] = val //assign the English words back to original matrix
				}
			}
		case "stateID": //same logic for the rest of the decoders
			for row := 1; row < len(Matrix); row++ {
				val, ok := decoder_state[Matrix[row][col]]
				if ok {
					Matrix[row][col] = val
				}
			}
		case "countyID":
			for row := 1; row < len(Matrix); row++ {
				val, ok := decoder_county[Matrix[row][col]]
				if ok {
					Matrix[row][col] = val
				}
			}
		case "pollutantID":
			for row := 1; row < len(Matrix); row++ {
				val, ok := decoder_pollutant[Matrix[row][col]]
				if ok {
					Matrix[row][col] = val
				}
			}
		case "processID":
			for row := 1; row < len(Matrix); row++ {
				val, ok := decoder_emissionProcess[Matrix[row][col]]
				if ok {
					Matrix[row][col] = val
				}
			}
		case "sourceTypeID":
			for row := 1; row < len(Matrix); row++ {
				val, ok := decoder_sourceUseType[Matrix[row][col]]
				if ok {
					Matrix[row][col] = val
				}
			}
		case "regClassID":
			for row := 1; row < len(Matrix); row++ {
				val, ok := decoder_regulatoryClass[Matrix[row][col]]
				if ok {
					Matrix[row][col] = val
				}
			}
		case "fuelTypeID":
			for row := 1; row < len(Matrix); row++ {
				val, ok := decoder_fuelType[Matrix[row][col]]
				if ok {
					Matrix[row][col] = val
				}
			}
		case "fuelSubTypeID":
			for row := 1; row < len(Matrix); row++ {
				val, ok := decoder_fuelSubType[Matrix[row][col]]
				if ok {
					Matrix[row][col] = val
				}
			}
		case "roadTypeID":
			for row := 1; row < len(Matrix); row++ {
				val, ok := decoder_roadtype[Matrix[row][col]]
				if ok {
					Matrix[row][col] = val
				}
			}
		case "engTechID":
			for row := 1; row < len(Matrix); row++ {
				val, ok := decoder_engineTech[Matrix[row][col]]
				if ok {
					Matrix[row][col] = val
				}
			}
		case "sectorID":
			for row := 1; row < len(Matrix); row++ {
				val, ok := decoder_sector[Matrix[row][col]]
				if ok {
					Matrix[row][col] = val
				}
			}
		case "hpID":
			for row := 1; row < len(Matrix); row++ {
				val, ok := decoder_nrHpRangeBin[Matrix[row][col]]
				if ok {
					Matrix[row][col] = val
				}
			}
		case "activityTypeID":
			for row := 1; row < len(Matrix); row++ {
				val, ok := decoder_activityType[Matrix[row][col]]
				if ok {
					Matrix[row][col] = val
				}
			}
		default:
			//DO NOTHING for non-ID fields
			fmt.Println("decode NOTHING for non-ID fields, column name: ", Matrix[0][col])
		}
	}
}
