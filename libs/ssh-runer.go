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
	script := confing.Script
	if confing.SetSecretsToScript != nil && *confing.SetSecretsToScript {
		envCommand := r.loadEnvToScript(confing)
		script = envCommand + "\n" + script
	}

	// Since DockerComunication.ExecuteCommand already wraps with bash -c,
	// we just need to properly escape the script for SSH execution
	if strings.Contains(script, "\n") || strings.Contains(script, "'") || strings.Contains(script, "\"") {
		// For complex scripts with special characters, escape them properly
		escapedScript := strings.ReplaceAll(script, "'", "'\"'\"'")
		command := r.createCommand(confing, fmt.Sprintf("'%s'", escapedScript))
		return command, nil
	}

	// For simple commands, pass directly
	command := r.createCommand(confing, script)
	return command, nil
}

func (r *SSHRuner) createCommand(confing *SSHRunerConfig, command string) string {
	// Simple command generation without excessive escaping - matching working manual command
	if confing.SSHKey != nil {
		return fmt.Sprintf("sshpass -p %s ssh -o StrictHostKeyChecking=no -i %s %s@%s %s",
			confing.Password, *confing.SSHKey, confing.User, confing.IP, command)
	}
	return fmt.Sprintf("sshpass -p %s ssh -o StrictHostKeyChecking=no %s@%s %s",
		confing.Password, confing.User, confing.IP, command)
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
			command += fmt.Sprintf("-e %s=%s ", key, value)
		}
	}
	return command
}

func (r *SSHRuner) loadEnvToScript(confing *SSHRunerConfig) string {
	command := ""
	if confing.Env != nil {
		for key, value := range *confing.Env {
			command += fmt.Sprintf("export %s=%s\n", key, value)
		}
	}
	return command
}
