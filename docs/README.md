# Swagger Documentation

## Генерация Swagger документации

Для генерации Swagger документации используйте команду:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/server/main.go
```

Это создаст файлы `docs/swagger.json` и `docs/swagger.yaml`, которые будут использоваться Swagger UI.

## Доступ к Swagger UI

После запуска приложения, Swagger UI будет доступен по адресу:
- http://localhost:8080/swagger/index.html

## Обновление документации

После изменения аннотаций в коде, необходимо заново запустить `swag init` для обновления документации.

