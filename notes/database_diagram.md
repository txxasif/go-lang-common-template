# Database Diagram Documentation

## Overview

This document contains the dbdiagram.io syntax for creating a comprehensive database diagram. dbdiagram.io is a simple tool to draw database diagrams by just writing code.

## Basic Syntax

### Table Definition

```sql
Table table_name {
  column_name data_type [constraints]
  column_name data_type [constraints]
  ...
}
```

### Data Types

- `integer` - Whole numbers
- `string` - Text data
- `text` - Long text data
- `boolean` - True/False values
- `date` - Date values
- `timestamp` - Date and time values
- `float` - Decimal numbers
- `double` - Double precision decimal numbers
- `json` - JSON data
- `jsonb` - Binary JSON data

### Constraints

- `pk` - Primary key
- `increment` - Auto-incrementing
- `null` - Can be null
- `not null` - Cannot be null
- `unique` - Must be unique
- `default` - Default value
- `ref` - Reference to another table

### Relationships

```sql
Table posts {
  id integer [pk, increment]
  user_id integer [ref: > users.id]
  title string
  content text
}

Table users {
  id integer [pk, increment]
  name string
  email string [unique]
}
```

### Indexes

```sql
Table users {
  id integer [pk, increment]
  email string [unique]
  name string
  Index idx_name (name)
}
```

### Enums

```sql
Table posts {
  id integer [pk, increment]
  status enum('draft', 'published', 'archived')
}
```

### Notes and Comments

```sql
Table users {
  id integer [pk, increment]
  name string
  email string [unique]
  Note: 'Stores user information'
}
```

## Example Complete Diagram

```sql
// Example Database Schema

Table users {
  id integer [pk, increment]
  username string [unique, not null]
  email string [unique, not null]
  password_hash string [not null]
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  status enum('active', 'inactive', 'suspended')
  Note: 'Stores user account information'
}

Table posts {
  id integer [pk, increment]
  user_id integer [ref: > users.id]
  title string [not null]
  content text
  status enum('draft', 'published', 'archived')
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  Index idx_user_id (user_id)
  Note: 'Blog posts created by users'
}

Table comments {
  id integer [pk, increment]
  post_id integer [ref: > posts.id]
  user_id integer [ref: > users.id]
  content text [not null]
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  Index idx_post_id (post_id)
  Index idx_user_id (user_id)
  Note: 'Comments on blog posts'
}

Table categories {
  id integer [pk, increment]
  name string [not null]
  slug string [unique, not null]
  description text
  created_at timestamp [default: `now()`]
  Note: 'Post categories'
}

Table post_categories {
  post_id integer [ref: > posts.id]
  category_id integer [ref: > categories.id]
  Index idx_post_category (post_id, category_id)
  Note: 'Many-to-many relationship between posts and categories'
}
```

## Best Practices

1. **Naming Conventions**

   - Use singular nouns for table names
   - Use snake_case for column names
   - Be consistent with naming patterns

2. **Primary Keys**

   - Always include a primary key
   - Use auto-incrementing integers for simple cases
   - Consider UUID for distributed systems

3. **Foreign Keys**

   - Always define relationships explicitly
   - Use consistent naming for foreign key columns
   - Add appropriate indexes for foreign keys

4. **Data Types**

   - Choose appropriate data types for each column
   - Consider storage requirements
   - Plan for future growth

5. **Indexes**

   - Add indexes for frequently queried columns
   - Consider composite indexes for common query patterns
   - Don't over-index as it affects write performance

6. **Documentation**
   - Add notes to explain table purposes
   - Document complex relationships
   - Include examples of common queries

## Common Patterns

### Soft Deletes

```sql
Table users {
  id integer [pk, increment]
  name string
  deleted_at timestamp [null]
  Note: 'Uses soft delete pattern'
}
```

### Audit Trail

```sql
Table posts {
  id integer [pk, increment]
  title string
  content text
  created_by integer [ref: > users.id]
  created_at timestamp [default: `now()`]
  updated_by integer [ref: > users.id]
  updated_at timestamp [default: `now()`]
  Note: 'Includes audit trail fields'
}
```

### Polymorphic Relationships

```sql
Table comments {
  id integer [pk, increment]
  commentable_type string [not null]
  commentable_id integer [not null]
  content text
  Index idx_commentable (commentable_type, commentable_id)
  Note: 'Supports comments on multiple types of content'
}
```

## Export Options

dbdiagram.io supports exporting to:

- PNG
- PDF
- PostgreSQL
- MySQL
- SQLite
- Oracle
- SQL Server

## Tips and Tricks

1. Use the `Note` feature to document table purposes
2. Group related tables together in the code
3. Use consistent indentation for readability
4. Add indexes for performance optimization
5. Consider using enums for fixed sets of values
6. Plan for future scalability
7. Document complex relationships
8. Use appropriate data types for each column
9. Consider using soft deletes for data retention
10. Include audit trail fields where appropriate
