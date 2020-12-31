CREATE TABLE
IF NOT EXISTS blocks
(
    hash TEXT NOT NULL PRIMARY KEY,
    height INTEGER,
    prevBlockHash TEXT,
    merkleRoot TEXT,
    extraNonce INTEGER,
    timestamp INTEGER,
    bits INTEGER,
    nonce INTEGER,
    transactions BLOB,
    txCount INTEGER
);

CREATE TABLE
IF NOT EXISTS transactions
(
    txHash TEXT NOT NULL PRIMARY KEY,
    totalValue INTEGER,
    fee INTEGER,
    senderAddr TEXT,
    outCount INTEGER,
    outs BLOB,
    timestamp INTEGER
);

CREATE TABLE
IF NOT EXISTS accounts
(
    addr TEXT NOT NULL PRIMARY KEY, 
    balance INTEGER
);
