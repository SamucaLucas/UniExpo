package routers

import (
	controllers "modulo/controle"
	"net/http"
)

func CarregadoRotas() {
	// Nossas rotas usando os controllers
	http.HandleFunc("/", controllers.HomePage)
	http.HandleFunc("/aluno", controllers.AlunoPage)
	http.HandleFunc("/cadastro", controllers.CadastroPage)
}