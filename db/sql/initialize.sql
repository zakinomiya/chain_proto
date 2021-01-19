CREATE TABLE
IF NOT EXISTS blocks
(
    hash TEXT NOT NULL PRIMARY KEY,
    height INTEGER UNIQUE,
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
    txType TEXT, 
    totalValue TEXT,
    fee TEXT,
    senderAddr TEXT,
    outCount INTEGER,
    outs BLOB,
    signature TEXT,
    timestamp INTEGER
);

CREATE TABLE
IF NOT EXISTS accounts
(
    addr TEXT NOT NULL PRIMARY KEY, 
    balance TEXT
);
