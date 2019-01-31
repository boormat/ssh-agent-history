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

	"golang.org/x/crypto/ssh/agent"
)

type myExtendedAgent struct {
	agent.Agent
	// my stuff here
}

// func NewKeyring() Agent {

// unc (r *myExtendedAgent) RemoveAll() error {
// 	r.mu.Lock()

// func (r *keyring) Add(key AddedKey) error {

// Agent represents the capabilities of an ssh-agent.
//List returns the identities known to the agent.
// func (r *myExtendedAgent) List() ([]*agent.Key, error) {
// 	var a agent.Agent
// 	println(a)
// 	a = r
// 	return a.List()
// }

// // Sign has the agent sign the data using a protocol 2 key as defined
// // in [PROTOCOL.agent] section 2.6.2.
// func (r *myExtendedAgent) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) {
// 	return r.memKeyRing.Sign(key, data)
// }

// // Add adds a private key to the agent.
// func (r *myExtendedAgent) Add(key agent.AddedKey) error {
// 	return r.memKeyRing.Add(key)
// }

// // Remove removes all identities with the given public key.
// func (r *myExtendedAgent) Remove(key ssh.PublicKey) error {
// 	return r.memKeyRing.Remove(key)
// }

// // RemoveAll removes all identities.
// func (r *myExtendedAgent) RemoveAll() error {
// 	return r.memKeyRing.RemoveAll()
// }

// // Lock locks the agent. Sign and Remove will fail, and List will empty an empty list.
// func (r *myExtendedAgent) Lock(passphrase []byte) error {
// 	return r.memKeyRing.Lock(passphrase)
// }

// // Unlock undoes the effect of Lock
// func (r *myExtendedAgent) Unlock(passphrase []byte) error {
// 	return r.memKeyRing.Unlock(passphrase)
// }

// // Signers returns signers for all the known keys.
// func (r *myExtendedAgent) Signers() ([]ssh.Signer, error) {
// 	return r.memKeyRing.Signers()
// }

// // SignWithFlags signs like Sign, but allows for additional flags to be sent/received
// func (r *myExtendedAgent) SignWithFlags(key ssh.PublicKey, data []byte, flags agent.SignatureFlags) (*ssh.Signature, error) {
// 	// not using the extended flags
// 	return r.memKeyRing.Sign(key, data)
// }

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
func (r *myExtendedAgent) Extension(extensionType string, contents []byte) ([]byte, error) {
	log.Print("Extension fired:", extensionType, contents)
	return nil, agent.ErrExtensionUnsupported
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
	var a agent.Agent
	a = &myExt

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
		go agent.ServeAgent(a, con)
	}
}
