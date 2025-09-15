CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    balance NUMERIC(15,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,                -- 流水號
    account_id INT NOT NULL REFERENCES accounts(id), -- 對應哪個帳號
    type INT NOT NULL,                    -- 1=提款, 2=存款
    amount NUMERIC(20,2) NOT NULL,        -- 金額
    ref_id        VARCHAR(50) NOT NULL, -- 關聯 ID
    description VARCHAR(255),             -- 備註
    created_at TIMESTAMP NOT NULL DEFAULT NOW() -- 交易時間
);
