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
        "/api/check-code": {
            "post": {
                "description": "Check if the provided code is valid for the given username",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Check if a code is valid",
                "parameters": [
                    {
                        "description": "Request body for checking code",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/transport.CheckCodeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/transport.CheckCodeResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/check-username": {
            "get": {
                "description": "Check if a username exists in the database",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Check if a username exists",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username to check",
                        "name": "username",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/transport.CheckUsernameResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/users/login": {
            "post": {
                "description": "Logs in a user and returns the user ID, refresh token, and access token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Log in a user",
                "parameters": [
                    {
                        "description": "User details",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/transport.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/transport.LogInResponse"
                        }
                    }
                }
            }
        },
        "/api/users/refresh_token": {
            "get": {
                "description": "Refreshes the access token using the refresh token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Refresh access token",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Refresh token",
                        "name": "refresh",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/transport.LogInResponse"
                        }
                    }
                }
            }
        },
        "/api/users/signup": {
            "post": {
                "description": "Registers a new user and returns the user ID, refresh token, and access token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Sign up a new user",
                "parameters": [
                    {
                        "description": "User details",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/transport.SignUpRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/transport.LogInResponse"
                        }
                    }
                }
            }
        },
        "/api/users/{id}": {
            "get": {
                "description": "Retrieves a user by their ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get a user by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.User"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "transport.CheckCodeRequest": {
            "description": "Request structure for checking a code",
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "transport.CheckCodeResponse": {
            "description": "Response structure for checking a code",
            "type": "object",
            "properties": {
                "valid": {
                    "type": "boolean"
                }
            }
        },
        "transport.CheckUsernameResponse": {
            "description": "Check if a username exists",
            "type": "object",
            "properties": {
                "found": {
                    "type": "boolean"
                }
            }
        },
        "transport.LogInResponse": {
            "type": "object",
            "properties": {
                "authorization": {
                    "type": "string"
                },
                "refresh": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/types.User"
                }
            }
        },
        "transport.LoginRequest": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "transport.SignUpRequest": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "types.User": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "password": {
                    "type": "string"
                },
                "tg_username": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
