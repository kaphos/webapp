package swagger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/types"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

func Generate(appName, version string) OpenAPI {
	o := OpenAPI{
		OpenAPIVersion: "3.0.3",
		Info: Info{
			Title:   appName,
			Version: version,
		},
		Paths: make(map[string]Path, 0),
		Components: Components{
			SecuritySchemes: map[string]SecurityScheme{
				"keycloak": {
					Type:        "http",
					Description: "Keycloak authentication",
					Scheme:      "bearer",
				},
			},
		},
	}
	return o
}

func (o *OpenAPI) AddServer(url, description string) {
	o.Servers = append(o.Servers, Server{url, description})
}

// GenContent is a utility function to generate a Swagger-compatible "Content"
// object, given an interface. Automatically sets it to "application/json"
// content type.
func GenContent(t interface{}, hideEmptyBind bool) (*map[string]MediaType, []Parameter) {
	if t == nil {
		return nil, make([]Parameter, 0)
	}

	schema, queryParams := genSchema(reflect.TypeOf(t), hideEmptyBind)

	return &map[string]MediaType{
		"application/json": {
			Schema: schema,
		},
	}, queryParams
}

// buildRequestBody is a utility function to create a Swagger-compatible
// request body for a function that requires a given interface.
func buildRequestBody(t interface{}, hideEmptyBind bool) (*RequestBody, []Parameter) {
	queryParams := make([]Parameter, 0)

	if t == nil || t == *new(types.Nil) {
		return nil, queryParams
	}

	reflected := reflect.TypeOf(t)

	body := RequestBody{}
	body.Description = reflected.String()
	bodyContent, queryParams := GenContent(t, hideEmptyBind)
	body.Content = *bodyContent

	return &body, queryParams
}

func (val *Path) buildParams(params map[string]SimpleParam, pathParams []string) {
	for paramName, simpleParam := range params {
		found := false
		for _, x := range val.Parameters {
			if x.Name == paramName {
				found = true
				break
			}
		}
		if found {
			continue
		}

		param := Parameter{
			Name:        paramName,
			Description: simpleParam.Description,
			Schema:      Schema{Type: simpleParam.Type},
		}

		if slices.Contains(pathParams, paramName) {
			param.In = "path"
			param.Required = true
		} else {
			param.In = "query"
		}

		val.Parameters = append(val.Parameters, param)
	}

	// Include any parameters that we did not pass in,
	// but were found in the URL
	for _, param := range pathParams {
		if _, ok := params[param]; ok {
			// Already added
			continue
		}

		found := false
		for _, x := range val.Parameters {
			if x.Name == param {
				found = true
				break
			}
		}
		if found {
			continue
		}

		val.Parameters = append(val.Parameters, Parameter{
			Name:     param,
			In:       "path",
			Required: true,
		})
	}
}

func (o *OpenAPI) AddPath(t interface{}, repo, method, path, summary, description string,
	params map[string]SimpleParam, authGroups []string, responses map[int]Response) {
	requestBody, _ := buildRequestBody(t, method == "POST" || method == "PUT")
	cleanedPath, pathParams := processPath(path)

	val, ok := o.Paths[cleanedPath]
	if !ok {
		val = Path{}
	}

	val.buildParams(params, pathParams)

	operation := Operation{
		Summary:     summary,
		Description: description,
		Tags:        []string{repo},
		RequestBody: requestBody,
		Responses:   responses,
		Security:    make([]map[string][]string, 0),
	}

	if len(authGroups) > 0 {
		operation.Security = append(operation.Security, map[string][]string{"keycloak": authGroups})
	}

	switch method {
	case http.MethodGet:
		val.Get = &operation
	case http.MethodPut:
		val.Put = &operation
	case http.MethodPost:
		val.Post = &operation
	case http.MethodDelete:
		val.Delete = &operation
	}

	o.Paths[cleanedPath] = val
}

func (o *OpenAPI) Write(filename string) error {
	var file []byte
	var err error

	if strings.HasSuffix(filename, ".json") {
		file, err = json.MarshalIndent(o, "", "  ")
	} else if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		var b bytes.Buffer
		encoder := yaml.NewEncoder(&b)
		encoder.SetIndent(2)
		if err := encoder.Encode(o); err != nil {
			return err
		}

		if err := encoder.Close(); err != nil {
			return err
		}

		file = b.Bytes()
	} else {
		return fmt.Errorf("unrecognised file extension")
	}

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
