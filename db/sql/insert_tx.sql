-- TODO write tx insert sql
INSERT INTO transactions
    (
    txHash,
    blockHash,
    pendingNo,
    totalValue,
    fee,
    senderAddr,
    outCount,
    outs,
    timestamp
    )
VALUES
    (
        :txHash,
        :blockHash,
        :pendingNo,
        :totalValue,
        :fee,
        :senderAddr,
        :outCount,
        :outs,
        :timestamp
    );