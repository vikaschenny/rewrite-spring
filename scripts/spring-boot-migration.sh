#!/bin/bash

# Spring Boot Migration Script using rewrite-spring-go
# Usage: ./spring-boot-migration.sh <project-path> <from-version> <to-version>

PROJECT_PATH="$1"
FROM_VERSION="$2"
TO_VERSION="$3"

if [ -z "$PROJECT_PATH" ] || [ -z "$FROM_VERSION" ] || [ -z "$TO_VERSION" ]; then
    echo "Usage: $0 <project-path> <from-version> <to-version>"
    echo ""
    echo "Examples:"
    echo "  $0 ./my-spring-app 2.7 3.0"
    echo "  $0 /path/to/project 1.5 2.0"
    echo "  $0 ./microservice 2.0 2.1"
    exit 1
fi

echo "=== Spring Boot Migration: $FROM_VERSION → $TO_VERSION ==="
echo "Project: $PROJECT_PATH"
echo ""

# Build the tool if it doesn't exist
if [ ! -f "rewrite-spring-go" ]; then
    echo "Building rewrite-spring-go..."
    go build -o rewrite-spring-go cmd/rewrite-spring/main.go
    if [ $? -ne 0 ]; then
        echo "Failed to build rewrite-spring-go"
        exit 1
    fi
fi

# Function to run migration with logging
run_migration() {
    local old_key="$1"
    local new_key="$2"
    local description="$3"
    
    echo "🔄 $description"
    echo "   $old_key → $new_key"
    
    ./rewrite-spring-go -source "$PROJECT_PATH" -recipe change-property-key \
        -old-key "$old_key" \
        -new-key "$new_key" \
        -backup
    
    if [ $? -eq 0 ]; then
        echo "   ✅ Complete"
    else
        echo "   ❌ Failed"
    fi
    echo ""
}

# Function to add property
add_property() {
    local property="$1"
    local value="$2"
    local comment="$3"
    local description="$4"
    
    echo "➕ $description"
    echo "   Adding: $property = $value"
    
    ./rewrite-spring-go -source "$PROJECT_PATH" -recipe add-property \
        -property "$property" \
        -value "$value" \
        -comment "$comment"
    
    if [ $? -eq 0 ]; then
        echo "   ✅ Complete"
    else
        echo "   ❌ Failed"
    fi
    echo ""
}

# Spring Boot 2.7 → 3.0 Migration
if [ "$FROM_VERSION" == "2.7" ] && [ "$TO_VERSION" == "3.0" ]; then
    echo "Applying Spring Boot 2.7 → 3.0 migrations..."
    echo ""
    
    # Management/Actuator properties
    run_migration "management.metrics.binders.jvm.enabled" \
                  "management.metrics.enable.jvm" \
                  "Migrating JVM metrics binder property"
    
    run_migration "management.metrics.binders.logback.enabled" \
                  "management.metrics.enable.logback" \
                  "Migrating Logback metrics binder property"
    
    run_migration "management.metrics.binders.tomcat.enabled" \
                  "management.metrics.enable.tomcat" \
                  "Migrating Tomcat metrics binder property"
    
    # Server properties
    run_migration "server.servlet.context-path" \
                  "server.servlet.contextPath" \
                  "Migrating servlet context path property"
    
    # Database properties
    run_migration "spring.datasource.hikari.maximum-pool-size" \
                  "spring.datasource.hikari.max-pool-size" \
                  "Migrating Hikari maximum pool size"
    
    run_migration "spring.datasource.hikari.minimum-idle" \
                  "spring.datasource.hikari.min-idle" \
                  "Migrating Hikari minimum idle"
    
    # Add new Spring Boot 3.x properties
    add_property "management.observations.annotations.enabled" \
                 "true" \
                 "Enable observation annotations for Spring Boot 3.x" \
                 "Adding Spring Boot 3.x observability support"

# Spring Boot 1.x → 2.x Migration
elif [ "$FROM_VERSION" == "1.5" ] && [ "$TO_VERSION" == "2.0" ]; then
    echo "Applying Spring Boot 1.5 → 2.0 migrations..."
    echo ""
    
    # Actuator endpoints
    run_migration "endpoints.health.path" \
                  "management.endpoints.web.path-mapping.health" \
                  "Migrating health endpoint path"
    
    run_migration "endpoints.info.path" \
                  "management.endpoints.web.path-mapping.info" \
                  "Migrating info endpoint path"
    
    run_migration "management.security.enabled" \
                  "management.endpoints.web.exposure.include" \
                  "Migrating actuator security settings"
    
    # Database properties
    run_migration "spring.datasource.initialize" \
                  "spring.datasource.initialization-mode" \
                  "Migrating datasource initialization"
    
    run_migration "spring.jpa.hibernate.naming-strategy" \
                  "spring.jpa.hibernate.naming.strategy" \
                  "Migrating Hibernate naming strategy"

# Spring Boot 2.0 → 2.1 Migration
elif [ "$FROM_VERSION" == "2.0" ] && [ "$TO_VERSION" == "2.1" ]; then
    echo "Applying Spring Boot 2.0 → 2.1 migrations..."
    echo ""
    
    # Security properties
    run_migration "security.oauth2.resource.jwt.key-value" \
                  "spring.security.oauth2.resourceserver.jwt.key-value" \
                  "Migrating OAuth2 JWT key value"
    
    # Actuator properties
    add_property "management.endpoints.web.exposure.include" \
                 "health,info" \
                 "Default actuator endpoints for Spring Boot 2.1" \
                 "Adding default actuator endpoint exposure"

else
    echo "❌ Unsupported migration path: $FROM_VERSION → $TO_VERSION"
    echo ""
    echo "Supported migrations:"
    echo "  - 1.5 → 2.0"
    echo "  - 2.0 → 2.1" 
    echo "  - 2.7 → 3.0"
    echo ""
    echo "For custom migrations, use the tool directly:"
    echo "  ./rewrite-spring-go -source $PROJECT_PATH -recipe change-property-key \\"
    echo "    -old-key 'old.property' -new-key 'new.property'"
    exit 1
fi

echo "🎉 Migration completed!"
echo ""
echo "📋 Summary:"
echo "   - Project: $PROJECT_PATH"
echo "   - Migration: Spring Boot $FROM_VERSION → $TO_VERSION"
echo "   - Backup files created with .backup extension"
echo ""
echo "🔍 Next steps:"
echo "   1. Review the changes in your version control"
echo "   2. Update your build files (pom.xml/build.gradle)"
echo "   3. Test your application thoroughly"
echo "   4. Check for additional manual migration steps in Spring Boot docs"
echo ""
echo "📝 To verify changes:"
echo "   ./rewrite-spring-go -source $PROJECT_PATH -recipe change-property-key \\"
echo "     -old-key 'any.old.property' -new-key 'any.new.property' -dry-run" 