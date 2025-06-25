# Используем Alpine Linux для минимального размера
FROM alpine:3.19

# Устанавливаем необходимые пакеты
RUN apk add --no-cache \
    bash \
    curl \
    wget \
    git \
    openssh-client \
    rsync \
    docker-cli \
    python3 \
    py3-pip \
    nodejs \
    npm \
    jq \
    unzip \
    tar \
    gzip \
    ca-certificates \
    && rm -rf /var/cache/apk/*

# Создаем рабочие директории
RUN mkdir -p /workspace /scripts /secrets

# Создаем пользователя для безопасности
RUN addgroup -g 1000 deployer && \
    adduser -D -s /bin/bash -u 1000 -G deployer deployer

# Устанавливаем права доступа
RUN chown -R deployer:deployer /workspace /scripts /secrets

# Создаем скрипт для очистки workspace
COPY cleanup.sh /usr/local/bin/cleanup.sh
RUN chmod +x /usr/local/bin/cleanup.sh

# Переключаемся на пользователя deployer
USER deployer

# Устанавливаем рабочую директорию
WORKDIR /workspace

# Команда по умолчанию - ожидание команд
CMD ["tail", "-f", "/dev/null"]