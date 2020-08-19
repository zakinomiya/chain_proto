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
