package backendtest

var StandardTestQueries = []string{
	`{airport(code:"JFK"){code,name,city,state}}`,
	`{airport(code:"LAX"){code,name,city,state}}`,
	`{airportList(term:"vegas"){code,name,city,state}}`,
	`{airportList(term:"jack"){code,name,city,state}}`,
	`{flightStatsByAirline(origin:"JFK",destination:"LAX"){airline,onTimePercentage,lastFlight}}`,
	`{flightStatsByAirline(origin:"JFK",destination:"LAX"){airline,onTimePercentage,lastFlight},origin:airport(code:"JFK"){code,name,city,state},destination:airport(code:"LAX"){code,name,city,state}}`,
	//`{dailyFlightStats(origin:"JFK",destination:"LAX"){airline,rows{date,onTimePercentage}}}`,
	//`{monthlyFlightStats(origin:"JFK",destination:"LAX"){airline,rows{date,onTimePercentage}}}`,

	"",
	"{}",
	`{someGarbage(code:"LAX"){code,name,city,state}}`,
	`{airport(code:"FOUR"){city,state}}`,
	`{airportList(term:";"){city,state}}`,
	`{flightStatsByAirline(origin:"FOUR",destination:"LAX"){airline}}`,
	`{flightStatsByAirline(origin:"JFK",destination:"FOUR"){airline}}`,
	`{flightStatsByAirline(origin:"FOUR",destination:"ABCD"){airline}}`,
	`{flightStatsByAirline(origin:"FOUR",destination:"ABCD"){airline}}`,
}
