//simple implementation of binary trie for ipv4 lookup
package trie

import "fmt"

type PCTNode struct {
	lNode   *PCTNode
	rNode   *PCTNode
	skip    uint8
	segment uint32
	adjency string
}

type PCTrie struct {
	root *PCTNode
	adj  string
}

func (pc *PCTrie) AddDefault(def string) {
	pc.adj = def
}

func inOrderPCT(node *PCTNode) {
	if node != nil {
		fmt.Printf("%v\n", node)
		inOrderPCT(node.lNode)
		inOrderPCT(node.rNode)
	}
	return
}

func (pc *PCTrie) PrintTrie() {
	node := pc.root
	inOrderPCT(node)
}

func countNilPCT(node *PCTNode, nilNodes, total *int) {
	if node != nil {
		if node.adjency == "" {
			*nilNodes++
		}
		*total++
		countNilPCT(node.lNode, nilNodes, total)
		countNilPCT(node.rNode, nilNodes, total)
	}
	return
}

func (pc *PCTrie) AddPrefix(prefix Prefix) {
	if prefix.prefixLen == 0 {
		pc.adj = prefix.adj
		return
	}
	node := pc.root
	if node == nil {
		node = new(PCTNode)
		pc.root = node
	}
	shift := uint8(31)
	var prevNode *PCTNode
	_ = prevNode
	var leftNode bool
	_ = leftNode
mainLoop:
	for ; IPV4_LEN-shift <= prefix.prefixLen; shift-- {
		prevNode = node
		shiftHead := IPV4_LEN - shift
		if (prefix.prefix>>shift)&1 == 0 {
			leftNode = true
		} else {
			leftNode = false
		}
		//current node has child with same first(shift+1) bit as prefix
		if (leftNode && node.lNode != nil) ||
			(!leftNode && node.rNode != nil) {
			if leftNode {
				node = node.lNode
			} else {
				node = node.rNode
			}
			if prefix.prefixLen > shiftHead+node.skip {
				if ((prefix.prefix << shiftHead) >> (IPV4_LEN - node.skip)) == node.segment {
					shift -= node.skip
					continue
				}
				for i := node.skip - 1; i >= 0; i-- {
					if (prefix.prefix<<shiftHead)>>(IPV4_LEN-i) == (node.segment >> (node.skip - i)) {
						//trying to find first bit where current child and new prefix are not equal
						node.skip = node.skip - i - 1
						segment := node.segment
						node.segment = (node.segment << (IPV4_LEN - node.skip)) >> (IPV4_LEN - node.skip)
						newNode := new(PCTNode)
						newNode.skip = i
						newNode.segment = ((prefix.prefix << shiftHead) >> (IPV4_LEN - newNode.skip))
						if leftNode {
							prevNode.lNode = newNode
						} else {
							prevNode.rNode = newNode
						}
						if (segment>>(node.skip))&1 == 0 {
							newNode.lNode = node
						} else {
							newNode.rNode = node
						}
						node = newNode
						shift -= node.skip
						continue mainLoop
					}
				}
			} else {
				/*
					trying to find first bit where current child and new prefix are not equal.
					also we could has situation where they are equal, but prefix len of new prefix is less than current's childs
				*/
				for i := uint8(1); i <= prefix.prefixLen-shiftHead; i++ {
					if (prefix.prefix<<shiftHead)>>(IPV4_LEN-i) != (node.segment >> (node.skip - i)) {
						//found diff bit
						node.skip = node.skip - i
						node.segment = (node.segment << (IPV4_LEN - node.skip)) >> (IPV4_LEN - node.skip)
						newNode := new(PCTNode)
						newNode.skip = i - 1
						newNode.segment = ((prefix.prefix << shiftHead) >> (IPV4_LEN - newNode.skip))
						if leftNode {
							prevNode.lNode = newNode
						} else {
							prevNode.rNode = newNode
						}
						if ((prefix.prefix<<shiftHead)>>(IPV4_LEN-newNode.skip-1))&1 == 0 {
							//original node has diff val in last bit compare to newNode thats why 0  means rNode
							newNode.rNode = node
						} else {
							newNode.lNode = node
						}
						node = newNode
						shift -= node.skip
						continue mainLoop
					}
				}
				newNode := new(PCTNode)
				newNode.skip = prefix.prefixLen - shiftHead
				node.skip -= (newNode.skip + 1)
				segment := node.segment
				node.segment = (node.segment << (IPV4_LEN - node.skip)) >> (IPV4_LEN - node.skip)
				newNode.segment = ((prefix.prefix << shiftHead) >> (IPV4_LEN - newNode.skip))
				if leftNode {
					prevNode.lNode = newNode
				} else {
					prevNode.rNode = newNode
				}
				if (segment>>(node.skip))&1 == 0 {
					newNode.lNode = node
				} else {
					newNode.rNode = node
				}
				node = newNode
				break
			}
		} else {
			if leftNode {
				node.lNode = new(PCTNode)
				node = node.lNode
			} else {
				node.rNode = new(PCTNode)
				node = node.rNode
			}
			node.skip = prefix.prefixLen - shiftHead
			node.segment = ((prefix.prefix << shiftHead) >> (shiftHead + IPV4_LEN - prefix.prefixLen))
			break
		}

	}
	node.adjency = prefix.adj
}

func (pc *PCTrie) FindAdj(host uint32) string {
	adj := pc.adj
	if pc.root != nil {
		shift := uint8(31)
		node := pc.root
		for ; IPV4_LEN-shift <= 32; shift-- {
			shiftHead := IPV4_LEN - shift
			if (host>>shift)&1 == 0 {
				node = node.lNode
			} else {
				node = node.rNode
			}
			if node == nil {
				break
			}
			if node.skip != 0 {
				if (host<<shiftHead)>>(IPV4_LEN-node.skip) != node.segment {
					break
				}
				shift -= node.skip
			}
			if node.adjency != "" {
				adj = node.adjency
			}
		}

	}
	return adj
}
