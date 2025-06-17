package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/openrewrite/rewrite-spring-go/pkg/core"
	"github.com/openrewrite/rewrite-spring-go/pkg/recipes"
	"github.com/openrewrite/rewrite-spring-go/pkg/utils"
)

func main() {
	var (
		sourcePath  = flag.String("source", "", "Source directory to process")
		outputPath  = flag.String("output", "", "Output directory (optional, defaults to source)")
		recipe      = flag.String("recipe", "", "Recipe to apply (change-property-key, add-property)")
		oldKey      = flag.String("old-key", "", "Old property key (for change-property-key)")
		newKey      = flag.String("new-key", "", "New property key (for change-property-key)")
		property    = flag.String("property", "", "Property key (for add-property)")
		value       = flag.String("value", "", "Property value (for add-property)")
		comment     = flag.String("comment", "", "Comment for the property (optional)")
		exceptStr   = flag.String("except", "", "Comma-separated list of exceptions")
		patternsStr = flag.String("patterns", "", "Comma-separated list of file patterns")
		dryRun      = flag.Bool("dry-run", false, "Show what would be changed without modifying files")
		backup      = flag.Bool("backup", true, "Create backup files before modifying")
		debug       = flag.Bool("debug", false, "Enable debug logging")
		help        = flag.Bool("help", false, "Show help")
	)

	flag.Parse()

	if *help {
		showHelp()
		return
	}

	if *sourcePath == "" {
		fmt.Fprintf(os.Stderr, "Error: source path is required\n")
		os.Exit(1)
	}

	if *recipe == "" {
		fmt.Fprintf(os.Stderr, "Error: recipe is required\n")
		os.Exit(1)
	}

	// Set output path to source path if not specified
	if *outputPath == "" {
		*outputPath = *sourcePath
	}

	// Create logger
	logger := core.NewConsoleLogger(*debug)

	// Create execution context
	ctx := context.Background()

	// Parse exceptions
	var except []string
	if *exceptStr != "" {
		except = strings.Split(*exceptStr, ",")
		for i, ex := range except {
			except[i] = strings.TrimSpace(ex)
		}
	}

	// Parse patterns
	var patterns []string
	if *patternsStr != "" {
		patterns = strings.Split(*patternsStr, ",")
		for i, pattern := range patterns {
			patterns[i] = strings.TrimSpace(pattern)
		}
	}

	// Create recipe
	var recipeInstance core.Recipe
	switch *recipe {
	case "change-property-key":
		if *oldKey == "" || *newKey == "" {
			fmt.Fprintf(os.Stderr, "Error: old-key and new-key are required for change-property-key recipe\n")
			os.Exit(1)
		}
		recipeInstance = recipes.NewChangeSpringPropertyKeyRecipe(*oldKey, *newKey, except)
	case "add-property":
		if *property == "" || *value == "" {
			fmt.Fprintf(os.Stderr, "Error: property and value are required for add-property recipe\n")
			os.Exit(1)
		}
		recipeInstance = recipes.NewAddSpringPropertyRecipe(*property, *value, *comment, patterns)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown recipe '%s'\n", *recipe)
		os.Exit(1)
	}

	logger.Info("Starting rewrite-spring-go")
	logger.Info("Recipe: %s", recipeInstance.GetDisplayName())
	logger.Debug("Source: %s", *sourcePath)
	logger.Debug("Output: %s", *outputPath)

	// Find configuration files
	defaultPatterns := []string{
		"**/application.properties",
		"**/application.yml",
		"**/application.yaml",
		"**/*.java",
	}

	if len(patterns) == 0 {
		patterns = defaultPatterns
	}

	configFiles, err := utils.FindSpringConfigFiles(*sourcePath, patterns)
	if err != nil {
		logger.Error("Failed to find configuration files: %v", err)
		os.Exit(1)
	}

	logger.Info("Found %d configuration files", len(configFiles))

	// Process each file
	modifiedCount := 0
	for _, filePath := range configFiles {
		logger.Debug("Processing file: %s", filePath)

		// Load source file
		sourceFile, err := utils.LoadSourceFile(filePath)
		if err != nil {
			logger.Error("Failed to load file %s: %v", filePath, err)
			continue
		}

		// Apply recipe
		originalContent := sourceFile.GetContent()
		transformedFile, err := recipeInstance.Apply(ctx, sourceFile)
		if err != nil {
			logger.Error("Failed to apply recipe to %s: %v", filePath, err)
			continue
		}

		// Check if file was modified
		if transformedFile.GetContent() != originalContent {
			modifiedCount++
			logger.Info("Modified: %s", filePath)

			if *dryRun {
				logger.Info("DRY RUN: Would modify file %s", filePath)
				continue
			}

			// Create backup if requested
			if *backup {
				if err := utils.BackupFile(filePath); err != nil {
					logger.Warn("Failed to create backup for %s: %v", filePath, err)
				}
			}

			// Determine output file path
			outputFilePath := filePath
			if *outputPath != *sourcePath {
				relPath, err := filepath.Rel(*sourcePath, filePath)
				if err != nil {
					logger.Error("Failed to get relative path for %s: %v", filePath, err)
					continue
				}
				outputFilePath = filepath.Join(*outputPath, relPath)
			}

			// Save transformed file
			if err := utils.SaveSourceFile(transformedFile, outputFilePath); err != nil {
				logger.Error("Failed to save file %s: %v", outputFilePath, err)
				continue
			}

			logger.Debug("Saved transformed file: %s", outputFilePath)
		}
	}

	if *dryRun {
		logger.Info("DRY RUN completed. %d files would be modified", modifiedCount)
	} else {
		logger.Info("Processing completed. %d files modified", modifiedCount)
	}
}

