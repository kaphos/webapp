package main

import (
	"github.com/kaphos/webapp"
	"github.com/kaphos/webapp/internal/swagger"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"testing"
)

func TestSwagger(t *testing.T) {
	s := setupServer()
	err := s.GenDocs([]webapp.APIServer{{"http://localhost:5000", "Dev server"}}, "swagger.yml")
	assert.Nil(t, err)

	yamlFile, err := ioutil.ReadFile("swagger.yml")
	assert.Nil(t, err)

	var api swagger.OpenAPI
	err = yaml.Unmarshal(yamlFile, &api)
	assert.Nil(t, err)
	assert.Equal(t, api.OpenAPIVersion, "3.0.3")
	assert.ElementsMatch(t, api.Servers, []swagger.Server{{"http://localhost:5000", "Dev server"}})

	pathItems, found := api.Paths["/items/"]
	assert.True(t, found)

	assert.Equal(t, pathItems.Get.Summary, "Retrieves the list of items stored in the database.")
	assert.Equal(t, pathItems.Get.Description, "Simply fetches all items.")

	getItemsSuccess, found := pathItems.Get.Responses[200]
	assert.True(t, found)
	getItemsJson, found := getItemsSuccess.Content["application/json"]
	assert.True(t, found)
	assert.Equal(t, getItemsJson.Schema.Type, "array")
	properties := getItemsJson.Schema.Items.Properties
	assert.Equal(t, properties["id"].Type, "string")
	assert.Equal(t, properties["id"].Format, "uuid")
	assert.Equal(t, properties["created"].Type, "string")
	assert.Equal(t, properties["created"].Format, "date-time")
	assert.Equal(t, properties["created"].Nullable, false)
	assert.Equal(t, properties["edited"].Type, "string")
	assert.Equal(t, properties["edited"].Format, "date-time")
	assert.Equal(t, properties["edited"].Nullable, true)
	assert.Equal(t, properties["name"].Type, "string")
	assert.Equal(t, properties["name"].Nullable, false)
	assert.Equal(t, properties["owner"].Type, "string")
	assert.Equal(t, properties["owner"].Nullable, true)
	assert.Equal(t, properties["found"].Type, "boolean")
	assert.Equal(t, properties["found"].Nullable, true)
	assert.Equal(t, properties["count"].Type, "integer")
	assert.Equal(t, properties["count"].Format, "int64")
	assert.Equal(t, properties["count"].Nullable, true)
	assert.Equal(t, properties["price"].Type, "number")
	assert.Equal(t, properties["price"].Format, "float64")
	assert.Equal(t, properties["price"].Nullable, true)

	itemsExample := getItemsJson.Schema.Items.Example
	assert.Equal(t, "3fa85f64-5717-4562-b3fc-2c963f66afa6", itemsExample["id"])
	assert.Equal(t, "2023-05-21T17:32:28Z", itemsExample["created"])
	assert.Equal(t, "2023-05-21T17:32:28Z", itemsExample["edited"])
	assert.Equal(t, "string value", itemsExample["name"])
	assert.Equal(t, "string value", itemsExample["owner"])
	assert.Equal(t, true, itemsExample["found"])
	assert.Equal(t, 123, itemsExample["count"])
	assert.Equal(t, 12.3, itemsExample["price"])

	assert.Equal(t, pathItems.Post.Summary, "Creates a new item.")
	assert.Equal(t, pathItems.Post.Description, "Only allowed by authenticated users.")

	pathUsers, found := api.Paths["/users/"]
	assert.True(t, found)
	createUserSchema := pathUsers.Post.RequestBody.Content["application/json"].Schema
	assert.Equal(t, "integer", createUserSchema.Properties["id"].Type)
	assert.Equal(t, "int", createUserSchema.Properties["id"].Format)
	assert.Equal(t, "string", createUserSchema.Properties["name"].Type)
	assert.Equal(t, "string", createUserSchema.Properties["email"].Type)
	assert.Equal(t, "email", createUserSchema.Properties["email"].Format)
	assert.Equal(t, "boolean", createUserSchema.Properties["admin"].Type)
	assert.Equal(t, "integer", createUserSchema.Properties["groups"].Type)
	assert.Equal(t, "int", createUserSchema.Properties["groups"].Format)
	assert.Equal(t, "number", createUserSchema.Properties["age"].Type)

	assert.Equal(t, 123, createUserSchema.Example["id"])
	assert.Equal(t, "John Doe", createUserSchema.Example["name"])
	assert.Equal(t, true, createUserSchema.Example["admin"])
	assert.Equal(t, 31, createUserSchema.Example["groups"])
	assert.Equal(t, 12.3, createUserSchema.Example["age"])
}
