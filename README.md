# url-shortener

Minimal URL shortener built with Go and SQLite. No frameworks.

## Features

- Shorten any URL to a random 6-character hex code
- Redirect short codes to original URLs
- List all stored URLs
- Persistent storage via SQLite (`data.db`)

## Requirements

- Go 1.26+
- GCC (for `go-sqlite3` CGo compilation)

## Run

```bash
go run main.go
```

Server starts on `:8080`.

## API

### Shorten a URL

```
POST /shorten
Content-Type: application/json

{"url": "https://example.com"}
```

Response `201`:

```json
{
  "id": 1,
  "url": "https://example.com",
  "code": "a3f9c1",
  "createdAt": "2026-04-18T10:00:00Z"
}
```

### List all URLs

```
GET /urls
```

### Redirect

```
GET /{code}
```

Redirects to the original URL (`303`), or returns `404` if the code is unknown.

## Hot Reload

Uses [Air](https://github.com/air-verse/air). Run:

```bash
air
```
