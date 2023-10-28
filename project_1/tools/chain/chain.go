package chain

import (
	"blockchainSimulate/tools/block"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

type Node struct {
	SelfBlock *block.Block
	Index     int
	Nextnode  []*Node
	MinerID   int
}

type Chain struct {
	Head         *Node
	Taillist     []*Node
	Unmergedlist []*Node
}

// 创建创世块
func NewChain() *Chain {
	genesisBlock := &block.Block{
		Index:     0,
		PrevHash:  [32]byte{},
		Timestamp: time.Now().Unix(),
		Bits:      block.Difficulty,
		Nonce:     0,
		Tx_root:   sha256.Sum256([]byte("This is the Genesis Block")),
	}
	nchain := &Chain{
		Head: &Node{
			SelfBlock: genesisBlock,
			Index:     0,
			Nextnode:  nil,
			MinerID:   -1,
		},
		Taillist: []*Node{},
	}
	nchain.Taillist = append(nchain.Taillist, nchain.Head)
	return nchain
}

// 提交新区块
func (c *Chain) SubmitBlock(b *block.Block, tail *Node, minerid int) *Node {

	tailhash := block.CalculateHash(*tail.SelfBlock)
	if tailhash == b.PrevHash && tail.Index+1 == b.Index && tail.SelfBlock.Bits == b.Bits && uint64(block.CalculateHash(*b)[0]) < b.Bits {
		newnode := &Node{
			SelfBlock: b,
			Index:     tail.Index + 1,
			Nextnode:  nil,
			MinerID:   minerid,
		}
		tail.Nextnode = append(tail.Nextnode, newnode)
		c.Unmergedlist = append(c.Unmergedlist, newnode)
		return newnode
	}
	return nil
}

// 获取尾节点
func (c *Chain) GetTail() *Node {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	randomIndex := rng.Intn(len(c.Taillist))
	return c.Taillist[randomIndex]
}

func (n *Node) Height() int {
	if len(n.Nextnode) == 0 {
		return 1
	}

	maxChildHeight := 0

	for _, child := range n.Nextnode {
		childHeight := child.Height()
		if childHeight > maxChildHeight {
			maxChildHeight = childHeight
		}
	}

	return maxChildHeight + 1
}

// 合并最长链PoW
func Mergenode(root *Node) *Node {
	if len(root.Nextnode) == 0 {
		return root
	}

	maxHeight := -1
	highestnode := []*Node{}
	for _, child := range root.Nextnode {
		subtree := Mergenode(child)
		nowheight := subtree.Height()
		if nowheight > maxHeight {
			maxHeight = nowheight
			highestnode = []*Node{subtree}
		} else if nowheight == maxHeight {
			highestnode = append(highestnode, subtree)
		}
	}

	root.Nextnode = highestnode
	return root
}

func CalTailList(root *Node) []*Node {
	taillist := []*Node{}

	var dfs func(*Node)
	dfs = func(node *Node) {
		if len(node.Nextnode) == 0 {
			taillist = append(taillist, node)
			return
		}

		for _, child := range node.Nextnode {
			dfs(child)
		}
	}

	dfs(root)
	return taillist
}

func (c *Chain) Merge() {
	c.Head = Mergenode(c.Head)
	c.Taillist = CalTailList(c.Head)
	c.Unmergedlist = []*Node{}
}

func (n *Node) printNode() {
	if n == nil {
		return
	}

	fmt.Printf("Node ID: %d with miner: %d\n", n.Index, n.MinerID)
}

func (c *Chain) PrintChain() {
	if c == nil {
		return
	}
	fmt.Println("Now submitted chain is:")
	stack := []*Node{c.Head}
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		node.printNode()
		for i := len(node.Nextnode) - 1; i >= 0; i-- {
			stack = append(stack, node.Nextnode[i])
		}
	}
	fmt.Println("chain has taillist:")
	for _, node := range c.Taillist {
		node.printNode()
	}
	fmt.Println("chain has unmergedlist:")
	for _, node := range c.Unmergedlist {
		node.printNode()
	}
}

func (c *Chain) CalBranchLength() []int {
	branchNode := c.Head
	for {
		if branchNode == nil {
			return []int{0}
		} else if len(branchNode.Nextnode) < 2 {
			branchNode = branchNode.Nextnode[0]
		} else {
			break
		}
	}
	n := len(branchNode.Nextnode)
	lenlist := make([]int, n)
	for i := 0; i < n; i++ {
		lenlist[i] = branchNode.Nextnode[i].Height()
	}
	return lenlist
}

// 根据链上区块计算矿工收益
func (c *Chain) CalMinerReward() map[int]int {
	reward := make(map[int]int)
	ptr := c.Head
	for ptr != nil {
		if _, exists := reward[ptr.MinerID]; !exists {
			reward[ptr.MinerID] = 1
		} else {
			reward[ptr.MinerID]++
		}

		if len(ptr.Nextnode) != 1 {
			break
		} else {
			ptr = ptr.Nextnode[0]
		}
	}
	return reward
}
