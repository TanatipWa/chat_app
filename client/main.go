// Chat Client
package main

import (
	"io"
	"log"
	"net"
	"os"
)

// ^Z = EOF
func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		io.Copy(os.Stdout, conn)
	}()
	io.Copy(conn, os.Stdin) // waiting...
}
