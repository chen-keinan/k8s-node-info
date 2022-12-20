package collector

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

const (
	WorkerNode   = "worker"
	MasterNode   = "master"
	ShellCommand = "sh"
)

// Shell command interface to preform shell exec commands
type Shell interface {
	Execute(commandArgs string) (string, error)
	FindNodeType() (string, error)
}

// NewShellCmd instansiate new shell command
func NewShellCmd() Shell {
	return &cmd{}
}

type cmd struct {
}

// Execute execute a shell command and retun it output or error
func (e *cmd) Execute(commandArgs string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cm := exec.Command(ShellCommand, "-c", commandArgs)
	cm.Stdout = &stdout
	cm.Stderr = &stderr
	err := cm.Run()
	if err != nil {
		return "", err
	}
	if len(stderr.String()) > 0 {
		return "", errors.New(stderr.String())
	}
	return strings.TrimSuffix(stdout.String(), "\n"), nil
}

func (e *cmd) FindNodeType() (string, error) {
	masterConfigFiles := []string{
		"ls /etc/kubernetes/controller-manager.conf",
		"ls /etc/kubernetes/manifests/kube-apiserver.yaml",
		"ls /etc/kubernetes/scheduler.conf",
	}
	for _, path := range masterConfigFiles {
		output, err := e.Execute(path)
		if !strings.Contains(path, output) || err != nil {
			return WorkerNode, nil
		}
	}
	return MasterNode, nil
}
