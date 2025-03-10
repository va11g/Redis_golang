package memdb

import (
	"fmt"
	"resp/logger"
)

type Node[T Val] struct {
	Value  T
	left   *Node[T]
	right  *Node[T]
	height int64
}

type Val interface {
	Cmp(val Val) int64
	SetScore(score float64)
	GetScore() float64
	GetNames() map[string]struct{}
	AddName(name string)
	DeleteName(name string)
	Empty()
	IsNameExist(name string) bool
}

type ATree[T Val] struct {
	root   *Node[T]
	values []map[string]float64
	len    int
	dict   map[string]*Node[T]
}

// ----------------------------------------------------------------------------------------------------
//
//       _____                    _____                    _____
//      /\    \                  /\    \                  /\    \
//     /::\____\                /::\    \                /::\____\
//    /:::/    /               /::::\    \              /:::/    /
//   /:::/    /               /::::::\    \            /:::/    /
//  /:::/    /               /:::/\:::\    \          /:::/    /
// /:::/____/               /:::/__\:::\    \        /:::/    /
// |::|    |               /::::\   \:::\    \      /:::/    /
// |::|    |     _____    /::::::\   \:::\    \    /:::/    /
// |::|    |    /\    \  /:::/\:::\   \:::\    \  /:::/    /
// |::|    |   /::\____\/:::/  \:::\   \:::\____\/:::/____/
// |::|    |  /:::/    /\::/    \:::\  /:::/    /\:::\    \
// |::|    | /:::/    /  \/____/ \:::\/:::/    /  \:::\    \
// |::|____|/:::/    /            \::::::/    /    \:::\    \
// |:::::::::::/    /              \::::/    /      \:::\    \
// \::::::::::/____/               /:::/    /        \:::\    \
//  ~~~~~~~~~~                    /:::/    /          \:::\    \
//                               /:::/    /            \:::\    \
//                              /:::/    /              \:::\____\
//                              \::/    /                \::/    /
//                               \/____/                  \/____/
//
// ----------------------------------------------------------------------------------------------------

type SortedSet[T Val] struct {
	*ATree[T]
}

func NewSortedSet() *SortedSet[*SortedSetNode] {
	return &SortedSet[*SortedSetNode]{NewATree[*SortedSetNode]()}
}

type SortedSetNode struct {
	Names map[string]struct{} // allow multiple member with the same value
	Score float64
}

func (n *SortedSetNode) Cmp(val Val) int64 {
	n2 := val.(*SortedSetNode)
	if n.Score > n2.Score {
		return 1
	} else if n.Score < n2.Score {
		return -1
	} else if n.Score == n2.Score {
		return 0
	} else {
		logger.Error("cant compare values in sorted set Cmp")
		return 0
	}
}

func (n *SortedSetNode) SetScore(score float64) {
	n.Score = score
}

func (n *SortedSetNode) GetScore() float64 {
	return n.Score
}

func (n *SortedSetNode) GetNames() map[string]struct{} {
	return n.Names
}

func (n *SortedSetNode) AddName(name string) {
	n.Names[name] = struct{}{}
}

func (n *SortedSetNode) DeleteName(name string) {
	delete(n.Names, name)
}

func (n *SortedSetNode) Empty() {
	n.Names = nil
}

func (n *SortedSetNode) IsNameExist(name string) bool {
	_, ok := n.Names[name]
	return ok
}

// ----------------------------------------------------------------------------------------------------
//
//   _______   _____    ______   ______
//  |__   __| |  __ \  |  ____| |  ____|
//     | |    | |__) | | |__    | |__
//     | |    |  _  /  |  __|   |  __|
//     | |    | | \ \  | |____  | |____
//     |_|    |_|  \_\ |______| |______|
//
// ----------------------------------------------------------------------------------------------------

func NewATree[T Val]() *ATree[T] {
	return new(ATree[T]).Init()
}

func (t *ATree[T]) Init() *ATree[T] {
	t.root = nil
	t.values = nil
	t.len = 0
	t.dict = make(map[string]*Node[T])
	return t
}

func NewBtree[T Val]() *ATree[T] {
	return new(ATree[T]).Init()
}

func (t *ATree[T]) String() string {
	return fmt.Sprint()
}

func (t *ATree[T]) IsEmpty() bool {
	return t.root == nil
}

// ----------------------------------------------------------------------------------------------------
//
//   _   _    ____    _____    ______
//  | \ | |  / __ \  |  __ \  |  ____|
//  |  \| | | |  | | | |  | | | |__
//  | . ` | | |  | | | |  | | |  __|
//  | |\  | | |__| | | |__| | | |____
//  |_| \_|  \____/  |_____/  |______|
//
// ----------------------------------------------------------------------------------------------------

func (n *Node[T]) Init() *Node[T] {
	n.height = 1
	n.left = nil
	n.right = nil
	return n
}

func (n *Node[T]) String() string {
	return fmt.Sprint(n.Value)
}

func height[T Val](n *Node[T]) int64 {
	if n != nil {
		return n.height
	}
	return 0
}

func balance[T Val](n *Node[T]) int64 {
	if n == nil {
		return 0
	}
	return height(n.left) - height(n.right)
}

func (n *Node[T]) Debug() {
	var info string
	if n.left == nil && n.right == nil {
		info = "no children |"
	} else if n.left != nil && n.right != nil {
		info = fmt.Sprint("left child:", n.left.String(), " right child:", n.right.String())
	} else if n.right != nil {
		info = fmt.Sprint("right child:", n.right.String())
	} else {
		info = fmt.Sprint("left child:", n.left.String())
	}
	fmt.Println(n.String(), "|", "height", n.height, "|", "balance", balance(n), "|", info)
}

func (n *Node[T]) get(target T) *Node[T] {
	node := new(Node[T])
	c := target.Cmp(n.Value)
}
