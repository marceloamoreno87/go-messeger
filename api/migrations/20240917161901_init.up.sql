-- cria tabela de devices onde tenha numero de telefone jid e id
CREATE TABLE devices (
    id SERIAL PRIMARY KEY,
    jid VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- cria index 
CREATE UNIQUE INDEX devices_jid_index ON devices (jid);

CREATE UNIQUE INDEX devices_phone_number_index ON devices (phone_number);