#!/bin/bash

# manage-containers.sh - Скрипт для управления множественными deployment контейнерами
# с учетом ресурсов сервера

set -e

# Конфигурация
CONTAINER_NAME_PREFIX="deploy-worker"
IMAGE_NAME="deployment-worker"
DOCKERFILE_PATH="./Dockerfile"
NETWORK_NAME="deployment-network"

# Ресурсы (можно настроить)
RESERVED_CPU_CORES=1          # Резерв для основного приложения
RESERVED_MEMORY_GB=2          # Резерв памяти для основного приложения
CONTAINER_CPU_LIMIT="0.5"     # CPU лимит на один контейнер
CONTAINER_MEMORY_LIMIT="512m" # Память на один контейнер
MIN_CONTAINERS=1              # Минимальное количество контейнеров
MAX_CONTAINERS=10             # Максимальное количество контейнеров

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Функции логирования
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Проверка зависимостей
check_dependencies() {
    log_info "Проверка зависимостей..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker не установлен"
        exit 1
    fi
    
    if ! docker info &> /dev/null; then
        log_error "Docker daemon не запущен"
        exit 1
    fi
    
    log_success "Все зависимости в порядке"
}

# Получение информации о ресурсах системы
get_system_resources() {
    log_info "Анализ ресурсов системы..."
    
    # Получаем количество CPU ядер
    TOTAL_CPU_CORES=$(nproc)
    
    # Получаем общую память в GB
    TOTAL_MEMORY_KB=$(grep MemTotal /proc/meminfo | awk '{print $2}')
    TOTAL_MEMORY_GB=$((TOTAL_MEMORY_KB / 1024 / 1024))
    
    # Вычисляем доступные ресурсы для контейнеров
    AVAILABLE_CPU_CORES=$((TOTAL_CPU_CORES - RESERVED_CPU_CORES))
    AVAILABLE_MEMORY_GB=$((TOTAL_MEMORY_GB - RESERVED_MEMORY_GB))
    
    log_info "Всего CPU ядер: $TOTAL_CPU_CORES"
    log_info "Всего памяти: ${TOTAL_MEMORY_GB}GB"
    log_info "Доступно CPU ядер для контейнеров: $AVAILABLE_CPU_CORES"
    log_info "Доступно памяти для контейнеров: ${AVAILABLE_MEMORY_GB}GB"
    
    # Проверка минимальных требований
    if [ $AVAILABLE_CPU_CORES -lt 1 ] || [ $AVAILABLE_MEMORY_GB -lt 1 ]; then
        log_error "Недостаточно ресурсов для запуска контейнеров"
        exit 1
    fi
}

# Вычисление оптимального количества контейнеров
calculate_optimal_containers() {
    log_info "Вычисление оптимального количества контейнеров..."
    
    # Вычисляем максимальное количество контейнеров по CPU
    # Используем awk для работы с дробными числами
    CPU_CONTAINER_LIMIT=$(awk "BEGIN {printf \"%.0f\", $AVAILABLE_CPU_CORES / $CONTAINER_CPU_LIMIT}")
    
    # Вычисляем максимальное количество контейнеров по памяти
    CONTAINER_MEMORY_MB=$(echo $CONTAINER_MEMORY_LIMIT | sed 's/[^0-9]//g')
    AVAILABLE_MEMORY_MB=$((AVAILABLE_MEMORY_GB * 1024))
    MEMORY_CONTAINER_LIMIT=$((AVAILABLE_MEMORY_MB / CONTAINER_MEMORY_MB))
    
    # Берем минимальное значение
    if [ $CPU_CONTAINER_LIMIT -lt $MEMORY_CONTAINER_LIMIT ]; then
        OPTIMAL_CONTAINERS=$CPU_CONTAINER_LIMIT
    else
        OPTIMAL_CONTAINERS=$MEMORY_CONTAINER_LIMIT
    fi
    
    # Применяем ограничения
    if [ $OPTIMAL_CONTAINERS -lt $MIN_CONTAINERS ]; then
        OPTIMAL_CONTAINERS=$MIN_CONTAINERS
    fi
    
    if [ $OPTIMAL_CONTAINERS -gt $MAX_CONTAINERS ]; then
        OPTIMAL_CONTAINERS=$MAX_CONTAINERS
    fi
    
    log_info "Ограничение по CPU: $CPU_CONTAINER_LIMIT контейнеров"
    log_info "Ограничение по памяти: $MEMORY_CONTAINER_LIMIT контейнеров"
    log_success "Оптимальное количество контейнеров: $OPTIMAL_CONTAINERS"
}

