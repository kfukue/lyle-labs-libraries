SELECT *
FROM trades

ALTER SEQUENCE positions_id_seq RESTART WITH 1
ALTER SEQUENCE trades_id_seq RESTART WITH 1
ALTER SEQUENCE accounts_id_seq RESTART WITH 1

SELECT *
FROM accounts
ORDER BY alternate_name
WHERE address = ''

SELECT *
FROM trades
where trade_date between '2017-12-27' AND '2017-12-29'
-- AND asset_id = 35
ORDER BY trade_date, asset_id

from_account_id = 87
SELECT * FRom accounts
where id in
(4,38)
SELECT *
FROM positions
WHERE
account_id = 87
  AND base_asset_id = 35
-- -- AND quantity > 0
ORDER BY start_date desc, base_asset_id


SELECT *
From trades
where id = 1795

SELECT *
FROM positions
WHERE
account_id = 94
  AND base_asset_id = 35
  AND start_date BETWEEN '2017-11-1' AND '2018-1-1'
ORDER BY start_date

SELECT id, alternate_name
FROm accounts
order by alternate_name asc
SELECT *
from assets
where id = 29


SELECT account_id, quantity, *
FROM positions
WHERE
account_id IN (85,	88,	96,	83,	35,	86,	138,	36,	84,	78,	38,	37,	87,	70,	39,	99,	54,	55,	56,	57,	42,	118,69,	4)


  AND base_asset_id = 35
  AND start_date = '2020-12-08'
-- -- AND quantity > 0
ORDER BY account_id
SELECT *
FRom accounts
where id =69

ORDER BY start_date desc, base_asset_id
 LIMIT 100
 
 SELECT * FROM accounts
WHERE alternate_name LIKE '%trezor%'
 ORDER BY alternate_name

SELECT *
FROM accounts
ORDER BY alternate_name
SELECT *
FROM trades
where 
from_account_id in (69)
  OR to_account_id in (69)

SELECT trade_date, to_quantity, *
FROM trades
where 
from_account_id in (4)
  AND (asset_id IN ( select id
  from assets
  where base_asset_id = 35)OR asset_id IN ( 35))
  and is_active = true
-- OR to_account_id in (69)



SELECT trade_date, to_quantity, *
FROM trades
where 
 to_account_id in (4)
  AND (asset_id IN ( select id
  from assets
  where base_asset_id = 35)OR asset_id IN ( 35))
  and is_active = true

SELECT *
FROM assets
SELECT *
From trades
where
trade_date BETWEEN '2017-11-30' AND '2017-12-1'

SELECT trade_date, from_quantity, *
FROM trades
where 
from_account_id in (69)
  AND (asset_id IN ( 35))
  and is_active = true
ORDER BY trade_date
-- OR to_account_id in (69)


SELECT trade_date, to_quantity, *
FROM trades
where 
 to_account_id in (69)
  AND (asset_id IN ( 35))
  and is_active = true


TRUNCATE TABLE accounts
TRUNCATE TABLE trades
TRUNCATE TABLE positions

SELECT *
FROM pg_settings
WHERE pending_restart = true;

SELECT *
FROM pg_stat_user_tables

SELECT *
FROm trades

EXPLAIN
ANALYZE
SELECT *
FROM positions;
