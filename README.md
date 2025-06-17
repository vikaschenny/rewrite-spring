# rewrite-spring-go

A Go port of the OpenRewrite Spring transformation recipes. This tool provides automated refactoring and migration capabilities for Spring/Spring Boot applications, focusing on configuration property transformations.

## Features

- **Change Property Keys**: Rename Spring property keys across `.properties`, `.yml`, `.yaml`, and Java `@Value` annotations
- **Add Properties**: Add new properties to Spring configuration files
- **Multiple File Format Support**: Works with Properties files, YAML files, and Java source files
- **Glob Pattern Matching**: Flexible file selection using glob patterns
- **Dry Run Mode**: Preview changes before applying them
- **Backup Creation**: Automatically create backups of modified files
- **Comprehensive Logging**: Debug mode for detailed operation logs

## Installation

### From Source

```bash
git clone https://github.com/openrewrite/rewrite-spring-go.git
cd rewrite-spring-go
go build -o bin/rewrite-spring-go cmd/rewrite-spring/main.go
```

### Using Go Install

```bash
go install github.com/openrewrite/rewrite-spring-go/cmd/rewrite-spring@latest
```

## Usage

### Basic Usage

```bash
rewrite-spring-go -source /path/to/project -recipe RECIPE_NAME [OPTIONS]
```

### Available Recipes

#### 1. Change Property Key

Changes property keys across all Spring configuration files and Java `@Value` annotations.

```bash
rewrite-spring-go -source ./myproject -recipe change-property-key \
  -old-key "management.metrics.binders.*.enabled" \
  -new-key "management.metrics.enable.process.files"
```

**Options:**
- `-old-key`: The property key to rename (required)
- `-new-key`: The new property key name (required)
- `-except`: Comma-separated list of exceptions (optional)

#### 2. Add Property

Adds new properties to Spring configuration files.

```bash
rewrite-spring-go -source ./myproject -recipe add-property \
  -property "server.port" \
  -value "8080" \
  -comment "Server port configuration"
```

**Options:**
- `-property`: The property key to add (required)
- `-value`: The property value (required)
- `-comment`: Optional comment for the property

### Common Options

- `-source`: Source directory to process (required)
- `-output`: Output directory (optional, defaults to source)
- `-patterns`: Comma-separated list of file patterns (optional)
- `-dry-run`: Show what would be changed without modifying files
- `-backup`: Create backup files before modifying (default: true)
- `-debug`: Enable debug logging
- `-help`: Show help message

### Examples

#### Change deprecated property keys

```bash
# Change deprecated Spring Boot 2.x property to 3.x equivalent
rewrite-spring-go -source ./spring-boot-app -recipe change-property-key \
  -old-key "spring.datasource.hikari.maximum-pool-size" \
  -new-key "spring.datasource.hikari.max-pool-size"
```

#### Add monitoring configuration

```bash
# Add actuator endpoint configuration
rewrite-spring-go -source ./microservice -recipe add-property \
  -property "management.endpoints.web.exposure.include" \
  -value "health,info,metrics" \
  -comment "Expose actuator endpoints"
```

#### Dry run to preview changes

```bash
# See what would be changed without modifying files
rewrite-spring-go -source ./my-app -recipe change-property-key \
  -old-key "old.deprecated.property" \
  -new-key "new.recommended.property" \
  -dry-run
```

#### Process specific file patterns

```bash
# Only process test configuration files
rewrite-spring-go -source ./project -recipe add-property \
  -property "spring.profiles.active" \
  -value "test" \
  -patterns "**/application-test.properties,**/application-test.yml"
```

## File Support

### Properties Files
- `application.properties`
- `application-{profile}.properties`
- `bootstrap.properties`
- Any `.properties` file matching the specified patterns

### YAML Files
- `application.yml`
- `application.yaml`
- `application-{profile}.yml`
- `application-{profile}.yaml`
- `bootstrap.yml`
- `bootstrap.yaml`

### Java Files
- Transformations in `@Value` annotations
- Support for `@ConditionalOnProperty` annotations
- Property references in string literals

## Architecture

The tool is structured with a modular architecture:

```
pkg/
├── core/           # Core interfaces and types
├── recipes/        # Transformation recipes
└── utils/          # Utility functions

cmd/
└── rewrite-spring/ # CLI application
```

### Key Components

- **Recipe Interface**: Plugin-based transformation system
- **SourceFile Interface**: Abstraction for different file types
- **Logger Interface**: Configurable logging system
- **Utility Functions**: File operations and pattern matching

## Differences from Java Version

This Go port provides a simplified but functional subset of the original Java rewrite-spring library:

### Similarities
- Core property transformation functionality
- Support for multiple file formats
- Dry run and backup capabilities
- Glob pattern matching

### Differences
- Simplified YAML processing (no full AST manipulation)
- Limited Java annotation support
- Fewer built-in recipes
- No integration with build tools (Maven/Gradle)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the same license as the original OpenRewrite project.

## Related Projects

- [OpenRewrite](https://github.com/openrewrite/rewrite) - The original Java-based refactoring tool
- [Rewrite Spring](https://github.com/openrewrite/rewrite-spring) - The original Java Spring transformation recipes

## Support

For issues and questions:
- Check the [GitHub Issues](https://github.com/openrewrite/rewrite-spring-go/issues)
- Refer to the [OpenRewrite Documentation](https://docs.openrewrite.org/)
- Join the [OpenRewrite Community](https://join.slack.com/t/rewriteoss/shared_invite/zt-nj42n3ea-b~62rIHzb3Vo0E1APKCXEA)
