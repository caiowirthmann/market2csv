package main

// TODO:
// Adicionar agentes para fezer a rotação do srape e não cair no limitador também
// Convertes valores que estão em string para numerico

import (
	"bufio"
	"fmt"
	"os"
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
		anuncios: []mercadolivre.Anuncio{},
	}

	scrapper.OnRequest(func(r *colly.Request) {
		fmt.Println("\nRealizando a busca no Mercado Livre por:", termoBuscaML, "\nURL da pesquisa:", r.URL)
	})
	scrapper.OnResponse(func(r *colly.Response) {
		fmt.Println("\nColetando dados...")
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
		}

		if err := e.DOM.Find("ul.ui-search-top-keywords__list").Text(); len(err) != 0 {
			e.ForEach("section.ui-search-top-keywords ul.ui-search-top-keywords__list a", func(i int, keyword *colly.HTMLElement) {
				resultadoScrapper.buscaRelacionada = append(resultadoScrapper.buscaRelacionada, keyword.Text)
			})
		}
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

		// coleta dados do anuncio e "monta" anuncio
		anuncio := mercadolivre.NovoAnuncio(prod)
		resultadoScrapper.anuncios = append(resultadoScrapper.anuncios, anuncio)
	})

	scrapper.OnHTML("li.andes-pagination__button--next a", func(e *colly.HTMLElement) {
		if resultadoScrapper.anunciosColetados >= resultadoScrapper.limiteAnuncios {
			return
		}
		proximaPagina := e.Request.AbsoluteURL(e.Attr("href"))
		fmt.Println("Visitando próxima página")
		e.Request.Visit(proximaPagina)
	})

	scrapper.Visit(queryPesquisa)

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

	fmt.Printf("Tempo de execução: %.2s\n", fim)

}
