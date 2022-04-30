package macaddr

import (
	"bytes"
	"testing"
)

func TestMacAddrOUINIC(t *testing.T) {
	fixtures := []struct {
		notation    string
		expectedOUI []byte
		expectedNIC []byte
	}{
		{
			notation:    "01:23:45:67:89:ab",
			expectedOUI: []byte{0x01, 0x23, 0x45},
			expectedNIC: []byte{0x67, 0x89, 0xab},
		},
		{
			notation:    "01:23:45:67:89:AB",
			expectedOUI: []byte{0x01, 0x23, 0x45},
			expectedNIC: []byte{0x67, 0x89, 0xab},
		},
		{
			notation:    "01-23-45-67-89-ab",
			expectedOUI: []byte{0x01, 0x23, 0x45},
			expectedNIC: []byte{0x67, 0x89, 0xab},
		},
		{
			notation:    "01-23-45-67-89-AB",
			expectedOUI: []byte{0x01, 0x23, 0x45},
			expectedNIC: []byte{0x67, 0x89, 0xab},
		},
	}

	for _, fixture := range fixtures {
		addr, err := New(fixture.notation)
		if err != nil {
			t.Errorf("macaddr.New doesn't work with %s", fixture.notation)
		}
		if bytes.Compare(addr.OUI, fixture.expectedOUI) != 0 {
			t.Errorf("OUI expected: %v, actual: %v", fixture.expectedOUI, fixture.notation)
		}
		if bytes.Compare(addr.NIC, fixture.expectedNIC) != 0 {
			t.Errorf("NIC expected: %v, actual: %v", fixture.expectedNIC, fixture.notation)
		}
		if addr.Reserved {
			t.Errorf("Reserved is not false for %v", fixture.notation)
		}
	}
}

func TestMacAddrBroadcast(t *testing.T) {
	s := "ff-ff-ff-ff-ff-ff"
	addr, err := New(s)
	if err != nil {
		t.Error(err)
	}
	if !addr.Reserved {
		t.Errorf("macaddr.New doens't set reserved for %s", s)
	}
}

func TestIsEqual(t *testing.T) {
	addr1, _ := New("00:11:22:33:44:55")
	addr2, _ := New("11:22:33:44:55:66")
	if IsEqual(addr1, addr2) {
		t.Errorf("%v and %v is not equal", addr1, addr2)
	}

	addr3, _ := New("00:11:22:33:44:55")
	if !IsEqual(addr1, addr3) {
		t.Errorf("%v and %v is equal", addr1, addr3)
	}
}

func TestIsOUISame(t *testing.T) {
	addr1, _ := New("00:11:22:33:44:55")
	addr2, _ := New("11:22:33:44:55:66")
	if IsOUISame(addr1, addr2) {
		t.Errorf("%v and %v 's OUIs are not equal", addr1, addr2)
	}

	addr3, _ := New("00:11:22:44:55:66")
	if !IsOUISame(addr1, addr3) {
		t.Errorf("%v and %v 's OUIs are equal", addr1, addr3)
	}
}

func TestMacAddrIllegalFormat(t *testing.T) {
	s := "02:00:5e:10:00:00:00:01"
	_, err := New(s)
	if err == nil {
		t.Errorf("macaddr.New doesn't deny %s", s)
	}
}
