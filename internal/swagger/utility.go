package swagger

import (
	"reflect"
	"strconv"
	"strings"
)

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
	}

	return fieldType, ""
}

// genExampleValue returns an example value for a given struct field,
// if one is provided using the "example" tag.
func genExampleValue(field reflect.StructField, fieldType string) interface{} {
	eg := field.Tag.Get("example")
	if eg == "" {
		return nil
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
func genSchema(reflected reflect.Type) Schema {
	schema := Schema{}

	if reflected.String() == "types.Nil" {
		return schema
	}

	if reflected.Kind() == reflect.Slice {
		schema.Type = "array"
		itemsSchema := genSchema(reflected.Elem())
		schema.Items = &itemsSchema
		return schema
	} else if reflected.Kind() == reflect.Map {
		schema.Type = "object"
		additionalSchema := genSchema(reflected.Elem())
		schema.AdditionalProperties = &additionalSchema
		return schema
	}

	schema.Properties = map[string]*Schema{}
	example := map[string]interface{}{}

	for i := 0; i < reflected.NumField(); i++ {
		schemaProperty := Schema{}
		field := reflected.Field(i)
		fieldName := field.Tag.Get("json")
		if fieldName == "-" {
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
			schemaProperty = genSchema(reflected.Field(i).Type)
		}

		schema.Properties[fieldName] = &schemaProperty
	}

	if len(example) == len(schema.Properties) {
		schema.Example = example
	}

	return schema
}
