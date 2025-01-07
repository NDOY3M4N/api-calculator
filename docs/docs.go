// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag/v2"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "components": {"schemas":{"main.APIError":{"properties":{"error":{"type":"string"}},"type":"object"},"main.APISuccess":{"properties":{"result":{"type":"integer"}},"type":"object"},"main.Payload":{"properties":{"number1":{"example":6,"type":"integer"},"number2":{"example":9,"type":"integer"}},"type":"object"}}},
    "info": {"contact":{"email":"pa.ndoye@outlook.com","name":"Abdoulaye NDOYE","url":"https://github.com/NDOY3M4N"},"description":"{{escape .Description}}","license":{"name":"MIT","url":"https://github.com/NDOY3M4N/api-calculator/blob/main/LICENSE.md"},"title":"{{.Title}}","version":"{{.Version}}"},
    "externalDocs": {"description":"","url":""},
    "paths": {"/add":{"post":{"description":"Add two numbers together","requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/main.Payload"}}},"description":"Numbers needed for the operation","required":true},"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/main.APISuccess"}}},"description":"OK"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/main.APIError"}}},"description":"Bad Request"}},"summary":"Add two numbers","tags":["Math"]}},"/divide":{"post":{"description":"Divide two numbers together","requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/main.Payload"}}},"description":"Numbers needed for the operation","required":true},"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/main.APISuccess"}}},"description":"OK"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/main.APIError"}}},"description":"Bad Request"}},"summary":"Divide two numbers","tags":["Math"]}},"/multiply":{"post":{"description":"Multiply two numbers together","requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/main.Payload"}}},"description":"Numbers needed for the operation","required":true},"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/main.APISuccess"}}},"description":"OK"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/main.APIError"}}},"description":"Bad Request"}},"summary":"Multiply two numbers","tags":["Math"]}},"/substract":{"post":{"description":"Substract two numbers together","requestBody":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/main.Payload"}}},"description":"Numbers needed for the operation","required":true},"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/main.APISuccess"}}},"description":"OK"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/main.APIError"}}},"description":"Bad Request"}},"summary":"Substract two numbers","tags":["Math"]}},"/sum":{"post":{"description":"Add all numbers in an array","requestBody":{"content":{"application/json":{"schema":{"items":{"type":"integer"},"type":"array"}}},"description":"Array of numbers needed for the operation","required":true},"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/main.APISuccess"}}},"description":"OK"},"400":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/main.APIError"}}},"description":"Bad Request"}},"summary":"Sum numbers","tags":["Math"]}}},
    "openapi": "3.1.0",
    "servers": [
        {"description":"Development server","url":"http://localhost:3000/api/v1"}
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Title:            "Calculator API",
	Description:      "This is a simple server for Calculator API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
