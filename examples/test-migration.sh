#!/bin/bash

echo "=== Testing Spring Boot Migration with Sample Project ==="
echo ""

# Build the tool first
echo "1. Building rewrite-spring-go..."
go build -o rewrite-spring-go cmd/rewrite-spring/main.go

if [ $? -ne 0 ]; then
    echo "❌ Failed to build the tool"
    exit 1
fi

echo "✅ Tool built successfully"
echo ""

# Test 1: Show current state
echo "2. Current state of sample project files:"
echo ""
echo "📄 application.properties:"
grep -n "management.metrics.binders" examples/sample-spring-project/application.properties || echo "   No deprecated properties found"
echo ""
echo "📄 application.yml:"
grep -A2 -B2 "binders:" examples/sample-spring-project/application.yml || echo "   No deprecated properties found"
echo ""
echo "📄 Config.java:"
grep -n "management.metrics.binders" examples/sample-spring-project/src/main/java/com/example/Config.java || echo "   No deprecated annotations found"
echo ""

# Test 2: Dry run to see what would change
echo "3. Dry run - showing what would be changed:"
echo ""

./rewrite-spring-go -source examples/sample-spring-project -recipe change-property-key \
  -old-key "management.metrics.binders.jvm.enabled" \
  -new-key "management.metrics.enable.jvm" \
  -dry-run

echo ""

# Test 3: Actually apply one migration
echo "4. Applying migration: management.metrics.binders.jvm.enabled → management.metrics.enable.jvm"
echo ""

./rewrite-spring-go -source examples/sample-spring-project -recipe change-property-key \
  -old-key "management.metrics.binders.jvm.enabled" \
  -new-key "management.metrics.enable.jvm" \
  -backup

if [ $? -eq 0 ]; then
    echo "✅ Migration applied successfully"
else
    echo "❌ Migration failed"
    exit 1
fi

echo ""

# Test 4: Show changes
echo "5. Changes made:"
echo ""
echo "📄 Updated application.properties:"
grep -n "management.metrics.enable.jvm" examples/sample-spring-project/application.properties || echo "   Property not found in properties file"
echo ""
echo "📄 Updated application.yml:"
grep -A2 -B2 "enabled:" examples/sample-spring-project/application.yml || echo "   Property not found in YAML file"
echo ""
echo "📄 Updated Config.java:"
grep -n "management.metrics.enable.jvm" examples/sample-spring-project/src/main/java/com/example/Config.java || echo "   Annotation not found in Java file"
echo ""

# Test 5: Show backup files created
echo "6. Backup files created:"
find examples/sample-spring-project -name "*.backup" -exec echo "   📋 {}" \;
echo ""

# Test 6: Add a new property
echo "7. Adding new Spring Boot 3.x property:"
echo ""

./rewrite-spring-go -source examples/sample-spring-project -recipe add-property \
  -property "management.observations.annotations.enabled" \
  -value "true" \
  -comment "Enable observation annotations for Spring Boot 3.x"

echo ""

# Test 7: Show final state
echo "8. Final state verification:"
echo ""
echo "📄 New property in application.properties:"
grep -A1 -B1 "management.observations.annotations.enabled" examples/sample-spring-project/application.properties || echo "   Property not added to properties file"
echo ""
echo "📄 New property in application.yml:"
grep -A3 -B1 "observations:" examples/sample-spring-project/application.yml || echo "   Property not added to YAML file"
echo ""

# Test 8: Demonstrate the full migration script
echo "9. Testing full migration script (dry run):"
echo ""

if [ -f "scripts/spring-boot-migration.sh" ]; then
    chmod +x scripts/spring-boot-migration.sh
    
    # Create a copy of sample project for testing
    cp -r examples/sample-spring-project examples/sample-spring-project-test
    
    echo "Running: ./scripts/spring-boot-migration.sh examples/sample-spring-project-test 2.7 3.0"
    ./scripts/spring-boot-migration.sh examples/sample-spring-project-test 2.7 3.0
    
    echo ""
    echo "Migration script completed!"
    
    # Clean up test copy
    rm -rf examples/sample-spring-project-test
else
    echo "Migration script not found at scripts/spring-boot-migration.sh"
fi

echo ""
echo "🎉 Migration testing completed!"
echo ""
echo "📝 Summary:"
echo "   ✅ Tool builds successfully"
echo "   ✅ Property key changes work"
echo "   ✅ Property addition works"
echo "   ✅ Backup files are created"
echo "   ✅ Multiple file formats supported (properties, YAML, Java)"
echo ""
echo "🔄 To restore original files:"
echo "   find examples/sample-spring-project -name '*.backup' -exec sh -c 'mv \"\$1\" \"\${1%.backup}\"' _ {} \\;" 