package libs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type Container struct {
	ID        string
	Name      string
	Status    string
	CreatedAt string
}

type ContainerStats struct {
	CPUPercent    float64
	MemoryUsage   uint64
	MemoryLimit   uint64
	MemoryPercent float64
	NetworkRx     uint64
	NetworkTx     uint64
}

type DockerComunication struct {
	client                     *client.Client
	activeContainers           []Container
	deploymentWorkerContainers []Container
	deploymentWorkerMutex      sync.RWMutex
	lastCacheUpdate            time.Time
	cacheExpiration            time.Duration
}

// NewDockerCommunication создает новый экземпляр Docker клиента
func NewDockerCommunication() (*DockerComunication, error) {
	fmt.Println("DEBUG: Creating Docker client...")

	// Выводим информацию о Docker окружении
	fmt.Printf("DEBUG: DOCKER_HOST=%s\n", os.Getenv("DOCKER_HOST"))
	fmt.Printf("DEBUG: DOCKER_API_VERSION=%s\n", os.Getenv("DOCKER_API_VERSION"))
	fmt.Printf("DEBUG: DOCKER_CERT_PATH=%s\n", os.Getenv("DOCKER_CERT_PATH"))
	fmt.Printf("DEBUG: DOCKER_TLS_VERIFY=%s\n", os.Getenv("DOCKER_TLS_VERIFY"))

	// Сначала пробуем стандартное подключение
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Printf("DEBUG: Failed to create Docker client with FromEnv: %v\n", err)

		// Пробуем подключиться через Unix socket (для контейнеров)
		fmt.Println("DEBUG: Trying to connect via Unix socket...")
		cli, err = client.NewClientWithOpts(
			client.WithHost("unix:///var/run/docker.sock"),
			client.WithAPIVersionNegotiation(),
		)
		if err != nil {
			fmt.Printf("DEBUG: Failed to create Docker client with Unix socket: %v\n", err)
			return nil, fmt.Errorf("failed to create Docker client: %w", err)
		}
	}

	fmt.Println("DEBUG: Docker client created successfully")

	// Проверяем подключение
	ctx := context.Background()
	version, err := cli.ServerVersion(ctx)
	if err != nil {
		fmt.Printf("DEBUG: Failed to get Docker server version: %v\n", err)
		return nil, fmt.Errorf("failed to connect to Docker daemon: %w", err)
	}

	fmt.Printf("DEBUG: Connected to Docker daemon - Version: %s, API Version: %s\n",
		version.Version, version.APIVersion)

	return &DockerComunication{
		client:          cli,
		cacheExpiration: 30 * time.Second, // кэш действителен 30 секунд
	}, nil
}

// SetCacheExpiration устанавливает время жизни кэша
func (dc *DockerComunication) SetCacheExpiration(duration time.Duration) {
	dc.deploymentWorkerMutex.Lock()
	defer dc.deploymentWorkerMutex.Unlock()
	dc.cacheExpiration = duration
}

// Close закрывает соединение с Docker
func (dc *DockerComunication) Close() error {
	return dc.client.Close()
}

// ListContainers получает список всех контейнеров
func (dc *DockerComunication) ListContainers(ctx context.Context, all bool) ([]Container, error) {
	containers, err := dc.client.ContainerList(ctx, container.ListOptions{All: all})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	var result []Container
	for _, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = c.Names[0]
		}

		result = append(result, Container{
			ID:        c.ID[:12], // короткий ID
			Name:      name,
			Status:    c.Status,
			CreatedAt: time.Unix(c.Created, 0).Format(time.RFC3339),
		})
	}

	dc.activeContainers = result
	return result, nil
}

// isCacheExpired проверяет, истек ли кэш
func (dc *DockerComunication) isCacheExpired() bool {
	return time.Since(dc.lastCacheUpdate) > dc.cacheExpiration
}

// ListDeploymentWorkerContainers получает список контейнеров с именем deployment-worker
func (dc *DockerComunication) ListDeploymentWorkerContainers(ctx context.Context, all bool) ([]Container, error) {
	fmt.Printf("DEBUG: Starting to list containers (all=%v)\n", all)

	containers, err := dc.client.ContainerList(ctx, container.ListOptions{All: all})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	fmt.Printf("DEBUG: Found %d total containers\n", len(containers))

	var result []Container
	for i, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/") // убираем префикс "/"
		}

		fmt.Printf("DEBUG: Container [%d]: ID=%s, Name='%s', Image=%s, Status=%s\n",
			i+1, c.ID[:12], name, c.Image, c.Status)

		// Фильтруем только контейнеры с именем deployment-worker
		if name == "deployment-worker" || strings.Contains(name, "deployment-worker") {
			fmt.Printf("DEBUG: ✓ Container '%s' matches deployment-worker filter\n", name)
			result = append(result, Container{
				ID:        c.ID[:12], // короткий ID
				Name:      name,
				Status:    c.Status,
				CreatedAt: time.Unix(c.Created, 0).Format(time.RFC3339),
			})
		} else {
			fmt.Printf("DEBUG: ✗ Container '%s' does not match deployment-worker filter\n", name)
		}
	}

	fmt.Printf("DEBUG: Found %d deployment-worker containers\n", len(result))

	// Обновляем кэш
	dc.updateDeploymentWorkerCache(result)

	return result, nil
}

