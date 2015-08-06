package trie

import (
	"errors"
	"strconv"
	"strings"
)

const (
	IPV4_LEN          = 32
	IPV4_TOTAL_OCTETS = 4
	BITS_IN_OCTET     = 8
	ALL_ONES_32BIT    = uint32(4294967295)
)

func IPv4ToUint32(ipv4 string) (uint32, error) {
	octets := strings.Split(ipv4, ".")
	if len(octets) != IPV4_TOTAL_OCTETS {
		return 0, errors.New("cant convert ipv4 to string")

	}
	ipv4int := uint32(0)
	for cntr := 0; cntr < IPV4_TOTAL_OCTETS; cntr++ {
		tmpVal, err := strconv.Atoi(octets[cntr])
		if err != nil {
			return 0, errors.New("cant convert ipv4 to string")

		}
		ipv4int += uint32(tmpVal << uint((3-cntr)*BITS_IN_OCTET))

	}
	return ipv4int, nil

}

func IPv4ToUint32NoError(ipv4 string) uint32 {
	addr, _ := IPv4ToUint32(ipv4)
	return addr
}

func StringToPrefix(stringPrefix, adjency string) Prefix {
	var prefix Prefix
	prefixWithMask := strings.Split(stringPrefix, "/")
	if len(prefixWithMask) != 2 {
		return prefix
	}
	net := prefixWithMask[0]
	mask := prefixWithMask[1]
	netU32, err := IPv4ToUint32(net)
	if err != nil {
		return Prefix{}
	}
	maskU8, err := strconv.ParseUint(mask, 10, 8)
	if err != nil {
		return Prefix{}
	}
	prefix.prefix = uint32(netU32)
	prefix.prefixLen = uint8(maskU8)
	prefix.adj = adjency
	return prefix
}
