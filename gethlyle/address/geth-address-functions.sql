CREATE OR REPLACE FUNCTION get_distinct_addresses_by_asset_id(asset_id integer)
-- RETURNS TABLE(distinct_address id)
RETURNS setof geth_addresses
LANGUAGE plpgsql
AS
	$$
	BEGIN
	CREATE TEMP TABLE IF NOT EXISTS temp_table_address AS
			SELECT DISTINCT sender_address_id as address  
			FROM geth_transfers
			WHERE token_address=(SELECT contract_address from assets where id = $1 );
		INSERT INTO temp_table_address
		SELECT DISTINCT to_address_id as address  FROM geth_transfers
		WHERE  token_address= (SELECT contract_address from assets where id = $1);

		  RETURN QUERY SELECT * FROM geth_addresses 
		  WHERE 
		  	id IN (SELECT DISTINCT address from temp_table_address)
		  	-- only eoa address types
		  	AND address_type_id = 84;
		  DROP TABLE temp_table_address;
	END;
$$ 