// updateDeploymentWorkerCache обновляет кэш deployment-worker контейнеров
func (dc *DockerComunication) updateDeploymentWorkerCache(containers []Container) {
	dc.deploymentWorkerMutex.Lock()
	defer dc.deploymentWorkerMutex.Unlock()

	dc.deploymentWorkerContainers = containers
	dc.lastCacheUpdate = time.Now()
}

// GetCachedDeploymentWorkers возвращает кэшированные deployment-worker контейнеры
func (dc *DockerComunication) GetCachedDeploymentWorkers() []Container {
	dc.deploymentWorkerMutex.RLock()
	defer dc.deploymentWorkerMutex.RUnlock()

	// Возвращаем копию, чтобы избежать изменений извне
	result := make([]Container, len(dc.deploymentWorkerContainers))
	copy(result, dc.deploymentWorkerContainers)

	return result
}

// GetDeploymentWorkersWithCache возвращает deployment-worker контейнеры из кэша или обновляет кэш
func (dc *DockerComunication) GetDeploymentWorkersWithCache(ctx context.Context, forceRefresh bool) ([]Container, error) {
	// Если кэш не истек и не требуется принудительное обновление, возвращаем из кэша
	if !forceRefresh && !dc.isCacheExpired() {
		cached := dc.GetCachedDeploymentWorkers()
		if len(cached) > 0 {
			return cached, nil
		}
	}

	// Обновляем кэш
	return dc.ListDeploymentWorkerContainers(ctx, false)
}

// RefreshDeploymentWorkersCache принудительно обновляет кэш
func (dc *DockerComunication) RefreshDeploymentWorkersCache(ctx context.Context) error {
	_, err := dc.ListDeploymentWorkerContainers(ctx, false)
	return err
}

// GetDeploymentWorkerByName находит deployment-worker контейнер по имени в кэше
func (dc *DockerComunication) GetDeploymentWorkerByName(name string) *Container {
	dc.deploymentWorkerMutex.RLock()
	defer dc.deploymentWorkerMutex.RUnlock()

	for _, container := range dc.deploymentWorkerContainers {
		if container.Name == name || strings.Contains(container.Name, name) {
			return &container
		}
	}

	return nil
}

// GetDeploymentWorkerByID находит deployment-worker контейнер по ID в кэше
func (dc *DockerComunication) GetDeploymentWorkerByID(id string) *Container {
	dc.deploymentWorkerMutex.RLock()
	defer dc.deploymentWorkerMutex.RUnlock()

	for _, container := range dc.deploymentWorkerContainers {
		if container.ID == id || strings.HasPrefix(container.ID, id) {
			return &container
		}
	}

	return nil
}

// GetCacheInfo возвращает информацию о состоянии кэша
func (dc *DockerComunication) GetCacheInfo() map[string]interface{} {
	dc.deploymentWorkerMutex.RLock()
	defer dc.deploymentWorkerMutex.RUnlock()

	return map[string]interface{}{
		"cached_containers_count": len(dc.deploymentWorkerContainers),
		"last_update":             dc.lastCacheUpdate.Format(time.RFC3339),
		"cache_expiration":        dc.cacheExpiration.String(),
		"is_expired":              dc.isCacheExpired(),
		"time_until_expiry":       dc.cacheExpiration - time.Since(dc.lastCacheUpdate),
	}
}

// StartAutoRefresh запускает автоматическое обновление кэша
func (dc *DockerComunication) StartAutoRefresh(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := dc.RefreshDeploymentWorkersCache(ctx); err != nil {
					fmt.Printf("Error refreshing deployment-worker cache: %v\n", err)
				}
			}
		}
	}()
}

