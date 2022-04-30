package main

import (
	"bufio"
	"bytes"
	"flag"
	"math/rand"
	"os"
	"time"

	"github.com/goura/macaddrobfus/internal/stream"
)

func randomString() string {
	var buf bytes.Buffer

	// Initialize seed with time
	rand.Seed(time.Now().UnixMicro())

	// Pick 32 random letters
	const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := 0; i < 32; i++ {
		idx := rand.Intn(len(letters))
		buf.WriteByte(letters[idx])
	}
	return buf.String()
}

func main() {
	// Define flags
	flagSecret := flag.String("secret", "not a secret", "secret key used to obfuscate strings")
	flag.Parse()

	// If secret is not specified, generate it randomely
	var secret string
	if *flagSecret == "not a secret" {
		secret = randomString()
	} else {
		secret = *flagSecret
	}

	// Do the magic
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	st := stream.New(in, out, secret)
	st.ReadWrite()
}
