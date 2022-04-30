package stream

import (
	"bufio"
	"bytes"
	"io"

	"github.com/rgoura/macaddrobfus/internal/macaddr"
	"github.com/rgoura/macaddrobfus/internal/obfus"
)

var MacAddrLength int = len("xx:xx:xx:xx:xx:xx")

type StreamObfuscator struct {
	Reader     *bufio.Reader
	Writer     *bufio.Writer
	Obfuscator *obfus.Obfuscator
}

func New(r io.Reader, w io.Writer, secret string) *StreamObfuscator {
	s := StreamObfuscator{
		Reader:     bufio.NewReader(r),
		Writer:     bufio.NewWriter(w),
		Obfuscator: obfus.New(secret),
	}
	return &s
}

func (st *StreamObfuscator) ReadWrite() error {
	var buf bytes.Buffer
	buf.Grow(32)

	for {
		// Try to read 17 chars
		buflen := buf.Len()
		for i := 0; i < (MacAddrLength - buflen); i++ {
			b, err := st.Reader.ReadByte()
			if err == io.EOF {
				st.Writer.Write(buf.Bytes())
				buf.Reset()
				goto escape
			}
			buf.WriteByte(b)
		}

		if buf.Len() != MacAddrLength {
			panic("Couldn't read from the input for some reason")
		}

		// 17 chars were read, now test whether it's a MAC address
		s := buf.String()
		addr, err := macaddr.New(s)
		if err == nil {
			// Those 17 chars were a MAC address
			// Let's obfuscate it
			buf.Reset() // Clear the buffer
			newAddr, err := st.Obfuscator.Obfuscate(addr)
			if err != nil {
				return err
			}
			st.Writer.Write([]byte(newAddr.HardwareAddr.String()))
			continue
		}

		// 17 chars were read but were not a MAC address
		// Write the first byte but keep the rest 16
		b, err := buf.ReadByte()
		if err == io.EOF {
			panic("TODO: this shouldn't happen")
		}
		st.Writer.WriteByte(b)
	}

escape:
	st.Writer.Flush()
	return nil
}
