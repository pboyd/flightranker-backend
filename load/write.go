package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

var columns = []string{
	"id",
	"date",
	"departure_time",
	"scheduled_departure_time",
	"arrival_time",
	"scheduled_arrival_time",
	"carrier",
	"flight_number",
	"tail_number",
	"origin",
	"destination",
	"cancelled",
	"cancellation_code",
	"diverted",
	"elapsed_time",
	"schedule_time",
	"air_time",
	"taxi_in_time",
	"taxi_out_time",
	"wheels_off_time",
	"wheels_on_time",
	"arrival_delay",
	"departure_delay",
	"carrier_delay",
	"weather_delay",
	"nas_delay",
	"security_delay",
	"late_aircraft_delay",
}

func connect(address, user, pass, dbName string) (*sql.DB, error) {
	dsn := (&mysql.Config{
		User:   user,
		Passwd: pass,
		Net:    "tcp",
		Addr:   address,
		DBName: dbName,

		AllowNativePasswords: true,
		ParseTime:            true,
	}).FormatDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func loadRecords(db *sql.DB, records <-chan record) error {
	fh, csv, err := makeCSV()
	if err != nil {
		return err
	}
	defer cleanupTemp(fh)

	i := 0
	for record := range records {
		err = writeSingleRecord(csv, &record)
		if err != nil {
			log.Printf("failed to write record: %v %v", err, record)
		}

		i++
		if i > 0 && i%batchSize == 0 {
			csv.Flush()
			err = insertTempFile(db, fh.Name())
			if err != nil {
				return err
			}
			cleanupTemp(fh)

			fh, csv, err = makeCSV()
			if err != nil {
				return err
			}
		}
	}

	csv.Flush()
	return insertTempFile(db, fh.Name())
}

func makeCSV() (*os.File, *csv.Writer, error) {
	fh, err := ioutil.TempFile("", "flights-load-")
	if err != nil {
		return nil, nil, err
	}

	w := csv.NewWriter(fh)
	err = w.Write(columns)
	if err != nil {
		cleanupTemp(fh)
		return nil, nil, err
	}

	return fh, w, nil
}

func cleanupTemp(fh *os.File) {
	fh.Close()
	os.Remove(fh.Name())
}

func insertTempFile(db *sql.DB, path string) error {
	mysql.RegisterLocalFile(path)
	_, err := db.Exec(fmt.Sprintf(`
		LOAD DATA LOCAL INFILE '%s' INTO TABLE flights
			FIELDS TERMINATED BY ','
			IGNORE 1 LINES
		`, path))
	mysql.DeregisterLocalFile(path)
	return err
}

func writeSingleRecord(csv *csv.Writer, record *record) error {
	return csv.Write([]string{
		"0",
		record.Date,                               // date
		strconv.Itoa(record.DepTime),              // departure_time
		strconv.Itoa(record.ScheduledDepTime),     // scheduled_departure_time
		strconv.Itoa(record.ArrTime),              // arrival_time
		strconv.Itoa(record.ScheduledArrTime),     // scheduled_arrival_time
		record.Airline,                            // carrier
		record.FlightNum,                          // flight_number
		record.TailNum,                            // tail_number
		record.Origin,                             // origin
		record.Dest,                               // destination
		strconv.FormatBool(record.Cancelled),      // cancelled
		record.CancellationCode,                   // cancellation_code
		strconv.FormatBool(record.Diverted),       // diverted
		strconv.Itoa(record.ActualElapsedTime),    // elapsed_time
		strconv.Itoa(record.ScheduledElapsedTime), // schedule_time
		strconv.Itoa(record.AirTime),              // air_time
		strconv.Itoa(record.TaxiIn),               // taxi_in_time
		strconv.Itoa(record.TaxiOut),              // taxi_out_time
		strconv.Itoa(record.WheelsOff),            // wheels_off_time
		strconv.Itoa(record.WheelsOn),             // wheels_on_time
		strconv.Itoa(record.ArrDelay),             // arrival_delay
		strconv.Itoa(record.DepDelay),             // departure_delay
		strconv.Itoa(record.CarrierDelay),         // carrier_delay
		strconv.Itoa(record.WeatherDelay),         // weather_delay
		strconv.Itoa(record.NASDelay),             // nas_delay
		strconv.Itoa(record.SecurityDelay),        // security_delay
		strconv.Itoa(record.LateAircraftDelay),    // late_aircraft_delay
	})
}
