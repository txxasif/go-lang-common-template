# Simplified Database Design for Multiple Authentication Methods

## Core Tables

### Users Table

```sql
Table users {
  id integer [pk, increment]
  username string [unique]
  email string [unique]
  password_hash string [null]  // Nullable for social auth users
  auth_provider string [null]  // e.g., 'email', 'google', 'github', 'facebook'
  provider_user_id string [null]  // ID from social provider
  provider_email string [null]  // Email from social provider
  access_token string [null]  // For social auth
  refresh_token string [null]  // For social auth
  token_expires_at timestamp [null]  // For social auth
  is_email_verified boolean [default: false]
  verification_token string [null]
  verification_token_expires_at timestamp [null]
  reset_token string [null]
  reset_token_expires_at timestamp [null]
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  last_login_at timestamp
  status enum('active', 'inactive', 'suspended') [default: 'active']
  Note: 'Combined user information for all auth methods'
}
```

## Design Considerations

### 1. Authentication Flow

- Email/Password:

  1. User registers with email and password
  2. Verification email sent (stored in verification_token)
  3. User verifies email
  4. User can log in

- Social Auth:
  1. User clicks social login button
  2. Redirected to provider
  3. Returns with provider info
  4. Account created with provider details

### 2. Security Considerations

- Passwords are always hashed
- Social tokens are encrypted
- Email verification required for email/password
- Password reset flow available
- Token expiration management

### 3. Data Relationships

```
Users 1---1 AuthProvider (if social auth)
```

## Implementation Notes

1. **User Creation**

   - Check if email exists
   - For social auth, check provider_user_id
   - Store all auth-related data in one record

2. **Login Flow**

   - Try email/password first
   - Check social auth if no password
   - Update last_login_at

3. **Account Management**

   - Allow switching auth methods
   - Handle email changes
   - Manage token expiration

4. **Security Features**
   - Rate limiting
   - IP tracking
   - Suspicious activity detection
   - Session management

## Example Queries

### Find User by Email

```sql
SELECT * FROM users WHERE email = ?;
```

### Find User by Social Provider

```sql
SELECT * FROM users
WHERE auth_provider = ? AND provider_user_id = ?;
```

### Get User's Auth Details

```sql
SELECT
  id,
  email,
  auth_provider,
  is_email_verified,
  status
FROM users
WHERE id = ?;
```

## Best Practices

1. **Data Privacy**

   - Encrypt sensitive data
   - Regular token rotation
   - Secure token storage

2. **Performance**

   - Index frequently queried fields
   - Cache user data
   - Optimize token queries

3. **Maintenance**

   - Regular token cleanup
   - Provider config updates
   - Security audits

4. **Scalability**
   - Horizontal scaling possible
   - Stateless authentication
   - Caching strategies

## Advantages of This Design

1. **Simpler Schema**

   - Single table for all user data
   - Easier to query and maintain
   - Reduced joins

2. **Flexible Authentication**

   - Easy to add new auth providers
   - Simple to switch auth methods
   - Clear user state

3. **Better Performance**

   - Fewer table joins
   - Faster queries
   - Simpler caching

4. **Easier Management**
   - Single point of update
   - Clearer data relationships
   - Simpler migrations
