// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/lukpueh/git-remote-foo/third_party/pktline"

	"github.com/go-git/go-git/v5/utils/trace"
)

type pkt struct {
	raw     []byte
	header  []byte
	payload []byte
}

// read until flush
func readPkts(w io.Reader) []pkt {
	pkts := []pkt{}
	s := pktline.NewScanner(w)
	for s.Scan() {
		pkts = append(pkts, pkt{s.Bytes(), s.Header(), s.Payload()})
		if s.IsFlushPkt() {
			break
		}
	}
	return pkts
}

// write
func writePkts(w io.Writer, pkts []pkt) {
	for _, pkt := range pkts {
		w.Write(pkt.raw)
	}
}

func transport() error {

	if len(os.Args) < 3 {
		return fmt.Errorf("usage: %s <remote-name> <url>", os.Args[0])
	}

	// Parse URL
	url := strings.TrimLeft(os.Args[2], "foo://")
	parts := strings.Split(url, ":")
	host := parts[0]
	repo := parts[1]

	// Setup logging
	logFilePath := os.Getenv("FOO_LOG_FILE")
	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer logFile.Close()

	// Write standard log to FOO_LOG_FILE
	log.SetOutput(logFile)

	// ... and packet trace from go-git too.
	trace.SetLogger(log.Default())
	trace.SetTarget(trace.Packet)

	// Enter text based top-level menu
	gitScanner := bufio.NewScanner(os.Stdin)
	for gitScanner.Scan() {
		gitLine := gitScanner.Text()

		switch {
		case gitLine == "capabilities":
			os.Stdout.WriteString("stateless-connect\npush\n\n")

		case strings.HasPrefix(gitLine, "stateless-connect"):
			service := "git-upload-pack"
			// TODO: We only support stateless-connect fetch
			// write "fallback", if git wants to do a stateless-connect push,
			// and fail if it's something else

			// https://git-scm.com/docs/gitprotocol-pack#_ssh_transport
			remoteCmd := fmt.Sprintf("%s '%s'", service, repo)

			// TODO: Add this to test container ssh config
			//	 ~/.ssh/config
			//	 Host *
			//	 StrictHostKeyChecking accept-new
			sshCmd := exec.Command("ssh", "-o", "SendEnv=GIT_PROTOCOL", "-o", "StrictHostKeyChecking=accept-new", host, remoteCmd)
			sshCmd.Stderr = os.Stderr
			sshStdin, _ := sshCmd.StdinPipe()
			sshStdout, _ := sshCmd.StdoutPipe()
			sshCmd.Env = append(os.Environ(), "GIT_PROTOCOL=version=2")
			sshCmd.Start()

			// Tell git that connection was established
			// TODO: or failed (exit with error)
			os.Stdout.WriteString("\n")

			// Enter packet line communication

			// 1. Capability Advertisement
			pkts := readPkts(sshStdout)
			writePkts(os.Stdout, pkts)

			// 2. Command Request
			pkts = readPkts(os.Stdin)
			writePkts(sshStdin, pkts)

			// 3. List refs
			pkts = readPkts(sshStdout)
			// TODO: Do gittuf stuff here
			// writePkts(os.Stdout, pkts)

			// 4. Fetch
			// ...

			// 5. Send packfile
			// ...

			return nil
		}
	}
	return nil
}

func main() {
	if err := transport(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
