package controle

import (
	"log"
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

func CadastroProjetoPage(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        // Para o método GET, precisamos buscar todos os alunos para popular o formulário
        alunos, err := modelo.GetAllAlunos() // Você precisará ter essa função no seu modelo de aluno
        if err != nil {
            http.Error(w, "Erro ao buscar alunos", http.StatusInternalServerError)
            return
        }
        templates.ExecuteTemplate(w, "cadastro_projeto.html", alunos)
        return
    }

    if r.Method == http.MethodPost {
        r.ParseForm()

        // 1. Cria o objeto Projeto com os dados do formulário
        projeto := modelo.Projeto{
            Titulo:      r.FormValue("titulo"),
            Descricao:   r.FormValue("descricao"),
            ImagemURL:   r.FormValue("imagem_url"),
            LinkProjeto: r.FormValue("link_projeto"),
        }

        // 2. Cria a lista de membros da equipe
        var equipe []modelo.MembroEquipe
        // O formulário enviará os IDs e funções como listas (slices)
        alunoIDs := r.Form["aluno_id"]
        funcoes := r.Form["funcao"]

        // Garante que temos a mesma quantidade de IDs e funções
        if len(alunoIDs) == len(funcoes) {
            for i, idStr := range alunoIDs {
                alunoID, _ := strconv.Atoi(idStr)
                equipe = append(equipe, modelo.MembroEquipe{
                    AlunoID: alunoID,
                    Funcao:  funcoes[i],
                })
            }
        }

        // 3. Chama a função do modelo para salvar tudo no banco
        err := modelo.CreateProjeto(&projeto, equipe)
        if err != nil {
            log.Printf("Erro ao criar projeto: %v", err)
            http.Error(w, "Erro ao salvar o projeto", http.StatusInternalServerError)
            return
        }

        // 4. Redireciona para a página de listagem de projetos
        http.Redirect(w, r, "/projetos", http.StatusSeeOther)
    }
}
