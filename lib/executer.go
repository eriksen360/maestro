package lib

import (
	"bytes"
	"context"
	"fmt"
	scp "github.com/bramvdbogaerde/go-scp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
)

/*
 SSH Compression
 SSH Transport Protocol
 SSH Channels for parallel connection to same host


*/

func ConfigureClient() *ssh.Client                       { return nil }
func EstablishConnection(client *ssh.Client) *ssh.Client { return nil }

func Orchestrate() {

	/* The brain of everything. Handles all connections and states for those connections. Orchestrates the
	processes and communication.

	Maestro can consume a configuration file without applying it

	*/

	/*
		while ACTIVE {

			if user.startsTask(host, task, **kwargs) {

				create go routine
				Find host and task from saved instruction sets

				If connection for host already open, use new session/channel (Potentiel for regulating bandwidth for each client)
				Apply instruction set on session. If scheduled, sleep async until ready (close session meanwhile)

				We must at the same time use events and signals to resolve dependencies

				The trick here is not the functionality but high availability, performance and bandwidth.
				There is plenty of oppurtunity for optimization and gain here if done correctly

			}




		}



	*/

}

/* Implement a MsgQueue */

type Handler interface {
}

func UserInterfaceEventHandle() {}

func ExecuteSingleCommand(ipAddr string) {

	// var hostKey ssh.PublicKey
	var signer = GetSignedPrivateKey("/home/mathias/.ssh/bfs_thinkpad")
	var network = "tcp"

	clientConf := &ssh.ClientConfig{
		// Config: ssh.Config{},
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: must work with 'ssh.FixedHostKey(hostKey)'
	}

	client, err := ssh.Dial(network, ipAddr, clientConf)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer func(client *ssh.Client) {
		err := client.Close()
		if err != nil {
			log.Fatal("Failed to close client connection", err)
		}
	}(client)

	/* Once a client connection has been established a session can be created */
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create client session: ", err)
	}
	defer func(session *ssh.Session) {
		err := session.Close()
		if err != nil {
		}
	}(session)

	var b bytes.Buffer
	session.Stdout = &b

	if err := session.Run("ls"); err != nil {
		log.Fatal("Failed to run", err.Error())
	}
	fmt.Println(b.String())

}

func TransferFileToRemote(client *ssh.Client) (int, error) {

	var localPath = "/home/mathias/payload.txt"
	var remotePath = "~/payload.txt"

	_client, err := scp.NewClientBySSH(client)
	if err != nil {
		log.Fatal("Could not establish scp client", err)
	}

	f, _ := os.Open(localPath)
	err = _client.CopyFromFile(context.Background(), *f,
		remotePath, "0655")

	return 0, nil
}

func TransferFileFromRemote(client *ssh.Client) (int, error) {

	// scp.CopyFromRemote

	return 0, nil
}

func GetSignedPrivateKey(keyPath string) ssh.Signer {

	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("Could not read from key file: ", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("unable to parse private key: ", err)
	}

	return signer
}
