`load` reads data from the US Bureau of Transportation Statistics On-Time
Flight Performance export files and imports them into a MySQL database.

The files can be downloaded
[here](https://www.transtats.bts.gov/DL_SelectFields.asp?Table_ID=236). Use the
"Prezipped Files".

To run:

```
go build && ./load -address 127.0.0.1 -user flightdb-user -pass flightdb-pass -db flightdb file1.zip file2.zip ...
```

It works by reading the data from the zip files in parallel and merging them
into a CSV file. When the CSV gets large enough, it's pushed to MySQL via `LOAD
DATA LOCAL INFILE`. This is fairly efficient, but it still takes a while.

The program doesn't output anything except errors. Don't be concerned. You can
check the progress by executing `SELECT COUNT(*) FROM flights`, and pass the
time by updating the README.