func showHelp() {
	fmt.Println("rewrite-spring-go - Spring configuration transformation tool")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  rewrite-spring-go [OPTIONS]")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -source string")
	fmt.Println("        Source directory to process (required)")
	fmt.Println("  -output string")
	fmt.Println("        Output directory (optional, defaults to source)")
	fmt.Println("  -recipe string")
	fmt.Println("        Recipe to apply: change-property-key, add-property (required)")
	fmt.Println("  -old-key string")
	fmt.Println("        Old property key (required for change-property-key)")
	fmt.Println("  -new-key string")
	fmt.Println("        New property key (required for change-property-key)")
	fmt.Println("  -property string")
	fmt.Println("        Property key (required for add-property)")
	fmt.Println("  -value string")
	fmt.Println("        Property value (required for add-property)")
	fmt.Println("  -comment string")
	fmt.Println("        Comment for the property (optional)")
	fmt.Println("  -except string")
	fmt.Println("        Comma-separated list of exceptions")
	fmt.Println("  -patterns string")
	fmt.Println("        Comma-separated list of file patterns")
	fmt.Println("  -dry-run")
	fmt.Println("        Show what would be changed without modifying files")
	fmt.Println("  -backup")
	fmt.Println("        Create backup files before modifying (default: true)")
	fmt.Println("  -debug")
	fmt.Println("        Enable debug logging")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  # Change property key in all Spring config files")
	fmt.Println("  rewrite-spring-go -source ./myproject -recipe change-property-key \\")
	fmt.Println("    -old-key management.metrics.binders.*.enabled \\")
	fmt.Println("    -new-key management.metrics.enable.process.files")
	fmt.Println()
	fmt.Println("  # Add a new property to Spring config files")
	fmt.Println("  rewrite-spring-go -source ./myproject -recipe add-property \\")
	fmt.Println("    -property server.port -value 8080 -comment \"Server port configuration\"")
	fmt.Println()
	fmt.Println("  # Dry run to see what would be changed")
	fmt.Println("  rewrite-spring-go -source ./myproject -recipe change-property-key \\")
	fmt.Println("    -old-key old.property -new-key new.property -dry-run")
}
