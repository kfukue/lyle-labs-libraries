
SELECT agg.start_date, assets.name, assets.id, agg.quantity
from assets
  LEFT JOIN (
SELECT SUM(positions.quantity) as quantity, positions.base_asset_id as base_asset_id, positions.start_date
  FROM positions

  WHERE
-- account_id IN (85,	88,	96,	83,	35,	86,	138,	36,	84,	78,	38,	37,	87,	70,	39,	99,	54,	55,	56,	57,	42,	118,69,	4)


-- AND base_asset_id = 35
 start_date = '2020-12-08'
    AND quantity > 0
  GROUP BY positions.base_asset_id, start_date) agg
  on assets.id = agg.base_asset_id
WHERE agg.quantity <> 0
ORDER BY assets.name