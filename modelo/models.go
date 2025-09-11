// internal/models/aluno.go
package models

import (
	"database/sql"
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
	query := `SELECT id, nome, periodo, instagram, github, foto_url, linguagens, frameworks, ferramentas, bio, linkedin
			FROM alunos WHERE id = $1`

	err := database.DB.QueryRow(query, id).Scan(
		&a.ID, &a.Nome, &a.Periodo, &a.Instagram, &a.Github, &a.FotoURL,
		pq.Array(&a.Linguagens), pq.Array(&a.Frameworks), pq.Array(&a.Ferramentas), &a.Bio, &a.Linkedin,
	)
	if err == sql.ErrNoRows {
		return a, fmt.Errorf("aluno com ID %d não encontrado", id)
	}
	return a, err
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