# Проверка существующих контейнеров
check_existing_containers() {
    log_info "Проверка существующих контейнеров..."
    
    EXISTING_CONTAINERS=$(docker ps -a --filter "name=${CONTAINER_NAME_PREFIX}" --format "{{.Names}}" | wc -l)
    RUNNING_CONTAINERS=$(docker ps --filter "name=${CONTAINER_NAME_PREFIX}" --format "{{.Names}}" | wc -l)
    
    log_info "Существующих контейнеров: $EXISTING_CONTAINERS"
    log_info "Запущенных контейнеров: $RUNNING_CONTAINERS"
    
    if [ $EXISTING_CONTAINERS -gt 0 ]; then
        log_info "Список существующих контейнеров:"
        docker ps -a --filter "name=${CONTAINER_NAME_PREFIX}" --format "table {{.Names}}\t{{.Status}}\t{{.CreatedAt}}"
    fi
}

# Создание сети Docker (если не существует)
create_network() {
    if ! docker network ls | grep -q $NETWORK_NAME; then
        log_info "Создание Docker сети: $NETWORK_NAME"
        docker network create $NETWORK_NAME
        log_success "Сеть $NETWORK_NAME создана"
    else
        log_info "Сеть $NETWORK_NAME уже существует"
    fi
}

# Сборка Docker образа
build_image() {
    log_info "Сборка Docker образа..."
    
    if [ ! -f "$DOCKERFILE_PATH" ]; then
        log_error "Dockerfile не найден: $DOCKERFILE_PATH"
        exit 1
    fi
    
    docker build -t $IMAGE_NAME .
    log_success "Образ $IMAGE_NAME собран успешно"
}

# Запуск контейнера
start_container() {
    local container_id=$1
    local container_name="${CONTAINER_NAME_PREFIX}-${container_id}"
    
    log_info "Запуск контейнера: $container_name"
    
    docker run -d \
        --name "$container_name" \
        --network "$NETWORK_NAME" \
        --cpus="$CONTAINER_CPU_LIMIT" \
        --memory="$CONTAINER_MEMORY_LIMIT" \
        --restart=unless-stopped \
        -v /var/run/docker.sock:/var/run/docker.sock \
        -v "${PWD}/workspace-${container_id}:/workspace" \
        -e "CONTAINER_ID=${container_id}" \
        "$IMAGE_NAME"
    
    log_success "Контейнер $container_name запущен"
}

# Остановка и удаление контейнера
stop_container() {
    local container_name=$1
    
    log_info "Остановка контейнера: $container_name"
    
    if docker ps -q --filter "name=$container_name" | grep -q .; then
        docker stop "$container_name"
    fi
    
    if docker ps -aq --filter "name=$container_name" | grep -q .; then
        docker rm "$container_name"
    fi
    
    log_success "Контейнер $container_name остановлен и удален"
}

