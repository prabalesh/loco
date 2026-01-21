# Loco Database Migration (v2)

This project has been updated with v2 schema extensions to support automated problem creation, custom data types, and reference solution validation.

## Migration Steps

To apply the v2 migrations, follow these steps:

1.  **Backup your database**: Always perform a backup before running migrations on production data.
2.  **Import the migration function**: In your server initialization (e.g., `cmd/server/main.go`), import the `migrations` package.
3.  **Call `MigrateV2Schema`**: After initializing the GORM DB connection, call the migration function:

```go
import "github.com/prabalesh/loco/backend/internal/migrations"

// ... after db initialization ...
if err := migrations.MigrateV2Schema(db.DB); err != nil {
    log.Fatal("Failed to run v2 migrations", zap.Error(err))
}
```

## Changes Summary

### Existing Tables Extended
- `problems`: Added fields for function signatures, validation status, and complexity limits.
- `test_cases`: Added support for multiple valid answers and resource limits.
- `users`: Added `is_bot` flag for automated testing users.

### New Tables Created
- `problem_boilerplates`: Caches language-specific starter code.
- `custom_types`: Definitions for `TreeNode`, `ListNode`, etc.
- `type_implementations`: Language-specific serializers for custom types.
- `problem_reference_solutions`: Reference implementations for validation.
