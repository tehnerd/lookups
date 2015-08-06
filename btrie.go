//simple implementation of binary trie for ipv4 lookup
package trie

import "fmt"

type Prefix struct {
	prefix    uint32
	prefixLen uint8
	adj       string
}

type BTNode struct {
	lNode   *BTNode
	rNode   *BTNode
	adjency string
}

type BTrie struct {
	root *BTNode
	adj  string
}

func (bt *BTrie) AddDefault(def string) {
	bt.adj = def
}

func inOrderBT(node *BTNode) {
	if node != nil {
		fmt.Printf("%v\n", node)
		inOrderBT(node.lNode)
		inOrderBT(node.rNode)
	}
	return
}

func (bt *BTrie) PrintTrie() {
	node := bt.root
	inOrderBT(node)
}

func countNilBT(node *BTNode, nilNodes, total *int) {
	if node != nil {
		if node.adjency == "" {
			*nilNodes++
		}
		*total++
		countNilBT(node.lNode, nilNodes, total)
		countNilBT(node.rNode, nilNodes, total)
	}
	return
}

func (bt *BTrie) AddPrefix(prefix Prefix) {
	if prefix.prefixLen == 0 {
		bt.adj = prefix.adj
		return
	}
	shift := uint8(31)
	node := bt.root
	if node == nil {
		node = new(BTNode)
		bt.root = node
	}
	for ; IPV4_LEN-shift <= prefix.prefixLen; shift-- {
		if (prefix.prefix>>shift)&1 == 0 {
			if node.lNode != nil {
				node = node.lNode
			} else {
				node.lNode = new(BTNode)
				node = node.lNode
			}
		} else {
			if node.rNode != nil {
				node = node.rNode
			} else {
				node.rNode = new(BTNode)
				node = node.rNode
			}
		}
	}
	node.adjency = prefix.adj
}

func (bt *BTrie) FindAdj(host uint32) string {
	adj := bt.adj
	if bt.root != nil {
		shift := uint8(31)
		node := bt.root
		for ; IPV4_LEN-shift <= 32; shift-- {
			if node.adjency != "" {
				adj = node.adjency
			}
			if (host>>shift)&1 == 0 {
				node = node.lNode
			} else {
				node = node.rNode
			}
			if node == nil {
				break
			}
		}
		if node != nil {
			adj = node.adjency
		}

	}
	return adj
}