# Масштабирование контейнеров
scale_containers() {
    local target_count=$1
    
    log_info "Масштабирование до $target_count контейнеров..."
    
    # Получаем список запущенных контейнеров
    local running_containers=($(docker ps --filter "name=${CONTAINER_NAME_PREFIX}" --format "{{.Names}}" | sort))
    local current_count=${#running_containers[@]}
    
    if [ $target_count -gt $current_count ]; then
        # Добавляем контейнеры
        local containers_to_add=$((target_count - current_count))
        log_info "Добавление $containers_to_add контейнеров"
        
        for ((i=current_count+1; i<=target_count; i++)); do
            start_container $i
            sleep 2  # Небольшая пауза между запусками
        done
        
    elif [ $target_count -lt $current_count ]; then
        # Удаляем лишние контейнеры
        local containers_to_remove=$((current_count - target_count))
        log_info "Удаление $containers_to_remove контейнеров"
        
        for ((i=0; i<containers_to_remove; i++)); do
            local container_name=${running_containers[$((current_count-1-i))]}
            stop_container "$container_name"
        done
    else
        log_info "Количество контейнеров уже соответствует целевому"
    fi
}

# Очистка ресурсов
cleanup() {
    log_info "Очистка неиспользуемых ресурсов..."
    
    # Удаляем остановленные контейнеры
    docker container prune -f
    
    # Удаляем неиспользуемые образы
    docker image prune -f
    
    log_success "Очистка завершена"
}

# Мониторинг контейнеров
monitor() {
    log_info "Статус контейнеров:"
    docker ps --filter "name=${CONTAINER_NAME_PREFIX}" --format "table {{.Names}}\t{{.Status}}\t{{.CPU}}\t{{.MemUsage}}"
    
    echo ""
    log_info "Использование ресурсов системы:"
    echo "CPU: $(top -bn1 | grep "Cpu(s)" | awk '{print $2}' | awk -F'%' '{print $1}')% использовано"
    echo "Память: $(free -h | awk '/^Mem:/ {print $3 "/" $2}')"
    echo "Диск: $(df -h / | awk 'NR==2 {print $3 "/" $2 " (" $5 " использовано)"}')"
}

# Главная функция
main() {
    local action=${1:-"auto"}
    local count=${2:-"$OPTIMAL_CONTAINERS"}
    
    check_dependencies
    get_system_resources
    calculate_optimal_containers
    check_existing_containers
    create_network
    
    case $action in
        "build")
            build_image
            ;;
        "start")
            build_image
            scale_containers $count
            ;;
        "stop")
            log_info "Остановка всех контейнеров..."
            for container in $(docker ps --filter "name=${CONTAINER_NAME_PREFIX}" --format "{{.Names}}"); do
                stop_container "$container"
            done
            ;;
        "restart")
            main stop
            sleep 5
            main start $count
            ;;
        "scale")
            if [ -z "$count" ] || ! [[ "$count" =~ ^[0-9]+$ ]]; then
                log_error "Укажите корректное количество контейнеров"
                exit 1
            fi
            scale_containers $count
            ;;
        "monitor")
            monitor
            ;;
        "cleanup")
            cleanup
            ;;
        "auto")
            build_image
            scale_containers $OPTIMAL_CONTAINERS
            ;;
        "help"|"-h"|"--help")
            echo "Использование: $0 [ДЕЙСТВИЕ] [КОЛИЧЕСТВО]"
            echo ""
            echo "ДЕЙСТВИЯ:"
            echo "  auto      - Автоматическое масштабирование (по умолчанию)"
            echo "  build     - Сборка образа"
            echo "  start     - Запуск контейнеров"
            echo "  stop      - Остановка всех контейнеров"
            echo "  restart   - Перезапуск контейнеров"
            echo "  scale     - Масштабирование до указанного количества"
            echo "  monitor   - Мониторинг статуса"
            echo "  cleanup   - Очистка ресурсов"
            echo "  help      - Показать эту справку"
            echo ""
            echo "ПРИМЕРЫ:"
            echo "  $0                    # Автоматическое масштабирование"
            echo "  $0 start 5           # Запуск 5 контейнеров"
            echo "  $0 scale 3           # Масштабирование до 3 контейнеров"
            echo "  $0 monitor           # Мониторинг статуса"
            ;;
        *)
            log_error "Неизвестное действие: $action"
            main help
            exit 1
            ;;
    esac
    
    if [ "$action" != "help" ] && [ "$action" != "monitor" ]; then
        echo ""
        monitor
    fi
}

# Обработка сигналов
trap cleanup EXIT

# Запуск
main "$@" 