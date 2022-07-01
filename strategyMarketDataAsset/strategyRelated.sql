SELECT * FROm structured_value_types

SELECT * FROm structured_values
where structured_value_type_id = 8]\


SELECT* FROM assets
SELECT * FROM 
(SELECT DISTINCT positions.start_date from positions ) x
LEFT JOIN (SELECT start_date, COUNT(*) as count_date FROM market_data
WHERE asset_id =35
GROUP BY start_date) y on x.start_date = y.start_date
WHERE y.start_date is null
order by x.start_date asc

SELECT * FROM market_data 
WHERE asset_id =35 

SELECT * FROM