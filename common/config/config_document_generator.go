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

// Generates documentation for a configuration struct by parsing the configuration. Output is determinsitic.
func DocumentConfig[T any]( // TODO T VerifiableConfig
	constructor func() T,
	envPrefix string,
	packagePaths []string,
	outputPath string,
) error {

	if envPrefix == "" {
		return fmt.Errorf("envPrefix may not be empty")
	}
	defaultConfig := constructor()

	// Unwrap pointer to get the named type
	t := reflect.TypeOf(defaultConfig)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Name() == "" {
		return fmt.Errorf("target type must be a named type, got %v", t)
	}

	fields, err := gatherConfigFieldData(defaultConfig, envPrefix, packagePaths)
	if err != nil {
		return fmt.Errorf("failed to gather config field data: %w", err)
	}

	markdownString := generateMarkdownDoc(t.Name(), fields)

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
		return "", 0, false, err
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
	// Name of the field.
	FieldName string
	// Type of the field as a string.
	FieldType string
	// The default value of the field as a string.
	DefaultValue string
	// GoDoc comment associated with the field.
	Godoc string
}

func gatherConfigFieldData(
	target any,
	prefix string,
	packagePaths []string,
) ([]*configFieldData, error) {

	// Find the source file and line number where the target type is defined.
	structFile, line, err := findTypeDefLocation(packagePaths, reflect.TypeOf(target))
	if err != nil {
		return nil, fmt.Errorf("failed to find source file for target type %T: %w", target, err)
	}

	// Extract GoDoc comments for the struct fields.
	godocs, err := parseStructGodocs(structFile, line)
	if err != nil {
		return nil, fmt.Errorf("failed to parse struct godocs: %w", err)
	}

	var fields []*configFieldData

	// Handle pointer to struct
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() == reflect.Ptr {
		targetValue = targetValue.Elem()
	}
	targetType := targetValue.Type()

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
			nestedPrefix := prefix + "_" + strings.ToUpper(field.Name)
			nestedFieldData, err := gatherConfigFieldData(tmp, nestedPrefix, packagePaths)
			if err != nil {
				return nil, fmt.Errorf("failed to gather field data for field %s: %w", field.Name, err)
			}
			fields = append(fields, nestedFieldData...)
		case reflect.Ptr:
			// Handle pointer to struct
			if field.Type.Elem().Kind() == reflect.Struct {
				tmp := reflect.New(field.Type.Elem()).Interface()

				nestedPrefix := prefix + "_" + strings.ToUpper(field.Name)
				nestedFieldData, err := gatherConfigFieldData(tmp, nestedPrefix, packagePaths)
				if err != nil {
					return nil, fmt.Errorf("failed to gather field data for field %s: %w", field.Name, err)
				}
				fields = append(fields, nestedFieldData...)
			} else {
				// Pointer to non-struct type, treat as regular field.
				// TODO be sure to unit test this
				fields = append(fields, &configFieldData{
					FieldName:    prefix + "_" + strings.ToUpper(field.Name),
					FieldType:    field.Type.String(),
					DefaultValue: fmt.Sprintf("%v", targetValue.Field(i).Interface()),
					Godoc:        godocs[field.Name],
				})
			}
		default:
			fields = append(fields, &configFieldData{
				FieldName:    prefix + "_" + strings.ToUpper(field.Name),
				FieldType:    field.Type.String(),
				DefaultValue: fmt.Sprintf("%v", targetValue.Field(i).Interface()),
				Godoc:        godocs[field.Name],
			})
		}
	}

	// Alphabetically sort fields by FieldName for deterministic output.
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].FieldName < fields[j].FieldName
	})

	return fields, nil
}

func generateMarkdownDoc(
	componentName string,
	fields []*configFieldData,
) string {

	var sb strings.Builder

	sb.WriteString("<!-- Code generated by config_document_generator.go. DO NOT EDIT BY HAND. -->\n\n")

	sb.WriteString(fmt.Sprintf("# %s Configuration\n\n", componentName))
	sb.WriteString("Configuration is provided via environment variables.\n")
	sb.WriteString("| Environment Variable | Type | Default | Description |\n")
	sb.WriteString("|----------------------|------|---------|-------------|\n")

	for _, f := range fields {
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
			f.FieldName, f.FieldType, f.DefaultValue, f.Godoc))
	}
	return sb.String()
}
