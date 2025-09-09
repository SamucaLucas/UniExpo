package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)
var DB *sql.DB

// InitDB inicializa a conexão com o banco de dados PostgreSQL.
func InitDB(){
	connStr := "user=teste_doe dbname=teste_doenet host=teste-doenet.postgres.uhserver.com password=Samuca!2004} sslmode=disable search_path=teste_doenet"
	
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
