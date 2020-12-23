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
WHERE blocks.height >= :start AND blocks.height<:
end;