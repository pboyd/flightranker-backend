package main

import (
	"flag"
	"log"
	"sync"
)

const batchSize = 1000000

func main() {
	var (
		address string
		user    string
		pass    string
		db      string
	)

	flag.StringVar(&address, "address", "127.0.0.1", "MySQL address")
	flag.StringVar(&user, "user", "flightdb", "MySQL user")
	flag.StringVar(&pass, "pass", "flightdb", "MySQL password")
	flag.StringVar(&db, "db", "flightdb", "MySQL database name")
	flag.Parse()

	conn, err := connect(address, user, pass, db)
	if err != nil {
		log.Fatalf("could not connect to mysql: %v", err)
	}

	recordChs := make([]<-chan record, flag.NArg())
	for i, path := range flag.Args() {
		recordChs[i], err = readZipFile(path)
		if err != nil {
			log.Fatalf("error reading %q: %v", path, err)
		}
	}
	ch := merge(recordChs)

	err = loadRecords(conn, ch)
	if err != nil {
		log.Fatalf("failed to load records: %v", err)
	}
}

func merge(chs []<-chan record) <-chan record {
	out := make(chan record)
	var wg sync.WaitGroup
	wg.Add(len(chs))

	for _, ch := range chs {
		go func(ch <-chan record) {
			defer wg.Done()

			for r := range ch {
				out <- r
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
