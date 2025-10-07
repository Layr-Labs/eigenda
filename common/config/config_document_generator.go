package config

import (
	"fmt"
	"go/types"
	"os"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

// A tag on struct fields used by this framework to generate documentation.
const DocsTag = "docs"

// Use this tag value to indicate that a field is required, e.g. `docs:"required"`.
// Note that this tag does not enforce that the field is actually required, it is only
// used for documentation generation.
const RequiredTag = "required"

// Use this tag value to indicate that a field is deprecated, e.g. `docs:"deprecated"`.
// Note that this tag does not enforce that the field is actually deprecated, it is only
// used for documentation generation. Fields that are deprecated will not show up in the
// "required" or "optional" lists in the generated documentation.
const DeprecatedTag = "deprecated"

// Use this tag value to indicate that a field is unsafe, e.g. `docs:"unsafe"`.
// Note that this tag does not enforce that the field is actually unsafe, it is only
// used for documentation generation. Fields that are unsafe will be listed in a
// separate "unsafe" section in the generated documentation.
const UnsafeTag = "unsafe"

// Generates documentation for a configuration struct by parsing the configuration. Output is deterministic.
func DocumentConfig[T any]( // TODO T VerifiableConfig
	// The name of the thing being documented, e.g. "Validator".
	componentName string,
	// The default constructor for the config struct. Default values will be extracted from the returned struct.
	constructor func() T,
	// The prefix for environment variables that will set fields in this config struct.
	envPrefix string,
	// A list of package import paths to search for the source file defining the config struct.
	// The package for the struct and the packages for any nested structs must be included in this list.
	packagePaths []string,
	// The path to write the generated markdown document to.
	outputPath string,
	// If true, fields without GoDoc comments will cause this method to return an error.
	requireDocs bool,
) error {

	if envPrefix == "" {
		return fmt.Errorf("envPrefix may not be empty")
	}
	defaultConfig := constructor()

	// TODO
	fmt.Printf("Default config: %+v\n", defaultConfig)

	// Unwrap pointer to get the named type
	t := reflect.TypeOf(defaultConfig)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Name() == "" {
		return fmt.Errorf("target type must be a named type, got %v", t)
	}

	fields, err := gatherConfigFieldData(defaultConfig, envPrefix, "", packagePaths)
	if err != nil {
		return fmt.Errorf("failed to gather config field data: %w", err)
	}

	if requireDocs {
		for _, f := range fields {
			if f.Deprecated {
				// Deprecated fields don't need docs
				continue
			}
			if f.Godoc == "" {
				return fmt.Errorf("field %q is missing GoDoc comments", f.TOML)
			}
		}
	}

	markdownString := generateMarkdownDoc(componentName, fields)

	if outputPath == "" {
		fmt.Println(markdownString)
		return nil
	}

	if err := os.WriteFile(outputPath, []byte(markdownString), 0o644); err != nil {
		return fmt.Errorf("failed to write config doc to %q: %w", outputPath, err)
	}

	return nil
}

// Find the file path and line number where the given type is defined, searching in the given package paths.
func findTypeDefLocation(packagePaths []string, t reflect.Type) (string, int, error) {
	for _, pkgPath := range packagePaths {
		if file, line, found, err := findInPackage(pkgPath, t); err != nil {
			return "", 0, fmt.Errorf("failed to search package %q: %w", pkgPath, err)
		} else if found {
			return file, line, nil
		}
	}

	return "", 0, fmt.Errorf("could not find source file for target type %s in provided package paths %v",
		t.String(), packagePaths)
}

// Look for the given type in the given package, returning its file and line number if found.
func findInPackage(pkgImportPath string, t reflect.Type) (string, int, bool, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedModule,
	}
	pkgs, err := packages.Load(cfg, pkgImportPath)
	if err != nil {
		return "", 0, false, fmt.Errorf("failed to load packages: %w", err)
	}
	if packages.PrintErrors(pkgs) > 0 || len(pkgs) == 0 {
		return "", 0, false, fmt.Errorf("failed to load package %q", pkgImportPath)
	}

	typeName := t.Name()
	wantPkgPath := t.PkgPath()

	for _, pkg := range pkgs {
		for _, obj := range pkg.TypesInfo.Defs {
			tn, ok := obj.(*types.TypeName)
			if !ok || tn == nil {
				continue
			}
			if tn.Name() != typeName {
				continue
			}
			// Check package path match for safety
			if obj.Pkg() == nil || obj.Pkg().Path() != wantPkgPath {
				continue
			}
			pos := pkg.Fset.Position(obj.Pos())
			return pos.Filename, pos.Line, true, nil
		}
	}

	return "", 0, false, nil
}

