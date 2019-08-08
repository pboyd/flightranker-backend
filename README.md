`flightranker-backend` is the back-end code for
[flightranker.com](https://flightranker.com). This was largely an experiment in
software organization, so there are actually two back-ends:

* `backendA` - a flat structure without much attention paid to writing good code
* `backendB` - a hierarchical structure where each dependency is isolated

Both are functionally identical. You can read more about it on my blog (once I
get around to writing it, that is).

# Running

To run, you'll need a MySQL database. Run the files in the `sql` directory to
set up the schema and populate the `airlines` and `airports` tables. Then
import flight data via the program in the `load` directory (see
`load/README.md` for details).

# Credits

On-Time performance data are provided by the [US Bureau of Transportation Statistics](https://www.transtats.bts.gov).

The list of airports was taken from [stat-computing.org](http://stat-computing.org/dataexpo/2009/supplemental-data.html).
