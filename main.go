package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/caiowirthmann/market2csv/assets"
	"github.com/caiowirthmann/market2csv/scraper/mercadolivre"
	"github.com/caiowirthmann/market2csv/utils"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

// resultados busca
type resultadoPesquisaInicial struct {
	quantidadeResultados int64                  // quantidade de anuncios encontrados com o termo da busca
	anunciosColetados    int                    // quantidade de anuncios que passaram pelo scrape
	limiteAnuncios       int                    // quantidade de anuncios que serao analisados (definida pelo usuario)
	capacidadeAnuncio    int                    // contador para calculo troca de pagina
	indexPesquisa        int                    // numero do link "desde_x" no link da pesquisa
	pagina               int                    // pagina atual
	arvoreCategoria      []string               // armazena os nós da categoria da busca
	anuncios             []mercadolivre.Anuncio // slice de anuncios da pagina de resultado
}

// const linkBase = "https://lista.mercadolivre.com.br/xxx-xxxxx-xxxxxx-xxxxxx#D[A:xxx xx]"

// Gera o link de pesquisa no ML, que será usado no crawler
func criarQueryPesquisa(termoPesquisa string) string {
	var pesquisaURL, parteInicial string
	parteInicial = strings.ReplaceAll(termoPesquisa, " ", "-")
	pesquisaURL = fmt.Sprintf("https://lista.mercadolivre.com.br/%s#D[A:%s]", parteInicial, termoPesquisa)
	return pesquisaURL
}

