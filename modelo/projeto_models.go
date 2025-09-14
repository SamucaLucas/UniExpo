// modelo/projeto.go
package modelo

import (
	"log"
	database "modulo/db"
)

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

func GetAllAlunos() ([]Aluno, error) {
    // Selecionamos as colunas mais importantes e ordenamos por nome.
    // Ordenar por nome (ORDER BY nome ASC) é crucial para a usabilidade do formulário.
    query := `SELECT id, nome, periodo, foto_url FROM alunos ORDER BY nome ASC`

    rows, err := database.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var alunos []Aluno
    // O loop vai ler cada linha do resultado do banco
    for rows.Next() {
        var a Aluno
        // O Scan vai atribuir os valores das colunas para os campos da struct 'a'
        // A ordem dos &a.Campo deve ser a MESMA das colunas no SELECT.
        err := rows.Scan(&a.ID, &a.Nome, &a.Periodo, &a.FotoURL)
        if err != nil {
            // Se houver um erro ao ler uma linha, logamos e continuamos para a próxima
            log.Println("Erro ao escanear aluno:", err)
            continue
        }
        alunos = append(alunos, a)
    }

    return alunos, nil
}

func CreateProjeto(projeto *Projeto, equipe []MembroEquipe) error {
    // Inicia a transação
    tx, err := database.DB.Begin()
    if err != nil {
        return err
    }
    // Garante que a transação será desfeita (ROLLBACK) se algo der errado
    defer tx.Rollback()

    // 1. Insere os dados na tabela 'projetos' e obtém o ID do novo projeto
    stmtProjeto, err := tx.Prepare(`INSERT INTO projetos (titulo, descricao, imagem_url, link_projeto)
                                VALUES ($1, $2, $3, $4) RETURNING id`)
    if err != nil {
        return err
    }
    defer stmtProjeto.Close()

    var projetoID int
    err = stmtProjeto.QueryRow(projeto.Titulo, projeto.Descricao, projeto.ImagemURL, projeto.LinkProjeto).Scan(&projetoID)
    if err != nil {
        return err
    }

    // 2. Insere cada membro na tabela de junção 'equipe_projetos'
    stmtEquipe, err := tx.Prepare(`INSERT INTO equipe_projetos (projeto_id, aluno_id, funcao)
                                VALUES ($1, $2, $3)`)
    if err != nil {
        return err
    }
    defer stmtEquipe.Close()

    for _, membro := range equipe {
        // Ignora se o ID do aluno for 0 (opção "Selecione..." do formulário)
        if membro.AlunoID == 0 {
            continue
        }
        _, err := stmtEquipe.Exec(projetoID, membro.AlunoID, membro.Funcao)
        if err != nil {
            return err
        }
    }

    // 3. Se tudo correu bem, confirma a transação (COMMIT)
    return tx.Commit()
}