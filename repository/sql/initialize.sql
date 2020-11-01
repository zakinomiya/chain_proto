CREATE TABLE
IF NOT EXISTS blocks
(
    height INTEGER NOT NULL PRIMARY KEY,
    hash STRING,
    prevBlockHash BLOB,
    merkleRoot BLOB,
    extraNonce INTEGER,
    timestamp INTEGER,
    bits INTEGER,
    nonce INTEGER,
);

CREATE TABLE
IF NOT EXISTS transactions
(
    txHash BLOB NOT NULL PRIMARY KEY,
    blockHash BLOB,
    index INTEGER,
    pendingNo TEXT,
    totalValue INTEGER,
    fee INTEGER,
    senderAddr BLOB,
    outs BLOB,
    timestamp INTEGER
);

CREATE TABLE
IF NOT EXISTS accounts
(
    addr BLOB NOT NULL PRIMARY KEY, 
    balance INTEGER
);
