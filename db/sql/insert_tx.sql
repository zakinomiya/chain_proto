-- TODO write tx insert sql
INSERT INTO transactions
    (
    txHash,
    blockHash,
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
        :pendingNo,
        :totalValue,
        :fee,
        :senderAddr,
        :outCount,
        :outs,
        :timestamp
    );
