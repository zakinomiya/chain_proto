# Go Blockchain

## Block

まずは Block の Struct を定義

## Blockchain

次に Blockchain の Struct を定義

Block の配列をプロパティにもつ

Blockchain を New するときにブロックを渡す（続きの場合）

もしブロックが渡されない場合は Genesis ブロックを作成

Blockchain は Runtime 中に一つしか作れないようにする

## Transaction

Transaction の Struct を作成

vin, vout は index と hash のみ持つようにする手数料などはデータとしては持たず、都度計算するようにする。その方が最初はわかりやすい。

vin, vout はもう少し考える。

BlockModel のなかに Transaction を入れるかどうか。DB に Transaction 情報をしまうかどうか

しまう場合：　
　 Transaction を ByteSlice にしてしまう
　　 → 　複数しまうのがめんどくさい。区切り文字？もしくは Transaction の Byte サイズを固定するか
　 Transcsaction の Hash だけしまう
　　 → 　これが一番現実的
　　なので blockModel は Transaction の Hash だけもつけど Block はどうするか
　　 → 　 1. Block has all the transaction data, meaning each block struct has all the transaction struct in it on memory
　　 → 　一番良いのは、Block は Transaction の Struct をもつ、BlockModel は TxHash のみ（string）
　　 → 　 Block → 　 BlockModel の変換のときに Transaction の Hash 以外の情報は落とす
　　 → 　 BlockModel 　 → 　 Block のときは、TransactionHash を元に全部 Tranaction を引っ張ってくる
　　　　　　 → 　最初は Transaction の BlockHash で取ってきて、もし見つからない場合（MerkleRoot が合わない場合）

しまわない場合：
　 Transaction 側から BLockHeight を元にひっぱってくる
　　 → 　これだともし Transaction 側で欠けているデータがあるとどうしようもない
　　 → 　しまわないはない
