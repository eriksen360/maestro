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



 */

func ConfigureClient() *ssh.Client                       { return nil }
func EstablishConnection(client *ssh.Client) *ssh.Client { return nil }
func ParseConfFile(f string) {

	// Should parse a .yaml file and create a set of instructions to be executed
}

func Orchestrate() {

	/* The brain of everything. Handles all connections and states for those connections. Orchestrates the
	processes and communication



	*/

}

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
