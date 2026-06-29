-- 3. Select / Stored Procedure Script
-- Select: Get country by person name
SELECT Country FROM Person WHERE Name = 'Adam';
SELECT Country FROM Person WHERE Name = 'John';
SELECT Country FROM Person WHERE Name = 'Henry';
SELECT Country FROM Person WHERE Name = 'Dominic';

-- Select all persons
SELECT * FROM Person;

-- Stored Procedure equivalent: GetCountry by Name
-- SQLite doesn't support stored procedures, so we use a parameterized query instead.
-- In PostgreSQL / MySQL / SQL Server, the equivalent would be:
--
-- PostgreSQL:
-- CREATE OR REPLACE FUNCTION GetCountry(p_name VARCHAR)
-- RETURNS VARCHAR AS $$
--     SELECT Country FROM Person WHERE Name = p_name;
-- $$ LANGUAGE sql;
--
-- MySQL:
-- DELIMITER //
-- CREATE PROCEDURE GetCountry(IN p_name VARCHAR(100))
-- BEGIN
--     SELECT Country FROM Person WHERE Name = p_name;
-- END //
-- DELIMITER ;
--
-- SQL Server:
-- CREATE PROCEDURE GetCountry
--     @Name NVARCHAR(100)
-- AS
--     SELECT Country FROM Person WHERE Name = @Name;
-- GO;
