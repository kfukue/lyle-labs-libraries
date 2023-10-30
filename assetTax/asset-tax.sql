DROP TABLE IF EXISTS asset_taxes CASCADE;

CREATE TABLE asset_taxes
(
  tax_id INT NOT NULL,
  asset_id INT NOT NULL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  tax_rate_override NUMERIC NULL,
  tax_rate_type_id INT NULL,
  description TEXT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(tax_id, asset_id),
  CONSTRAINT fk_tax_rate_type_id FOREIGN KEY(tax_rate_type_id) REFERENCES structured_values(id),
  CONSTRAINT fk_tax FOREIGN KEY(tax_id) REFERENCES taxes(id),
  CONSTRAINT fk_asset FOREIGN KEY(asset_id) REFERENCES assets(id)
);

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-user";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-user";

GRANT USAGE ON ALL SEQUENCES IN SCHEMA public TO "asset-tracker-api";
GRANT SELECT,INSERT,UPDATE,DELETE  ON ALL TABLES IN SCHEMA public TO "asset-tracker-api";