package main

// TODO:
// Adicionar agentes para fezer a rotação do srape e não cair no limitador também
// Convertes valores que estão em string para numerico

import (
	"bufio"
	"fmt"
	"market2csv/assets"
	"market2csv/scraper/mercadolivre"
	"market2csv/utils"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

// resultados busca
type resultadoPesquisaInicial struct {
	quantidadeResultados int64                  // quantidade de anuncios encontrados com o termo da busca
	anunciosColetados    int                    // quantidade de anuncios que passaram pelo scrape
	limiteAnuncios       int                    // quantidade de anuncios que serao analisados (definida pelo usuario)
	buscaRelacionada     []string               // campo superior indicando outras palavras-chaves relacionadas a pesquisa
	anuncios             []mercadolivre.Anuncio // slice de anuncios da pagina de resultado
}

// const linkBase = "https://lista.mercadolivre.com.br/xxx-xxxxx-xxxxxx-xxxxxx#D[A:xxx xx]"

// Gera o link de pesquisa no ML, que será usado no crawler
func criarQueryPesquisa(termoPesquisa string) string {
	var pesquisaURL, parteInicial string
	parteInicial = strings.Replace(termoPesquisa, " ", "-", -1)
	pesquisaURL = fmt.Sprintf("https://lista.mercadolivre.com.br/%s#D[A:%s]", parteInicial, termoPesquisa)
	return pesquisaURL
}

func main() {

	assets.PrintTelaTerminal()
	var inputSolicitado bool = false // condicional para que o não rode mais de uma vez o input solicitando quantos anucnios quer analisar

	começo := time.Now()

	reader := bufio.NewReader(os.Stdin)
	var termoBuscaML string

	for {
		fmt.Printf("Digite a sua pesquisa para gerar um csv com os dados dos anúncios buscados no:\n")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Erro ao ler a busca. Tente novamente...")
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

	// 2 - contruir a URL da pesquisa. Args[0] é o nome do programa, por isso é [1:]
	// pesquisaML := criarArg(args[1:])
	queryPesquisa := criarQueryPesquisa(termoBuscaML)
	// fmt.Println(queryPesquisa)

	// 3 - Scrape dos seguintes dados
	scrapper := colly.NewCollector()
	scrapperDetalhado := colly.NewCollector(
		colly.Async(true), // async = true permite rodar em paralelo e evita race condition
	)
	extensions.RandomUserAgent(scrapperDetalhado)

	// define limite de chamados em paralelo para esse coletor
	scrapperDetalhado.Limit(&colly.LimitRule{
		RandomDelay: 1 * time.Second, // delay para não bater no rate limit do site (ainda não descobri qual o limite, mas é bom não forçar, pra não dar merda quando começar a pegar multi pagina), ou sobrecarregar o site
		// que não tenha uma forma de proteção ou gargalo de request
	})

	resultadoScrapper := resultadoPesquisaInicial{
		anuncios: []mercadolivre.Anuncio{},
	}

	scrapper.OnRequest(func(r *colly.Request) {
		fmt.Println("Realizando a busca no Mercado Livre por:", termoBuscaML, "\nURL da pesquisa:", r.URL)
	})
	scrapper.OnResponse(func(r *colly.Response) {
		fmt.Println("Coletando dados ...")
	})
	scrapper.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Erro: %s", err)
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
			fmt.Printf("Foram encontrados %v resultados para a busca\n\n", s)
			for {
				fmt.Printf("Quantos anuncios deseja analisar? -- Quantidade NÃO pode ser maior que %v --\n\nDigite 0 para analisar todos os anuncios encontrados OU Digite outra quantidade:\n", resultadoScrapper.quantidadeResultados)
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
		}
		if err := e.DOM.Find("ul.ui-search-top-keywords__list").Text(); len(err) != 0 {
			e.ForEach("section.ui-search-top-keywords ul.ui-search-top-keywords__list a", func(i int, keyword *colly.HTMLElement) {
				resultadoScrapper.buscaRelacionada = append(resultadoScrapper.buscaRelacionada, keyword.Text)
			})
		}

		// e.ForEach(".ui-search-layout__item", func(i int, h *colly.HTMLElement) {
		// })
		// fmt.Printf("%#v\n", resultadoScrapper.buscaRelacionada)
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

	})

	scrapperDetalhado.OnHTML("body main", func(prod *colly.HTMLElement) {

		// fmt.Printf("Buscando dados do anuncio [%s]\n", prod.Request.URL.String())
		// começoAnuncio := time.Now()
		anuncio := mercadolivre.NovoAnuncio(prod)

		resultadoScrapper.anuncios = append(resultadoScrapper.anuncios, anuncio)
		// fmt.Printf("Tempo de scrape: %v para o anuncio [%s]\n", time.Since(começoAnuncio), anuncio.link)
	})

	scrapper.OnHTML("li.andes-pagination__button--next a", func(e *colly.HTMLElement) {
		if resultadoScrapper.anunciosColetados >= resultadoScrapper.limiteAnuncios {
			return
		}
		proximaPagina := e.Request.AbsoluteURL(e.Attr("href"))
		fmt.Println("Visitando próxima página")
		e.Request.Visit(proximaPagina)
	})

	// printf com %#v printa no formato field:data para struct
	scrapper.Visit(queryPesquisa)

	// fmt.Printf("%#v\n\n", resultadoScrapper.anuncios[:8])
	fim := time.Since(começo)

	fmt.Printf("Tempo de execução: %v\n", fim)

	err := mercadolivre.ExportarCSV(termoBuscaML, resultadoScrapper.anuncios)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("CSV com os dados dos anuncios do Merado Livre criado com sucesso")
	}
}
