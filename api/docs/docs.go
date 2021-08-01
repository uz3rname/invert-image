// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/get_last_images": {
            "get": {
                "description": "Get last images",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "images"
                ],
                "summary": "returns last images from db",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Number of images to return",
                        "name": "count",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.ImagePairListDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorDTO"
                        }
                    }
                }
            }
        },
        "/get_task_status": {
            "get": {
                "description": "Get task status",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "images"
                ],
                "summary": "returns status of background task",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Task ID",
                        "name": "taskId",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.TaskStatusDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorDTO"
                        }
                    }
                }
            }
        },
        "/negative_image": {
            "post": {
                "description": "Create negative of image",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "images"
                ],
                "summary": "creates inversion of image",
                "parameters": [
                    {
                        "description": "Base64 encoded image",
                        "name": "image",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.UploadImageDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.InvertImageResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorDTO"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.EncodedImage": {
            "type": "object",
            "properties": {
                "base64": {
                    "type": "string"
                },
                "mimeType": {
                    "type": "string"
                }
            }
        },
        "api.ErrorDTO": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "string",
                    "enum": [
                        "error"
                    ]
                }
            }
        },
        "api.ImagePair": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "hash": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "negative": {
                    "$ref": "#/definitions/api.EncodedImage"
                },
                "original": {
                    "$ref": "#/definitions/api.EncodedImage"
                }
            }
        },
        "api.ImagePairListDTO": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "items": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/api.ImagePair"
                    }
                },
                "status": {
                    "type": "string",
                    "enum": [
                        "ok"
                    ]
                }
            }
        },
        "api.InvertImageResponse": {
            "type": "object",
            "properties": {
                "pair": {
                    "$ref": "#/definitions/api.ImagePair"
                },
                "status": {
                    "type": "string",
                    "enum": [
                        "ok",
                        "defered"
                    ]
                },
                "taskId": {
                    "type": "string"
                }
            }
        },
        "api.TaskStatusDTO": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string",
                    "enum": [
                        "ok"
                    ]
                },
                "taskStatus": {
                    "type": "string",
                    "enum": [
                        "new",
                        "running",
                        "failed",
                        "done",
                        "canceled"
                    ]
                }
            }
        },
        "api.UploadImageDTO": {
            "type": "object",
            "required": [
                "data"
            ],
            "properties": {
                "data": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "/api",
	Schemes:     []string{},
	Title:       "Backend test task",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
