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
        "/news": {
            "get": {
                "description": "Fetch news articles with optional filters like ID, title, status, author_id, start_date, end_date, sort_by, and sort_order.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "News"
                ],
                "summary": "Fetch news articles",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Limit the number of results",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page number for pagination",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Filter by news ID",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by news title",
                        "name": "title",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by news status",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Filter by author ID",
                        "name": "author_id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter news starting from this date (RFC3339 format)",
                        "name": "start_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter news until this date (RFC3339 format)",
                        "name": "end_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Field to sort by",
                        "name": "sort_by",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Order of sorting (asc or desc)",
                        "name": "sort_order",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response with the list of news",
                        "schema": {
                            "$ref": "#/definitions/main.ResponseNews"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a news article with the provided details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "News"
                ],
                "summary": "Create a new news article",
                "parameters": [
                    {
                        "description": "Create news request",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.NewsRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Success message",
                        "schema": {
                            "$ref": "#/definitions/main.ResponseMessage"
                        }
                    }
                }
            }
        },
        "/news/{id}": {
            "get": {
                "description": "Fetch a news article by its unique ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "News"
                ],
                "summary": "Get news article by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID of the news article",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response with the news article",
                        "schema": {
                            "$ref": "#/definitions/main.NewsArticle"
                        }
                    }
                }
            },
            "put": {
                "description": "Update the news article with the provided ID using the details provided in the request body",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "News"
                ],
                "summary": "Update an existing news article",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "News ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update news request",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.NewsRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success message",
                        "schema": {
                            "$ref": "#/definitions/main.ResponseMessage"
                        }
                    }
                }
            },
            "delete": {
                "description": "Remove the news article with the specified ID",
                "tags": [
                    "News"
                ],
                "summary": "Delete a news article",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "News ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/topic": {
            "get": {
                "description": "Retrieve a list of topics with optional pagination and filtering",
                "tags": [
                    "Topic"
                ],
                "summary": "Fetch topics",
                "parameters": [
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "Limit the number of topics returned",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Page number for pagination",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Filter by topic ID",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by topic name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Field to sort by",
                        "name": "sort_by",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Sort order (asc or desc)",
                        "name": "sort_order",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response",
                        "schema": {
                            "$ref": "#/definitions/main.ResponseTopic"
                        }
                    }
                }
            },
            "post": {
                "description": "Add a new topic to the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Topic"
                ],
                "summary": "Create a new topic",
                "parameters": [
                    {
                        "description": "Create Topic Request",
                        "name": "topic",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.TopicRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Successful response message",
                        "schema": {
                            "$ref": "#/definitions/main.ResponseMessage"
                        }
                    }
                }
            }
        },
        "/topic/{id}": {
            "get": {
                "description": "Retrieve a topic by its ID",
                "tags": [
                    "Topic"
                ],
                "summary": "Get topic by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Topic ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response with the topic",
                        "schema": {
                            "$ref": "#/definitions/main.Topic"
                        }
                    }
                }
            },
            "put": {
                "description": "Update the details of a topic by its ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Topic"
                ],
                "summary": "Update an existing topic",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Topic ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Update Topic Request",
                        "name": "topic",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.TopicRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response message",
                        "schema": {
                            "$ref": "#/definitions/main.ResponseMessage"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete an existing topic by its ID",
                "tags": [
                    "Topic"
                ],
                "summary": "Delete a topic by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Topic ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Author": {
            "type": "object",
            "properties": {
                "id": {
                    "description": "ID of the author",
                    "type": "integer"
                },
                "name": {
                    "description": "Name of the author",
                    "type": "string"
                }
            }
        },
        "main.NewsArticle": {
            "type": "object",
            "properties": {
                "author": {
                    "description": "Author of the news article",
                    "allOf": [
                        {
                            "$ref": "#/definitions/main.Author"
                        }
                    ]
                },
                "content": {
                    "description": "Content of the news article (HTML)",
                    "type": "string"
                },
                "created_at": {
                    "description": "Created timestamp of the news article",
                    "type": "string"
                },
                "id": {
                    "description": "ID of the news article",
                    "type": "integer"
                },
                "status": {
                    "description": "Status of the news article",
                    "type": "string"
                },
                "title": {
                    "description": "Title of the news article",
                    "type": "string"
                },
                "topics": {
                    "description": "List of topics associated with the news article",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Topic"
                    }
                },
                "updated_at": {
                    "description": "Updated timestamp of the news article",
                    "type": "string"
                }
            }
        },
        "main.NewsRequest": {
            "type": "object",
            "required": [
                "author_id",
                "content",
                "status",
                "title",
                "topic_ids"
            ],
            "properties": {
                "author_id": {
                    "type": "integer"
                },
                "content": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "topic_ids": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "main.PaginationMeta": {
            "type": "object",
            "properties": {
                "current_page": {
                    "description": "Current page number",
                    "type": "integer"
                },
                "total_data": {
                    "description": "Total number of data entries",
                    "type": "integer"
                },
                "total_pages": {
                    "description": "Total number of pages",
                    "type": "integer"
                }
            }
        },
        "main.ResponseMessage": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "main.ResponseNews": {
            "type": "object",
            "properties": {
                "data": {
                    "description": "List of news articles",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.NewsArticle"
                    }
                },
                "meta": {
                    "description": "Metadata about the response",
                    "allOf": [
                        {
                            "$ref": "#/definitions/main.PaginationMeta"
                        }
                    ]
                }
            }
        },
        "main.ResponseTopic": {
            "type": "object",
            "properties": {
                "data": {
                    "description": "List of news articles",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Topic"
                    }
                },
                "meta": {
                    "description": "Metadata about the response",
                    "allOf": [
                        {
                            "$ref": "#/definitions/main.PaginationMeta"
                        }
                    ]
                }
            }
        },
        "main.Topic": {
            "type": "object",
            "properties": {
                "id": {
                    "description": "ID of the topic",
                    "type": "integer"
                },
                "name": {
                    "description": "Name of the topic",
                    "type": "string"
                }
            }
        },
        "main.TopicRequest": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "name": {
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
	Title:            "Swagger Example API",
	Description:      "This is a sample server Petstore server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}