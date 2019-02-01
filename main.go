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
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	// "golang.org/x/crypto/ssh"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type myExtendedAgent struct {
	agent.Agent
	histfile io.Writer // file to write bash history-ish format to.
}

// Check that my struct implements the interfaces properly.  (otherwise will be ignored silently in agent module)
func testTypes() {
	var _ agent.Agent = myExtendedAgent{}               // Verify that T implements I.
	var _ agent.Agent = (*myExtendedAgent)(nil)         // Verify that *T implements I.
	var _ agent.ExtendedAgent = myExtendedAgent{}       // Verify that T implements I.
	var _ agent.ExtendedAgent = (*myExtendedAgent)(nil) // Verify that *T implements I.
}

// Don't care pass to base impl
func (r myExtendedAgent) SignWithFlags(key ssh.PublicKey, data []byte, flags agent.SignatureFlags) (*ssh.Signature, error) {
	log.Println("Hey got some (sign request)", key)
	return r.Sign(key, data)
}

type historyMsgContent struct {
	command  string
	hostname string
	user     string
}

// In the case of success, since [PROTOCOL.agent] section 4.7 specifies that the contents
// of the response are unspecified (including the type of the message), the complete
// response will be returned as a []byte slice, including the "type" byte of the message.
func (r myExtendedAgent) Extension(extensionType string, contents []byte) ([]byte, error) {
	log.Print("Extension fired:", extensionType, contents)
	switch extensionType {
	case "HISTORY":
		// parse the contents.  3 Strings is expected, command, hostname, userid.
		var req historyMsgContent
		if err := ssh.Unmarshal(contents, &req); err != nil {
			log.Print("Bad Extension Message:", extensionType, contents)
			return nil, err
		}

		ts := time.Now()

		// write timestamp and other meta

		if _, err := r.histfile.Write([]byte(fmt.Sprintf("#%d %s %s\n", ts.Unix(), req.hostname, req.user))); err != nil {
			log.Fatal(err)
		}
		if _, err := r.histfile.Write([]byte(fmt.Sprintln(req.command))); err != nil {
			log.Fatal(err)
		}

	default:
		log.Print("Unknown Extension Type:", extensionType, contents)
		return nil, agent.ErrExtensionUnsupported
	}

	// histfile io.Writer
	return []byte("yo Bogus MEssage"), nil
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

	// config time
	histfilename, _ := os.LookupEnv("AGENT_HISTFILE")
	hf, err := os.OpenFile(histfilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal("Could not create/open history file", histfilename, err)
	}
	defer hf.Close()

	// if _, err := f.Write([]byte("appended some data\n")); err != nil {
	// 	log.Fatal(err)
	// }
	myExt := myExtendedAgent{agent.NewKeyring(), hf}
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
		handleClient(myExt, con) // should be a go sub?
	}
}

func handleClient(a myExtendedAgent, connection net.Conn) {
	defer println("closed", connection)
	defer connection.Close()
	agent.ServeAgent(a, connection)
}
