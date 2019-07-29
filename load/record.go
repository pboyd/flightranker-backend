package main

import (
	"math"
	"strconv"
)

type record struct {
	Date                 string
	DepTime              int
	ScheduledDepTime     int
	ArrTime              int
	ScheduledArrTime     int
	Airline              string
	FlightNum            string
	TailNum              string
	ActualElapsedTime    int
	ScheduledElapsedTime int
	AirTime              int
	ArrDelay             int
	DepDelay             int
	Origin               string
	Dest                 string
	TaxiIn               int
	TaxiOut              int
	WheelsOff            int
	WheelsOn             int
	Cancelled            bool
	CancellationCode     string
	Diverted             bool
	CarrierDelay         int
	WeatherDelay         int
	NASDelay             int
	SecurityDelay        int
	LateAircraftDelay    int
}

func newRecord(row, header []string) record {
	if len(row) != len(header) {
		// This is supposed to be handled by the caller
		panic("row length != header length")
	}

	r := record{}

	for i, col := range header {
		switch col {
		case "FlightDate":
			r.Date = row[i]
		case "DepTime":
			r.DepTime, _ = strconv.Atoi(row[i])
		case "CRSDepTime":
			r.ScheduledDepTime, _ = strconv.Atoi(row[i])
		case "ArrTime":
			r.ArrTime, _ = strconv.Atoi(row[i])
		case "CRSArrTime":
			r.ScheduledArrTime, _ = strconv.Atoi(row[i])
		case "Reporting_Airline":
			r.Airline = row[i]
		case "Flight_Number_Reporting_Airline":
			r.FlightNum = row[i]
		case "Tail_Number":
			r.TailNum = row[i]
		case "ActualElapsedTime":
			r.ActualElapsedTime = readDecimalInt(row[i])
		case "CRSElapsedTime":
			r.ScheduledElapsedTime = readDecimalInt(row[i])
		case "AirTime":
			r.AirTime = readDecimalInt(row[i])
		case "ArrDelay":
			r.ArrDelay = readDecimalInt(row[i])
		case "DepDelay":
			r.DepDelay = readDecimalInt(row[i])
		case "Origin":
			r.Origin = row[i]
		case "Dest":
			r.Dest = row[i]
		case "TaxiIn":
			r.TaxiIn = readDecimalInt(row[i])
		case "TaxiOut":
			r.TaxiOut = readDecimalInt(row[i])
		case "Cancelled":
			r.Cancelled = readBool(row[i])
		case "CancellationCode":
			r.CancellationCode = row[i]
		case "Diverted":
			r.Diverted = readBool(row[i])
		case "CarrierDelay":
			r.CarrierDelay = readDecimalInt(row[i])
		case "WeatherDelay":
			r.WeatherDelay = readDecimalInt(row[i])
		case "NASDelay":
			r.NASDelay = readDecimalInt(row[i])
		case "SecurityDelay":
			r.SecurityDelay = readDecimalInt(row[i])
		case "LateAircraftDelay":
			r.LateAircraftDelay = readDecimalInt(row[i])
		case "WheelsOff":
			r.WheelsOff, _ = strconv.Atoi(row[i])
		case "WheelsOn":
			r.WheelsOn, _ = strconv.Atoi(row[i])
		default:
			//log.Print("unhandled column: " + col)
		}
	}
	return r
}

func readBool(v string) bool {
	b, _ := strconv.Atoi(v)
	return b == 1
}

func readDecimalInt(v string) int {
	f, _ := strconv.ParseFloat(v, 64)
	return int(math.Round(f))
}
