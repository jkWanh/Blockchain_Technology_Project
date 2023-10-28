package miner

import (
	"blockchainSimulate/tools/block"
	"blockchainSimulate/tools/chain"
	"fmt"
	"math/rand"
)

var Q int = 100

type Miner struct {
	MinerID int
}

func NewMiner(id int) Miner {
	return Miner{
		MinerID: id,
	}
}

func (m Miner) MineBlock(chain *chain.Chain) bool {
	target := chain.GetTail()
	ran_message := rand.Intn(1000000)
	message := fmt.Sprintf("This is block %d from miner %d, rand_num is %d", target.Index+1, m.MinerID, ran_message)
	newBlock := block.MineBlock_inOneRound(target.SelfBlock, message, Q)
	if newBlock != nil {
		newnode := chain.SubmitBlock(newBlock, target, m.MinerID)
		if newnode != nil {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

type Vir_Miner struct {
	Miner
	SelfBlocklist []*block.Block
	InitTarget    *chain.Node
}

func NewVir_Miner(id int) Vir_Miner {
	return Vir_Miner{
		Miner: Miner{
			MinerID: id + 1000,
		},
		SelfBlocklist: []*block.Block{},
		InitTarget:    nil,
	}
}

func (vm *Vir_Miner) MineBlock_SelfishAttack(chain *chain.Chain) bool {
	if len(vm.SelfBlocklist) == 0 {
		vm.InitTarget = chain.GetTail()
		vm.SelfBlocklist = append(vm.SelfBlocklist, vm.InitTarget.SelfBlock)
	}
	targetBlock := vm.SelfBlocklist[len(vm.SelfBlocklist)-1]
	ran_message := rand.Intn(1000000)
	message := fmt.Sprintf(("This is block %d from vitriolic miner %d, rand_num is %d"), targetBlock.Index+1, vm.MinerID, ran_message)
	newBlock := block.MineBlock_inOneRound(targetBlock, message, vm.MinerID)
	if newBlock != nil {
		vm.SelfBlocklist = append(vm.SelfBlocklist, newBlock)
		fmt.Print(vm.SelfBlocklist)
		return true
	} else {
		return false
	}
}

func (vm *Vir_Miner) SubmitBlock(chain *chain.Chain) {
	targetNode := vm.InitTarget
	if len(vm.SelfBlocklist) > 0 {
		for i := 1; i < len(vm.SelfBlocklist); i++ {
			targetNode = chain.SubmitBlock(vm.SelfBlocklist[i], targetNode, vm.MinerID)
			fmt.Printf("Miner %d submit block %d\n", vm.MinerID, i)
		}
	}
	vm.InitTarget = nil
	vm.SelfBlocklist = []*block.Block{}
}
