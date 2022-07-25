package swagger

import (
	"encoding/json"
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

func BuildRequestBody(t interface{}) *RequestBody {
	if t == nil {
		return nil
	}

	reflected := reflect.TypeOf(t)

	body := RequestBody{}
	body.Description = reflected.String()
	body.Content = map[string]MediaType{}

	mediaType := MediaType{Schema{
		//Example:    map[string]interface{}{},
		Properties: map[string]SchemaProperty{},
	}}

	for i := 0; i < reflected.NumField(); i++ {
		field := reflected.Field(i)
		fieldName := field.Tag.Get("json")
		if fieldName == "-" {
			continue
		}
		if fieldName == "" {
			fieldName = field.Name
		}

		fieldType := field.Type.String()

		//if eg := field.Tag.Get("example"); eg != "" {
		//	mediaType.Schema.Example[fieldName] = eg
		//} else {
		//	switch fieldType {
		//	case "string":
		//		mediaType.Schema.Example[fieldName] = "string"
		//	case "int":
		//		mediaType.Schema.Example[fieldName] = 0
		//	}
		//}

		mediaType.Schema.Properties[fieldName] = SchemaProperty{Type: fieldType}
	}

	body.Content["application/json"] = mediaType

	return &body
}

func (o *OpenAPI) AddPath(repo, method, path string, requestBody *RequestBody, responses map[string]Response) {
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

	// TODO: Handle authentication (e.g., using keycloak)

	o.Paths[path] = val
}

func (o *OpenAPI) Write(filename string) error {
	file, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
