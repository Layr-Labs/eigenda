package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"golang.org/x/crypto/ssh"
)

// TODO consider dropping password support

// SSHSession encapsulates an SSH session with a remote host.
type SSHSession struct {
	logger  logging.Logger
	client  *ssh.Client
	session *ssh.Session
	user    string
	host    string
	port    uint64
	keyPath string
	verbose bool
}

// Create a new SSH session to a remote host.
func NewSSHSession(
	logger logging.Logger,
	user string,
	host string,
	port uint64,
	keyPath string,
	verbose bool,
) (*SSHSession, error) {
	config := &ssh.ClientConfig{
		User: user,
	}

	key, err := ssh.ParsePrivateKey([]byte(keyPath))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}
	config.Auth = []ssh.AuthMethod{
		ssh.PublicKeys(key),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s:%d: %v", host, port, err)
	}

	session, err := client.NewSession()
	if err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("failed to create SSH session: %v", err)
	}

	return &SSHSession{
		logger:  logger,
		client:  client,
		session: session,
		user:    user,
		host:    host,
		port:    port,
		keyPath: keyPath,
		verbose: verbose,
	}, nil
}

// Close the SSH session.
func (s *SSHSession) Close() error {
	err := s.session.Close()
	if err != nil {
		return fmt.Errorf("failed to close SSH session: %v", err)
	}
	err = s.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close SSH client: %v", err)
	}

	return nil
}

// RemoteLs executes "ls" on the remote machine and returns the list of files in the specified path.
func (s *SSHSession) Ls(path string) ([]string, error) {
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	s.session.Stdout = &stdoutBuf
	s.session.Stderr = &stderrBuf

	command := fmt.Sprintf("ls '%s'", path)
	if s.verbose {
		s.logger.Infof("Executing remotely: %s", command)
	}

	if err := s.session.Run(command); err != nil {
		return nil, fmt.Errorf("failed to execute command '%s': %v, stderr: %s",
			command, err, stderrBuf.String())
	}

	files := strings.Split(stdoutBuf.String(), "\n")
	return files, nil
}

// Search for all files matching a regex inside a file tree at the specified root path.
func (s *SSHSession) FindRegex(root string, regex string) ([]string, error) {

	command := fmt.Sprintf("find '%s' -type f | grep -E '%s'", root, regex)
	if s.verbose {
		s.logger.Infof("Executing remotely: %s", command)
	}

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	s.session.Stdout = &stdoutBuf
	s.session.Stderr = &stderrBuf

	if err := s.session.Run(command); err != nil {
		return nil, fmt.Errorf("failed to execute command '%s': %v, stderr: %s",
			command, err, stderrBuf.String())
	}

	files := strings.Split(stdoutBuf.String(), "\n")

	return files, nil
}

// Mkdirs creates the specified directory on the remote machine, including any necessary parent directories.
func (s *SSHSession) Mkdirs(path string) error {
	command := fmt.Sprintf("mkdir -p '%s'", path)
	if s.verbose {
		s.logger.Infof("Executing remotely: %s", command)
	}

	var stderrBuf bytes.Buffer
	s.session.Stderr = &stderrBuf

	if err := s.session.Run(command); err != nil {
		return fmt.Errorf("failed to create directory '%s': %v, stderr: %s",
			path, err, stderrBuf.String())
	}

	return nil
}

// Rsync transfers files from the local machine to the remote machine using rsync.
func (s *SSHSession) Rsync(sourceFile string, destFile string) error {
	sshCmd := fmt.Sprintf("ssh -i %s -p %d", s.keyPath, s.port)
	target := fmt.Sprintf("%s@%s:%s", s.user, s.host, destFile)

	arguments := []string{
		"rsync",
		"-avz", // TODO look into flag use
		"-e", "'" + sshCmd + "'",
		sourceFile,
		target,
	}

	if s.verbose {
		s.logger.Infof("Executing remotely: %s", strings.Join(arguments, " "))
	}

	cmd := exec.Command(arguments[0], arguments[1:]...)

	if s.verbose {
		cmd.Stdout = os.Stdout
	}
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to rsync data: %v", err)
	}

	return nil
}
