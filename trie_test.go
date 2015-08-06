package trie

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestAddPrefix(t *testing.T) {
	var bt BTrie
	bt.AddPrefix(StringToPrefix("192.168.1.0/24", "test1"))
	bt.AddPrefix(StringToPrefix("192.168.1.1/32", "test2"))
	//bt.PrintTrie()
}

func TestAddPrefixPCT(t *testing.T) {
	var pct PCTrie
	pct.AddPrefix(StringToPrefix("192.168.1.0/24", "test1"))
	pct.PrintTrie()
	fmt.Println("###")
	pct.AddPrefix(StringToPrefix("192.168.1.1/32", "test2"))
	pct.PrintTrie()
	fmt.Println("###")
	pct.AddPrefix(StringToPrefix("192.168.1.1/16", "test3"))
	pct.PrintTrie()
	fmt.Println("###")
	pct.AddPrefix(StringToPrefix("192.168.1.1/15", "test4"))
	pct.PrintTrie()
	fmt.Println("###")
	pct.AddPrefix(StringToPrefix("192.168.1.1/14", "test5"))
	pct.PrintTrie()
	fmt.Println("###")
}

func TestFindAdj(t *testing.T) {
	prefix := Prefix{prefix: 123, prefixLen: 24, adj: "test prefix"}
	var bt BTrie
	bt.AddPrefix(prefix)
	bt.AddPrefix(Prefix{prefix: 123, prefixLen: 32, adj: "/32 test prefix"})
	bt.AddPrefix(Prefix{prefix: 0, prefixLen: 1, adj: "/1 test prefix"})
	adj := bt.FindAdj(uint32(123))
	if adj != "/32 test prefix" {
		fmt.Println("error in adjency: got ", adj, " expected is /32 test prefix")
		t.Errorf("find adjency for existing route doesnt work")
	}
	fmt.Println(adj)
	adj = bt.FindAdj(uint32(2147483647))
	if adj != "/1 test prefix" {
		fmt.Println("error in adjecncy: got ", adj, " expected is /1 test prefix")
		t.Errorf("find adjency for existing route doesnt work")
	}
	fmt.Println(adj)

}

func TestFindAdjPCT(t *testing.T) {
	var pct PCTrie
	pct.AddPrefix(StringToPrefix("192.168.1.0/24", "test1"))
	pct.AddPrefix(StringToPrefix("192.168.1.1/32", "test2"))
	pct.AddPrefix(StringToPrefix("192.168.1.1/16", "test3"))
	pct.AddPrefix(StringToPrefix("192.168.1.1/15", "test4"))
	pct.AddPrefix(StringToPrefix("192.168.1.1/14", "test5"))
	pct.AddPrefix(StringToPrefix("8.8.8.0/24", "test8"))
	pct.AddPrefix(StringToPrefix("8.8.32.0/24", "test8.8.32.0/24 "))
	pct.AddPrefix(StringToPrefix("8.8.33.0/24", "test8.8.33.0/24 "))
	pct.AddPrefix(StringToPrefix("8.8.39.0/24", "test8.8.39.0/24 "))
	pct.AddPrefix(StringToPrefix("8.8.8.0/24", "test8.8.8.0/24 "))
	pct.AddPrefix(StringToPrefix("8.8.4.0/24", "test8.8.4.0/24 "))
	pct.AddPrefix(StringToPrefix("8.8.65.0/24", "test8.8.65.0/24 "))
	pct.AddPrefix(StringToPrefix("8.8.178.0/24", "test8.8.178.0/24 "))
	pct.AddPrefix(StringToPrefix("8.8.128.0/21", "test8.8.128.0/21 "))
	pct.PrintTrie()
	ipv4 := IPv4ToUint32NoError("192.168.1.1")
	adj := pct.FindAdj(ipv4)
	fmt.Println("adj is ", adj)
	if adj != "test2" {
		t.Errorf("findadj not working. adj suppose to be \"test2\"\n")
	}
}

func TestFullViewPCTf8888(t *testing.T) {
	var pc PCTrie
	fd, err := os.Open("./fv.out")
	if err != nil {
		t.Errorf("cant find file with full view. filename must be fv.out in local dir")
		return
	}
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 2 {
			pc.AddPrefix(StringToPrefix(fields[0], fields[1]))
		}
	}
	fd.Close()
	ipv4 := IPv4ToUint32NoError("8.8.8.8")
	adj := pc.FindAdj(ipv4)
	if adj != "google_dns" {
		t.Errorf("adjency must be google_dns instead of %s\n", adj)
	}
}

