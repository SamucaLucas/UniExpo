package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)
var DB *sql.DB

// InitDB inicializa a conexão com o banco de dados PostgreSQL.
func InitDB() {
    connStr := "user=samucael dbname=uniexpo host=72.60.149.222 password=2784 sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro ao abrir a conexão com o banco de dados:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Erro ao conectar com o banco de dados:", err)
	}

	fmt.Println("Conexão com o PostgreSQL estabelecida com sucesso!")
}
