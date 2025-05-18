# Тестовое задание в Effective Mobile

[![testing](https://github.com/P1xart/effective_mobile_service/actions/workflows/test.yaml/badge.svg)](https://github.com/P1xart/effective_mobile_service/actions/workflows/test.yaml)

## Запуск
Если сервис запускается впервые, нужно запустить миграции
```
docker compose down -v # Только вперые, для очистки данных докера
docker compose up postgres
make compose-migrate
```
Запустить сервис можно с помощью команды
```
make prod
# или
docker compose up --build
```
Это запустит базу данных и сам сервис.

Документацию после запуска сервиса можно посмотреть по адресу http://localhost:8080/swagger/index.html с портом 8080 по умолчанию

Перед запуском интеграционных тестов нужно запустить сервис командой выше и не выключать его до конца тестирования. При этом вы должны быть в корне проекта

Запуск интеграционных тестов
```
make integration-test
```

## Configuration

Сконфигурировать приложение можно используя config.yaml, указав путь до файла в переменной CONFIG_PATH
example
```
export CONFIG_PATH=config/config.yaml
```
```
http:
  host: localhost
  port: 8080
postgresql:
  user: user
  password: password
  host: localhost
  port: 5432
  database: postgres
  ssl_mode: disable
  auto_create: false
```
вместе с файлом приложение можно настроить используя перемнные окружения
```
CONFIG_PATH=path - настройка расположения yaml конфиг файла; дефолт значение = "config.yaml"
APP_ENV=prod/dev - настройка окружения приложения
LOG_LEVEL=debug/info/warn/error - настройка уровня логирования; дефолт значение = "debug"
```
для всех полей из yaml файла есть переменные окружения для конфигурации
```
PORT
HOST
PG_USER
PG_PASSWORD
PG_HOST
PG_PORT
PG_DATABASE
PG_SSL
PG_AUTO_CREATE
AGE_URL
GENDER_URL
NATION_URL
```
Перменные содержащие URL в конце предполагают внешние API, которые обогащают полученные данные для базы данных