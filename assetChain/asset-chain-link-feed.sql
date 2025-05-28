-- SQL schema for asset_chain_link_data_feed association table
CREATE TABLE IF NOT EXISTS asset_chain (
    asset_id INT NOT NULL,
    chain_id INT NOT NULL,
    chainlink_data_feed_contract_address TEXT NOT NULL,
    created_by TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by TEXT NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (asset_id, chain_id),
    FOREIGN KEY (asset_id) REFERENCES assets(id),
    FOREIGN KEY (chain_id) REFERENCES chains(id)
); 


CREATE INDEX asset_chain_asset_id ON asset_chain(asset_id);
CREATE INDEX asset_chain_chain_id ON asset_chain(chain_id);




GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-user";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-user";

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-api";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-api";

