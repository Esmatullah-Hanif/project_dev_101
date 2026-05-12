# Database Guide

## Overview

This project uses PostgreSQL as the primary data store, accessed via Supabase or a local PostgreSQL instance. The database design prioritizes data integrity, performance, and scalability.

## Connection

### Local PostgreSQL
```bash
createdb microservices
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/microservices"
```

### Supabase
```bash
# Get connection string from Supabase dashboard
export DATABASE_URL="postgresql://postgres:[password]@[project].supabase.co:5432/postgres"
```

### Environment Variable
```bash
# .env file
DATABASE_URL=postgresql://user:password@host:port/database
```

## Schema

### Users Table

The core table for storing user information:

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  first_name TEXT DEFAULT '',
  last_name TEXT DEFAULT '',
  bio TEXT DEFAULT '',
  avatar_url TEXT DEFAULT '',
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE
);
```

#### Columns

| Column | Type | Constraints | Purpose |
|--------|------|-----------|---------|
| id | UUID | PRIMARY KEY | Unique user identifier |
| email | TEXT | UNIQUE NOT NULL | Login email, must be unique |
| password_hash | TEXT | NOT NULL | Bcrypt hashed password |
| first_name | TEXT | DEFAULT '' | User's first name |
| last_name | TEXT | DEFAULT '' | User's last name |
| bio | TEXT | DEFAULT '' | User profile bio |
| avatar_url | TEXT | DEFAULT '' | URL to user's avatar image |
| is_active | BOOLEAN | DEFAULT TRUE | Account status |
| created_at | TIMESTAMP WITH TIME ZONE | DEFAULT NOW() | Account creation time |
| updated_at | TIMESTAMP WITH TIME ZONE | DEFAULT NOW() | Last profile update |
| deleted_at | TIMESTAMP WITH TIME ZONE | NULL | Soft delete timestamp |

#### Constraints

- **Primary Key (id)**: Ensures each user is uniquely identified
- **Unique (email)**: Prevents duplicate email registrations
- **Not Null (email, password_hash)**: Ensures required data is present

#### Indexes

```sql
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_created_at ON users(created_at DESC);
```

| Index | Purpose | Usage |
|-------|---------|-------|
| idx_users_email | Fast email lookup | Login queries |
| idx_users_deleted_at | Soft delete filtering | Exclude deleted users |
| idx_users_created_at | Timeline queries | Sort by creation date |

## Migrations

### Running Migrations

#### Method 1: psql
```bash
psql $DATABASE_URL < migrations/001_create_users_table.sql
```

#### Method 2: Supabase Dashboard
1. Go to SQL Editor
2. Copy migration SQL
3. Execute in editor
4. Verify schema in Table Editor

#### Method 3: Make Command
```bash
make migrate-up
```

### Migration Files

Located in `migrations/` directory:
- `001_create_users_table.sql` - Initial users table schema

### Adding New Migrations

1. Create file: `migrations/002_your_change.sql`
2. Start with a descriptive comment block
3. Use `IF NOT EXISTS` or `IF EXISTS` for safety
4. Include rollback notes in comments
5. Test locally before applying to production

Example:
```sql
/*
  # Add phone number to users

  1. New Columns
    - `phone` (text, optional)

  2. Migration
    - Add phone column to users table
*/

ALTER TABLE users ADD COLUMN IF NOT EXISTS phone TEXT DEFAULT '';
```

## Data Access

### pgx Connection Pool

The project uses `jackc/pgx` for high-performance database access:

```go
config, _ := pgxpool.ParseConfig(databaseURL)
config.MaxConns = 25
config.MinConns = 5
pool, _ := pgxpool.NewWithConfig(ctx, config)
```

### Repository Pattern

Each service has a repository for data access:

```go
// Auth Service
type UserRepository struct {
    db *pgxpool.Pool
}

// Query methods
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error)
func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error
func (r *UserRepository) UpdatePassword(ctx context.Context, userID string, hash string) error
```

### Parameterized Queries

Always use parameterized queries to prevent SQL injection:

```go
// Safe - parameterized
query := "SELECT * FROM users WHERE email = $1"
row := db.QueryRow(ctx, query, email)

