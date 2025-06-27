package libs

import (
	"fmt"
	"strings"
)

type SSHRuner struct{}

type SSHRunerConfig struct {
	IP                    string
	User                  string
	Password              string
	Script                string
	SSHKey                *string
	Env                   *map[string]string
	DockerUser            *string
	DockerPassword        *string
	DockerImage           *string
	DockerRegistry        *string
	DockerContainerName   *string
	DockerTag             *string
	SetSecretsToScript    *bool
	SetSecretsToContainer *bool
}

func NewSSHRuner() *SSHRuner {
	return &SSHRuner{}
}

func (r *SSHRuner) LoginDockerCommand(confing *SSHRunerConfig) (string, error) {
	command := r.loginDockerCommand(confing)
	return command, nil
}

func (r *SSHRuner) PullDockerCommand(confing *SSHRunerConfig) (string, error) {
	command := r.pullDockerCommand(confing)
	return command, nil
}

func (r *SSHRuner) RunDockerCommand(confing *SSHRunerConfig) (string, error) {
	command := r.runDockerCommand(confing)
	return command, nil
}

func (r *SSHRuner) CreateScriptRunner(confing *SSHRunerConfig) (string, error) {
	if confing.SetSecretsToScript != nil && *confing.SetSecretsToScript {
		envCommand := r.loadEnvToScript(confing)
		confing.Script = envCommand + "\n" + confing.Script
	}
	command := r.createCommand(confing, confing.Script)
	return command, nil
}

// escapeShellArg properly escapes shell arguments for safe execution
func (r *SSHRuner) escapeShellArg(arg string) string {
	// Use single quotes and escape any single quotes within the argument
	escaped := strings.ReplaceAll(arg, "'", "'\"'\"'")
	return "'" + escaped + "'"
}

func (r *SSHRuner) createCommand(confing *SSHRunerConfig, command string) string {
	// Escape password and command to prevent shell injection
	escapedPassword := r.escapeShellArg(confing.Password)
	escapedCommand := r.escapeShellArg(command)

	if confing.SSHKey != nil {
		escapedSSHKey := r.escapeShellArg(*confing.SSHKey)
		return fmt.Sprintf("sshpass -p %s ssh -o StrictHostKeyChecking=no -o ConnectTimeout=30 -o ServerAliveInterval=60 -o ServerAliveCountMax=3 -i %s %s@%s %s",
			escapedPassword, escapedSSHKey, confing.User, confing.IP, escapedCommand)
	}
	return fmt.Sprintf("sshpass -p %s ssh -o StrictHostKeyChecking=no -o ConnectTimeout=30 -o ServerAliveInterval=60 -o ServerAliveCountMax=3 %s@%s %s",
		escapedPassword, confing.User, confing.IP, escapedCommand)
}

func (r *SSHRuner) loginDockerCommand(confing *SSHRunerConfig) string {
	command := ""
	if confing.DockerUser != nil && confing.DockerPassword != nil {
		command += fmt.Sprintf("docker login -u %s -p %s", *confing.DockerUser, *confing.DockerPassword)
	}
	execCommand := r.createCommand(confing, command)
	return execCommand
}

func (r *SSHRuner) pullDockerCommand(confing *SSHRunerConfig) string {
	command := ""
	if confing.DockerImage != nil && confing.DockerRegistry != nil && confing.DockerTag != nil {
		command += fmt.Sprintf("docker pull %s/%s:%s", *confing.DockerRegistry, *confing.DockerImage, *confing.DockerTag)
	}
	execCommand := r.createCommand(confing, command)
	return execCommand
}

func (r *SSHRuner) runDockerCommand(confing *SSHRunerConfig) string {
	command := ""
	envCommand := ""
	if confing.SetSecretsToContainer != nil && *confing.SetSecretsToContainer {
		envCommand = r.createEnvContainerCommand(confing)
	}
	if confing.DockerImage != nil && confing.DockerRegistry != nil && confing.DockerTag != nil {
		command += fmt.Sprintf("docker run -d --name %s %s --network gateway_network --ip 172.30.0.20 %s/%s:%s",
			*confing.DockerContainerName, envCommand, *confing.DockerRegistry, *confing.DockerImage, *confing.DockerTag)
	}
	execCommand := r.createCommand(confing, command)
	return execCommand
}

func (r *SSHRuner) createEnvContainerCommand(confing *SSHRunerConfig) string {
	command := ""
	if confing.Env != nil {
		for key, value := range *confing.Env {
			escapedValue := r.escapeShellArg(value)
			command += fmt.Sprintf("-e %s=%s ", key, escapedValue)
		}
	}
	return command
}

func (r *SSHRuner) loadEnvToScript(confing *SSHRunerConfig) string {
	command := ""
	if confing.Env != nil {
		for key, value := range *confing.Env {
			escapedValue := r.escapeShellArg(value)
			command += fmt.Sprintf("export %s=%s\n", key, escapedValue)
		}
	}
	return command
}
