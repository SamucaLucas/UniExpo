// cmd/web/main.go
package main

import (
	"log"
	database "modulo/db"
	routers "modulo/rotas"
	"net/http"
)

func main() {

	routers.CarregadoRotas()
	// Inicializa a conexão com o banco de dados
	database.InitDB()
	defer database.DB.Close()
	// Servir arquivos estáticos (CSS, JS, imagens)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))



	log.Println("Servidor iniciado em http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}