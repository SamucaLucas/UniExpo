package controle

import (
	"log"
	models "modulo/modelo"
	"net/http"
	"strconv"
	"strings"
)

// 2. Carregamos os templates e injetamos nossas funções com .Funcs(funcMap)


func HomePage(w http.ResponseWriter, r *http.Request) {
	periodoFilter, _ := strconv.Atoi(r.URL.Query().Get("periodo"))
	skillFilter := r.URL.Query().Get("skill")

	alunos, err := models.GetAll(periodoFilter, skillFilter)
	if err != nil {
		log.Println("Erro ao buscar alunos:", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	err = templates.ExecuteTemplate(w, "index.html", alunos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AlunoPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	aluno, err := models.GetByID(id)
	if err != nil {
		log.Printf("Aluno com ID %d não encontrado: %v", id, err)
		http.NotFound(w, r)
		return
	}

	err = templates.ExecuteTemplate(w, "aluno.html", aluno)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CadastroPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templates.ExecuteTemplate(w, "cadastro.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		periodo, _ := strconv.Atoi(r.FormValue("periodo"))
		novoAluno := models.Aluno{
			Nome:        r.FormValue("nome"),
			Periodo:     periodo,
			Instagram:   r.FormValue("instagram"),
			Github:      r.FormValue("github"),
			Linkedin:      r.FormValue("linkedin"),
			FotoURL:     r.FormValue("foto_url"),
			Linguagens:  splitAndTrim(r.FormValue("linguagens")),
			Frameworks:  splitAndTrim(r.FormValue("frameworks")),
			Ferramentas: splitAndTrim(r.FormValue("ferramentas")),
			Bio:         r.FormValue("bio"),
		}

		err := novoAluno.Create()
		if err != nil {
			log.Println("Erro ao cadastrar novo aluno:", err)
			http.Error(w, "Erro ao salvar os dados", http.StatusInternalServerError)
			return
		}

		log.Printf("Novo aluno cadastrado com ID: %d", novoAluno.ID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// splitAndTrim é uma função utilitária para limpar os dados do formulário.
func splitAndTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i, v := range parts {
		parts[i] = strings.TrimSpace(v)
	}
	return parts
}