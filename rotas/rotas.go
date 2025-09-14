package routers

import (
	controle "modulo/controle"
	"net/http"
)

func CarregadoRotas() {
	// Nossas rotas usando os controllers
	http.HandleFunc("/", controle.HomePage)
	http.HandleFunc("/aluno", controle.AlunoPage)
	http.HandleFunc("/cadastro", controle.CadastroPage)

	http.HandleFunc("/projetos", controle.ProjetosListPage)
	http.HandleFunc("/projeto", controle.ProjetoDetailPage)
	http.HandleFunc("/cadastro-projeto", controle.CadastroProjetoPage)

	http.HandleFunc("/adicionar-avaliacao", controle.AdicionarAvaliacaoHandler)
}