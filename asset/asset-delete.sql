SELECT * FROM assets
where id = 510
SELECT * FROM asset_sources
WHERE asset_id = 510 and source_id = 3
ROLLBACK 
BEGIN TRANSACTION;
DELETE FROM market_data_quotes
WHERE base_asset_id = 510;
DELETE FROM market_data_jobs mj
WHERE mj.market_data_id IN (SELECT id FROM  market_data m
WHERE m.asset_id = 510);
DELETE FROM market_data
WHERE asset_id = 510;
DELETE FROM assets
WHERE id = 510;
COMMIT
SELECT * FROM market_data_jobs mj
JOIN market_data m ON mj.market_data_id = m.id
WHERE m.asset_id = 510