func TestConvertFuncsAdnFind(t *testing.T) {
	var bt BTrie
	bt.AddPrefix(StringToPrefix("192.168.1.0/24", "1918"))
	bt.AddPrefix(StringToPrefix("10.10.0.0/16", "1918again"))
	adj := bt.FindAdj(IPv4ToUint32NoError("192.168.1.244"))
	if adj != "1918" {
		fmt.Println("wrong adj is ", adj, " expected is 1918")
		t.Errorf("error in find adj and/or string to prefix")
	}
	adj = bt.FindAdj(IPv4ToUint32NoError("10.11.1.1"))
	if adj != "" {
		fmt.Println("wrong adj is ", adj, " expected is <empty>")
		t.Errorf("error in find adj and/or string to prefix")

	}
	bt.AddDefault("default")
	adj = bt.FindAdj(IPv4ToUint32NoError("10.11.1.1"))
	if adj != "default" {
		fmt.Println("wrong adj is ", adj, " expected is default")
		t.Errorf("error in find adj and/or string to prefix")
	}

}

func TestFullViewNillNodesBT(t *testing.T) {
	var bt BTrie
	var ms runtime.MemStats
	fd, err := os.Open("./fv.out")
	if err != nil {
		t.Errorf("cant find file with full view. filename must be fv.out in local dir")
		return
	}
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 2 {
			bt.AddPrefix(StringToPrefix(fields[0], fields[1]))
		}
	}
	fd.Close()
	runtime.GC()
	runtime.ReadMemStats(&ms)
	fmt.Println("bt heap alloc: ", ms.HeapAlloc)
	node := bt.root
	var nilCount, total int
	countNilBT(node, &nilCount, &total)
	fmt.Println("bt trie - nil nodes: ", nilCount, " total nodes: ", total)
	if (total - nilCount) != 563280 {
		t.Errorf("not all prefixes from fv in trie")
	}
}

func TestFullViewNillNodesPCT(t *testing.T) {
	var pc PCTrie
	var ms runtime.MemStats
	fd, err := os.Open("./fv.out")
	if err != nil {
		t.Errorf("cant find file with full view. filename must be fv.out in local dir")
		return
	}
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 2 {
			pc.AddPrefix(StringToPrefix(fields[0], fields[1]))
		}
	}
	fd.Close()
	runtime.GC()
	runtime.ReadMemStats(&ms)
	fmt.Println("pct heap alloc: ", ms.HeapAlloc)
	node := pc.root
	var nilCount, total int
	countNilPCT(node, &nilCount, &total)
	fmt.Println("pct trie - nil nodes: ", nilCount, " total nodes: ", total)
	if (total - nilCount) != 563280 {
		t.Errorf("not all prefixes from fv in trie")
	}
}

func BenchmarkFullViewLookup(b *testing.B) {
	var bt BTrie
	//var ms runtime.MemStats
	fd, err := os.Open("./fv.out")
	if err != nil {
		b.Errorf("cant find file with full view. filename must be fv.out in local dir")
		return
	}
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 2 {
			bt.AddPrefix(StringToPrefix(fields[0], fields[1]))
		}
	}
	fd.Close()
	bt.AddDefault("default")
	ipv4 := IPv4ToUint32NoError("8.8.8.8")
	for i := 0; i < b.N; i++ {
		bt.FindAdj(ipv4)
		//runtime.ReadMemStats(&ms)
	}
	//fmt.Println(bt.FindAdj(ipv4))
	/*
		fmt.Println(ms.HeapAlloc)
		fmt.Println(ms.TotalAlloc)
	*/
}

func BenchmarkFullViewLookupPCT(b *testing.B) {
	var pc PCTrie
	//var ms runtime.MemStats
	fd, err := os.Open("./fv.out")
	if err != nil {
		b.Errorf("cant find file with full view. filename must be fv.out in local dir")
		return
	}
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 2 {
			pc.AddPrefix(StringToPrefix(fields[0], fields[1]))
		}
	}
	fd.Close()
	pc.AddDefault("default")
	ipv4 := IPv4ToUint32NoError("8.8.8.8")
	for i := 0; i < b.N; i++ {
		pc.FindAdj(ipv4)
		//runtime.ReadMemStats(&ms)
	}
	//fmt.Println(pc.FindAdj(ipv4))
	/*
		fmt.Println(ms.HeapAlloc)
		fmt.Println(ms.TotalAlloc)
	*/
}
