UPDATE airports SET is_active=0;
UPDATE airports SET is_active=1
    WHERE
        code IN (SELECT DISTINCT origin FROM flights) OR
        code IN (SELECT DISTINCT destination FROM flights);
