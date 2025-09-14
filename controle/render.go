// arquivo: controle/render.go
package controle

import (
	"html/template"
)
var templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("templete/**/*.html"))
// funcMap contém TODAS as funções customizadas que qualquer template possa precisar.
var funcMap = template.FuncMap{
    // Funções que você já criou
    "add":      func(a, b int) int { return a + b },
    "loop":     func(n int) []int { return make([]int, n) },
    "truncate": func(s string, length int) string {
        runes := []rune(s)
        if len(runes) > length {
            return string(runes[:length]) + "..."
        }
        return s
    },
}

