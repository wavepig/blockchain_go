package block

import (
	"blockchain_go/pkg/utils"
	"bytes"
)

// TXOutput 交易输出结构体信息 包含两部分
// Value: 有多少币，就是存储在 Value 里面
// ScriptPubKey: 对输出进行锁定
// 在当前实现中，ScriptPubKey 将仅用一个字符串来代替
type TXOutput struct {
	Value      int
	PubKeyHash []byte
}

func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := utils.Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}

// CanUnlockOutputWith 检查是否可以使用提供的数据解锁
//func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
//	return in.ScriptSig == unlockingData
//}

// CanBeUnlockedWith 检查是否可以使用提供的数据解锁
//func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
//	return out.ScriptPubKey == unlockingData
//}
