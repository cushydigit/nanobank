CREATE TYPE transaction_status AS ENUM ('pending', 'confirmed', 'canceled');


CREATE TABLE IF NOT EXISTS transactions (
  id TEXT PRIMARY KEY,
  from_user_id TEXT NOT NULL,
  to_user_id TEXT NOT NULL,
  amount BIGINT NOT NULL,
  status transaction_status NOT NULL,
  confirmation_token UUID NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now() 
);  
