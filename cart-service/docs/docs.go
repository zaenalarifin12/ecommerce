// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/carts": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "List all products in the cart",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Carts"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/domain.Cart"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Add a product to the cart",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Carts"
                ],
                "parameters": [
                    {
                        "description": "Product details",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.CartRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Product added successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/carts/tx": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Tx cart",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "TransactionCarts"
                ],
                "responses": {
                    "200": {
                        "description": "Product tx successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/carts/tx-rollback": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Tx Cart Rollback",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "TransactionCarts"
                ],
                "responses": {
                    "200": {
                        "description": "Rollback Tx Successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/carts/{uuid}": {
            "put": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Update quantity of a product in the cart",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Carts"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Product quantity",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.CartRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Product quantity updated successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Remove a product from the cart",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Carts"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Product UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Product removed successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/domain.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Cart": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "productUUID": {
                    "type": "string"
                },
                "quantity": {
                    "type": "integer"
                },
                "updatedAt": {
                    "type": "string"
                },
                "userUUID": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "domain.ErrorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "object",
                    "properties": {
                        "code": {
                            "type": "integer"
                        },
                        "detail": {},
                        "message": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "dto.CartRequest": {
            "type": "object",
            "properties": {
                "product_id": {
                    "description": "Product ID of the item to be added to the cart.",
                    "type": "string",
                    "example": "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
                },
                "quantity": {
                    "description": "Quantity of the item to be added.",
                    "type": "integer",
                    "example": 1
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    },
    "tags": [
        {
            "name": "Carts"
        },
        {
            "name": "TransactionCarts"
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:5002",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Cart Service API",
	Description:      "This is a sample server celler server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
