// modelo/projeto.go
package modelo

import database "modulo/db"

// MembroEquipe combina dados do aluno com sua função no projeto
type MembroEquipe struct {
    AlunoID   int
    Nome      string
    FotoURL   string
    Funcao    string // Ex: "Desenvolvedor Backend", "UI/UX Designer"
}

// Projeto armazena todas as informações de um projeto, incluindo a equipe
type Projeto struct {
    ID          int
    Titulo      string
    Descricao   string
    ImagemURL   string
    LinkProjeto string
    Equipe      []MembroEquipe // Uma lista de membros
}

// ---- Funções de Interação com o Banco de Dados ----

// GetProjetoByID busca um projeto e sua equipe
func GetProjetoByID(id int) (Projeto, error) {
    var p Projeto
    // Primeiro, busca os dados do projeto
    queryProjeto := `SELECT id, titulo, descricao, imagem_url, link_projeto FROM projetos WHERE id = $1`
    err := database.DB.QueryRow(queryProjeto, id).Scan(&p.ID, &p.Titulo, &p.Descricao, &p.ImagemURL, &p.LinkProjeto)
    if err != nil {
        return p, err
    }

    // Depois, busca os membros da equipe
    queryEquipe := `SELECT a.id, a.nome, a.foto_url, ep.funcao
                    FROM alunos a
                    JOIN equipe_projetos ep ON a.id = ep.aluno_id
                    WHERE ep.projeto_id = $1`
    rows, err := database.DB.Query(queryEquipe, id)
    if err != nil {
        return p, err
    }
    defer rows.Close()

    for rows.Next() {
        var membro MembroEquipe
        if err := rows.Scan(&membro.AlunoID, &membro.Nome, &membro.FotoURL, &membro.Funcao); err != nil {
            continue // Pula membros com erro
        }
        p.Equipe = append(p.Equipe, membro)
    }
    return p, nil
}

// GetAllProjetos busca todos os projetos (sem detalhes da equipe, para a lista principal)
func GetAllProjetos() ([]Projeto, error) {
    query := `SELECT id, titulo, descricao, imagem_url FROM projetos ORDER BY created_at DESC`
    rows, err := database.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var projetos []Projeto
    for rows.Next() {
        var p Projeto
        if err := rows.Scan(&p.ID, &p.Titulo, &p.Descricao, &p.ImagemURL); err != nil {
            continue
        }
        projetos = append(projetos, p)
    }
    return projetos, nil
}