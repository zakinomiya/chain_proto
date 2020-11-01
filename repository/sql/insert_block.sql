INSERT INTO blocks
    (
    height ,
    hash,
    prevBlockHash,
    merkleRoot,
    extraNonce,
    timestamp,
    bits,
    nonce,
    transactions
    )
VALUES
    (
        :height,
        :hash,
        :prevBlockHash,
        :merkleRoot,
        :extraNonce,
        :timestamp,
        :bits,
        :nonce,
        :transactions
);