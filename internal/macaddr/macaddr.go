package macaddr

import (
	"bytes"
	"errors"
	"net"
)

type MacAddr struct {
	net.HardwareAddr
	OUI      []byte
	NIC      []byte
	Reserved bool
}

func New(s string) (*MacAddr, error) {
	hw, err := net.ParseMAC(s)
	if err != nil {
		return nil, err
	}
	if len(hw) > 6 {
		return nil, errors.New("not a 6 bytes MAC address")
	}
	oui := hw[0:3]
	nic := hw[3:6]
	reserved := isReserved(hw)
	addr := MacAddr{hw, oui, nic, reserved}

	return &addr, nil
}

func IsEqual(addr1 *MacAddr, addr2 *MacAddr) bool {
	return IsOUISame(addr1, addr2) && IsNICSame(addr1, addr2)
}

func IsOUISame(addr1 *MacAddr, addr2 *MacAddr) bool {
	return bytes.Compare(addr1.OUI, addr2.OUI) == 0
}

func IsNICSame(addr1 *MacAddr, addr2 *MacAddr) bool {
	return bytes.Compare(addr1.NIC, addr2.NIC) == 0
}

func isReserved(hw net.HardwareAddr) bool {
	// Broadcast address ff:ff:ff:ff:ff:ff
	if bytes.Compare(hw, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}) == 0 {
		return true
	}
	return false
}

func Helloworld() string {
	return "helloworld"
}
