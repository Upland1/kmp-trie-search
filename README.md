# Autocomplete and pattern search in text

Pattern search and autocomplete project based on KMP and Trie.

This repository contains a small application written in Go that allows:
- Pattern searching in a text corpus using the KMP algorithm.
- Autocomplete suggestions by prefix (using a Trie structure) and combining suggestions with the search.
- Simple web interface to test searches and view highlighted contexts.

## Requirements
- Go (recommended >= 1.20). Download it at: https://go.dev/dl/
- A web browser for the interface.

## Run in development

```bash
go mod tidy
```

Run the app:

```bash
go run ./app
```

Open in browser:

```bash
http://localhost:8080
```

The server will serve a simple page to enter the pattern and view occurrences with context.
