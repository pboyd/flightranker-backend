`flightranker-backend` is the back-end code for
[flightranker.com](https://flightranker.com). This was largely an experiment in
software organization, so there are actually three back-ends:

* `backendA` - a flat structure without much attention paid to writing good code
* `backendB` - a hierarchical structure where each dependency is isolated
* `backendC` - based on ideas from [John Ousterhout's A Philosophy of Software Design](https://books.google.com/books?id=pD6-swEACAAJ)

All three are functionally identical. You can read more about it on my
[blog](https://pboyd.io/posts/code-structure-experiment/).

# Set up

## Database

You will need to set up and populate a MySQL database before running either
backend. For local development, try:

```sh
docker run -d -p 3306:3306 --name=flightdb -e MYSQL_ROOT_PASSWORD=flightdb -e MYSQL_USER=flightdb -e MYSQL_PASSWORD=flightdb -e MYSQL_DATABASE=flightdb mysql/mysql-server:5.7
```

The files in the `sql` directory will set up the schema and populate the
`airports` and `carriers` tables:

```sh
cat sql/*.sql | mysql -uflightdb -pflightdb -h 127.0.0.1 flightdb
```

The `load` program will populate the `flights` table. See `load/README.md` for
details. The test data was generated from a database loaded with only the
following months:

* [January 2019](https://transtats.bts.gov/PREZIP/On_Time_Reporting_Carrier_On_Time_Performance_1987_present_2019_1.zip)
* [February 2019](https://transtats.bts.gov/PREZIP/On_Time_Reporting_Carrier_On_Time_Performance_1987_present_2019_2.zip)
* [March 2019](https://transtats.bts.gov/PREZIP/On_Time_Reporting_Carrier_On_Time_Performance_1987_present_2019_3.zip)

After the data has been loaded, run the SQL files in `sql/updates`:

```sh
cat sql/updates/*.sql | mysql -uflightdb -pflightdb -h 127.0.0.1 flightdb
```

## Configuration

All configuration is read from environment variables, each backend reads the
following:

* `MYSQL_ADDRESS`: Network address for the database (e.g. `127.0.0.1:3306`)
* `MYSQL_DATABASE`: Database name
* `MYSQL_USER`: Username for MySQL
* `MYSQL_PASS`: Password for MySQL
* `CORS_ALLOW_ORIGIN`: Value to return in the `Access-Control-Allow-Origin`
  header. If this variable is not set, the header is omitted.

## Running

Each backend can be run in the usual Go way:

```sh
cd backendA && go install && backendA
```

```sh
cd backendB && go install && backendB
```

```sh
cd backendC && go install && backendC
```

## Tests

Database tests in all the backends require the same set of environment
variables as the backend server.

```
MYSQL_USER=flightdb MYSQL_PASS=flightdb MYSQL_ADDRESS=127.0.0.1:3306 MYSQL_DATABASE=flightdb go test ./...
```

## Docker

There is a `Dockerfile` in root of the repository that can be used for either backend.

For `backendA`:

```sh
docker build . -t flightranker-backend-a --build-arg which=backendA
```

For `backendB`:

```sh
docker build . -t flightranker-backend-b --build-arg which=backendB
```

For `backendC`:

```sh
docker build . -t flightranker-backend-c --build-arg which=backendC
```

# Credits

On-Time performance data are provided by the [US Bureau of Transportation Statistics](https://www.transtats.bts.gov).

The list of airports was taken from [stat-computing.org](http://stat-computing.org/dataexpo/2009/supplemental-data.html).
