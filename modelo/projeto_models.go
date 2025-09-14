// modelo/projeto.go
package modelo

import (
	"log"
	database "modulo/db"
	"time"
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
    ID              int
    Titulo          string
    Descricao       string
    ImagemURL       string
    LinkProjeto     string
    Equipe          []MembroEquipe 
    Avaliacoes      []Avaliacao    // <-- NOVO CAMPO
    NotaMedia       float64        // <-- NOVO CAMPO
    TotalAvaliacoes int            // <-- NOVO CAMPO
}

// ---- Funções de Interação com o Banco de Dados ----

// GetProjetoByID busca um projeto e sua equipe
func GetProjetoByID(id int) (Projeto, error) {
    var p Projeto

    // --- PARTE 1: Busca os dados principais do projeto ---
    queryProjeto := `SELECT id, titulo, descricao, imagem_url, link_projeto FROM projetos WHERE id = $1`
    err := database.DB.QueryRow(queryProjeto, id).Scan(&p.ID, &p.Titulo, &p.Descricao, &p.ImagemURL, &p.LinkProjeto)
    if err != nil {
        return p, err
    }

    // --- PARTE 2: Busca os membros da equipe ---
    queryEquipe := `SELECT a.id, a.nome, a.foto_url, ep.funcao
                    FROM alunos a
                    JOIN equipe_projetos ep ON a.id = ep.aluno_id
                    WHERE ep.projeto_id = $1`
    rowsEquipe, err := database.DB.Query(queryEquipe, id)
    if err != nil {
        return p, err
    }
    defer rowsEquipe.Close()

    for rowsEquipe.Next() {
        var membro MembroEquipe
        if err := rowsEquipe.Scan(&membro.AlunoID, &membro.Nome, &membro.FotoURL, &membro.Funcao); err != nil {
            log.Println("Erro ao escanear membro da equipe:", err)
            continue
        }
        p.Equipe = append(p.Equipe, membro)
    }
    // Verifica se ocorreu um erro durante a iteração do loop da equipe
    if err = rowsEquipe.Err(); err != nil {
        return p, err
    }


    // --- PARTE 3: Busca as avaliações e calcula a média ---
    queryMedia := `SELECT COALESCE(AVG(nota), 0), COUNT(nota) FROM avaliacoes WHERE projeto_id = $1`
    // Usamos "=" aqui para reatribuir à variável 'err' já existente
    err = database.DB.QueryRow(queryMedia, id).Scan(&p.NotaMedia, &p.TotalAvaliacoes)
    if err != nil {
        return p, err
    }

    queryAvaliacoes := `SELECT id, nota, comentario, nome_avaliador, created_at 
                        FROM avaliacoes 
                        WHERE projeto_id = $1 ORDER BY created_at DESC`
    rowsAvaliacoes, err := database.DB.Query(queryAvaliacoes, id)
    if err != nil {
        return p, err
    }
    defer rowsAvaliacoes.Close()

    for rowsAvaliacoes.Next() {
        var a Avaliacao
        if err := rowsAvaliacoes.Scan(&a.ID, &a.Nota, &a.Comentario, &a.NomeAvaliador, &a.CreatedAt); err != nil {
            log.Println("Erro ao escanear avaliação:", err)
            continue
        }
        p.Avaliacoes = append(p.Avaliacoes, a)
    }
    // Verifica se ocorreu um erro durante a iteração do loop de avaliações
    if err = rowsAvaliacoes.Err(); err != nil {
        return p, err
    }

    return p, nil
}

// GetAllProjetos busca todos os projetos (sem detalhes da equipe, para a lista principal)
func GetAllProjetos() ([]Projeto, error) {
    // Esta query usa LEFT JOIN para incluir projetos mesmo que não tenham nenhuma avaliação.
    // COALESCE(AVG(a.nota), 0) garante que a nota média seja 0 em vez de nula se não houver avaliações.
    query := `
        SELECT
            p.id,
            p.titulo,
            p.descricao,
            p.imagem_url,
            COALESCE(AVG(a.nota), 0) as nota_media,
            COUNT(a.nota) as total_avaliacoes
        FROM
            projetos p
        LEFT JOIN
            avaliacoes a ON p.id = a.projeto_id
        GROUP BY
            p.id
        ORDER BY
            nota_media DESC, p.created_at DESC`

    rows, err := database.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var projetos []Projeto
    for rows.Next() {
        var p Projeto
        // O Scan deve seguir a ordem exata das colunas no SELECT
        if err := rows.Scan(&p.ID, &p.Titulo, &p.Descricao, &p.ImagemURL, &p.NotaMedia, &p.TotalAvaliacoes); err != nil {
            log.Println("Erro ao escanear projeto:", err)
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

type Avaliacao struct {
    ID             int
    ProjetoID      int
    Nota           int
    Comentario     string
    NomeAvaliador  string
    CreatedAt      time.Time
}

// CreateAvaliacao salva uma nova avaliação no banco de dados.
func CreateAvaliacao(avaliacao *Avaliacao) error {
    query := `INSERT INTO avaliacoes (projeto_id, nota, comentario, nome_avaliador)
            VALUES ($1, $2, $3, $4)`
    
    _, err := database.DB.Exec(query, avaliacao.ProjetoID, avaliacao.Nota, avaliacao.Comentario, avaliacao.NomeAvaliador)
    return err
}