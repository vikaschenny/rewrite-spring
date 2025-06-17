# Spring Boot Migration Guide using rewrite-spring-go

This guide shows how to use `rewrite-spring-go` to migrate between Spring Boot versions by automating configuration changes.

## Spring Boot 2.x to 3.x Migration

### Common Property Changes

#### 1. Management/Actuator Properties

```bash
# Change deprecated metrics binder properties
./rewrite-spring-go -source ./myproject -recipe change-property-key \
  -old-key "management.metrics.binders.jvm.enabled" \
  -new-key "management.metrics.enable.jvm"

./rewrite-spring-go -source ./myproject -recipe change-property-key \
  -old-key "management.metrics.binders.logback.enabled" \
  -new-key "management.metrics.enable.logback"

./rewrite-spring-go -source ./myproject -recipe change-property-key \
  -old-key "management.metrics.binders.tomcat.enabled" \
  -new-key "management.metrics.enable.tomcat"
```

#### 2. Server Properties

```bash
# Change servlet context path property
./rewrite-spring-go -source ./myproject -recipe change-property-key \
  -old-key "server.servlet.context-path" \
  -new-key "server.servlet.contextPath"

# Change server error properties
./rewrite-spring-go -source ./myproject -recipe change-property-key \
  -old-key "server.error.include-exception" \
  -new-key "server.error.include-exception"
```

#### 3. Database Properties

```bash
# Change Hikari connection pool properties
./rewrite-spring-go -source ./myproject -recipe change-property-key \
  -old-key "spring.datasource.hikari.maximum-pool-size" \
  -new-key "spring.datasource.hikari.max-pool-size"

./rewrite-spring-go -source ./myproject -recipe change-property-key \
  -old-key "spring.datasource.hikari.minimum-idle" \
  -new-key "spring.datasource.hikari.min-idle"
```

#### 4. Logging Properties

```bash
# Change logging pattern properties
./rewrite-spring-go -source ./myproject -recipe change-property-key \
  -old-key "logging.pattern.dateformat" \
  -new-key "logging.pattern.dateFormat"
```

### Adding New Spring Boot 3.x Properties

```bash
# Add new observability properties
./rewrite-spring-go -source ./myproject -recipe add-property \
  -property "management.observations.annotations.enabled" \
  -value "true" \
  -comment "Enable observation annotations"

# Add new security properties
./rewrite-spring-go -source ./myproject -recipe add-property \
  -property "spring.security.filter.dispatcher-types" \
  -value "request,error,async" \
  -comment "Configure security filter dispatcher types"
```

## Spring Boot 1.x to 2.x Migration

### Actuator Endpoints

```bash
# Change actuator endpoint paths
./rewrite-spring-go -source ./myproject -recipe change-property-key \
  -old-key "endpoints.health.path" \
  -new-key "management.endpoints.web.path-mapping.health"

./rewrite-spring-go -source ./myproject -recipe change-property-key \
  -old-key "endpoints.info.path" \
  -new-key "management.endpoints.web.path-mapping.info"

# Change actuator security
./rewrite-spring-go -source ./myproject -recipe change-property-key \
  -old-key "management.security.enabled" \
  -new-key "management.endpoints.web.exposure.include"
```

### Database Configuration

```bash
# Change database initialization
./rewrite-spring-go -source ./myproject -recipe change-property-key \
  -old-key "spring.datasource.initialize" \
  -new-key "spring.datasource.initialization-mode"

./rewrite-spring-go -source ./myproject -recipe change-property-key \
  -old-key "spring.jpa.hibernate.naming-strategy" \
  -new-key "spring.jpa.hibernate.naming.strategy"
```

## Batch Migration Script

Create a comprehensive migration script:

