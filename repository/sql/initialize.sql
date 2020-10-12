CREATE TABLE  IF NOT EXISTS block
(
    height INTEGER NOT NULL PRIMARY KEY,
    hash BLOB,
    prevBlockHash BLOB,
    merkleRoot BLOB,
    transactions BLOB,
    extraNone INTEGER,
    timestamp INTEGER,
    bits INTEGER,
    nonce INTEGER
);

CREATE TABLE IF NOT EXISTS transactions
(
    txHash BLOB NOT NULL PRIMARY KEY,
    totalValue INTEGER,
    fee INTEGER,
    senderAddr BLOB,
    outs BLOB,
    timestamp INTEGER
);

CREATE TABLE IF NOT EXISTS accounts (
    addr BLOB NOT NULL PRIMARY KEY, 
    balance INTEGER
);
