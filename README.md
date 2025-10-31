# Autocomplete and pattern search in text

Pattern search and autocomplete project based on KMP and Trie.

Este repositorio contiene una pequeña aplicación escrita en Go que permite:
- Buscar patrones en un corpus de texto usando el algoritmo KMP.
- Sugerir autocompletado por prefijo (estructura Trie) y combinar la sugerencia con la búsqueda.
- Interfaz web sencilla para probar búsquedas y ver contextos resaltados.

## Requisitos
- Go (recomendado >= 1.20). Descárgalo en: https://go.dev/dl/
- Navegador para la interfaz web.

## Ejecutar en desarrollo

```bash
go mod tidy
```

Run with app:

```bash
go run ./app
```

Open with browser:

```bash
http://localhost:8080
```

El servidor servirá una página simple para ingresar el patrón y ver las ocurrencias con contexto.