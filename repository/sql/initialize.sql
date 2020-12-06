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
    transactions TEXT,
    txCount INTEGER
);

CREATE TABLE
IF NOT EXISTS transactions
(
    txHash TEXT NOT NULL PRIMARY KEY,
    blockHash TEXT,
    pendingNo TEXT,
    totalValue INTEGER,
    fee INTEGER,
    senderAddr TEXT,
    outCount INTEGER,
    outs TEXT,
    timestamp INTEGER
);

CREATE TABLE
IF NOT EXISTS accounts
(
    addr TEXT NOT NULL PRIMARY KEY, 
    balance INTEGER
);
