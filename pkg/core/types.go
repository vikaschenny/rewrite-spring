package core

import (
	"context"
	"io"
)

// Recipe represents a transformation recipe that can be applied to source files
type Recipe interface {
	GetDisplayName() string
	GetDescription() string
	Apply(ctx context.Context, sourceFile SourceFile) (SourceFile, error)
}

// SourceFile represents a source file that can be transformed
type SourceFile interface {
	GetPath() string
	GetContent() string
	SetContent(content string)
	GetType() FileType
	Save(writer io.Writer) error
}

// FileType represents the type of configuration file
type FileType int

const (
	Properties FileType = iota
	YAML
	Java
)

// ExecutionContext provides context and configuration for recipe execution
type ExecutionContext struct {
	Options map[string]interface{}
	Logger  Logger
}

// Logger interface for logging recipe execution
type Logger interface {
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// PropertyTransformation represents a property key transformation
type PropertyTransformation struct {
	OldKey    string
	NewKey    string
	NewValue  string
	Comment   string
	Operation OperationType
}

// OperationType represents the type of operation to perform
type OperationType int

const (
	ChangeKey OperationType = iota
	ChangeValue
	AddProperty
	DeleteProperty
	CommentOut
)

// SpringConfigFile represents a Spring configuration file
type SpringConfigFile struct {
	Path    string
	Content string
	Type    FileType
}

// GetPath returns the file path
func (f *SpringConfigFile) GetPath() string {
	return f.Path
}

// GetContent returns the file content
func (f *SpringConfigFile) GetContent() string {
	return f.Content
}

// SetContent sets the file content
func (f *SpringConfigFile) SetContent(content string) {
	f.Content = content
}

// GetType returns the file type
func (f *SpringConfigFile) GetType() FileType {
	return f.Type
}

// Save writes the file content to the provided writer
func (f *SpringConfigFile) Save(writer io.Writer) error {
	_, err := writer.Write([]byte(f.Content))
	return err
}

// BaseRecipe provides common functionality for all recipes
type BaseRecipe struct {
	DisplayName string
	Description string
}

// GetDisplayName returns the recipe display name
func (r *BaseRecipe) GetDisplayName() string {
	return r.DisplayName
}

// GetDescription returns the recipe description
func (r *BaseRecipe) GetDescription() string {
	return r.Description
}