// Unsafe - string interpolation (DO NOT USE)
query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
```

## Common Queries

### Get User by Email (Login)
```go
query := `
    SELECT id, email, password_hash, created_at, updated_at
    FROM users
    WHERE email = $1 AND deleted_at IS NULL
`
row := db.QueryRow(ctx, query, email)
```

### Create New User
```go
query := `
    INSERT INTO users (id, email, password_hash, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5)
`
_, err := db.Exec(ctx, query, userID, email, passwordHash, now, now)
```

### Update User Profile
```go
query := `
    UPDATE users
    SET first_name = $1, last_name = $2, bio = $3, avatar_url = $4, updated_at = $5
    WHERE id = $6 AND deleted_at IS NULL
    RETURNING id, email, first_name, last_name, bio, avatar_url, is_active, created_at, updated_at
`
```

### List Users with Pagination
```go
query := `
    SELECT id, email, first_name, last_name, bio, avatar_url, is_active, created_at, updated_at
    FROM users
    WHERE deleted_at IS NULL
    ORDER BY created_at DESC
    LIMIT $1 OFFSET $2
`
```

### Soft Delete User
```go
query := `
    UPDATE users
    SET deleted_at = $1, updated_at = $2
    WHERE id = $3 AND deleted_at IS NULL
`
```

## Performance Optimization

### Indexing Strategy

Current indexes optimize for:
- **Authentication**: Fast email lookup
- **Data Integrity**: Soft delete filtering
- **Sorting**: Chronological user listings

### Connection Pooling

Configured in `shared/pkg/database/database.go`:
- **MaxConns**: 25 (adjustable based on load)
- **MinConns**: 5 (keeps minimum connections ready)
- **MaxConnLifetime**: 1 hour (prevent stale connections)
- **MaxConnIdleTime**: 10 minutes (recycle idle connections)

### Query Optimization Tips

1. **Always filter soft-deleted rows**:
   ```sql
   WHERE deleted_at IS NULL
   ```

2. **Use indexes for WHERE clauses**:
   ```sql
   SELECT * FROM users WHERE email = $1  -- Uses idx_users_email
   ```

3. **Limit result sets for pagination**:
   ```sql
   LIMIT 10 OFFSET 0  -- Never fetch all rows
   ```

4. **Use EXPLAIN to verify index usage**:
   ```sql
   EXPLAIN ANALYZE SELECT * FROM users WHERE email = $1;
   ```

## Backup Strategy

### Supabase Backups
- Automatic daily backups
- 30-day retention
- Point-in-time recovery available
- Access via Supabase dashboard

### Manual Backups
```bash
# Dump entire database
pg_dump $DATABASE_URL > backup.sql

# Restore from backup
psql $DATABASE_URL < backup.sql
```

## Security Considerations

### Password Storage
- Bcrypt hashing (cost factor: 10)
- Salt generated automatically per password
- Never store plain text passwords
- 60-character hash column sufficient

### Data Privacy
- No sensitive data in logs
- Consider encryption at rest for sensitive fields
- Soft deletes preserve audit trail
- Regular security updates for PostgreSQL

### Access Control
- Database credentials in environment variables
- Different credentials for dev/staging/prod
- Read-only replicas for analytics (future)
- Row-level security (Supabase feature)

## Monitoring

### Connection Pool Health
```go
// Monitor from logs
// Check for connection timeouts or exhaustion
```

### Query Performance
```sql
-- Find slow queries
SELECT query, calls, total_time, mean_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;
```

### Table Size
```sql
-- Check table sizes
SELECT schemaname, tablename,
  pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

## Future Enhancements

1. **Full-Text Search**: Add GIN indexes for user search
2. **Partitioning**: Partition users table by date for scale
3. **Read Replicas**: Separate read traffic for analytics
4. **JSON Fields**: Add flexible metadata columns
5. **Audit Logging**: Track all user modifications
6. **Encryption**: Encrypt sensitive fields at rest
7. **Replication**: Multi-region replication for DR

## Troubleshooting

### Connection Refused
```bash
# Check PostgreSQL is running
pg_isready -h localhost -p 5432

# Verify connection string
echo $DATABASE_URL
```

### Slow Queries
```bash
# Check slow query log
EXPLAIN ANALYZE SELECT ...;

# Consider adding index if missing
```

### Permission Denied
```bash
# Verify user has proper roles
\du  -- in psql

# Grant permissions if needed
GRANT ALL ON users TO postgres;
```

### Disk Space Issues
```bash
# Check disk usage
SELECT pg_size_pretty(pg_database_size('microservices'));

# Clean up old backups or WAL files
```
