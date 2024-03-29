SELECT
    height,
    hash,
    prevBlockHash,
    merkleRoot,
    transactions,
    extraNonce,
    timestamp,
    bits,
    nonce,
    transactions
FROM blocks
LIMIT :limit OFFSET :offset;
