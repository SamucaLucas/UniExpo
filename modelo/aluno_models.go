package modelo

import (
	"fmt"
	"log"
	database "modulo/db"

	"github.com/lib/pq" // Pacote necessário para lidar com arrays do PostgreSQL
)

type Aluno struct {
	ID          int
	Nome        string
	Periodo     int
	Instagram   string
	Github      string
	Linkedin	string
	FotoURL     string
	Linguagens  []string
	Frameworks  []string
	Ferramentas []string
	Bio         string
	ProjetosParticipa []Projeto
}

// GetAll busca todos os alunos no banco de dados, com filtros.
func GetAll(periodoFilter int, skillFilter string) ([]Aluno, error) {
	query := `SELECT id, nome, periodo, foto_url, linguagens FROM alunos WHERE 1=1`
	args := []interface{}{}
	argID := 1

	if periodoFilter > 0 {
		query += fmt.Sprintf(" AND periodo = $%d", argID)
		args = append(args, periodoFilter)
		argID++
	}

	if skillFilter != "" {
		query += fmt.Sprintf(" AND $%d = ANY(linguagens)", argID)
		args = append(args, skillFilter)
	}

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alunos []Aluno
	for rows.Next() {
		var a Aluno
		err := rows.Scan(&a.ID, &a.Nome, &a.Periodo, &a.FotoURL, pq.Array(&a.Linguagens))
		if err != nil {
			log.Println("Erro ao escanear aluno:", err)
			continue
		}
		alunos = append(alunos, a)
	}

	return alunos, nil
}

// GetByID busca um único aluno pelo seu ID.
func GetByID(id int) (Aluno, error) {
    var a Aluno

    // PARTE 1: Busca os dados principais do aluno (como já fazia antes)
    queryAluno := `SELECT id, nome, periodo, instagram, github, linkedin, foto_url, linguagens, frameworks, ferramentas, bio
                FROM alunos WHERE id = $1`

    err := database.DB.QueryRow(queryAluno, id).Scan(
        &a.ID, &a.Nome, &a.Periodo, &a.Instagram, &a.Github, &a.Linkedin, &a.FotoURL,
        pq.Array(&a.Linguagens), pq.Array(&a.Frameworks), pq.Array(&a.Ferramentas), &a.Bio,
    )
    if err != nil {
        return a, err
    }

    // PARTE 2 (NOVA): Busca os projetos associados a este aluno
    queryProjetos := `SELECT p.id, p.titulo, p.imagem_url
                    FROM projetos p
                    JOIN equipe_projetos ep ON p.id = ep.projeto_id
                    WHERE ep.aluno_id = $1`

    rows, err := database.DB.Query(queryProjetos, id)
    if err != nil {
        // Se der erro aqui, não quebra a página, apenas logamos o erro.
        // O perfil do aluno ainda será exibido, mas sem os projetos.
        log.Printf("Erro ao buscar projetos do aluno %d: %v", id, err)
        return a, nil // Retornamos 'nil' para o erro, pois os dados principais do aluno foram carregados
    }
    defer rows.Close()

    // Itera sobre os resultados da busca de projetos
    for rows.Next() {
        var p Projeto // Usamos a struct Projeto, mas só preenchemos alguns campos
        if err := rows.Scan(&p.ID, &p.Titulo, &p.ImagemURL); err != nil {
            log.Println("Erro ao escanear projeto do aluno:", err)
            continue
        }
        // Adiciona o projeto encontrado à lista de projetos do aluno
        a.ProjetosParticipa = append(a.ProjetosParticipa, p)
    }

    return a, nil
}

// Create insere um novo aluno no banco de dados.
func (a *Aluno) Create() error {
	query := `INSERT INTO alunos (nome, periodo, instagram, github, foto_url, linguagens, frameworks, ferramentas, bio, linkedin)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id`

	err := database.DB.QueryRow(
		query,
		a.Nome, a.Periodo, a.Instagram, a.Github, a.FotoURL,
		pq.Array(a.Linguagens), pq.Array(a.Frameworks), pq.Array(a.Ferramentas), a.Bio, a.Linkedin,
	).Scan(&a.ID)

	return err
}