package domain

type GenericType string

const (
	TypeInteger      GenericType = "integer"
	TypeString       GenericType = "string"
	TypeBoolean      GenericType = "boolean"
	TypeIntegerArray GenericType = "integer_array"
	TypeStringArray  GenericType = "string_array"
	// Optional: Custom types could also be represented here or handled separately
)

type SchemaParameter struct {
	Name     string      `json:"name"`
	Type     GenericType `json:"type"`
	IsCustom bool        `json:"is_custom"`
}

type ProblemSchema struct {
	FunctionName string            `json:"function_name"`
	Parameters   []SchemaParameter `json:"parameters"`
	ReturnType   GenericType       `json:"return_type"`
}
