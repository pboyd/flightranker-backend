INSERT INTO flights_day (date, carrier, origin, destination, total_flights, delayed_flights)
    SELECT totals.date, totals.carrier, totals.origin, totals.destination, totals.count, delays.count
    FROM (
        SELECT date, carrier, origin, destination, count(*) AS count FROM flights GROUP BY date, carrier, origin, destination
    ) AS totals
    LEFT OUTER JOIN (
        SELECT date, carrier, origin, destination, count(*) AS count FROM flights WHERE scheduled_departure_time <= departure_time AND scheduled_arrival_time <= arrival_time GROUP BY date, carrier, origin, destination
    ) AS delays ON totals.date=delays.date AND totals.carrier=delays.carrier AND totals.origin=delays.origin AND totals.destination=delays.destination;
