# Defining Indexes in dbdiagram.io

## Basic Index Syntax

```sql
Table table_name {
  column_name data_type [constraints]
  Index index_name (column_name)
}
```

## Multiple Column Index

```sql
Table table_name {
  column1 data_type
  column2 data_type
  Index idx_name (column1, column2)
}
```

## Unique Index

```sql
Table table_name {
  column_name data_type
  Index idx_name (column_name) [unique]
}
```

## Example with Users Table

```sql
Table users {
  id integer [pk, increment]
  username string [unique]
  email string [unique]
  password_hash string [null]
  auth_provider string [null]
  provider_user_id string [null]
  provider_email string [null]
  access_token string [null]
  refresh_token string [null]
  token_expires_at timestamp [null]
  is_email_verified boolean [default: false]
  verification_token string [null]
  verification_token_expires_at timestamp [null]
  reset_token string [null]
  reset_token_expires_at timestamp [null]
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  last_login_at timestamp
  status enum('active', 'inactive', 'suspended') [default: 'active']

  // Single column indexes
  Index idx_email (email)
  Index idx_username (username)
  Index idx_auth_provider (auth_provider)
  Index idx_provider_user_id (provider_user_id)
  Index idx_verification_token (verification_token)
  Index idx_reset_token (reset_token)

  // Composite indexes
  Index idx_provider_auth (auth_provider, provider_user_id)
  Index idx_email_status (email, status)
  Index idx_created_status (created_at, status)

  // Unique indexes (alternative to [unique] constraint)
  Index idx_unique_email (email) [unique]
  Index idx_unique_username (username) [unique]
}
```

## Common Index Patterns

1. **Primary Key Index**

   - Automatically created with `[pk]` constraint
   - No need to explicitly define

2. **Foreign Key Index**

   - Should be created for all foreign key columns
   - Improves join performance

3. **Search Index**

   - Columns frequently used in WHERE clauses
   - Columns used in ORDER BY

4. **Composite Index**
   - Multiple columns used together in queries
   - Order matters (most selective first)

## Best Practices

1. **Index Selection**

   - Index columns used in WHERE clauses
   - Index columns used in JOIN conditions
   - Index columns used in ORDER BY
   - Consider composite indexes for common query patterns

2. **Index Naming**

   - Use consistent prefix (e.g., `idx_`)
   - Include table name for clarity
   - Describe the indexed columns

3. **Index Maintenance**

   - Monitor index usage
   - Remove unused indexes
   - Update statistics regularly

4. **Performance Considerations**
   - Balance between read and write performance
   - Consider index size
   - Watch for duplicate indexes

## Example Queries and Their Indexes

### Find by Email

```sql
SELECT * FROM users WHERE email = ?;
```

Required index: `Index idx_email (email)`

### Find by Auth Provider

```sql
SELECT * FROM users
WHERE auth_provider = ? AND provider_user_id = ?;
```

Required index: `Index idx_provider_auth (auth_provider, provider_user_id)`

### Find Active Users

```sql
SELECT * FROM users
WHERE status = 'active'
ORDER BY created_at DESC;
```

Required index: `Index idx_created_status (created_at, status)`
