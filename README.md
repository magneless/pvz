# Установка


## Требования
- Go 1.21+
- PostgreSQL 12+
- Утилита `migrate` для работы с миграциями
- Утилита `yq` для парсинга YAML конфигов

- Склонируйте
```bash
git clone https://github.com/magneless/pvz
cd pvz
```


- Установите зависимости
```bash
go mod download
```

- Создайте файл .env в корне проекта:

```bash
touch .env
Заполните .env файл (пример):
```

```bash
DB_PASSWORD=your_db_password
CONFIG_PATH=./configs/local.yaml
```
## Настройка базы данных
- Запустите PostgreSQL сервер

- Измените параметры configs/local.yaml

- Выполните миграции:

```bash
chmod +x scripts/scripts.sh
./scripts/scripts.sh up
```
## Запуск сервера
- Запустите сервер:

```bash
go run cmd/main/main.go
```