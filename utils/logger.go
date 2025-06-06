package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	LogErro *log.Logger
)

// TODO: Quando for implementado os outros marketplaces, sera adicionado um parametro na função para identificar o marketplace
func InitLogErro(pesquisaML string) {
	// cria a pasta de logs
	var pastaLogs = "logs"

	caminhoExecutavel, err := os.Executable()
	if err != nil {
		fmt.Println("Erro ao identificar o caminho do executável:", err)
	}
	caminhoReal, err := filepath.EvalSymlinks(caminhoExecutavel)
	if err != nil {
		fmt.Println("Erro ao resolver symlink do executável:", err)
	}

	var dirBase string
	if strings.HasPrefix(caminhoReal, os.TempDir()) {
		// caso de execução da ferramente NÃO compilada
		dirBase = "."
	} else {
		dirBase = filepath.Dir(caminhoReal)
	}

	dirLogs := filepath.Join(dirBase, pastaLogs)
	if err := os.MkdirAll(dirLogs, os.ModePerm); err != nil {
		log.Fatalf("Erro ao criar a pasta de logs. Erro [%v]", err)
	}
	// cria arquivo de log unico da execução
	nomePesquisa := strings.ReplaceAll(pesquisaML, " ", "-")
	dataExecucao := time.Now().Format("02-01-2006_15-04")
	nomeArquivo := fmt.Sprintf("%s_%s.log", nomePesquisa, dataExecucao)

	caminhoCompleto := filepath.Join(dirLogs, nomeArquivo)
	arquivo, err := os.OpenFile(caminhoCompleto, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Erro ao gerar o arquivo de log da ferramenta")
	}
	LogErro = log.New(arquivo, "ERRO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// funcao para detalhar o log com nome funcao e os parametros que foram passador para a função
// Depois pode ser testado no debugger para entender o que rolou
func LogarErroFunc(nomeFuncao string, parametrosFunc map[string]any, err error) {
	if LogErro == nil {
		fmt.Println("Log não inicializado")
		return
	}
	// construtor
	parametrosStr := ""
	for k, v := range parametrosFunc {
		parametrosStr += fmt.Sprintf("[%s=%s]|| ", k, v)
	}

	LogErro.Printf("Função: %s || Parametros: [%s] || Erro %v\n", nomeFuncao, parametrosStr, err)
}
