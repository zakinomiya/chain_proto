SELECT
    txHash,
    blockHash,
    pendingNo,
    totalValue,
    fee,
    senderAddr,
    outCount,
    outs,
    timestamp
FROM transactions
WHERE transactions.blockHash=:blockHash;