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
func parseFieldType(field reflect.StructField, binding string) (string, string) {
	fieldType := field.Type.String()
	if strings.HasPrefix(fieldType, "float") {
		return "number", fieldType
	} else if fieldType == "null.Float" {
		return "number", "float64"
	} else if strings.HasPrefix(fieldType, "int") {
		return "integer", fieldType
	} else if fieldType == "null.Int" {
		return "integer", "int64"
	} else if fieldType == "bool" || fieldType == "null.Bool" {
		return "boolean", ""
	} else if fieldType == "time.Time" || fieldType == "null.Time" {
		return "string", "date-time"
	} else if fieldType == "uuid.UUID" {
		return "string", "uuid"
	} else if fieldType == "string" || fieldType == "null.String" {
		if strings.Contains(binding, "email") {
			return "string", "email"
		}
		return "string", ""
	}

	return fieldType, ""
}

// genExampleValue returns an example value for a given struct field,
// if one is provided using the "example" tag.
func genExampleValue(field reflect.StructField, fieldType, fieldFmt string) interface{} {
	eg := field.Tag.Get("example")
	if eg == "" {
		switch fieldType {
		case "integer":
			return 123
		case "number":
			return 12.3
		case "boolean":
			return true
		case "string":
			switch fieldFmt {
			case "email":
				return "johndoe@email.com"
			case "date":
				return "2023-05-21"
			case "date-time":
				return "2023-05-21T17:32:28Z"
			case "uuid":
				return "3fa85f64-5717-4562-b3fc-2c963f66afa6"
			}

			return "string value"
		default:
			return "undefined-value"
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

	if typeName == "types.Nil" || typeName == "string" || typeName == "int" || typeName == "uuid.UUID" {
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

		fieldType, fieldFmt := parseFieldType(field, binding)

		if fieldType == "number" || fieldType == "integer" || fieldType == "string" || fieldType == "boolean" {
			if egVal := genExampleValue(field, fieldType, fieldFmt); egVal != nil {
				example[fieldName] = egVal
			}
			schemaProperty.Type = fieldType
			schemaProperty.Format = fieldFmt
			schemaProperty.Nullable = strings.HasPrefix(field.Type.String(), "null.")
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