// GetContainerStats получает статистики контейнера в реальном времени
func (dc *DockerComunication) GetContainerStats(ctx context.Context, containerID string) (*ContainerStats, error) {
	stats, err := dc.client.ContainerStats(ctx, containerID, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get container stats: %w", err)
	}
	defer stats.Body.Close()

	// Читаем JSON из потока
	var v container.StatsResponse
	decoder := json.NewDecoder(stats.Body)
	if err := decoder.Decode(&v); err != nil {
		return nil, fmt.Errorf("failed to decode stats: %w", err)
	}

	// Вычисляем CPU процент
	cpuDelta := float64(v.CPUStats.CPUUsage.TotalUsage - v.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(v.CPUStats.SystemUsage - v.PreCPUStats.SystemUsage)
	cpuPercent := 0.0
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}

	// Вычисляем использование памяти
	memoryUsage := v.MemoryStats.Usage
	memoryLimit := v.MemoryStats.Limit
	memoryPercent := float64(memoryUsage) / float64(memoryLimit) * 100.0

	// Получаем сетевую статистику
	var networkRx, networkTx uint64
	for _, network := range v.Networks {
		networkRx += network.RxBytes
		networkTx += network.TxBytes
	}

	return &ContainerStats{
		CPUPercent:    cpuPercent,
		MemoryUsage:   memoryUsage,
		MemoryLimit:   memoryLimit,
		MemoryPercent: memoryPercent,
		NetworkRx:     networkRx,
		NetworkTx:     networkTx,
	}, nil
}

// ExecuteCommand выполняет команду в контейнере
func (dc *DockerComunication) ExecuteCommand(ctx context.Context, containerID string, cmd []string) (string, error) {
	execConfig := container.ExecOptions{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	}

	execResp, err := dc.client.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create exec: %w", err)
	}

	attachResp, err := dc.client.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to attach to exec: %w", err)
	}
	defer attachResp.Close()

	// Читаем вывод команды
	output, err := io.ReadAll(attachResp.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to read exec output: %w", err)
	}

	return string(output), nil
}

// GetContainerLogs получает логи контейнера
func (dc *DockerComunication) GetContainerLogs(ctx context.Context, containerID string, tail string) (string, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
		Timestamps: true,
	}

	logs, err := dc.client.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer logs.Close()

	// Читаем логи
	output, err := io.ReadAll(logs)
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return string(output), nil
}

// StartContainer запускает контейнер
func (dc *DockerComunication) StartContainer(ctx context.Context, containerID string) error {
	err := dc.client.ContainerStart(ctx, containerID, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}
	return nil
}

// StopContainer останавливает контейнер
func (dc *DockerComunication) StopContainer(ctx context.Context, containerID string, timeoutSeconds *int) error {
	err := dc.client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: timeoutSeconds})
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}
	return nil
}

// RestartContainer перезапускает контейнер
func (dc *DockerComunication) RestartContainer(ctx context.Context, containerID string, timeoutSeconds *int) error {
	err := dc.client.ContainerRestart(ctx, containerID, container.StopOptions{Timeout: timeoutSeconds})
	if err != nil {
		return fmt.Errorf("failed to restart container: %w", err)
	}
	return nil
}

// CreateAndRunContainer создает и запускает новый контейнер
func (dc *DockerComunication) CreateAndRunContainer(ctx context.Context, imageName, containerName string, cmd []string) (string, error) {
	// Создаем контейнер
	resp, err := dc.client.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   cmd,
	}, nil, nil, nil, containerName)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	// Запускаем контейнер
	if err := dc.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	return resp.ID, nil
}

// PullImage скачивает образ
func (dc *DockerComunication) PullImage(ctx context.Context, imageName string) error {
	reader, err := dc.client.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	defer reader.Close()

	// Читаем вывод процесса скачивания
	_, err = io.Copy(io.Discard, reader)
	if err != nil {
		return fmt.Errorf("failed to read pull output: %w", err)
	}

	return nil
}

// MonitorContainer запускает мониторинг контейнера в горутине
func (dc *DockerComunication) MonitorContainer(ctx context.Context, containerID string, interval time.Duration, callback func(*ContainerStats)) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				stats, err := dc.GetContainerStats(ctx, containerID)
				if err != nil {
					fmt.Printf("Error getting stats for container %s: %v\n", containerID, err)
					continue
				}
				callback(stats)
			}
		}
	}()
}

// MonitorDeploymentWorkers запускает мониторинг всех deployment-worker контейнеров используя кэш
func (dc *DockerComunication) MonitorDeploymentWorkers(ctx context.Context, interval time.Duration, callback func(string, *ContainerStats)) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// Используем кэш для получения списка контейнеров
				containers, err := dc.GetDeploymentWorkersWithCache(ctx, false)
				if err != nil {
					fmt.Printf("Error getting deployment-worker containers: %v\n", err)
					continue
				}

				for _, cont := range containers {
					stats, err := dc.GetContainerStats(ctx, cont.ID)
					if err != nil {
						fmt.Printf("Error getting stats for container %s: %v\n", cont.ID, err)
						continue
					}
					callback(cont.Name, stats)
				}
			}
		}
	}()
}
