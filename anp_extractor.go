package main

import (
	"AnpExtractor/anp"
	"AnpExtractor/sheet_file"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"strings"
	"time"
)

var dateFormat = "02/01/2006"

func weekFormat(t time.Time) string {
	return t.Format(dateFormat) + "-" + t.Add(6*24*time.Hour).Format(dateFormat)
}

func main() {

	skip := false
	if len(os.Args) > 1 {
		args := os.Args[1]
		skip = args == "--today"
		if args == "--help" {
			fmt.Printf("Usage %s\n [--today]\n  --today     skips week prompt", os.Args[0])
			os.Exit(0)
		}
	}

	weekNum := anp.TimeToWeek(time.Now())
	if !skip {
		weekNum = promptWeek(weekNum, weekNum)
	}

	queryCount := (len(anp.StateList) + len(anp.CityList)) * len(anp.FuelList)
	bar := pb.StartNew(queryCount)
	cityChan := make(chan [][]string, 10)
	stationChan := make(chan [][]string, 10)

	for fuelId, fuel := range anp.FuelList {
		for _, stateName := range anp.StateList {
			go createCityRows(cityChan, weekNum, fuelId, fuel, stateName)
		}
		for cityId, cityName := range anp.CityList {
			go createStationRows(stationChan, weekNum, fuelId, fuel, cityId, cityName)
		}
	}

	file := sheet_file.CreateFile()
	citySheet := file.Sheets[0]
	stationSheet := file.Sheets[1]
	for i := 0; i < queryCount; i++ {
		select {
		case data := <-cityChan:
			addAllRows(citySheet, data)
		case data := <-stationChan:
			addAllRows(stationSheet, data)
		}
		bar.Increment()
	}
	bar.Finish()

	fileName := fmt.Sprintf("anp-preco-%s.xlsx", strings.ReplaceAll(weekFormat(anp.WeekToTime(weekNum)), "/", "_"))
	if _, err := os.Stat(fileName); err == nil {
		_ = os.Remove(fileName)
	}
	err := file.Save(fileName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File saved as ", fileName)
}

func addAllRows(sheet *xlsx.Sheet, rowList [][]string) {
	for _, col := range rowList {
		row := sheet.AddRow()
		for _, data := range col {
			cell := row.AddCell()
			cell.Value = data
		}
	}
}

func createCityRows(cityChan chan [][]string, week int, fuelId int, fuelName string, state string) {
	rows, err := anp.ReportByState(week, state, fuelId, fuelName)
	if err != nil && rows != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\rCould not get state %s %s pair\n", state, fuelName)
		cityChan <- [][]string{}
	} else {
		cityChan <- rows
	}
}

func createStationRows(stationChan chan [][]string, week int, fuelId int, fuelName string, cityId int, cityName string) {
	rows, err := anp.ReportByCity(week, cityId, cityName, fuelId, fuelName)
	if err != nil && rows != nil {
		_, _ = fmt.Fprintf(os.Stderr, "\rCould not get city %s %s pair\n", cityName, fuelName)
		stationChan <- [][]string{}
	} else {
		stationChan <- rows
	}
}
