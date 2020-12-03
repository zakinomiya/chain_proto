SELECT
    height,
    hash,
    prevBlockHash,
    merkleRoot,
    extraNonce,
    timestamp,
    bits,
    nonce,
    transactions
FROM blocks
ORDER BY height DESC
LIMIT 1