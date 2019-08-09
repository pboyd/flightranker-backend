-- Sets is_active on airports with flight data.
--
-- Consequently, only run this _after_ importing flight data!
UPDATE airports SET is_active=0;
UPDATE airports SET is_active=1
    WHERE
        code IN (SELECT DISTINCT origin FROM flights) OR
        code IN (SELECT DISTINCT destination FROM flights);
