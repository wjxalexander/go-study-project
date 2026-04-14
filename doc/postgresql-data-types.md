# PostgreSQL Data Types

## BYTEA

`BYTEA` is a PostgreSQL data type for storing **raw binary data** (byte arrays).

### Common Use Cases

- Password hashes (e.g. bcrypt, argon2)
- Cryptographic tokens and keys
- SHA-256 / SHA-512 digests
- Small binary files (prefer object storage for large files)

### Comparison with TEXT

| Type | Stores | Example Value |
|------|--------|---------------|
| `TEXT` / `VARCHAR` | Human-readable strings | `"hello world"` |
| `BYTEA` | Raw byte sequences | `\x243261243132...` |

### Example: Token Hash Column

```sql
CREATE TABLE IF NOT EXISTS tokens (
    token_hash BYTEA NOT NULL,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expiry     TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    scope      TEXT NOT NULL
);
```

### Go Mapping

`BYTEA` maps to `[]byte` in Go.

```go
type Token struct {
    Hash   []byte    `json:"-"`
    UserID int64     `json:"user_id"`
    Expiry time.Time `json:"expiry"`
    Scope  string    `json:"scope"`
}
```

### Storage Notes

- `BYTEA` uses hex format (`\x...`) by default for input/output.
- More space-efficient than encoding binary data as a hex or base64 `TEXT` string.
- PostgreSQL handles escaping automatically; no manual encoding needed in queries.

### Useful Commands

```sql
-- Check byte length of a BYTEA column
SELECT LENGTH(token_hash) FROM tokens;

-- Insert raw bytes using hex notation
INSERT INTO tokens (token_hash) VALUES ('\xDEADBEEF');

-- Compare with a known hash
SELECT * FROM tokens WHERE token_hash = '\xABCDEF1234';
```
