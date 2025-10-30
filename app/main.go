package main

import (
	"fmt"
	"html"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"textSearching/pkg/kmp"
	"textSearching/pkg/normalize"
)

func bringText(text string) string {
	textPath := text
	contentBytes, err := os.ReadFile(textPath)

	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return err.Error()
	}

	fileContent := string(contentBytes)

	return fileContent
}

type PageData struct {
	TextPath, Pattern, Result string
	Positions                 []int
	Contexts                  []template.HTML
	TimeElapsed               time.Duration
}

func main() {
	filePath := "./Books/dracula.txt"
	text := bringText(filePath)
	normalizedText := normalize.NormalizeText(text)

	tmpl := template.Must(template.New("index").Parse(`
	<!DOCTYPE html>
	<html lang="es">
	<head>
		<meta charset="UTF-8">
		<title>Buscador de Patrones (KMP)</title>
		<style>
			body { font-family: sans-serif; background: #f7f7f7; padding: 2rem; }
			h1 { color: #333; }
			form { background: white; padding: 1rem; border-radius: 8px; box-shadow: 0 2px 6px rgba(0,0,0,0.1); }
			input, textarea, button { width: 100%; margin-top: 10px; padding: 8px; }
			button { background: #007BFF; color: white; border: none; border-radius: 4px; cursor: pointer; }
			.result { margin-top: 1.5rem; background: #fff; padding: 1rem; border-radius: 8px; box-shadow: 0 2px 6px rgba(0,0,0,0.1); }
		</style>
	</head>
	<body>
		<h1>Buscador de Patrones con KMP</h1>
		<form method="POST" action="/">
			<label>Patrón a buscar:</label>
			<input name="pattern" value="{{.Pattern}}" placeholder="Ingresa una palabra o frase...">
			<button type="submit">Buscar</button>
		</form>

		{{if .Positions}}
		<div class="result">
			<h3>Resultados:</h3>
			<p><strong>Patrón:</strong> {{.Pattern}}</p>
			<p><strong>Tiempo:</strong> {{.TimeElapsed}}</p>
			<p><strong>Ocurrencias encontradas:</strong> {{len .Positions}}</p>
			<ul>
			{{range $idx, $pos := .Positions}}
				<li><strong>Posición {{$pos}}:</strong> {{index $.Contexts $idx}}</li>
			{{end}}
			</ul>
		</div>
		{{else if .Result}}
		<div class="result">
			<h3>Resultados:</h3>
			<p><strong>Patrón:</strong> {{.Pattern}}</p>
			<p><strong>Tiempo:</strong> {{.TimeElapsed}}</p>
			<p>{{.Result}}</p>
		</div>
		{{end}}
	</body>
	</html>
	`))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{}

		if r.Method == http.MethodPost {
			r.ParseForm()
			pattern := r.Form.Get("pattern")
			data.Pattern = pattern

			if strings.TrimSpace(pattern) == "" {
				data.Result = "Por favor, ingrese un patron valido."
				tmpl.Execute(w, data)
				return
			}

			normalizedPattern := normalize.NormalizeText(pattern)

			start := time.Now()
			positions := kmp.Kmp(normalizedPattern, normalizedText)
			elapsed := time.Since(start)
			data.TimeElapsed = elapsed

			if len(positions) == 0 {
				data.Result = "No se encontraron ocurrencias del patron."
			} else {
				// Guardar todas las posiciones y contexto alrededor de cada una
				data.Positions = positions
				contexts := make([]template.HTML, len(positions))
				escPattern := html.EscapeString(normalizedPattern)
				for i, pos := range positions {
					startIdx := pos - 20
					if startIdx < 0 {
						startIdx = 0
					}
					endIdx := pos + len(normalizedPattern) + 20
					if endIdx > len(normalizedText) {
						endIdx = len(normalizedText)
					}
					raw := normalizedText[startIdx:endIdx]
					// Escape the context, then replace escaped pattern with bold tag
					escaped := html.EscapeString(raw)
					highlighted := strings.ReplaceAll(escaped, escPattern, "<strong>"+escPattern+"</strong>")
					contexts[i] = template.HTML(highlighted)
				}
				data.Contexts = contexts
			}

			tmpl.Execute(w, data)
			return
		}
		tmpl.Execute(w, data)
	})

	fmt.Println("Servidor en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
