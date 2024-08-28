package test

import (
	"github.com/Layr-Labs/eigenda/configuration"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReadingFile(t *testing.T) {
	configFile := "basic-config1.json"
	config := DefaultBasicConfig()

	err := configuration.ParseJsonFile(&config, configFile)
	assert.NoError(t, err)

	assert.Equal(t, "asdf", config.Foo)
	assert.Equal(t, 1234, config.Bar)
	assert.Equal(t, true, config.Baz)
}

func TestReadingMultipleFiles(t *testing.T) {
	configFile1 := "basic-config1.json"
	configFile2 := "basic-config2.json"
	config := DefaultBasicConfig()

	err := configuration.ParseJsonFiles(&config, configFile1, configFile2)
	assert.NoError(t, err)

	assert.Equal(t, "qwerty", config.Foo)
	assert.Equal(t, 4321, config.Bar)
	assert.Equal(t, true, config.Baz)
}

func TestReadingZeroFiles(t *testing.T) {
	config := DefaultBasicConfig()

	err := configuration.ParseJsonFiles(&config)
	assert.NoError(t, err)

	assert.Equal(t, "this is a default value", config.Foo)
	assert.Equal(t, 1337, config.Bar)
	assert.Equal(t, false, config.Baz)
}

// Test the parsing of a simple config that has all values fully defined.
func TestWithAllValues(t *testing.T) {
	jsonString :=
		`{
			"Foo": "asdf",
			"Bar": 1234,
			"Baz": true
		}`
	config := DefaultBasicConfig()

	err := configuration.ParseJsonString(&config, jsonString)
	assert.NoError(t, err)

	assert.Equal(t, "asdf", config.Foo)
	assert.Equal(t, 1234, config.Bar)
	assert.Equal(t, true, config.Baz)
}

// Test the parsing of a config that has an extra value not in the struct. We want this to throw an error,
// since this means somebody probably mistyped a field name.
func TestWithValueNotInStruct(t *testing.T) {
	jsonString :=
		`{
  			"Foo": "asdf",
  			"Bar": 1234,
  			"Baz": true,
  			"ThisIsNotAField": "This is not a field"
		}`

	config := DefaultBasicConfig()

	err := configuration.ParseJsonString(&config, jsonString)
	assert.Error(t, err)
}

// Test a config that is missing an entry. The struct should not be modified for this value.
func TestDefaultValue(t *testing.T) {
	configJson :=
		`{
  			"Foo": "asdf",
  			"Baz": false
		}`
	config := DefaultBasicConfig()

	err := configuration.ParseJsonString(&config, configJson)
	assert.NoError(t, err)

	assert.Equal(t, "asdf", config.Foo)
	assert.Equal(t, 1337, config.Bar)
	assert.Equal(t, false, config.Baz)
}

// Test a config that just contains open and close parens.
func TestEmptyConfig(t *testing.T) {
	configJson := "{}"
	config := DefaultBasicConfig()

	err := configuration.ParseJsonString(&config, configJson)
	assert.NoError(t, err)

	assert.Equal(t, "this is a default value", config.Foo)
	assert.Equal(t, 1337, config.Bar)
	assert.Equal(t, false, config.Baz)
}

// Test configuration with nested structs.
func TestNestedStructs(t *testing.T) {
	configString :=
		`{
  			"RecursiveConfig": {
    			"RecursiveConfig": {
      				"BasicConfig": {
        				"Foo": "xxxx",
        				"Bar": 42
      				}
    			},
    			"BasicConfig": {
      				"Foo": "qwerty",
      				"Bar": 4321,
      				"Baz": false
    			}
  			},
  			"BasicConfig": {
    			"Foo": "asdf",
    			"Bar": 1234,
    			"Baz": true
  			}
		}`
	config := NestedConfig{}

	err := configuration.ParseJsonString(&config, configString)
	assert.NoError(t, err)

	assert.NotNil(t, config.BasicConfig)
	assert.Equal(t, "asdf", config.BasicConfig.Foo)
	assert.Equal(t, 1234, config.BasicConfig.Bar)
	assert.Equal(t, true, config.BasicConfig.Baz)

	assert.NotNil(t, config.RecursiveConfig)

	assert.NotNil(t, config.RecursiveConfig.BasicConfig)
	assert.Equal(t, "qwerty", config.RecursiveConfig.BasicConfig.Foo)
	assert.Equal(t, 4321, config.RecursiveConfig.BasicConfig.Bar)
	assert.Equal(t, false, config.RecursiveConfig.BasicConfig.Baz)

	assert.NotNil(t, config.RecursiveConfig.RecursiveConfig)

	assert.NotNil(t, config.RecursiveConfig.RecursiveConfig.BasicConfig)
	assert.Equal(t, "xxxx", config.RecursiveConfig.RecursiveConfig.BasicConfig.Foo)
	assert.Equal(t, 42, config.RecursiveConfig.RecursiveConfig.BasicConfig.Bar)
	assert.Equal(t, false, config.RecursiveConfig.RecursiveConfig.BasicConfig.Baz)

	assert.Nil(t, config.RecursiveConfig.RecursiveConfig.RecursiveConfig)
}

func TestAllPrimitiveTypes(t *testing.T) {
	configString :=
		`{
  			"Bool": true,
  			"Int": 1234,
  			"Int8": 123,
  			"Int16": 1234,
  			"Int32": 12345,
  			"Int64": 123456,
  			"Uint": 1234,
  			"Uint8": 123,
  			"Uint16": 1234,
  			"Uint32": 12345,
  			"Uint64": 123456,
  			"Float32": 123.456,
  			"Float64": 1234.5678,
  			"String": "asdf",
  			"Time": "2000-01-02T03:04:05Z",
  			"Duration": 12345
		}`
	config := AllPrimitiveTypes{}

	err := configuration.ParseJsonString(&config, configString)
	assert.NoError(t, err)

	assert.Equal(t, true, config.Bool)
	assert.Equal(t, 1234, config.Int)
	assert.Equal(t, int8(123), config.Int8)
	assert.Equal(t, int16(1234), config.Int16)
	assert.Equal(t, int32(12345), config.Int32)
	assert.Equal(t, int64(123456), config.Int64)
	assert.Equal(t, uint(1234), config.Uint)
	assert.Equal(t, uint8(123), config.Uint8)
	assert.Equal(t, uint16(1234), config.Uint16)
	assert.Equal(t, uint32(12345), config.Uint32)
	assert.Equal(t, uint64(123456), config.Uint64)
	assert.Equal(t, float32(123.456), config.Float32)
	assert.Equal(t, float64(1234.5678), config.Float64)
	assert.Equal(t, "asdf", config.String)

	expectedTime, err := time.Parse(time.RFC3339, "2000-01-02T03:04:05Z")

	assert.Equal(t, expectedTime, config.Time)
	assert.Equal(t, time.Duration(12345), config.Duration)
}

func TestReadingMultipleStrings(t *testing.T) {
	configString1 :=
		`{
  			"Foo": "asdf",
  			"Bar": 1234,
  			"Baz": true
		}`
	configString2 :=
		`{
  			"Foo": "qwerty",
  			"Bar": 4321
		}`
	config := DefaultBasicConfig()

	err := configuration.ParseJsonStrings(&config, configString1, configString2)
	assert.NoError(t, err)

	assert.Equal(t, "qwerty", config.Foo)
	assert.Equal(t, 4321, config.Bar)
	assert.Equal(t, true, config.Baz)
}

// TODO maps
// TODO multiple config strings
