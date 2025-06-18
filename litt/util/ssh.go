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
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO
	}

	exists, err := Exists(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to check if key exists: %v", err)
	}
	if !exists {
		return nil, fmt.Errorf("private key does not exist at path: %s", keyPath)
	}

	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %v", err)
	}

	key, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}
	config.Auth = []ssh.AuthMethod{
		ssh.PublicKeys(key),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s port %d: %v", host, port, err)
	}

	return &SSHSession{
		logger:  logger,
		client:  client,
		user:    user,
		host:    host,
		port:    port,
		keyPath: keyPath,
		verbose: verbose,
	}, nil
}

// Close the SSH session.
func (s *SSHSession) Close() error {
	err := s.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close SSH client: %v", err)
	}

	return nil
}

// RemoteLs executes "ls" on the remote machine and returns the list of files in the specified path.
func (s *SSHSession) Ls(path string) ([]string, error) {
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	session, err := s.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer func() {
		_ = session.Close()
	}()

	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	command := fmt.Sprintf("ls '%s'", path)
	if s.verbose {
		s.logger.Infof("Executing remotely: %s", command)
	}

	if err := session.Run(command); err != nil {
		return nil, fmt.Errorf("failed to execute command '%s': %v, stderr: %s",
			command, err, stderrBuf.String())
	}

	files := strings.Split(stdoutBuf.String(), "\n")
	return files, nil
}

// Search for all files matching a regex inside a file tree at the specified root path.
func (s *SSHSession) FindFiles(root string, extensions []string) ([]string, error) {

	session, err := s.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer func() {
		_ = session.Close()
	}()

	command := fmt.Sprintf("find \"%s\" -type f", root)
	if s.verbose {
		s.logger.Infof("Executing remotely: %s", command)
	}

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	err = session.Run(command)
	if err != nil {
		errString := stderrBuf.String()
		if !strings.Contains(errString, "No such file or directory") {
			return nil, fmt.Errorf("failed to execute command '%s': %v, stderr: %s",
				command, err, errString)
		}
		// There are no files since the directory does not exist.
		return []string{}, nil
	}

	files := strings.Split(stdoutBuf.String(), "\n")

	filteredFiles := make([]string, 0, len(files))
	for _, file := range files {
		if file == "" {
			continue // Skip empty lines
		}
		for _, ext := range extensions {
			if strings.HasSuffix(file, ext) {
				filteredFiles = append(filteredFiles, file)
				break // Stop checking other extensions once a match is found
			}
		}
	}

	return filteredFiles, nil
}

// Mkdirs creates the specified directory on the remote machine, including any necessary parent directories.
func (s *SSHSession) Mkdirs(path string) error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %v", err)
	}
	defer func() {
		_ = session.Close()
	}()

	command := fmt.Sprintf("mkdir -p '%s'", path)
	if s.verbose {
		s.logger.Infof("Executing remotely: %s", command)
	}

	var stderrBuf bytes.Buffer
	session.Stderr = &stderrBuf

	if err := session.Run(command); err != nil {
		return fmt.Errorf("failed to create directory '%s': %v, stderr: %s",
			path, err, stderrBuf.String())
	}

	return nil
}

// Rsync transfers files from the local machine to the remote machine using rsync.
func (s *SSHSession) Rsync(sourceFile string, destFile string) error {
	sshCmd := fmt.Sprintf("ssh -i %s -p %d", s.keyPath, s.port)
	target := fmt.Sprintf("%s@%s:%s", s.user, s.host, destFile)

	// If the source file is a symlink, we actually want to send the thing the symlink points to.
	fileInfo, err := os.Lstat(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %w", sourceFile, err)
	}
	isSymlink := fileInfo.Mode()&os.ModeSymlink != 0

	if isSymlink {
		// Resolve the symlink to get the actual file it points to
		sourceFile, err = os.Readlink(sourceFile)
		if err != nil {
			return fmt.Errorf("failed to resolve symlink %s: %w", sourceFile, err)
		}
	}

	arguments := []string{
		"rsync",
		"-z",
		"-e", sshCmd,
		sourceFile,
		target,
	}

	if s.verbose {
		s.logger.Infof("Executing: %s", strings.Join(arguments, " "))
	}

	cmd := exec.Command(arguments[0], arguments[1:]...)
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to rsync data: %v", err)
	}

	return nil
}
