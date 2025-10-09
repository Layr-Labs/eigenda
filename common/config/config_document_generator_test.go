package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToScreamingSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single word lowercase",
			input:    "field",
			expected: "FIELD",
		},
		{
			name:     "single word uppercase",
			input:    "FIELD",
			expected: "FIELD",
		},
		{
			name:     "camelCase",
			input:    "myFieldName",
			expected: "MY_FIELD_NAME",
		},
		{
			name:     "PascalCase",
			input:    "MyFieldName",
			expected: "MY_FIELD_NAME",
		},
		{
			name:     "HTTPServer - consecutive uppercase at start",
			input:    "HTTPServer",
			expected: "HTTP_SERVER",
		},
		{
			name:     "APIKey - consecutive uppercase at start",
			input:    "APIKey",
			expected: "API_KEY",
		},
		{
			name:     "ServerHTTP - consecutive uppercase at end",
			input:    "ServerHTTP",
			expected: "SERVER_HTTP",
		},
		{
			name:     "single character",
			input:    "X",
			expected: "X",
		},
		{
			name:     "with numbers",
			input:    "Field123Name",
			expected: "FIELD123_NAME",
		},
		{
			name:     "already snake_case",
			input:    "my_field_name",
			expected: "MY_FIELD_NAME",
		},
		{
			name:     "XMLParser - consecutive uppercase followed by word",
			input:    "XMLParser",
			expected: "XML_PARSER",
		},
		{
			name:     "MyYAMLParser - user example",
			input:    "MyYAMLParser",
			expected: "MY_YAML_PARSER",
		},
		{
			name:     "IPAddress - user example",
			input:    "IPAddress",
			expected: "IP_ADDRESS",
		},
		{
			name:     "URLPath",
			input:    "URLPath",
			expected: "URL_PATH",
		},
		{
			name:     "HTTPAPI",
			input:    "HTTPAPI",
			expected: "HTTPAPI",
		},
		{
			name:     "HTTPSConnection",
			input:    "HTTPSConnection",
			expected: "HTTPS_CONNECTION",
		},
		{
			name:     "two letter acronym",
			input:    "IOReader",
			expected: "IO_READER",
		},
		{
			name:     "NodeID - uppercase sequence at end",
			input:    "NodeID",
			expected: "NODE_ID",
		},
		{
			name:     "GetUUID - uppercase sequence at end",
			input:    "GetUUID",
			expected: "GET_UUID",
		},
		{
			name:     "single letter followed by uppercase sequence at end",
			input:    "AHTTP",
			expected: "AHTTP",
		},
		{
			name:     "RequestID",
			input:    "RequestID",
			expected: "REQUEST_ID",
		},
		{
			name:     "UserAPI - uppercase sequence at end",
			input:    "UserAPI",
			expected: "USER_API",
		},
		{
			name:     "MySQLDatabase - mixed case with uppercase sequence",
			input:    "MySQLDatabase",
			expected: "MY_SQL_DATABASE",
		},
		{
			name:     "Single Trailing lower case",
			input:    "EthRpcURLs",
			expected: "ETH_RPC_URLS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := toScreamingSnakeCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