func main() {

	assets.PrintTelaTerminal()
	// condicional para que não rode mais de uma vez o input solicitando quantos anucnios quer analisar
	// se passar de mais de uma pagina a quantidade
	var inputSolicitado bool = false

	começo := time.Now()

	reader := bufio.NewReader(os.Stdin)
	var termoBuscaML string

	// garante que input do usuario é valido ou que foi passado algo (não pode ser vazio)
	for {
		fmt.Printf("Digite o que quer pesquisar nos marketplaces:\n")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Erro ao ler a pesquisa. Tente novamente...")
			continue
		}
		// é necessario essa manipulação para limpar o '\n' do texto recebido do terminal
		// se for passado nada (""), a função ReadString vai adicionar o '\n' no final do texto, e a len NÃO vai ser 0
		// com o TrimSpace é "limpado" a quebra de linha
		input = strings.TrimSpace(input)
		if len(input) == 0 {
			fmt.Println("Não foi passado nada para ser pesquisado....")
			continue
		}
		termoBuscaML = input
		break
	}

	utils.InitLogErro(termoBuscaML)

	queryPesquisa := criarQueryPesquisa(termoBuscaML)

	scrapper := colly.NewCollector()
	scrapperDetalhado := colly.NewCollector(
		colly.Async(true), // async = true permite rodar em paralelo e evita race condition
	)
	extensions.RandomUserAgent(scrapperDetalhado)

	// define limite de chamados em paralelo para esse coletor
	scrapperDetalhado.Limit(&colly.LimitRule{
		RandomDelay: 2 * time.Second, // delay para não bater no rate limit do site (ainda não descobri qual o limite, mas é bom não forçar, pra não dar merda quando começar a pegar multi pagina), ou sobrecarregar o site
		// que não tenha uma forma de proteção ou gargalo de request
	})

	resultadoScrapper := resultadoPesquisaInicial{
		anuncios:      []mercadolivre.Anuncio{},
		indexPesquisa: 49,
		pagina:        1,
	}

	scrapper.OnRequest(func(r *colly.Request) {
		fmt.Println("\nRealizando a busca no Mercado Livre por:", termoBuscaML, "\nURL da pesquisa:", r.URL)
	})
	scrapper.OnResponse(func(r *colly.Response) {
		fmt.Println("\nColetando dados...")
	})
	scrapper.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Erro: %s\n", err)
	})

	scrapperDetalhado.OnScraped(func(r *colly.Response) {
		fmt.Printf("Dados do anúncio coletado com sucesso...\n")
		fmt.Printf("[%v / %v] anuncios analisados\n\n", resultadoScrapper.anunciosColetados, resultadoScrapper.limiteAnuncios)
	})

	scrapper.OnHTML("body main div.ui-search", func(e *colly.HTMLElement) {

		// Quantidade de resultados encontrados para a pesquisa
		qtdResultados := e.ChildText("aside.ui-search-sidebar div.ui-search-search-result span.ui-search-search-result__quantity-results")
		s, err := mercadolivre.TratarQtdResultados(qtdResultados)
		if err != nil {
			fmt.Println(err)
		}
		resultadoScrapper.quantidadeResultados = s

		if !inputSolicitado {
			fmt.Printf("Foram encontrados --> %v <-- anuncios para [%s]\n\n", s, termoBuscaML)
			for {
				fmt.Printf("Quantos anuncios deseja analisar? Digite um valor entre 1 e %v\n--> Caso queira TODOS os anuncios encontrados, digite 0\n", resultadoScrapper.quantidadeResultados)
				_, err := fmt.Scanln(&resultadoScrapper.limiteAnuncios)
				if err != nil || resultadoScrapper.limiteAnuncios < 0 || resultadoScrapper.limiteAnuncios > int(resultadoScrapper.quantidadeResultados) {
					fmt.Println("Valor inválido. Digite outra opção")
					continue
				}
				if resultadoScrapper.limiteAnuncios == 0 {
					resultadoScrapper.limiteAnuncios = int(resultadoScrapper.quantidadeResultados)
				}
				break
			}
			inputSolicitado = true
			// armazena os links de cada nó da arvore de categoria
			// ultimo link é o link usado para montar a paginação
			// dentro do check para que não seja executado de novo quando trocar de página
			e.ForEach(".andes-breadcrumb__item", func(i int, h *colly.HTMLElement) {
				resultadoScrapper.arvoreCategoria = append(resultadoScrapper.arvoreCategoria, h.ChildAttr("a.andes-breadcrumb__link", "href"))
			})
		}

		// contador para identificar quando trocar de pagina e limite da paginação
		e.ForEach(".poly-component__title", func(i int, h *colly.HTMLElement) {
			resultadoScrapper.capacidadeAnuncio++
		})
	})

	// Coleta de dados do anuncio
	scrapper.OnHTML(".poly-component__title", func(h *colly.HTMLElement) {
		if resultadoScrapper.anunciosColetados >= resultadoScrapper.limiteAnuncios {
			return
		}
		linkAnuncio := h.Request.AbsoluteURL(h.Attr("href"))
		resultadoScrapper.anunciosColetados++
		scrapperDetalhado.Visit(linkAnuncio)
		scrapperDetalhado.Wait()

		// TODO: continuar analise
		// formula p/ pegar o indexPagina (exceto 1ª pagina) = 49+(pagina*48)
		if resultadoScrapper.anunciosColetados >= resultadoScrapper.capacidadeAnuncio {
			incremento := 48
			queryPaginação := strings.ReplaceAll(termoBuscaML, " ", "-")
			// 2a pagina de resultados (1a que é com link "construido") tem essa estrutura estranha
			// 2a é Desde_49_xxxx
			// 3a em diante é 49+(pagina*48
			// exemplo: vv
			// https://lista.mercadolivre.com.br/saude/suplementos-alimentares/whey-morango-dux_Desde_49_NoIndex_True
			// https://lista.mercadolivre.com.br/saude/suplementos-alimentares/whey-morango-dux_Desde_97_NoIndex_True
			if resultadoScrapper.pagina == 1 {
				linkProximaPagina := fmt.Sprintf("%s/%s_Desde_%v_NoIndex_True", resultadoScrapper.arvoreCategoria[len(resultadoScrapper.arvoreCategoria)-1],
					queryPaginação, resultadoScrapper.indexPesquisa)
				resultadoScrapper.pagina++
				fmt.Println("Visitando próxima página")
				scrapper.Visit(linkProximaPagina)
			} else {
				linkProximaPagina := fmt.Sprintf("%s/%s_Desde_%v_NoIndex_True", resultadoScrapper.arvoreCategoria[len(resultadoScrapper.arvoreCategoria)-1],
					queryPaginação, resultadoScrapper.indexPesquisa+(resultadoScrapper.pagina*incremento))
				resultadoScrapper.pagina++
				fmt.Println("Visitando próxima página")
				scrapper.Visit(linkProximaPagina)
			}
		}
	})

	scrapperDetalhado.OnHTML("body main", func(prod *colly.HTMLElement) {

		// coleta dados do anuncio e "monta" anuncio
		anuncio := mercadolivre.NovoAnuncio(prod)
		resultadoScrapper.anuncios = append(resultadoScrapper.anuncios, anuncio)
	})

	scrapper.Visit(queryPesquisa)
	// scrapperDetalhado.Wait()

	err := mercadolivre.ExportarCSV(termoBuscaML, resultadoScrapper.anuncios)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("CSV com os dados dos anuncios do Merado Livre criado com sucesso")
	}

	err2 := mercadolivre.ExportarFichaTecnica(termoBuscaML, resultadoScrapper.anuncios)
	if err2 != nil {
		fmt.Println(err2)
	} else {
		fmt.Println("Arquivo JSON com a ficha tecnica de cada anuncio criado com sucesso")
	}

	fim := time.Since(começo)

	fmt.Printf("Tempo de execução: %s\n", fim)

	// apenas para garantir que a janela do terminal não fecha automaticamente
	// após a execução da ferramente, já que no windows se rodar dando duplo-cliue no .exe
	// existe esse comportamento, já que teoricamente a aplicação encerrou
	// unix não rola porque é executado pelo terminal
	if runtime.GOOS == "windows" {
		fmt.Println("\nPressione ENTER para sair...")
		fmt.Scanln()
	}

}
