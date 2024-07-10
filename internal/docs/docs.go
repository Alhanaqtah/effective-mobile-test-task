// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/tasks/{task_id}": {
            "post": {
                "description": "Запускает задачу по ее UUID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Запуск задачи",
                "parameters": [
                    {
                        "type": "string",
                        "description": "UUID задачи",
                        "name": "task_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Запущенная задача",
                        "schema": {
                            "$ref": "#/definitions/models.Task"
                        }
                    },
                    "400": {
                        "description": "Неверный формат UUID или пустое тело запроса",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    },
                    "404": {
                        "description": "Задача не найдена",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/tasks/{task_id}/start": {
            "post": {
                "description": "Запускает задачу по ее UUID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Запуск задачи",
                "parameters": [
                    {
                        "type": "string",
                        "description": "UUID задачи",
                        "name": "task_uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Запущенная задача",
                        "schema": {
                            "$ref": "#/definitions/models.Task"
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    },
                    "404": {
                        "description": "Задача не найдена",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "description": "Получить список пользователей с возможностью фильтрации и пагинации",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Получить пользователей",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Номер страницы",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Строка фильтра",
                        "name": "filter",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.User"
                            }
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            },
            "post": {
                "description": "Создает нового пользователя по паспортным данным",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Создание нового пользователя",
                "parameters": [
                    {
                        "description": "Данные для создания пользователя",
                        "name": "CreateUser",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.CreateUser"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Пользователь создан успешно",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "Некорректные данные запроса",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    },
                    "409": {
                        "description": "Пользователь уже существует",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/users/{user_id}/tasks": {
            "get": {
                "description": "Возвращает задачи пользователя в заданном диапазоне дат",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tasks"
                ],
                "summary": "Получить задачи в диапазоне дат",
                "parameters": [
                    {
                        "type": "string",
                        "description": "UUID пользователя",
                        "name": "user_uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Дата начала в формате YYYY-MM-DD",
                        "name": "start_date",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Дата окончания в формате YYYY-MM-DD",
                        "name": "end_date",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список задач",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Task"
                            }
                        }
                    },
                    "400": {
                        "description": "Некорректный запрос",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        },
        "/users/{uuid}": {
            "put": {
                "description": "Обновить информацию о пользователе по UUID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Обновить пользователя",
                "parameters": [
                    {
                        "type": "string",
                        "description": "UUID пользователя",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Информация о пользователе",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "400": {
                        "description": "Неверный формат UUID или пустое тело запроса",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    },
                    "404": {
                        "description": "Пользователь не найден",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            },
            "delete": {
                "description": "Удалить пользователя по UUID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Удалить пользователя",
                "parameters": [
                    {
                        "type": "string",
                        "description": "UUID пользователя",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Пользователь успешно удалён",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    },
                    "400": {
                        "description": "Неверный формат UUID",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    },
                    "404": {
                        "description": "Пользователь не найден",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка",
                        "schema": {
                            "$ref": "#/definitions/response.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Task": {
            "type": "object",
            "properties": {
                "created_at": {
                    "description": "Время создания задачи",
                    "type": "string"
                },
                "description": {
                    "description": "Описание задачи",
                    "type": "string"
                },
                "done": {
                    "description": "Признак завершённости задачи",
                    "type": "boolean"
                },
                "done_at": {
                    "description": "Время завершения задачи (если задача завершена)",
                    "type": "string"
                },
                "duration": {
                    "description": "Продолжительность выполнения задачи в часах (если указано)",
                    "type": "number"
                },
                "id": {
                    "description": "Уникальный идентификатор задачи",
                    "type": "string"
                },
                "title": {
                    "description": "Заголовок задачи",
                    "type": "string"
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "address": {
                    "description": "Адрес пользователя",
                    "type": "string"
                },
                "id": {
                    "description": "Уникальный идентификатор пользователя",
                    "type": "string"
                },
                "name": {
                    "description": "Имя пользователя",
                    "type": "string"
                },
                "passport_number": {
                    "description": "Номер паспорта пользователя",
                    "type": "integer"
                },
                "passport_serie": {
                    "description": "Серия паспорта пользователя",
                    "type": "integer"
                },
                "patronymic": {
                    "description": "Отчество пользователя",
                    "type": "string"
                },
                "surname": {
                    "description": "Фамилия пользователя",
                    "type": "string"
                }
            }
        },
        "request.CreateUser": {
            "type": "object",
            "properties": {
                "passportNumber": {
                    "description": "Номер паспорта пользователя",
                    "type": "string"
                }
            }
        },
        "response.Response": {
            "type": "object",
            "properties": {
                "data": {
                    "description": "Данные ответа"
                },
                "error": {
                    "description": "Ошибка, если есть",
                    "type": "string"
                },
                "message": {
                    "description": "Сообщение, если есть",
                    "type": "string"
                },
                "status": {
                    "description": "Статус ответа",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Time Tracker API",
	Description:      "Test task for Effective-mobile.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
