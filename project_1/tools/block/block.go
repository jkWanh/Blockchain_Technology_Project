package block

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

// 难度系数：哈希结果前导0的个数
var Difficulty uint64 = 0x000010

// 定义区块结构
type Block struct {
	Index     int
	PrevHash  [32]byte
	Timestamp int64
	Bits      uint64
	Nonce     uint64
	Tx_root   [32]byte
}

// 计算块的哈希值
func CalculateHash(b Block) [32]byte {
	data := strconv.Itoa(b.Index) + string(b.PrevHash[:]) + strconv.FormatInt(b.Timestamp, 10) + strconv.FormatUint(b.Bits, 10) + strconv.FormatUint(b.Nonce, 10) + string(b.Tx_root[:])
	return sha256.Sum256([]byte(data))
}

func MineBlock_inOneRound(b *Block, data string, q int) *Block {
	newBlock := &Block{
		Index:     b.Index + 1,
		PrevHash:  CalculateHash(*b),
		Timestamp: time.Now().Unix(),
		Bits:      b.Bits,
		Nonce:     0,
		Tx_root:   sha256.Sum256([]byte(data)),
	}

	for i := 0; i < q; i++ {

		hash := CalculateHash(*newBlock)
		combined := uint64(hash[0])<<16 | uint64(hash[1])<<8 | uint64(hash[2])

		if combined < b.Bits {
			fmt.Printf("find new Block Success: %x\n", hash)
			return newBlock
		}
		newBlock.Nonce++

	}
	return nil
}
