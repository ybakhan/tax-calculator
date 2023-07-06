// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Yasser Khan",
            "url": "http://github.com/ybakhan",
            "email": "ybakhan@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/login": {
            "post": {
                "description": "returns api key for calling taxes api",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "taxes"
                ],
                "summary": "login to taxes api",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.loginResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/main.taxServerError"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.taxServerResponse"
                        }
                    },
                    "405": {
                        "description": "Method Not Allowed",
                        "schema": {
                            "$ref": "#/definitions/main.taxServerError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.taxServerError"
                        }
                    }
                }
            }
        },
        "/tax/{year}": {
            "get": {
                "description": "calculate taxes for given a salary and tax year",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "taxes"
                ],
                "summary": "calculate taxes",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "tax year",
                        "name": "year",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "salary",
                        "name": "s",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/taxcalculator.TaxCalculation"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/main.taxServerResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.taxServerError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.loginResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "main.taxServerError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "main.taxServerResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "taxbracket.Bracket": {
            "type": "object",
            "properties": {
                "max": {
                    "type": "number",
                    "example": 100392
                },
                "min": {
                    "type": "number",
                    "example": 50197
                },
                "rate": {
                    "type": "number",
                    "example": 0.205
                }
            }
        },
        "taxcalculator.BracketTax": {
            "type": "object",
            "properties": {
                "band": {
                    "$ref": "#/definitions/taxbracket.Bracket"
                },
                "tax": {
                    "type": "number",
                    "example": 984.62
                }
            }
        },
        "taxcalculator.TaxCalculation": {
            "type": "object",
            "properties": {
                "effective_rate": {
                    "type": "number",
                    "example": 0.15
                },
                "salary": {
                    "type": "number",
                    "example": 55000
                },
                "taxes_by_band": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/taxcalculator.BracketTax"
                    }
                },
                "total_taxes": {
                    "type": "number",
                    "example": 8514.17
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
	Title:            "Tax Calculator API",
	Description:      "REST API for calculating taxes",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
