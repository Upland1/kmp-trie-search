package main

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"os"
	"strings"
	"time"

	"textSearching/pkg/kmp"
	"textSearching/pkg/normalize"
	trie "textSearching/pkg/trie_autocomplete"
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

func buildTrieFromText(text string) *trie.Trie {
	t := trie.NewTrie()
	words := strings.Fields(text)
	for _, word := range words {
		clean := strings.ToLower(strings.Trim(word, ",.!?\"'"))
		if clean == "" {
			continue
		}
		t.Insert(clean)
	}
	return t
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
	trie := buildTrieFromText(normalizedText)

	tmpl := template.Must(template.ParseFiles("app/index.html"))

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

			tmpl.ExecuteTemplate(w, "index.html", data)
			return
		}
		tmpl.ExecuteTemplate(w, "index.html", data)
	})

	http.HandleFunc("/suggest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		prefix := r.URL.Query().Get("prefix")

		if prefix == "" {
			w.Write([]byte("[]"))
			return
		}

		results := trie.Suggest(strings.ToLower(prefix), 10)

		jsonData, err := json.Marshal(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("[]"))
			return
		}
		w.Write(jsonData)
	})

	fmt.Println("Servidor en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
