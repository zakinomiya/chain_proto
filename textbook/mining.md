# Mining

btcd の実装を参考にしている。

マイニングは、runtime.NumCPU()個の Goroutine を最大で動かしている。

それぞれの Goroutine で ExtraNonce をランダムで決めている。それに対して最大まで Nonce を回している。

考えること

- ExtraNonce いるか？
- Is uint32 enough? (I guess so)
- block をどこで作って Miner に渡すか（もしくは Miner が Block を作るか）
