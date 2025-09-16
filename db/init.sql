CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY, -- 帳號 ID (自動增加)
    name VARCHAR(100) NOT NULL, -- 帳號名稱
    balance NUMERIC(15,2) NOT NULL DEFAULT 0, -- 帳號餘額
    created_at TIMESTAMP DEFAULT NOW(), -- 建立時間
    updated_at TIMESTAMP DEFAULT NOW() -- 更新時間
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
