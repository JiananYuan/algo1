package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	PASS = 0
	BET = 1
	PLAYER_A = 3
	PLAYER_B = 4
	N = 1000
	INIT_PROFIT = 3000
	TRY_THRESHOLD = 0.1
)

var last_player int
var last_op_A int
var last_op_B int
var win_A int  // A赢的次数
var random_choice_A int // 自定义策略的硬币决定变量
var random_choice_B int // 自定义策略的硬币决定变量
var profit_A int
var profit_B int
var a_mp map[string] int
var b_mp map[string] int
var a_need_update bool
var b_need_update bool
var current_games int  // 当前游戏进行了几场


type node struct {
	left, right *node
	op_str string
}

// 发牌函数: 从三张牌取两张
func distribute() (int, int) {
	a := 0
	b := 0
	a = rand.Int() % 3 + 1
	for true {
		b = rand.Int() % 3 + 1
		if a != b {
			break
		}
	}
	return a, b
}

// 随机操作函数: 从过牌和加注里面随机选择一种操作
// 同时也作为二选一函数接口
func random_operate() int {
	op := rand.Int() % 2
	if op == PASS {
		return PASS
	} else {
		return BET
	}
}

// 固定操作函数
func fix_operate(player int, a int, b int) int {
	if player == PLAYER_A {
		if a == 1 {
			return PASS
		} else if a == 2 {
			if last_op_A == BET {
				last_op_A = PASS
				return PASS
			} else {
				last_op_A = BET
				return BET
			}
		} else {
			return BET
		}
	} else {
		if b == 1 {
			return PASS
		} else if b == 2 {
			if last_op_B == BET {
				last_op_B = PASS
				return PASS
			} else {
				last_op_B = BET
				return BET
			}
		} else {
			return BET
		}
	}
}

// note: 需要事先生成好random_choice随机变量
func my_scheme(player int, a int, b int, op_str string) int {
	if player == PLAYER_A {
		if a == 1 {
			return PASS
		} else if a == 2 {
			// 用不多于20%的钱进行试探 (随机策略)
			if current_games < N * TRY_THRESHOLD {
				if random_operate() == BET {
					return BET
				} else {
					return PASS
				}
			} else {
				a_need_update = false
				if a_mp[op_str + "P"] >= a_mp[op_str + "B"] {
					return PASS
				} else {
					return BET
				}
			}
		} else {
			return BET
		}
	} else {
		if b == 1 {
			return PASS
		} else if b == 2 {
			// 用不多于20%的钱进行试探
			if current_games < N * TRY_THRESHOLD {
				if random_operate() == BET {
					return BET
				} else {
					return PASS
				}
			} else {
				b_need_update = false
				if b_mp[op_str + "P"] >= b_mp[op_str + "B"] {
					return PASS
				} else {
					return BET
				}
			}
		} else {
			return BET
		}
	}
}

// 人工操作
func artificial_operate() int {
	op := 0
	fmt.Scan(&op)
	return op
}

/**
  博弈函数 (note: A先手)
  @param player: 当前玩家 0代表A玩家，1代表B玩家
  @param a: A拿到的牌
  @param b: B拿到的牌
  @param profit: 当前押注
  @return
*/
func vs(treenode *node, a int, b int) int {
	score := 0
	if treenode.op_str == "PP" {
		if a > b {
			score = 1
		} else {
			score = -1
		}
		if a == 2 && a_need_update {
			a_mp["P"] += score
			a_mp["PP"] += score
		}
		if b == 2 && b_need_update {
			b_mp["P"] -= score
			b_mp["PP"] -= score
		}
		return score
	} else if treenode.op_str == "BB" {
		if a > b {
			score = 2
		} else {
			score = -2
		}
		if a == 2 && a_need_update {
			a_mp["B"] += score
			a_mp["BB"] += score
		}
		if b == 2 && b_need_update {
			b_mp["B"] -= score
			b_mp["BB"] -= score
		}
		return score
	} else if treenode.op_str == "BP" {
		score = 1
		if a == 2 && a_need_update {
			a_mp["B"] += score
			a_mp["BP"] += score
		}
		if b == 2 && b_need_update {
			b_mp["B"] -= score
			b_mp["BP"] -= score
		}
		return score
	} else if treenode.op_str == "PBP" {
		score = -1
		if a == 2 && a_need_update {
			a_mp["P"] += score
			a_mp["PB"] += score
			a_mp["PBP"] += score
		}
		if b == 2 && b_need_update {
			b_mp["P"] -= score
			b_mp["PB"] -= score
			b_mp["PBP"] -= score
		}
		return score
	} else if treenode.op_str == "PBB" {
		if a > b {
			score = 2
		} else {
			score = -2
		}
		if a == 2 && a_need_update {
			a_mp["P"] += score
			a_mp["PB"] += score
			a_mp["PBB"] += score
		}
		if b == 2 && b_need_update {
			b_mp["P"] -= score
			b_mp["PB"] -= score
			b_mp["PBB"] -= score
		}
		return score
	}

	var op int
	if last_player == PLAYER_B {
		op = my_scheme(PLAYER_A, a, b, treenode.op_str)
		last_player = PLAYER_A
	} else {
		op = my_scheme(PLAYER_B, a, b, treenode.op_str)
		last_player = PLAYER_B
	}

	if op == PASS {
		treenode = treenode.left
	} else if op == BET {
		treenode = treenode.right
	}
	treenode.left = &node{op_str: treenode.op_str + "P"}
	treenode.right = &node{op_str: treenode.op_str + "B"}
	return vs(treenode, a, b)
}

// 游戏逻辑
func game() {
	// 控制先手
	last_player = PLAYER_B
	a, b := distribute()
	last_op_A = BET
	last_op_B = BET
	root := &node{op_str: ""}
	root.left = &node{op_str: "P"}
	root.right = &node{op_str: "B"}
	profit := 0
	if last_player == PLAYER_B {
		profit = vs(root, a, b)
	} else {
		profit = vs(root, b, a)
	}
	if profit > 0 {
		win_A += 1
	}
	profit_A += profit
	profit_B -= profit
	fmt.Println(profit_A, profit_B)
}

func statistic() {
	fmt.Println("先手胜率统计: ", win_A)
	fmt.Println("后手胜率统计: ", N - win_A)
	fmt.Println("先手收益: ", profit_A)
	fmt.Println("后手收益: ", profit_B)
}

func main() {
	// Go生成随机数: https://blog.csdn.net/u011304970/article/details/72721747
	rand.Seed(time.Now().Unix())
	profit_A = INIT_PROFIT
	profit_B = INIT_PROFIT
	// 我的策略需要的初始化操作
	a_need_update = true
	b_need_update = true
	a_mp = make(map[string] int)
	b_mp = make(map[string] int)
	// 进行N 轮博弈
	for i := 0; i < N; i += 1 {
		current_games = i
		game()
	}
	statistic()
}
