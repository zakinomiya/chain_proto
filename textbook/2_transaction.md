# Transaction

## How to verify the ownership of an output

Output には PubKey 情報がある。

Input は Output への Reference をもっている。参照されているアウトプットの所有権を自分が持つことの証明として、

自分の秘密鍵での電子署名が、Output 内の公開鍵で Verify できるかを調べる。