```bash
#!/bin/bash

echo "Starting Spring Boot migration..."

PROJECT_PATH="$1"
if [ -z "$PROJECT_PATH" ]; then
    echo "Usage: $0 <project-path>"
    exit 1
fi

# Build the tool
go build -o rewrite-spring-go cmd/rewrite-spring/main.go

echo "Migrating management properties..."
./rewrite-spring-go -source "$PROJECT_PATH" -recipe change-property-key \
  -old-key "management.metrics.binders.jvm.enabled" \
  -new-key "management.metrics.enable.jvm"

./rewrite-spring-go -source "$PROJECT_PATH" -recipe change-property-key \
  -old-key "management.metrics.binders.logback.enabled" \
  -new-key "management.metrics.enable.logback"

./rewrite-spring-go -source "$PROJECT_PATH" -recipe change-property-key \
  -old-key "management.metrics.binders.tomcat.enabled" \
  -new-key "management.metrics.enable.tomcat"

echo "Migrating server properties..."
./rewrite-spring-go -source "$PROJECT_PATH" -recipe change-property-key \
  -old-key "server.servlet.context-path" \
  -new-key "server.servlet.contextPath"

echo "Migrating database properties..."
./rewrite-spring-go -source "$PROJECT_PATH" -recipe change-property-key \
  -old-key "spring.datasource.hikari.maximum-pool-size" \
  -new-key "spring.datasource.hikari.max-pool-size"

echo "Adding new Spring Boot 3.x properties..."
./rewrite-spring-go -source "$PROJECT_PATH" -recipe add-property \
  -property "management.observations.annotations.enabled" \
  -value "true" \
  -comment "Enable observation annotations for Spring Boot 3.x"

echo "Migration completed!"
```

## Migration Verification

After running migrations, verify the changes:

```bash
# Dry run to see what would change
./rewrite-spring-go -source ./myproject -recipe change-property-key \
  -old-key "management.metrics.binders.jvm.enabled" \
  -new-key "management.metrics.enable.jvm" \
  -dry-run

# Check for any remaining deprecated properties
grep -r "management.metrics.binders" ./myproject/src/main/resources/
grep -r "server.servlet.context-path" ./myproject/src/main/resources/
```

## Advanced Migration Patterns

### 1. Conditional Migrations with Exceptions

```bash
# Migrate all metrics binder properties except specific ones
./rewrite-spring-go -source ./myproject -recipe change-property-key \
  -old-key "management.metrics.binders" \
  -new-key "management.metrics.enable" \
  -except "custom,legacy"
```

### 2. Profile-Specific Migrations

```bash
# Migrate only specific profile configurations
./rewrite-spring-go -source ./myproject -recipe change-property-key \
  -old-key "spring.datasource.url" \
  -new-key "spring.datasource.jdbc-url" \
  -patterns "**/application-prod.properties,**/application-prod.yml"
```

### 3. Java Code Migrations

The tool also updates `@Value` annotations in Java files:

```java
// Before migration
@Value("${management.metrics.binders.jvm.enabled:true}")
private boolean jvmMetricsEnabled;

// After migration (automatically updated)
@Value("${management.metrics.enable.jvm:true}")
private boolean jvmMetricsEnabled;
```

## Best Practices

1. **Always use dry-run first**:
   ```bash
   ./rewrite-spring-go -source ./myproject -recipe change-property-key \
     -old-key "old.property" -new-key "new.property" -dry-run
   ```

2. **Create backups** (enabled by default):
   ```bash
   ./rewrite-spring-go -source ./myproject -recipe change-property-key \
     -old-key "old.property" -new-key "new.property" -backup
   ```

3. **Use version control** to track changes

4. **Test after migration** to ensure functionality

5. **Run migrations in stages** rather than all at once

## Common Spring Boot Version Migration Scenarios

### Spring Boot 2.0 → 2.1
- Actuator endpoint changes
- Security configuration updates

### Spring Boot 2.1 → 2.2  
- Configuration property relocations
- Deprecation removals

### Spring Boot 2.7 → 3.0
- Major property restructuring
- Jakarta EE namespace changes
- Observability additions

### Spring Boot 3.0 → 3.1
- New configuration options
- Enhanced security properties 