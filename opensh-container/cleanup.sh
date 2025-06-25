#!/bin/bash
# cleanup.sh - Скрипт для очистки workspace после deployment

echo "Starting cleanup process..."

# Очищаем workspace, оставляя только базовые директории
find /workspace -mindepth 1 -maxdepth 1 -exec rm -rf {} +

# Очищаем временные файлы
rm -rf /tmp/*

# Очищаем SSH ключи и конфигурации
rm -rf ~/.ssh/known_hosts ~/.ssh/config

# Очищаем git конфигурации
git config --global --unset-all user.name 2>/dev/null || true
git config --global --unset-all user.email 2>/dev/null || true

# Очищаем переменные окружения (кроме системных)
unset $(env | grep -E '^[A-Z_]+=' | grep -v -E '^(PATH|HOME|USER|SHELL|PWD|TERM)=' | cut -d= -f1)

echo "Cleanup completed successfully"