package main

// import (
// 	// "strings"
// 	"golang.org/x/crypto/ssh"
// 	"golang.org/x/crypto/ssh/agent"
// 	"log"
// 	"net"
// 	"os"
// )
// go ServeAgent(NewKeyring(), c2)
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"

	// "golang.org/x/crypto/ssh"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type myExtendedAgent struct {
	agent.Agent
	// alsoAgent agent.ExtendedAgent
	// agent.ExtendedAgent
	// my stuff here
}

func testTypes() {
	var _ agent.Agent = myExtendedAgent{}               // Verify that T implements I.
	var _ agent.Agent = (*myExtendedAgent)(nil)         // Verify that *T implements I.
	var _ agent.ExtendedAgent = myExtendedAgent{}       // Verify that T implements I.
	var _ agent.ExtendedAgent = (*myExtendedAgent)(nil) // Verify that *T implements I.
}

// // SignWithFlags signs like Sign, but allows for additional flags to be sent/received
func (r myExtendedAgent) SignWithFlags(key ssh.PublicKey, data []byte, flags agent.SignatureFlags) (*ssh.Signature, error) {
	// not using the extended flags
	return nil, agent.ErrExtensionUnsupported
}

// Extension processes a custom extension request. Standard-compliant agents are not
// required to support any extensions, but this method allows agents to implement
// vendor-specific methods or add experimental features. See [PROTOCOL.agent] section 4.7.
// If agent extensions are unsupported entirely this method MUST return an
// ErrExtensionUnsupported error. Similarly, if just the specific extensionType in
// the request is unsupported by the agent then ErrExtensionUnsupported MUST be
// returned.
//
// In the case of success, since [PROTOCOL.agent] section 4.7 specifies that the contents
// of the response are unspecified (including the type of the message), the complete
// response will be returned as a []byte slice, including the "type" byte of the message.
func (r myExtendedAgent) Extension(extensionType string, contents []byte) ([]byte, error) {
	log.Print("Extension fired:", extensionType, contents)
	// []byte("Here is a string....")
	return []byte("yo Bogus MEssage"), nil
	// return nil, agent.ErrExtensionUnsupported
}

func main() {

	// Create the SSH Agent Named Pipe
	tmpDir, err := ioutil.TempDir("", "ssh-go")
	if err != nil {
		log.Fatal(err)
	}
	pid := os.Getpid()
	pipeFile := filepath.Join(tmpDir, fmt.Sprint("agent.", pid))
	listener, err := net.Listen("unix", pipeFile)
	defer listener.Close()
	defer os.Remove(pipeFile)
	defer os.RemoveAll(tmpDir)
	log.Println(listener)

	myExt := myExtendedAgent{agent.NewKeyring()}
	var ea agent.ExtendedAgent
	ea = &myExt
	var a agent.Agent
	a = ea

	if eaa, ok := a.(agent.ExtendedAgent); !ok {
		log.Fatal(a, ok)
	} else {
		eaa.Extension("yo", nil)
	}
	// ssh-agent has a UNIX socket under $SSH_AUTH_SOCK
	// Output the ssh-agent compatible env vars for evaling.
	fmt.Printf("SSH_AUTH_SOCK=%s; export SSH_AUTH_SOCK;\n", pipeFile)
	fmt.Printf("SSH_AGENT_PID=%d; export SSH_AGENT_PID;\n", pid)
	fmt.Printf("echo Agent pid %d;\n", pid)

	// wait for clients...
	for {
		con, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		println("new Connection", con)
		handleClient(myExt, con) // should be a go sub
	}
}

func handleClient(a myExtendedAgent, connection net.Conn) {
	defer println("closed", connection)
	defer connection.Close()
	agent.ServeAgent(a, connection)
}
