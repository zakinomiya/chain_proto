SELECT
    txHash,
    txType,
    totalValue,
    fee,
    senderAddr,
    outCount,
    outs,
    signature,
    timestamp
FROM transactions
WHERE transactions.txHash=:hash;
