CREATE TABLE carriers (
    code VARCHAR(6),
    name VARCHAR(128),

    PRIMARY KEY (code)
);

CREATE TABLE airports (
    code CHAR(3),
    name VARCHAR(64),
    city VARCHAR(64),
    state VARCHAR(2),
    lat DECIMAL(10, 8),
    lng DECIMAL(11, 8),
    is_active BOOLEAN,

    PRIMARY KEY (code)
);

CREATE TABLE flights (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    date DATE NOT NULL,
    departure_time TIME,
    scheduled_departure_time TIME NOT NULL,
    arrival_time TIME,
    scheduled_arrival_time TIME NOT NULL,

    carrier VARCHAR(6),
    flight_number CHAR(4),
    tail_number CHAR(6),

    origin CHAR(3),
    destination CHAR(3),

    cancelled BOOLEAN,
    cancellation_code CHAR(1),
    diverted BOOLEAN,

    elapsed_time SMALLINT,
    schedule_time SMALLINT,
    air_time SMALLINT,
    taxi_in_time SMALLINT,
    taxi_out_time SMALLINT,
    wheels_off_time TIME,
    wheels_on_time TIME,

    arrival_delay SMALLINT,
    departure_delay SMALLINT,
    carrier_delay SMALLINT,
    weather_delay SMALLINT,
    nas_delay SMALLINT,
    security_delay SMALLINT,
    late_aircraft_delay SMALLINT,

    PRIMARY KEY (id),
    FOREIGN KEY (carrier) REFERENCES carriers(code),
    FOREIGN KEY (origin) REFERENCES airports(code),
    FOREIGN KEY (destination) REFERENCES airports(code),

    INDEX flight_number_idx (flight_number),
    INDEX carrier_idx (carrier),
    INDEX origin_idx (origin),
    INDEX destination_idx (destination),
    INDEX date_idx (date)
);

CREATE TABLE flights_day (
    date DATE NOT NULL,
    carrier VARCHAR(6),
    origin CHAR(3),
    destination CHAR(3),

    total_flights SMALLINT,
    delayed_flights SMALLINT,

    PRIMARY KEY (date, carrier, origin, destination),
    FOREIGN KEY (carrier) REFERENCES carriers(code),
    FOREIGN KEY (origin) REFERENCES airports(code),
    FOREIGN KEY (destination) REFERENCES airports(code),

    INDEX carrier_idx (carrier),
    INDEX origin_idx (origin),
    INDEX destination_idx (destination),
    INDEX date_idx (date)
);
