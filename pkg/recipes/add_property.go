package recipes

import (
	"context"
	"strings"

	"github.com/openrewrite/rewrite-spring-go/pkg/core"
	"github.com/openrewrite/rewrite-spring-go/pkg/utils"
)

// AddSpringPropertyRecipe adds properties to Spring configuration files
type AddSpringPropertyRecipe struct {
	core.BaseRecipe
	Property        string
	Value           string
	Comment         string
	PathExpressions []string
}

// NewAddSpringPropertyRecipe creates a new AddSpringProperty recipe
func NewAddSpringPropertyRecipe(property, value, comment string, pathExpressions []string) *AddSpringPropertyRecipe {
	defaultPaths := []string{
		"**/application.properties",
		"**/application.yml",
		"**/application.yaml",
	}

	if len(pathExpressions) == 0 {
		pathExpressions = defaultPaths
	}

	return &AddSpringPropertyRecipe{
		BaseRecipe: core.BaseRecipe{
			DisplayName: "Add a spring configuration property",
			Description: "Add a spring configuration property to a configuration file if it does not already exist in that file.",
		},
		Property:        property,
		Value:           value,
		Comment:         comment,
		PathExpressions: pathExpressions,
	}
}

// Apply executes the recipe on the provided source file
func (r *AddSpringPropertyRecipe) Apply(ctx context.Context, sourceFile core.SourceFile) (core.SourceFile, error) {
	if !r.shouldProcessFile(sourceFile.GetPath()) {
		return sourceFile, nil
	}

	switch sourceFile.GetType() {
	case core.Properties:
		return r.addToProperties(sourceFile)
	case core.YAML:
		return r.addToYAML(sourceFile)
	default:
		return sourceFile, nil
	}
}

// shouldProcessFile checks if the file should be processed based on path expressions
func (r *AddSpringPropertyRecipe) shouldProcessFile(filePath string) bool {
	for _, pattern := range r.PathExpressions {
		if matched, _ := utils.MatchGlob(filePath, pattern); matched {
			return true
		}
	}
	return false
}

// addToProperties adds the property to a properties file
func (r *AddSpringPropertyRecipe) addToProperties(sourceFile core.SourceFile) (core.SourceFile, error) {
	content := sourceFile.GetContent()

	// Check if property already exists
	if r.propertyExistsInProperties(content) {
		return sourceFile, nil
	}

	var newContent strings.Builder
	newContent.WriteString(content)

	// Add newline if content doesn't end with one
	if !strings.HasSuffix(content, "\n") {
		newContent.WriteString("\n")
	}

	// Add comment if provided
	if r.Comment != "" {
		newContent.WriteString("# ")
		newContent.WriteString(r.Comment)
		newContent.WriteString("\n")
	}

	// Add the property
	newContent.WriteString(r.Property)
	newContent.WriteString("=")
	newContent.WriteString(r.Value)
	newContent.WriteString("\n")

	sourceFile.SetContent(newContent.String())
	return sourceFile, nil
}

// addToYAML adds the property to a YAML file
func (r *AddSpringPropertyRecipe) addToYAML(sourceFile core.SourceFile) (core.SourceFile, error) {
	content := sourceFile.GetContent()

	// Check if property already exists
	if r.propertyExistsInYAML(content) {
		return sourceFile, nil
	}

	// Convert property to YAML format
	yamlProperty := r.convertToYAML()

	var newContent strings.Builder
	newContent.WriteString(content)

	// Add newline if content doesn't end with one
	if !strings.HasSuffix(content, "\n") {
		newContent.WriteString("\n")
	}

	newContent.WriteString(yamlProperty)

	sourceFile.SetContent(newContent.String())
	return sourceFile, nil
}

// propertyExistsInProperties checks if the property already exists in properties file
func (r *AddSpringPropertyRecipe) propertyExistsInProperties(content string) bool {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		if strings.HasPrefix(line, r.Property+"=") {
			return true
		}
	}
	return false
}

// propertyExistsInYAML checks if the property already exists in YAML file
func (r *AddSpringPropertyRecipe) propertyExistsInYAML(content string) bool {
	// Simple check - in a full implementation, you'd use a proper YAML parser
	propertyParts := strings.Split(r.Property, ".")

	// Check for the leaf property key
	leafKey := propertyParts[len(propertyParts)-1]
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, leafKey+":") {
			return true
		}
	}
	return false
}

// convertToYAML converts a dot-notation property to YAML format
func (r *AddSpringPropertyRecipe) convertToYAML() string {
	propertyParts := strings.Split(r.Property, ".")

	var yamlBuilder strings.Builder
	indent := ""

	// Add comment if provided
	if r.Comment != "" {
		yamlBuilder.WriteString("# ")
		yamlBuilder.WriteString(r.Comment)
		yamlBuilder.WriteString("\n")
	}

	for i, part := range propertyParts {
		yamlBuilder.WriteString(indent)
		yamlBuilder.WriteString(part)
		yamlBuilder.WriteString(":")

		if i == len(propertyParts)-1 {
			// Last part - add the value
			yamlBuilder.WriteString(" ")
			if r.needsQuotes(r.Value) {
				yamlBuilder.WriteString("\"")
				yamlBuilder.WriteString(r.Value)
				yamlBuilder.WriteString("\"")
			} else {
				yamlBuilder.WriteString(r.Value)
			}
		}
		yamlBuilder.WriteString("\n")
		indent += "  "
	}

	return yamlBuilder.String()
}

// needsQuotes determines if a YAML value needs to be quoted
func (r *AddSpringPropertyRecipe) needsQuotes(value string) bool {
	// Simple heuristic - quote if contains special characters
	specialChars := []string{":", "[", "]", "{", "}", ",", "&", "*", "#", "?", "|", "-", "<", ">", "=", "!", "%", "@", "`"}

	for _, char := range specialChars {
		if strings.Contains(value, char) {
			return true
		}
	}

	// Quote if it looks like a number but should be treated as string
	if strings.Contains(value, ".") && len(strings.Split(value, ".")) > 2 {
		return true
	}

	return false
}
