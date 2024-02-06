// Package openapi Code generated by swaggo/swag. DO NOT EDIT
package openapi

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "license": {
            "name": "GNU General Public License v3.0",
            "url": "https://github.com/lz1marine/notification-service/?tab=GPL-3.0-1-ov-file"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/internal/notifications/{id}": {
            "post": {
                "description": "post a notification to a channel given an event id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notifications"
                ],
                "summary": "Post a notification to a channel",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Event ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "The request body.",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.ChannelNotificationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.ChannelNotificationResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        },
        "/api/v1/notifications": {
            "get": {
                "description": "gets the list of all channels, including whether they are enabled",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notifications"
                ],
                "summary": "Gets the list of all channels, including whether they are enabled",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.ChannelResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/notifications/sub/{id}": {
            "get": {
                "description": "gets the list of all channels subscribed to by the user, including whether they are enabled for the user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notifications"
                ],
                "summary": "Gets the list of all channels subscribed to by the user, including whether they are enabled for the user",
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
                            "$ref": "#/definitions/v1.ChannelResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            },
            "patch": {
                "description": "patch he channel list that the user has subscribed to",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "notifications"
                ],
                "summary": "Patch the channel list that the user has subscribed to",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "The request body.",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.SetChannelsRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/v1.SetChannelsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/httputil.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "httputil.HTTPError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "status bad request"
                }
            }
        },
        "v1.Channel": {
            "type": "object",
            "properties": {
                "is_enabled": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "v1.ChannelNotificationRequest": {
            "type": "object",
            "properties": {
                "channel": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "subject": {
                    "type": "string"
                },
                "template_id": {
                    "type": "string"
                },
                "topic_id": {
                    "type": "string"
                }
            }
        },
        "v1.ChannelNotificationResponse": {
            "type": "object"
        },
        "v1.ChannelResponse": {
            "type": "object",
            "properties": {
                "channels": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/v1.Channel"
                    }
                }
            }
        },
        "v1.SetChannelsRequest": {
            "type": "object",
            "properties": {
                "channels": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/v1.Channel"
                    }
                }
            }
        },
        "v1.SetChannelsResponse": {
            "type": "object"
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "basic"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Notification Server API",
	Description:      "This is the notification server API that handles both external user notification subscriptions and internal notifications",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
