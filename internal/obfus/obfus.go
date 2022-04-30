package obfus

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"github.com/goura/macaddrobfus/internal/macaddr"
)

type ObfusedMacAddr struct {
	*macaddr.MacAddr
}

type Obfuscator struct {
	Secret []byte
}

func New(secret string) *Obfuscator {
	o := Obfuscator{[]byte(secret)}
	return &o
}

func (o *Obfuscator) Obfuscate(addr *macaddr.MacAddr) (*ObfusedMacAddr, error) {
	// Obfusecate the NIC part of a MAC address
	// Roughly speaking we want to do something like this
	// NIC = mod(NIC + hash(OUI, secret), 256^3)

	// Skip if the address is a "reserved" address (ff:ff:ff:ff:ff:ff)
	if addr.Reserved {
		return &ObfusedMacAddr{addr}, nil
	}

	// Put secret then OUI to calculate the hash to use
	// (For one secret, there will be a hash used for each different OUI)
	// Just to be clear, we will generate a 256 bit hash,
	// but will only use the last 24 bits of it
	m := sha256.New()
	m.Write(o.Secret)
	m.Write(addr.OUI)
	hash := m.Sum(nil)

	// Calculate the sum of the NIC part and the last 3 bytes of this hash

	last3 := hash[len(hash)-3:] // last 3 bytes of the hash

	// hash converted to uint32
	var toAdd uint32 = 0
	toAdd = toAdd | uint32(last3[2])
	toAdd = toAdd | (uint32(last3[1]) << 8)
	toAdd = toAdd | (uint32(last3[0]) << 16)

	// NIC part converted to uint32
	var nicInt uint32 = 0
	nicInt = nicInt | uint32(addr.NIC[2])
	nicInt = nicInt | (uint32(addr.NIC[1]) << 8)
	nicInt = nicInt | (uint32(addr.NIC[0]) << 16)

	sum := nicInt + toAdd

	// Convert the sum back to a byte array
	buf := make([]byte, binary.MaxVarintLen32)
	binary.BigEndian.PutUint32(buf, sum)

	// Get the least significant 3 bytes and use as the new NIC part
	nic := buf[1:4]

	// Construct a new address and return it
	s := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", addr.OUI[0], addr.OUI[1], addr.OUI[2], nic[0], nic[1], nic[2])

	newAddr, err := macaddr.New(s)
	if err != nil {
		return nil, err
	}

	return &ObfusedMacAddr{newAddr}, nil
}
