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
		return err
	}
	server, err := s.ServersService.GetServer(serverId, userId, iv)
	if err != nil {
		return err
	}
	var envMap map[string]string
	if loadEnv {
		env, err := s.EnvsService.GetSecret(envId, userId, iv)
		if err != nil {
			return err
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
	command, err := s.SSHRuner.CreateScriptRunner(&config)
	if err != nil {
		return err
	}
	go func() {
		rs, err := s.Docker.ExecuteCommand(context.Background(), worker.ID, command)
		if err != nil {
			fmt.Println("Error executing command:", err)
		}
		fmt.Println("Command executed successfully:", rs)
	}()
	return nil
}

func (s *ExecuteService) getWorker() *libs.Container {
	workerList := s.Docker.GetCachedDeploymentWorkers()
	randomWorker := workerList[rand.Intn(len(workerList))]
	return &randomWorker
}