// Parse the fields of the struct for godocs for a struct defined at a specific line in a file.
func parseStructGodocs(filePath string, lineNumber int) (map[string]string, error) {

	fields := make(map[string]string)

	// Read the file.
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", filePath, err)
	}
	fileString := string(fileBytes)

	lines := strings.Split(fileString, "\n")
	if lineNumber < 1 || lineNumber > len(lines) {
		return nil, fmt.Errorf("line number %d out of range for file %q with %d lines",
			lineNumber, filePath, len(lines))
	}

	var godoc strings.Builder

	// Search for fields starting from the given line number (which should be the line where the struct is defined).
	for i := lineNumber - 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			// Skip blank lines, but reset the GoDoc accumulator. We should assume blank lines mean that the prior
			// GoDoc comments are not associated with the next field.
			godoc.Reset()
			continue
		}

		if strings.Contains(line, "}") {
			// Anonymous (i.e. nested) structs are prohibited, so we can assume this is the end of the struct.
			break
		}

		if strings.HasPrefix(line, "//") {
			// Accumulate GoDoc comments for the next field.
			if godoc.Len() > 0 {
				godoc.WriteString("\n")
			}
			godoc.WriteString(strings.TrimSpace(strings.TrimPrefix(line, "//")))
			continue
		}

		// We've found a line that isn't a comment or blank line, so it should be a struct field.
		// Extract the field name and the accumulated GoDoc comments.

		godocString := godoc.String()
		godoc.Reset()

		parts := strings.Split(line, " ")
		if len(parts) < 2 {
			return nil, fmt.Errorf("failed to parse struct field from line %q in file %q", line, filePath)
		}

		fieldName := strings.TrimSpace(parts[0])
		if fieldName == "" {
			return nil, fmt.Errorf("failed to parse struct field from line %q in file %q", line, filePath)
		}

		fields[fieldName] = godocString
	}

	return fields, nil
}

// All the data needed to document a config field.
type configFieldData struct {
	// Name of the environment variable that will set this field.
	EnvVar string
	// The toml tag that will set this field.
	TOML string
	// Type of the field as a string.
	FieldType string
	// The default value of the field as a string.
	DefaultValue string
	// GoDoc comment associated with the field.
	Godoc string

	// If true, this field is required.
	Required bool
	// If true, this field is deprecated.
	Deprecated bool
	// If true, this field is unsafe.
	Unsafe bool
}

// parseDocsTag parses the `docs` struct tag and returns whether the field is required, deprecated, or unsafe.
// Only one tag value is allowed per field.
func parseDocsTag(tag string) (required bool, deprecated bool, unsafe bool) {
	if tag == "" {
		return false, false, false
	}

	// Trim whitespace for flexibility
	tag = strings.TrimSpace(tag)

	switch tag {
	case RequiredTag:
		required = true
	case DeprecatedTag:
		deprecated = true
	case UnsafeTag:
		unsafe = true
	}
	return required, deprecated, unsafe
}

