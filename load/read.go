package main

import (
	"archive/zip"
	"encoding/csv"
	"io"
	"log"
	"strings"
)

func readZipFile(path string) (<-chan record, error) {
	zipReader, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}

	ch := make(chan record, batchSize)

	go func() {
		defer zipReader.Close()
		defer close(ch)

		for _, f := range zipReader.File {
			if strings.HasSuffix(f.Name, ".csv") {
				r, err := f.Open()
				if err != nil {
					log.Printf("%s %s: open failed: %v", path, f.Name, err)
					continue
				}

				err = emitCSVRecords(ch, r)
				if err != nil {
					log.Printf("%s %s: read failed: %v", path, f.Name, err)
				}

				r.Close()
			}
		}
	}()

	return ch, nil
}

func emitCSVRecords(out chan record, r io.Reader) error {
	csvReader := csv.NewReader(r)
	csvReader.ReuseRecord = true

	header, err := csvReader.Read()
	if err != nil {
		return err
	}
	header = dup(header)

	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("failed to read row: %v", err)
			break
		}

		out <- newRecord(row, header)
	}

	return nil
}

func dup(s []string) []string {
	s2 := make([]string, len(s))
	for i := range s {
		s2[i] = s[i]
	}
	return s2
}
