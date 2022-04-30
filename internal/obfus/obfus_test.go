package obfus

import (
	"testing"

	"github.com/goura/macaddrobfus/internal/macaddr"
)

func TestObfuscateRuns(t *testing.T) {
	secret := "abcdefg"
	obfus := New(secret)
	addr, err := macaddr.New("00:11:22:33:44:55")
	if err != nil {
		t.Error(err)
	}
	_, err = obfus.Obfuscate(addr)
	if err != nil {
		t.Error(err)
	}
}

func TestObfuscateSkipBroadcast(t *testing.T) {
	secret := "abcdefg"
	obfus := New(secret)
	addr, err := macaddr.New("ff:ff:ff:ff:ff:ff")
	if err != nil {
		t.Error(err)
	}
	oaddr, err := obfus.Obfuscate(addr)
	if err != nil {
		t.Error(err)
	}

	if !macaddr.IsEqual(addr, oaddr.MacAddr) {
		t.Errorf("Obfuscator should skip ff:ff:ff:ff:ff:ff")
	}
}

func TestObfuscation(t *testing.T) {
	secret1 := "abcdefg"
	secret2 := "xyx0123"
	obfusA := New(secret1)
	obfusB := New(secret2)

	addr1a, _ := macaddr.New("00:11:22:33:44:55")
	addr1b, _ := macaddr.New("00:11:22:44:55:66")
	addr2a, _ := macaddr.New("01:23:45:33:44:55")

	oaddrA1a, _ := obfusA.Obfuscate(addr1a)
	oaddrA1b, _ := obfusA.Obfuscate(addr1b)
	oaddrA2a, _ := obfusA.Obfuscate(addr2a)

	oaddrB1a, _ := obfusB.Obfuscate(addr1a)

	// Obfuscate should generate a different MAC address from the original one
	if macaddr.IsEqual(addr1a, oaddrA1a.MacAddr) {
		t.Errorf("Obfuscate should generate different MAC address")
	}

	// But it's OUI part should be the same
	if !macaddr.IsOUISame(addr1a, oaddrA1a.MacAddr) {
		t.Errorf("Obfuscate should generate a MAC address with a same OUI")
	}

	// Different addresses should result to different results
	if macaddr.IsNICSame(oaddrA1a.MacAddr, oaddrA1b.MacAddr) {
		t.Errorf("Different addresses should result to different results")
	}

	// Different OUIs should result to different results
	if macaddr.IsNICSame(oaddrA1a.MacAddr, oaddrA2a.MacAddr) {
		t.Errorf("Different OUIs should result to different results")
	}

	// Different obfuscators should generate different results
	if macaddr.IsEqual(oaddrA1a.MacAddr, oaddrB1a.MacAddr) {
		t.Errorf("Different Obfuscators should generate different results")
	}
}
