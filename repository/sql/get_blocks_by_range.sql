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
WHERE height >= :start AND height<:
end;