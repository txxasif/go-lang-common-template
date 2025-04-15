# Database Design for Multiple Authentication Methods

## Core Tables

### Users Table

```sql
Table users {
  id integer [pk, increment]
  username string [unique]
  email string [unique]
  password_hash string [null]  // Nullable for social auth users
  is_email_verified boolean [default: false]
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  last_login_at timestamp
  status enum('active', 'inactive', 'suspended') [default: 'active']
  Note: 'Core user information'
}
```

### Social Auth Providers

```sql
Table social_auth_providers {
  id integer [pk, increment]
  name string [unique]  // e.g., 'google', 'github', 'facebook'
  client_id string
  client_secret string
  redirect_uri string
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  Note: 'Configuration for each social auth provider'
}
```

### User Social Accounts

```sql
Table user_social_accounts {
  id integer [pk, increment]
  user_id integer [ref: > users.id]
  provider_id integer [ref: > social_auth_providers.id]
  provider_user_id string  // The ID from the social provider
  email string [null]  // Email from social provider
  access_token string
  refresh_token string [null]
  token_expires_at timestamp [null]
  created_at timestamp [default: `now()`]
  updated_at timestamp [default: `now()`]
  Index idx_user_provider (user_id, provider_id)
  Index idx_provider_user (provider_id, provider_user_id)
  Note: 'Links users to their social media accounts'
}
```

### Email Verification Tokens

```sql
Table email_verification_tokens {
  id integer [pk, increment]
  user_id integer [ref: > users.id]
  token string [unique]
  expires_at timestamp
  created_at timestamp [default: `now()`]
  Note: 'Stores email verification tokens'
}
```

### Password Reset Tokens

```sql
Table password_reset_tokens {
  id integer [pk, increment]
  user_id integer [ref: > users.id]
  token string [unique]
  expires_at timestamp
  created_at timestamp [default: `now()`]
  Note: 'Stores password reset tokens'
}
```

## Design Considerations

### 1. Authentication Flow

- Email/Password:

  1. User registers with email and password
  2. Verification email sent
  3. User verifies email
  4. User can log in

- Social Auth:
  1. User clicks social login button
  2. Redirected to provider
  3. Returns with provider user ID
  4. Account created/linked automatically

### 2. Security Considerations

- Passwords are always hashed
- Social tokens are encrypted
- Email verification required for email/password
- Password reset flow available
- Account linking management

### 3. Data Relationships

```
Users 1---* UserSocialAccounts
Users 1---* EmailVerificationTokens
Users 1---* PasswordResetTokens
SocialAuthProviders 1---* UserSocialAccounts
```

## Implementation Notes

1. **User Creation**

   - Check if email exists
   - For social auth, check provider_user_id
   - Allow linking multiple social accounts

2. **Login Flow**

   - Try email/password first
   - Check social accounts
   - Update last_login_at

3. **Account Management**

   - Allow unlinking social accounts
   - Prevent last password removal
   - Handle email changes

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
SELECT u.*
FROM users u
JOIN user_social_accounts usa ON u.id = usa.user_id
JOIN social_auth_providers sap ON usa.provider_id = sap.id
WHERE sap.name = ? AND usa.provider_user_id = ?;
```

### Get User's Social Accounts

```sql
SELECT sap.name, usa.email, usa.created_at
FROM user_social_accounts usa
JOIN social_auth_providers sap ON usa.provider_id = sap.id
WHERE usa.user_id = ?;
```

## Best Practices

1. **Data Privacy**

   - Encrypt sensitive data
   - Regular token rotation
   - Secure token storage

2. **Performance**

   - Index frequently queried fields
   - Cache social provider configs
   - Optimize token queries

3. **Maintenance**

   - Regular token cleanup
   - Provider config updates
   - Security audits

4. **Scalability**
   - Horizontal scaling possible
   - Stateless authentication
   - Caching strategies
