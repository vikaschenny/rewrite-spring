package recipes

import (
	"context"
	"regexp"
	"strings"

	"github.com/openrewrite/rewrite-spring-go/pkg/core"
	"github.com/openrewrite/rewrite-spring-go/pkg/utils"
)

// ChangeSpringPropertyKeyRecipe changes Spring property keys in configuration files
type ChangeSpringPropertyKeyRecipe struct {
	core.BaseRecipe
	OldPropertyKey  string
	NewPropertyKey  string
	Except          []string
	PathExpressions []string
}

// NewChangeSpringPropertyKeyRecipe creates a new ChangeSpringPropertyKey recipe
func NewChangeSpringPropertyKeyRecipe(oldKey, newKey string, except []string) *ChangeSpringPropertyKeyRecipe {
	return &ChangeSpringPropertyKeyRecipe{
		BaseRecipe: core.BaseRecipe{
			DisplayName: "Change the key of a Spring application property",
			Description: "Change Spring application property keys existing in either Properties or YAML files, and in @Value annotations.",
		},
		OldPropertyKey: oldKey,
		NewPropertyKey: newKey,
		Except:         except,
		PathExpressions: []string{
			"**/application.properties",
			"**/application.yml",
			"**/application.yaml",
		},
	}
}

// Apply executes the recipe on the provided source file
func (r *ChangeSpringPropertyKeyRecipe) Apply(ctx context.Context, sourceFile core.SourceFile) (core.SourceFile, error) {
	if !r.shouldProcessFile(sourceFile.GetPath()) {
		return sourceFile, nil
	}

	switch sourceFile.GetType() {
	case core.Properties:
		return r.applyToProperties(sourceFile)
	case core.YAML:
		return r.applyToYAML(sourceFile)
	case core.Java:
		return r.applyToJava(sourceFile)
	default:
		return sourceFile, nil
	}
}

// shouldProcessFile checks if the file should be processed based on path expressions
func (r *ChangeSpringPropertyKeyRecipe) shouldProcessFile(filePath string) bool {
	if len(r.PathExpressions) == 0 {
		return true
	}

	for _, pattern := range r.PathExpressions {
		if matched, _ := utils.MatchGlob(filePath, pattern); matched {
			return true
		}
	}
	return false
}

// applyToProperties applies the transformation to properties files
func (r *ChangeSpringPropertyKeyRecipe) applyToProperties(sourceFile core.SourceFile) (core.SourceFile, error) {
	content := sourceFile.GetContent()
	lines := strings.Split(content, "\n")
	modified := false

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if r.shouldTransformKey(key) {
			newKey := r.transformPropertyKey(key)
			if newKey != key {
				lines[i] = strings.Replace(lines[i], key+"=", newKey+"=", 1)
				modified = true
			}
		}
	}

	if modified {
		sourceFile.SetContent(strings.Join(lines, "\n"))
	}

	return sourceFile, nil
}

// applyToYAML applies the transformation to YAML files
func (r *ChangeSpringPropertyKeyRecipe) applyToYAML(sourceFile core.SourceFile) (core.SourceFile, error) {
	content := sourceFile.GetContent()

	// Simple YAML key replacement - in a full implementation, you'd use a proper YAML parser
	lines := strings.Split(content, "\n")
	modified := false

	for i, line := range lines {
		if strings.Contains(line, ":") && !strings.TrimSpace(line)[0:1] == "#" {
			// Extract the key part before the colon
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				keyPart := strings.TrimSpace(parts[0])

				// Handle nested YAML properties by building the full path
				indent := len(line) - len(strings.TrimLeft(line, " "))

				if r.shouldTransformKey(keyPart) {
					newKey := r.transformPropertyKey(keyPart)
					if newKey != keyPart {
						lines[i] = strings.Replace(line, keyPart+":", newKey+":", 1)
						modified = true
					}
				}
			}
		}
	}

	if modified {
		sourceFile.SetContent(strings.Join(lines, "\n"))
	}

	return sourceFile, nil
}

// applyToJava applies the transformation to Java files (for @Value annotations)
func (r *ChangeSpringPropertyKeyRecipe) applyToJava(sourceFile core.SourceFile) (core.SourceFile, error) {
	content := sourceFile.GetContent()

	// Pattern to match @Value("${property.key:defaultValue}")
	valuePattern := regexp.MustCompile(`@Value\("?\$\{([^}:]+)([^}]*)\}"?\)`)

	modified := false
	newContent := valuePattern.ReplaceAllStringFunc(content, func(match string) string {
		submatch := valuePattern.FindStringSubmatch(match)
		if len(submatch) >= 2 {
			propertyKey := submatch[1]
			remainder := ""
			if len(submatch) > 2 {
				remainder = submatch[2]
			}

			if r.shouldTransformKey(propertyKey) {
				newKey := r.transformPropertyKey(propertyKey)
				if newKey != propertyKey {
					modified = true
					return strings.Replace(match, propertyKey, newKey, 1)
				}
			}
		}
		return match
	})

	if modified {
		sourceFile.SetContent(newContent)
	}

	return sourceFile, nil
}

// shouldTransformKey checks if a key should be transformed
func (r *ChangeSpringPropertyKeyRecipe) shouldTransformKey(key string) bool {
	// Check if the key matches the old property key pattern
	if !strings.HasPrefix(key, r.OldPropertyKey) {
		return false
	}

	// Check exceptions
	if len(r.Except) > 0 {
		for _, except := range r.Except {
			if strings.Contains(key, r.OldPropertyKey+"."+except) {
				return false
			}
		}
	}

	return true
}

// transformPropertyKey transforms a property key from old to new
func (r *ChangeSpringPropertyKeyRecipe) transformPropertyKey(key string) string {
	if strings.HasPrefix(key, r.OldPropertyKey) {
		return strings.Replace(key, r.OldPropertyKey, r.NewPropertyKey, 1)
	}
	return key
}