func gatherConfigFieldData(
	target any,
	envVarPrefix string,
	tomlPrefix string,
	packagePaths []string,
) ([]*configFieldData, error) {

	// Handle pointer to struct
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() == reflect.Ptr {
		targetValue = targetValue.Elem()
	}
	targetType := targetValue.Type()

	// Find the source file and line number where the target type is defined.
	structFile, line, err := findTypeDefLocation(packagePaths, targetType)
	if err != nil {
		return nil, fmt.Errorf("failed to find source file for target type %T: %w", target, err)
	}

	// Extract GoDoc comments for the struct fields.
	godocs, err := parseStructGodocs(structFile, line)
	if err != nil {
		return nil, fmt.Errorf("failed to parse struct godocs: %w", err)
	}

	var fields []*configFieldData

	// For each field in the struct, gather its data.
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		if field.PkgPath != "" { // unexported
			continue
		}

		switch field.Type.Kind() { //nolint:exhaustive // only handling struct and pointer types

		case reflect.Struct:
			// Recurse for nested structs
			tmp := reflect.New(field.Type).Elem().Interface()
			nestedEnvVarPrefix := envVarPrefix + "_" + strings.ToUpper(field.Name)

			var nestedTomlPrefix string
			if tomlPrefix == "" {
				nestedTomlPrefix = field.Name
			} else {
				nestedTomlPrefix = tomlPrefix + "." + field.Name
			}

			nestedFieldData, err := gatherConfigFieldData(tmp, nestedEnvVarPrefix, nestedTomlPrefix, packagePaths)
			if err != nil {
				return nil, fmt.Errorf("failed to gather field data for field %s: %w", field.Name, err)
			}
			fields = append(fields, nestedFieldData...)
		case reflect.Ptr:
			// Handle pointer to struct
			if field.Type.Elem().Kind() == reflect.Struct {
				tmp := reflect.New(field.Type.Elem()).Interface()

				nestedEnvVarPrefix := envVarPrefix + "_" + strings.ToUpper(field.Name)

				var nestedTomlPrefix string
				if tomlPrefix == "" {
					nestedTomlPrefix = field.Name
				} else {
					nestedTomlPrefix = tomlPrefix + "." + field.Name
				}

				nestedFieldData, err := gatherConfigFieldData(tmp, nestedEnvVarPrefix, nestedTomlPrefix, packagePaths)
				if err != nil {
					return nil, fmt.Errorf("failed to gather field data for field %s: %w", field.Name, err)
				}
				fields = append(fields, nestedFieldData...)
			} else {
				// Pointer to non-struct type, treat as regular field.
				// TODO be sure to unit test this

				var toml string
				if tomlPrefix == "" {
					toml = field.Name
				} else {
					toml = tomlPrefix + "." + field.Name
				}

				docsTag := field.Tag.Get("docs")
				required, deprecated, unsafe := parseDocsTag(docsTag)

				fields = append(fields, &configFieldData{
					EnvVar:       envVarPrefix + "_" + strings.ToUpper(field.Name),
					TOML:         toml,
					FieldType:    field.Type.String(),
					DefaultValue: fmt.Sprintf("%v", targetValue.Field(i).Interface()),
					Godoc:        godocs[field.Name],
					Required:     required,
					Deprecated:   deprecated,
					Unsafe:       unsafe,
				})
			}
		default:
			// Regular field

			var toml string
			if tomlPrefix == "" {
				toml = field.Name
			} else {
				toml = tomlPrefix + "." + field.Name
			}

			docsTag := field.Tag.Get("docs")
			required, deprecated, unsafe := parseDocsTag(docsTag)

			fields = append(fields, &configFieldData{
				EnvVar:       envVarPrefix + "_" + strings.ToUpper(field.Name),
				TOML:         toml,
				FieldType:    field.Type.String(),
				DefaultValue: fmt.Sprintf("%v", targetValue.Field(i).Interface()),
				Godoc:        godocs[field.Name],
				Required:     required,
				Deprecated:   deprecated,
				Unsafe:       unsafe,
			})
		}
	}

	// Alphabetically sort fields by for deterministic output.
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].TOML < fields[j].TOML
	})

	return fields, nil
}

