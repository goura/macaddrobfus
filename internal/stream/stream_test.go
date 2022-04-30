package stream

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/rgoura/macaddrobfus/internal/macaddr"
	"github.com/rgoura/macaddrobfus/internal/obfus"
)

func TestNoAddrPlainStrings(t *testing.T) {
	var buf bytes.Buffer

	fixtures := []string{
		"hoge",
		"hogehoge",
		"hogehogehoge",
		"hogehogehogehoge",
		"hogehogehogehogehoge",
		"hogehogehogehogehogehoge",
	}

	for _, s := range fixtures {
		buf.Reset()
		r := strings.NewReader(s)
		st := New(r, &buf, "secret")

		err := st.ReadWrite()
		if err != nil {
			t.Error(err)
		}

		if s != buf.String() {
			t.Errorf("String didn't match actual:%s expected:%s", buf.String(), s)
		}
	}
}

func TestAddrPlainStrings(t *testing.T) {
	var buf bytes.Buffer
	buf.Grow(64)

	addrStr := "11:22:33:44:55:66"
	addr, _ := macaddr.New(addrStr)

	o := obfus.New("secret")
	oaddr, _ := o.Obfuscate(addr)

	fixtures := []string{
		"%[1]vhoge",
		/*
			"hoge%[1]vhoge",
			"hoge%[1]vhoge%[1]vhoge",
			"hogehogehoge%[1]vhoge",
		*/
	}

	for _, fstr := range fixtures {
		buf.Reset()

		s := fmt.Sprintf(fstr, addrStr)
		expectedStr := fmt.Sprintf(fstr, oaddr.HardwareAddr.String())

		r := strings.NewReader(s)

		st := New(r, &buf, "secret")
		err := st.ReadWrite()
		if err != nil {
			t.Error(err)
		}

		if expectedStr != buf.String() {
			t.Errorf("String didn't match actual:%s expected:%s", buf.String(), expectedStr)
		}
	}
}
