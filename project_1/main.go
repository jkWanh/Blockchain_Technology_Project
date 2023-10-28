package main

import (
	"blockchainSimulate/tools/chain"
	"blockchainSimulate/tools/miner"
	"fmt"
	"time"
)

// num: 矿工节点数量
var num int = 0

// v_num: 恶意矿工节点数量
var v_num int = 0

// dif: 每个轮次出块成功概率（单位: %, 最小分度：1%）
var dif int = 0

func SettingInit() {
	// 读取矿工节点数量
	for {
		fmt.Println("请输入矿工节点数量: ")
		_, err := fmt.Scan(&num)
		if err != nil {
			fmt.Println("输入错误")
			continue
		} else {
			fmt.Printf("输入成功, 当前矿工节点数量为%d\n", num)
			break
		}
	}

	// 读取每个轮次出块成功概率
	for {
		fmt.Println("请输入每个轮次出块成功概率（单位: %, 最小分度：1%）: ")
		_, err := fmt.Scan(&dif)
		if err != nil {
			fmt.Println("输入错误")
			continue
		} else {
			fmt.Printf("输入成功, 当前每轮次出块概率为%d %\n", dif)
			break
		}
	}
	miner.Q = int((dif * 1048576) / (100 * num))

}

func CalBlockGenerateSpeed() {
	SettingInit()
	blockchain := chain.NewChain()
	minerlist := []miner.Miner{}
	for i := 0; i < num; i++ {
		minerlist = append(minerlist, miner.NewMiner(i))
	}
	fmt.Println("开始挖矿,测量每个轮次出块速度")
	r := 1
	nb := 0
	for {
		var tag bool = false
		fmt.Printf("当前轮次：%d\n", r)
		for j := 0; j < len(minerlist); j++ {
			if minerlist[j].MineBlock(blockchain) {
				fmt.Printf("Miner %d find a new block\n", j)
				tag = true
				nb++
			}
		}
		if tag {
			fmt.Printf("\t本轮产生新区块，当前出块平均速率为%f块/轮\n", float64(nb)/float64(r))
			blockchain.Merge()
		} else {
			fmt.Printf("\t本轮未产生新区块，当前出块平均速率为%f块/轮\n", float64(nb)/float64(r))
		}
		r++
		time.Sleep(2 * time.Second)
	}
}

func SelfishAttackSimulation() {
	SettingInit()
	for {
		fmt.Println("请设置自私挖矿恶意节点比例（单位：%）")
		_, err := fmt.Scan(&v_num)
		if err != nil {
			fmt.Println("输入错误")
			continue
		} else {
			fmt.Printf("输入成功, 当前自私挖矿恶意节点比例为%d %\n", v_num)
			break
		}
	}
	v_num = int(float64(v_num) / 100 * float64(num))
	blockchain := chain.NewChain()
	minerlist := []miner.Miner{}
	for i := 0; i < num-v_num; i++ {
		minerlist = append(minerlist, miner.NewMiner(i))
	}
	v_minerlist := []miner.Vir_Miner{}
	for i := 0; i < v_num; i++ {
		v_minerlist = append(v_minerlist, miner.NewVir_Miner(i))
	}
	fmt.Println("挖掘100块，测量节点收益")
	for blockchain.GetTail().Index < 100 {
		var tag bool = false
		for j := 0; j < len(minerlist); j++ {
			if minerlist[j].MineBlock(blockchain) {
				tag = true
			}
		}
		for j := 0; j < len(v_minerlist); j++ {
			v_minerlist[j].MineBlock_SelfishAttack(blockchain)
		}
		if tag {
			for j := 0; j < len(v_minerlist); j++ {
				v_minerlist[j].SubmitBlock(blockchain)
			}
		}
		blockchain.Merge()

	}
	dic := blockchain.CalMinerReward()
	fmt.Println("矿工收益情况：")
	var v_total, n_total float32 = 0, 0
	for k, _ := range dic {
		if k >= 1000 {
			v_total += 1
		} else if k != -1 {
			n_total += 1
		}
	}
	fmt.Printf("恶意矿工收益期望：%f\n", v_total/(v_total+n_total))
}

func BranchAttackSimulation() {
	SettingInit()
	for {
		fmt.Println("请设置分叉攻击恶意节点比例（单位：%）")
		_, err := fmt.Scan(&v_num)
		if err != nil {
			fmt.Println("输入错误")
			continue
		} else {
			fmt.Printf("输入成功, 当前分叉攻击恶意节点比例为%d %\n", v_num)
			break
		}
	}
	v_num = int(float64(v_num) / 100 * float64(num))
	minerlist := []miner.Miner{}
	for i := 0; i < num-v_num; i++ {
		minerlist = append(minerlist, miner.NewMiner(i))
	}
	v_minerlist := []miner.Vir_Miner{}
	for i := 0; i < v_num; i++ {
		v_minerlist = append(v_minerlist, miner.NewVir_Miner(i))
	}
	fmt.Println("开始估计分叉攻击成功1-6块的期望, 对不同攻击长度进行10次模拟并以成功次数作为估计值")
	dic := make(map[int]float64)
	for target := 1; target <= 6; target++ {
		success := 0
		fmt.Printf("当前模拟攻击长度为%d\n", target)
		for i := 0; i < 10; i++ {
			blockchain := chain.NewChain()
			blockchain_v := chain.NewChain()
			l := 0
			for l < target {
				tag := false
				for j := 0; j < len(minerlist); j++ {
					minerlist[j].MineBlock(blockchain)
				}
				blockchain.Merge()
				for j := 0; j < len(v_minerlist); j++ {
					if v_minerlist[j].MineBlock(blockchain_v) {
						tag = true
					}
				}
				if tag {
					blockchain_v.Merge()
					l += 1
				}
			}
			if blockchain.GetTail().Index < blockchain_v.GetTail().Index {
				success += 1
			}
		}
		dic[target] = float64(success) / 10
	}
	for k, v := range dic {
		fmt.Printf("攻击长度为%d时，成功概率为%f\n", k, v)
	}
}

func main() {
	var op_num int = 0
	for {
		fmt.Println("PoW仿真程序，请输入模拟任务编号：\n1. 计算出块速度\n2. 自私挖矿攻击仿真\n3. 分叉攻击仿真\n4. 退出")
		_, err := fmt.Scan(&op_num)
		if err != nil {
			fmt.Println("输入错误")
			continue
		} else {
			switch op_num {
			case 1:
				CalBlockGenerateSpeed()
				continue
			case 2:
				SelfishAttackSimulation()
				continue
			case 3:
				BranchAttackSimulation()
				continue
			case 4:
				fmt.Println("退出成功")
				return
			}
		}
	}
}
