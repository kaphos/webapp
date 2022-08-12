package swagger

import (
	"encoding/json"
	"go/types"
	"io/ioutil"
	"net/http"
	"reflect"
)

func Generate(appName, version string) OpenAPI {
	o := OpenAPI{
		OpenAPIVersion: "3.0.3",
		Info: Info{
			Title:   appName,
			Version: version,
		},
		Paths: make(map[string]Path, 0),
	}
	return o
}

func (o *OpenAPI) AddServer(url, description string) {
	o.Servers = append(o.Servers, Server{url, description})
}

// GenContent is a utility function to generate a Swagger-compatible "Content"
// object, given an interface. Automatically sets it to "application/json"
// content type.
func GenContent(t interface{}) *map[string]MediaType {
	if t == nil {
		return nil
	}

	return &map[string]MediaType{
		"application/json": {
			Schema: genSchema(reflect.TypeOf(t)),
		},
	}
}

// BuildRequestBody is a utility function to create a Swagger-compatible
// request body for a function that requires a given interface.
func BuildRequestBody(t interface{}) *RequestBody {
	if t == nil || t == *new(types.Nil) {
		return nil
	}

	reflected := reflect.TypeOf(t)

	body := RequestBody{}
	body.Description = reflected.String()
	body.Content = *GenContent(t)

	// TODO: Generate "required" field depending on JSON validation tags

	return &body
}

func (o *OpenAPI) AddPath(repo, method, path string, requestBody *RequestBody, responses map[int]Response) {
	val, ok := o.Paths[path]
	if !ok {
		o.Paths[path] = Path{}
		val = o.Paths[path]
	}

	operation := Operation{
		Tags:        []string{repo},
		RequestBody: requestBody,
		Responses:   responses,
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

	o.Paths[path] = val
}

func (o *OpenAPI) Write(filename string) error {
	var file []byte
	var err error

	if strings.HasSuffix(filename, ".json") {
		file, err = json.MarshalIndent(o, "", "  ")
	} else if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		file, err = yaml.Marshal(o)
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
