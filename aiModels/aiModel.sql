CREATE TABLE ai_models
(
  id SERIAL,
  uuid uuid NOT NULL DEFAULT uuid_generate_v4() ,
  name VARCHAR(255) NOT NULL,
  alternate_name VARCHAR(255) NULL,
  url VARCHAR(255) NULL,
  ticker VARCHAR(255) NULL,
  description TEXT NULL,
  ollama_name VARCHAR(255) NULL,
  params_size BIGINT NULL,
  quantiized_size VARCHAR(255) NULL,
  base_model_id INT NULL,
  created_by VARCHAR(255) NOT NULL,
  created_at timestamp NOT NULL,
  updated_by VARCHAR(255) NOT NULL,
  updated_at timestamp NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_base_model
    FOREIGN KEY(base_model_id) 
	  REFERENCES ai_models(id)
)