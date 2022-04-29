package schema

import (
	"encoding/json"
	"os"
)

const (
	STRING    = "STRING"
	INTEGER   = "INTEGER"
	FLOAT     = "FLOAT"
	UUID      = "UUID"
	TIMESTAMP = "TIMESTAMP"
	DATETIME  = "DATETIME"
	OBJECT    = "OBJECT"

	PARAMS_DISTRIBUTION_NORMAL  = "NORMAL"
	PARAMS_DISTRIBUTION_UNIFORM = "UNIFORM"
	PARAMS_DISTRIBUTION_RANDOM  = "RANDOM"

	PARAMS_PRECISION_S  = "s"
	PARAMS_PRECISION_MS = "ms"
	PARAMS_PRECISION_US = "us"
	PARAMS_PRECISION_NS = "ns"
)

type Field struct {
	// UUID of the field. Used for cache
	uuid string

	// Name of the field
	Name string `json:"NAME"`

	// Hardcoded value of the field
	Value interface{} `json:"VALUE,omitempty"`

	// Type of the Field
	// Supported values: INTEGER/FLOAT/TIMESTAMP/DATETIME/STRING/UUID/OBJECT
	Type string `json:"TYPE"`

	// Treat field as an array
	Repeated bool `json:"REPEATED,omitempty"`

	// Fields inside the object
	Fields *[]Field `json:"FIELDS,omitempty"`

	// Chance of the field being nullable
	Nullable float64 `json:"NULLABLE,omitempty"`

	// Parameters used to create the value
	Params struct {
		// Use current timestamp
		Now bool `json:"NOW,omitempty"`

		// Format of the DATETIME
		Format string `json:"FORMAT,omitempty"`

		// Precision of the timestamp
		// Supported values: s, ms, us, ns
		Precision string `json:"PRECISION,omitempty"`

		// Increment values from previous runs
		// Supported for: INTEGER, FLOAT, TIMESTAMP, DATETIME
		Incremental bool `json:"INCREMENTAL,omitempty"`
		// Starting point for the increments
		Start interface{} `json:"START,omitempty"`
		// Step to add to each message
		Step interface{} `json:"STEP,omitempty"`

		// Type of distribution used to generate
		// INTEGERS or FLOATS
		// Supported values: NORMAL, UNIFORM, RANDOM
		Distribution string `json:"DISTRIBUTION,omitempty"`

		// Mu used for normal distribution
		Mu float64 `json:"MU,omitempty"`
		// Sigma used for normal distribution
		Sigma float64 `json:"SIGMA,omitempty"`

		// MIN used for uniform distribution
		Min float64 `json:"MIN,omitempty"`
		//MAX used for uniform distribution
		Max float64 `json:"MAX,omitempty"`

		// Multiplier used when generating random floats
		Scale float64 `json:"SCALE,omitempty"`

		// Number of decimals in the float
		Round int `json:"ROUND,omitempty"`
		// Regex format that strings will follow
		Regex string `json:"REGEX,omitempty"`

		// Chance of the string being empty
		Empty float64 `json:"EMPTY,omitempty"`
	} `json:"PARAMS,omitempty"`
}

func addUUIDToFields(fields *[]Field) *[]Field {
	for i, field := range *fields {
		(*fields)[i].uuid, _ = field.generateUUID()
		if field.Fields != nil {
			(*fields)[i].Fields = addUUIDToFields(field.Fields)
		}
	}
	return fields
}

func GetFields(path string) (*[]Field, error) {
	var fields *[]Field

	data, err := os.ReadFile(path)
	if err != nil {
		return fields, err
	}

	if err = json.Unmarshal(data, &fields); err != nil {
		return fields, err
	}
	fields = addUUIDToFields(fields)

	return fields, nil
}
