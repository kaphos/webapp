package swagger

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var pathParamRegexp = regexp.MustCompile(":([a-z]*)")

func processPath(path string) (string, []string) {
	matches := pathParamRegexp.FindAllStringSubmatch(path, -1)
	newPath := pathParamRegexp.ReplaceAllString(path, "{$1}")

	params := make([]string, len(matches))
	for i, match := range matches {
		params[i] = match[1]
	}

	return newPath, params
}

// parseFieldType returns the type and format of a given struct field,
// based on its type.
func parseFieldType(field reflect.StructField) (string, string) {
	fieldType := field.Type.String()
	if strings.HasPrefix(fieldType, "float") {
		return "number", fieldType
	} else if strings.HasPrefix(fieldType, "int") {
		return "integer", fieldType
	} else if fieldType == "bool" {
		return "boolean", ""
	} else if fieldType == "time.Time" {
		return "string", "date"
	}

	return fieldType, ""
}

// genExampleValue returns an example value for a given struct field,
// if one is provided using the "example" tag.
func genExampleValue(field reflect.StructField, fieldType string) interface{} {
	eg := field.Tag.Get("example")
	if eg == "" {
		switch fieldType {
		case "integer":
			return 1
		case "number":
			return 1.0
		case "boolean":
			return true
		default:
			return "string value"
		}
	}

	switch fieldType {
	case "integer":
		egNum, _ := strconv.Atoi(eg)
		return egNum
	case "number":
		egNum, _ := strconv.ParseFloat(eg, 64)
		return egNum
	case "boolean":
		return true
	default:
		return eg
	}
}

// genSchema creates a Schema for a given reflect.Type.
// It also recursively resolves types for slices, maps and structs.
func genSchema(reflected reflect.Type, hideEmptyBind bool) (Schema, []Parameter) {
	schema := Schema{Required: make([]string, 0)}
	queryParams := make([]Parameter, 0)

	typeName := reflected.String()
	typeKind := reflected.Kind()

	if typeName == "types.Nil" || typeName == "string" || typeName == "int" {
		return schema, queryParams
	}

	if typeKind == reflect.Slice {
		schema.Type = "array"
		itemsSchema, itemsParams := genSchema(reflected.Elem(), hideEmptyBind)
		schema.Items = &itemsSchema
		return schema, append(queryParams, itemsParams...)
	} else if typeKind == reflect.Map {
		schema.Type = "object"
		additionalSchema, additionalParams := genSchema(reflected.Elem(), hideEmptyBind)
		schema.AdditionalProperties = &additionalSchema
		return schema, append(queryParams, additionalParams...)
	} else if typeKind != reflect.Struct {
		panic("return type '" + typeKind.String() + "' is not supported")
	}

	schema.Properties = map[string]*Schema{}
	example := map[string]interface{}{}

	for i := 0; i < reflected.NumField(); i++ {
		schemaProperty := Schema{}
		field := reflected.Field(i)

		fieldName := field.Tag.Get("json")
		if fieldName == "" {
			fieldName = field.Tag.Get("form")
		}

		binding := field.Tag.Get("binding")
		if fieldName == "-" || (hideEmptyBind && binding == "-") {
			// Field should be excluded from JSON
			continue
		}
		if fieldName == "" {
			// JSON tag not provided; default to field name
			fieldName = field.Name
		}

		fieldType, fieldFmt := parseFieldType(field)

		if fieldType == "number" || fieldType == "integer" || fieldType == "string" || fieldType == "boolean" {
			if egVal := genExampleValue(field, fieldType); egVal != nil {
				example[fieldName] = egVal
			}
			schemaProperty.Type = fieldType
			schemaProperty.Format = fieldFmt
		} else {
			var schemaParams []Parameter
			schemaProperty, schemaParams = genSchema(reflected.Field(i).Type, hideEmptyBind)
			queryParams = append(queryParams, schemaParams...)
		}

		schema.Properties[fieldName] = &schemaProperty

		if strings.Contains(binding, "required") {
			schema.Required = append(schema.Required, fieldName)
		}

		if formTag := field.Tag.Get("form"); formTag != "" {
			queryParams = append(queryParams, Parameter{
				Name:     formTag,
				In:       "query",
				Required: strings.Contains(binding, "required"),
				Schema:   schema,
			})
		}
	}

	if len(example) == len(schema.Properties) {
		schema.Example = example
	}

	return schema, queryParams
}
