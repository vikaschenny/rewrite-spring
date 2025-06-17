#!/bin/bash

echo "=== rewrite-spring-go Demo ==="
echo

# Build the tool
echo "Building rewrite-spring-go..."
go build -o rewrite-spring-go cmd/rewrite-spring/main.go

echo "Built successfully!"
echo

# Demo 1: Change property key (dry run first)
echo "Demo 1: Changing deprecated property key (dry run)"
echo "Command: ./rewrite-spring-go -source examples/sample-spring-project -recipe change-property-key -old-key 'management.metrics.binders.jvm.enabled' -new-key 'management.metrics.enable.jvm' -dry-run"
echo

./rewrite-spring-go -source examples/sample-spring-project -recipe change-property-key \
  -old-key "management.metrics.binders.jvm.enabled" \
  -new-key "management.metrics.enable.jvm" \
  -dry-run

echo
echo "Press Enter to continue..."
read

# Demo 2: Add new property (dry run)
echo "Demo 2: Adding new property (dry run)"
echo "Command: ./rewrite-spring-go -source examples/sample-spring-project -recipe add-property -property 'spring.profiles.active' -value 'development' -comment 'Default active profile' -dry-run"
echo

./rewrite-spring-go -source examples/sample-spring-project -recipe add-property \
  -property "spring.profiles.active" \
  -value "development" \
  -comment "Default active profile" \
  -dry-run

echo
echo "Press Enter to continue..."
read

# Demo 3: Actually apply changes (with backup)
echo "Demo 3: Actually applying property key change"
echo "Command: ./rewrite-spring-go -source examples/sample-spring-project -recipe change-property-key -old-key 'management.metrics.binders.jvm.enabled' -new-key 'management.metrics.enable.jvm' -backup"
echo

./rewrite-spring-go -source examples/sample-spring-project -recipe change-property-key \
  -old-key "management.metrics.binders.jvm.enabled" \
  -new-key "management.metrics.enable.jvm" \
  -backup

echo
echo "Demo 4: Adding the new property"
echo "Command: ./rewrite-spring-go -source examples/sample-spring-project -recipe add-property -property 'spring.profiles.active' -value 'development' -comment 'Default active profile'"
echo

./rewrite-spring-go -source examples/sample-spring-project -recipe add-property \
  -property "spring.profiles.active" \
  -value "development" \
  -comment "Default active profile"

echo
echo "=== Demo completed! ==="
echo "Check the files in examples/sample-spring-project to see the changes."
echo "Backup files (.backup) have been created for modified files."
echo
echo "To restore original files:"
echo "find examples/sample-spring-project -name '*.backup' -exec sh -c 'mv \"\$1\" \"\${1%.backup}\"' _ {} \\;" 