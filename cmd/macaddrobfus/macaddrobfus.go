package main

import (
	"bufio"
	"os"

	"github.com/rgoura/macaddrobfus/internal/stream"
)

func main() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	st := stream.New(in, out, "secret")
	st.ReadWrite()
}
