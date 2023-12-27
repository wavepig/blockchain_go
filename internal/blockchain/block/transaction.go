package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

// 挖出新块的奖励金 目前设置常量
const subsidy = 10

// 在比特币中，最先有输出，然后才有输入。换而言之，第一笔交易只有输出，没有输入。
// Transaction 由交易 ID，输入和输出构成
type Transaction struct {
	ID   []byte
	Vin  []TXInput  // 输入
	Vout []TXOutput // 输出
}

// IsCoinbase 检查交易是否为Coinbase
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// SetID 设置交易的ID
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	// 交易ID Sum256(Transaction)
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// TXInput 交易输入结构体信息 包含 3 部分
// Txid: 一个交易输入引用了之前一笔交易的一个输出, ID 表明是之前哪笔交易
// Vout: 一笔交易可能有多个输出，Vout 为输出的索引
// ScriptSig: 提供解锁输出 Txid:Vout 的数据
type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

// TXOutput 交易输出结构体信息 包含两部分
// Value: 有多少币，就是存储在 Value 里面
// ScriptPubKey: 对输出进行锁定
// 在当前实现中，ScriptPubKey 将仅用一个字符串来代替
type TXOutput struct {
	Value        int
	ScriptPubKey string
}

// CanUnlockOutputWith 检查是否可以使用提供的数据解锁
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// CanBeUnlockedWith 检查是否可以使用提供的数据解锁
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

// NewCoinbaseTX 创建一个新的Coinbase交易 发行新币
// 在区块链的最初，也就是第一个块，叫做创世块。正是这个创世块，产生了区块链最开始的输出。
// 对于创世块，不需要引用之前的交易输出。因为在创世块之前根本不存在交易，也就没有不存在交易输出。
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	// 在比特币中，最先有输出，然后才有输入。换而言之，第一笔交易只有输出，没有输入。
	txin := TXInput{[]byte{}, -1, data} // 第一笔 没有输入
	txout := TXOutput{subsidy, to}      // 输出 10 第一个区块奖励
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()

	return &tx
}

// NewUTXOTransaction 创建新的交易
// from 发送者  to 接受者  amount 大小
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput   // 输入
	var outputs []TXOutput // 输出
	// 找到足够的未花费输出
	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// 构建一个输入列表
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// 构建一个输出列表
	outputs = append(outputs, TXOutput{amount, to})
	// 如果 UTXO 总数超过所需，则产生找零
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}
