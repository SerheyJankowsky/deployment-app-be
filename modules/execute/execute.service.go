package execute

import (
	"context"
	"fmt"
	"math/rand"

	"deployer.com/libs"
	"deployer.com/modules/containers"
	"deployer.com/modules/deployments"
	"deployer.com/modules/projects"
	"deployer.com/modules/scripts"
	"deployer.com/modules/secrets"
	"deployer.com/modules/servers"
)

type ExecuteService struct {
	Docker             *libs.DockerComunication
	SSHRuner           *libs.SSHRuner
	EncryptionService  *libs.EncryptionService
	ScriptsService     *scripts.ScriptsService
	ServersService     *servers.ServersService
	EnvsService        *secrets.SecretsService
	ContainersService  *containers.ContainersService
	DeploymentsService *deployments.DeploymentsService
	ProjectsService    *projects.ProjectsService
}

func NewExecuteService(scriptsService *scripts.ScriptsService,
	serversService *servers.ServersService,
	envsService *secrets.SecretsService,
	containersService *containers.ContainersService,
	deploymentsService *deployments.DeploymentsService,
	projectsService *projects.ProjectsService,
	docker *libs.DockerComunication,
) *ExecuteService {
	sshRuner := libs.NewSSHRuner()
	encryptionService := libs.NewEncryptionService()

	return &ExecuteService{
		ScriptsService:     scriptsService,
		ServersService:     serversService,
		EnvsService:        envsService,
		ContainersService:  containersService,
		DeploymentsService: deploymentsService,
		ProjectsService:    projectsService,
		Docker:             docker,
		SSHRuner:           sshRuner,
		EncryptionService:  encryptionService,
	}
}

func (s *ExecuteService) RunScript(id, userId, serverId, envId uint, iv string, loadEnv bool) error {
	script, err := s.ScriptsService.GetScript(id, userId, iv)
	if err != nil {
		return fmt.Errorf("failed to get script: %w", err)
	}

	server, err := s.ServersService.GetServer(serverId, userId, iv)
	if err != nil {
		return fmt.Errorf("failed to get server: %w", err)
	}

	var envMap map[string]string
	if loadEnv {
		env, err := s.EnvsService.GetSecret(envId, userId, iv)
		if err != nil {
			return fmt.Errorf("failed to get secrets: %w", err)
		}
		envMap = s.EnvsService.GetEnvMap(env)
	}

	config := libs.SSHRunerConfig{
		IP:       server.Host,
		User:     server.Username,
		Password: server.Password,
		Script:   script.Script,
		Env:      &envMap,
	}

	worker := s.getWorker()
	if worker == nil {
		return fmt.Errorf("no deployment worker containers available")
	}

	command, err := s.SSHRuner.CreateScriptRunner(&config)
	if err != nil {
		return fmt.Errorf("failed to create script runner: %w", err)
	}

	// Log the execution details for debugging (without sensitive data)
	fmt.Printf("DEBUG: Executing script on server %s@%s using worker %s\n",
		server.Username, server.Host, worker.ID)

	go func() {
		// Test SSH connectivity first
		testCommand := "echo 'SSH connection test successful'"
		testSSHCommand := fmt.Sprintf("sshpass -p '%s' ssh -o StrictHostKeyChecking=no -o ConnectTimeout=10 %s@%s %s",
			server.Password, server.Username, server.Host, testCommand)

		fmt.Printf("DEBUG: Testing SSH connectivity to %s@%s...\n", server.Username, server.Host)
		testResult, testErr := s.Docker.ExecuteCommand(context.Background(), worker.ID, testSSHCommand)
		if testErr != nil {
			fmt.Printf("ERROR: SSH connectivity test failed: %v\n", testErr)
			fmt.Printf("DEBUG: Test command output: %s\n", testResult)
			return
		}

		if testResult != "" && testResult != "SSH connection test successful\n" {
			fmt.Printf("WARNING: SSH test returned unexpected output: %s\n", testResult)
		} else {
			fmt.Printf("DEBUG: SSH connectivity test passed\n")
		}

		// Execute the actual script
		fmt.Printf("DEBUG: Executing actual script...\n")
		rs, err := s.Docker.ExecuteCommand(context.Background(), worker.ID, command)
		if err != nil {
			fmt.Printf("ERROR: Failed to execute command in container: %v\n", err)
			return
		}

		// Check if the output contains permission denied errors
		if rs != "" {
			fmt.Printf("Command executed. Output: %s\n", rs)
			if rs == "Permission denied, please try again.\n" ||
				rs == "&Permission denied, please try again.\n" {
				fmt.Printf("ERROR: SSH authentication failed. Please check:\n")
				fmt.Printf("  1. Username: %s\n", server.Username)
				fmt.Printf("  2. Server IP: %s\n", server.Host)
				fmt.Printf("  3. Password is correct\n")
				fmt.Printf("  4. SSH server allows password authentication\n")
				fmt.Printf("  5. Network connectivity from container to server\n")
			}
		} else {
			fmt.Printf("Command executed successfully with no output\n")
		}
	}()

	return nil
}

func (s *ExecuteService) getWorker() *libs.Container {
	workerList := s.Docker.GetCachedDeploymentWorkers()
	if len(workerList) == 0 {
		return nil
	}
	randomWorker := workerList[rand.Intn(len(workerList))]
	return &randomWorker
}
