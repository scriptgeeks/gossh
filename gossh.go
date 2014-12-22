package main

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Hosts []string
	Auth  Auth
	Tasks []map[string]string
}

type Auth struct {
	User string
	Pass string
}

func execute(config *ssh.ClientConfig, hostname string, taskcmd string) {
	client, err := ssh.Dial("tcp", hostname+":22", config)
	if err != nil {
		panic("unable to connect: %s" + err.Error())
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(taskcmd); err != nil {
		color.Red("[%s]", hostname)
		fmt.Print(err)
	} else {
		color.Green("[%s]", hostname)
		fmt.Print(b.String())
	}
}

func main() {
	flag.Parse()

	bytes, err := ioutil.ReadFile("gossh.yaml")
	if err != nil {
		panic("read config file opssh.yaml error: " + err.Error())
	}

	var config Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}

	sshconf := &ssh.ClientConfig{
		User: config.Auth.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Auth.Pass),
		},
	}

	tasks := make(map[string]map[string]string)
	for _, task := range config.Tasks {
		tasks[task["name"]] = task
	}

	for _, hostname := range config.Hosts {
		taskcmd := tasks[flag.Arg(0)]["cmd"]
		execute(sshconf, hostname, taskcmd)

	}

}
