package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)
var DB *sql.DB

// InitDB inicializa a conexão com o banco de dados PostgreSQL.
func InitDB() {
    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        log.Fatal("A variável de ambiente DATABASE_URL não foi definida.")
    }	
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