func generateMarkdownDoc(
	componentName string,
	fields []*configFieldData,
) string {

	var sb strings.Builder

	// Sort fields into required, optional, and unsafe lists.
	requiredFields := make([]*configFieldData, 0)
	optionalFields := make([]*configFieldData, 0)
	unsafeFields := make([]*configFieldData, 0)
	for _, f := range fields {
		if f.Deprecated {
			// Deprecated fields are not documented.
			continue
		}
		if f.Unsafe {
			unsafeFields = append(unsafeFields, f)
		} else if f.Required {
			requiredFields = append(requiredFields, f)
		} else {
			optionalFields = append(optionalFields, f)
		}
	}

	// Write the markdown document.

	sb.WriteString("<!-- Code generated by config_document_generator.go. DO NOT EDIT BY HAND. -->\n\n")

	sb.WriteString(fmt.Sprintf("# %s Configuration\n\n", componentName))

	if len(requiredFields) > 0 {
		sb.WriteString("## Required Fields\n\n")
		sb.WriteString("| TOML | Environment Variable | Type | Description |\n")
		sb.WriteString("|------|----------------------|------|-------------|\n")

		for _, f := range requiredFields {
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				escapeMarkdownTableCell(f.TOML),
				escapeMarkdownTableCell(f.EnvVar),
				escapeMarkdownTableCell(f.FieldType),
				escapeMarkdownTableCell(f.Godoc)))
		}
	}

	if len(optionalFields) > 0 {
		sb.WriteString("\n## Optional Fields\n\n")
		sb.WriteString("| TOML | Environment Variable | Type | Default | Description |\n")
		sb.WriteString("|------|----------------------|------|---------|-------------|\n")

		for _, f := range optionalFields {
			defaultString := f.DefaultValue
			if f.FieldType == "string" {
				defaultString = fmt.Sprintf(`"%s"`, f.DefaultValue)
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				escapeMarkdownTableCell(f.TOML),
				escapeMarkdownTableCell(f.EnvVar),
				escapeMarkdownTableCell(f.FieldType),
				escapeMarkdownTableCell(defaultString),
				escapeMarkdownTableCell(f.Godoc)))
		}

	}

	if len(unsafeFields) > 0 {
		sb.WriteString("\n## Unsafe Fields\n\n")
		sb.WriteString("These fields are generally unsafe to modify unless you know what you are doing.\n\n")
		sb.WriteString("| TOML | Environment Variable | Type | Default | Description |\n")
		sb.WriteString("|------|----------------------|------|---------|-------------|\n")

		for _, f := range unsafeFields {
			defaultString := f.DefaultValue
			if f.FieldType == "string" {
				defaultString = fmt.Sprintf(`"%s"`, f.DefaultValue)
			}

			sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				escapeMarkdownTableCell(f.TOML),
				escapeMarkdownTableCell(f.EnvVar),
				escapeMarkdownTableCell(f.FieldType),
				escapeMarkdownTableCell(defaultString),
				escapeMarkdownTableCell(f.Godoc)))
		}
	}

	return sb.String()
}

// escapeMarkdownTableCell escapes special characters in markdown table cells.
// This handles pipes, newlines, and backslashes which have special meaning in markdown.
func escapeMarkdownTableCell(s string) string {
	var sb strings.Builder
	for _, r := range s {
		switch r {
		case '|':
			// Escape pipe characters which are table delimiters
			sb.WriteString("\\|")
		case '\n':
			// Replace newlines with <br> for markdown line breaks within table cells
			sb.WriteString("<br>")
		case '\r':
			// Skip carriage returns
			continue
		case '\\':
			// Escape backslashes
			sb.WriteString("\\\\")
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
