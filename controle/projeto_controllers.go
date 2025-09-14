package controle

import (
	"modulo/modelo"
	"net/http"
	"strconv"
	// ... seus imports
)

// ProjetosListPage exibe todos os projetos
func ProjetosListPage(w http.ResponseWriter, r *http.Request) {
    projetos, err := modelo.GetAllProjetos()
    if err != nil {
        // Tratar erro
    }
    templates.ExecuteTemplate(w, "projetos.html", projetos)
}

// ProjetoDetailPage exibe os detalhes de um projeto
func ProjetoDetailPage(w http.ResponseWriter, r *http.Request) {
    id, _ := strconv.Atoi(r.URL.Query().Get("id"))
    projeto, err := modelo.GetProjetoByID(id)
    if err != nil {
        // Tratar erro
    }
    templates.ExecuteTemplate(w, "projeto_detail.html", projeto)

}



// Você precisará de uma página de cadastro para o projeto também.
// A lógica do POST para criar o projeto e a equipe em uma transação
// seria mais complexa e pode ser um próximo passo.