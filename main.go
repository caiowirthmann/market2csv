package main

// TODO:
// Adicionar agentes para fezer a rotação do srape e não cair no limitador também
// Convertes valores que estão em string para numerico

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"market2csv/assets"
	"market2csv/utils"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

// dados dos anuncios
type anuncioML struct {
	precoBase, precoAtual float64
	descricao             string
	titulo                string
	link                  string
	nota                  string
	quantidadeReviews     string
	patrocinado           string
	full                  string
	freteGratis           string // possivelmente cortar ou deixar para implementar futuramente. Depende do CEP, local vendedor, se está logado ou não
	estoque               string // trocado temporariamente para debug
	quantidadeVendas      string
	// marca                 string // Temporariamente desativado

	vendedorML

	// - Ficha tecnica ==> entender como implementar esse campo, mas por enquanto é só uma idéia

}

// dados vendedor
type vendedorML struct {
	nome         string
	linkVendedor string
}

// resultados busca
type resultadoInicialPesquisaML struct {
	quantidadeResultados int64       // quantidade de anuncios encontrados com o termo da busca
	anunciosColetados    int         // quantidade de anuncios que passaram pelo scrape
	limiteAnuncios       int         // quantidade de anuncios que serao analisados (definida pelo usuario)
	buscaRelacionada     []string    // campo superior indicando outras palavras-chaves relacionadas a pesquisa
	anuncios             []anuncioML // slice de anuncios da pagina de resultado
}

// const linkBase = "https://lista.mercadolivre.com.br/xxx-xxxxx-xxxxxx-xxxxxx#D[A:xxx xx]"

// Gera o link de pesquisa no ML, que será usado no crawler
func criarQueryPesquisa(termoPesquisa string) string {
	var pesquisaURL, parteInicial string
	parteInicial = strings.Replace(termoPesquisa, " ", "-", -1)
	pesquisaURL = fmt.Sprintf("https://lista.mercadolivre.com.br/%s#D[A:%s]", parteInicial, termoPesquisa)
	return pesquisaURL
}

// une os valores do anuncio. Caso não tenha centavos, passar a segunda string vazia que a função converte em 00
func ConverterPrecoFloat(s1, s2 string) (vlrAnuncio float64, err error) {
	if s2 == "" {
		s2 = "00"
	}
	//para casos em que o valor é maior que 999,99, em que o separador de milhar local é ponto
	s1 = strings.Replace(s1, ".", "", -1)

	s := s1 + "." + s2
	if vlrAnuncio, err := strconv.ParseFloat(s, 64); err == nil {
		return vlrAnuncio, nil
	}
	utils.LogarErroFunc("ConverterPrecoFloat", map[string]any{
		"s1":         s1,
		"s2":         s2,
		"vlrAnuncio": vlrAnuncio,
	}, err)
	return 0.0, fmt.Errorf("não foi possível converter as STRINGS %s e %s do preço para um FLOAT. Cheque o log gerado para mais detalhes", s1, s2)
}

// formata e transforma string (xxx resultados) em int
func TratarQtdResultados(resultados string) (qtdResult int64, err error) {
	s := strings.Replace(resultados, ".", "", -1)
	s2 := strings.Split(s, " ")
	x, err := strconv.ParseInt(s2[0], 0, 0)
	if err != nil {
		utils.LogarErroFunc("TratarQtdResultados", map[string]any{
			"resultados": resultados,
			"s":          s,
			"s2":         s2,
		}, err)
		return 0, fmt.Errorf("erro ao tratar a string [%s] para coletar somente o numero. Cheque o log gerado para mais detalhes", resultados)
	}
	return x, nil
}

// funcao que trata a string que contem a qtd de vendas --> vem no formato {CONDIÇÃO | QTD VENDAS}. Ex: "Novo  |  +1000 vendidos". Remove texto
// e retorna string xxx vendidos. Como o ml fornece a qtd de vendas por uma range, não faz muito sentido cortar o +... e converter em um int
// já que ex: Um anuncio com +25 vendas (que pode ser 25 até 49), se convertido ficaria 25, não seria "preciso" por conta da range
//
// Melhor um dado qualitativo preciso do que um quantitativo impreciso
// TODO: converter função para retornar um int
func tratarQtdVendas(textoQtdVendas string) (qtdVendas string) {
	// quando não tem vendas, fica no formato {CONDICAO}, não tem |
	if strings.Contains(textoQtdVendas, "|") {
		s := strings.Split(textoQtdVendas, "|")
		s[1] = strings.Replace(s[1], "+", "", -1)
		corte := strings.Index(s[1], "v") //indice da palavra "vendidas"
		s[1] = strings.Replace(s[1][:corte], " ", "", -1)
		return s[1]
	}
	return "0"
}

// check se é patrocinado pela URL. Na query tem a tag is_advertising=true, indicando que teve impulsionamento pelo mercado ADS
func (a *anuncioML) Patrocinado() {
	if strings.Contains(a.link, "is_advertising=true") {
		a.patrocinado = "sim"
	} else {
		a.patrocinado = "não"
	}
}

// check se tem frete gratis --> não faz muito sentido, já que o frete pode variar por região e outros fatores que podem variar se estiver logado ou não, se tem item no carrinho
// se tem algum cupom, promoção ativa... Mas por enquanto pelo menos da um "norte" sobre o frete do anuncio
// talvez essa função seja cortada, por enquanto ta aqui
func (a *anuncioML) GetFrete(prod colly.HTMLElement) {
	if temFreteGratis := prod.ChildText(".poly-card__content div.poly-component__shipping"); len(temFreteGratis) != 0 {
		a.freteGratis = "Sim"
	} else {
		a.freteGratis = "Não"
	}
}

// check se tem full pela existencia do texto "enviado pelo", já texto full é um .svg que precede o texto
func (a *anuncioML) Full(prod colly.HTMLElement) {
	if temFull := prod.ChildText(".ui-pdp-promotions-pill-label__text"); len(temFull) != 0 {
		a.full = "sim"
	} else {
		a.full = "Não"
	}
}

// check para nota do anuncio, já que anuncio pode não ter rating disponivel
func (a *anuncioML) NotaAvaliacao(prod colly.HTMLElement) {
	if rating := prod.ChildText(".ui-pdp-review__rating"); len(rating) == 0 {
		a.nota = "Sem nota"
	} else {
		a.nota = rating
	}
}

// check para quantidade reviews, mesmo caso da rating
func (a *anuncioML) GetQtdReview(prod colly.HTMLElement) {
	if qtdReviews := prod.ChildText(".ui-pdp-review__amount"); len(qtdReviews) == 0 {
		a.quantidadeReviews = "Sem reviews"
	} else {
		a.quantidadeReviews = qtdReviews[1 : len(qtdReviews)-1] // para remover () da string da qtd de reviews
	}
}

// Check inicial de PREÇO COM DESCONTO ou PREÇO ATUAL DO ANUNCIO (anuncio que NÃO tem desconto).
// Valor é "construido" na pagina do ML por 2 elementos: MONEY-AMOUNT_FRACTION e MONEY-AMOUNT_CENTS. Se for um preço "cheio", não tem o cents.
// Por isso a função constrói o valor primeiro pegando o FRACTION e depois checando a existencia do cents, criando a string e tratando ela com a função ConverterPrecoFloat()
func (a *anuncioML) PrecoAtual(prod colly.HTMLElement) {
	precoAtual := prod.DOM.Find(".ui-pdp-price__second-line span.andes-money-amount__fraction").First().Text()
	if checkprecoAtualCentavos := prod.DOM.Find(".ui-pdp-price__second-line span.andes-money-amount__cents").First().Text(); len(checkprecoAtualCentavos) != 0 { //caso tenha centavos
		precoAtualConvertido, err := ConverterPrecoFloat(precoAtual, checkprecoAtualCentavos)
		if err != nil {
			utils.LogarErroFunc("PrecoAtual", map[string]any{
				"precoAtual":              precoAtual,
				"checkPrecoAtualCentavos": checkprecoAtualCentavos,
				"precoAtualConvertido":    precoAtualConvertido,
				"linkAnuncio":             a.link,
			}, err)
			fmt.Printf("Erro ao extrair o valor atual do anuncio. Cheque o log gerado para mais detalhes")
		} else {
			a.precoAtual = precoAtualConvertido
		}
	} else {
		a.precoAtual, _ = ConverterPrecoFloat(precoAtual, "")
	}
}

// Check inicial de PREÇO ORIGINAL DO ANUNCIO (anuncio que tem desconto). Por "original", entende-se o preço base do anuncio ANTES dos descontos.
// O Mercado Livre só mostra essa linha quando existe algum desconto no anuncio. Por isso em casos no qual NÃO existe desconto, esse valor é igual ao precoAtual
// Para fins analíticos, faz mais sentido manter os dois valores iguais do que colocar 0(zero). Ex: caso for calculado um percentual de desconto, o calculo seria feito errado se não rolasse tratamento na função de desconto.
// Valor é "construido" na pagina do ML por 2 elementos: MONEY-AMOUNT_FRACTION e MONEY-AMOUNT_CENTS. Se for um preço "cheio", não tem o cents.
// Por isso a função constrói o valor primeiro pegando o FRACTION e depois checando a existencia do cents, criando a string e tratando ela com a função ConverterPrecoFloat()
func (a *anuncioML) PrecoBase(prod colly.HTMLElement) {
	if checkPrecoBase := prod.DOM.Find(".ui-pdp-price__original-value span.andes-money-amount__fraction").First().Text(); len(checkPrecoBase) != 0 {
		precoBaseConvertido, err := ConverterPrecoFloat(checkPrecoBase, prod.DOM.Find(".ui-pdp-price__original-value span.andes-money-amount__cents").First().Text())
		if err != nil {
			utils.LogarErroFunc("PrecoBase", map[string]any{
				"checkPrecoBase":      checkPrecoBase,
				"precoBaseConvertido": precoBaseConvertido,
				"linkanuncio":         a.link,
			}, err)
			fmt.Printf("Erro ao extrair o valor base do anuncio. Cheque o log gerado para mais detalhes")
		} else {
			a.precoBase = precoBaseConvertido
		}
	} else {
		a.precoBase = a.precoAtual
	}
}

// Tratar e remover o {CONDICAO} | ... do texto da venda
// É mantido como string porque o Mercado Livre só disponibiliza a quantidade de vendas por uma range.
// Olhar função tratarQtdvendas() para explicação das ranges
func (a *anuncioML) QtdVendas(prod colly.HTMLElement) {
	s := tratarQtdVendas(prod.ChildText("span.ui-pdp-subtitle"))
	a.quantidadeVendas = strings.Replace(s, "mil", "000", -1)
}

// Trata string e remove "Vendido por "
func (a *anuncioML) GetVendedor(prod colly.HTMLElement) {
	prefixo := "Vendido por "
	vendedor := prod.ChildText(".ui-seller-data-header__title-container")
	if strings.Contains(vendedor, prefixo) {
		vendedor = strings.Replace(vendedor, prefixo, "", -1)
		a.vendedorML.nome = vendedor
	}
	a.vendedorML.nome = vendedor
}

// Pega link do vendedor do produto no ML
func (a *anuncioML) VendedorLink(prod colly.HTMLElement) {
	a.vendedorML.linkVendedor = prod.Request.AbsoluteURL(prod.ChildAttr("div.ui-seller-data-footer__container a", "href"))
}

// Alem do numero, pode aparecer "Ultimo disponível" --> Nesse caso irá ser transformado para 1
// Como o texto está envolvido por ( ), é removido por filtrar o 1º e ultimo caracter da string. E isso só acontece se não for o ultimo disponível
// Função busca as duas tags já que ML traz em lugares diferentes a informação caso seja o ultimo em estoque (genial isso kkkk)
func (a *anuncioML) Estoque(prod colly.HTMLElement) {
	estoqueNaoUltimo := prod.ChildText(".ui-pdp-buybox__quantity__available") // caso estoque > 1, vai ter a string (x disponíveis), e acima de 5 começa a mostrar por range com +x disponível (kkkk)
	estoqueUltimo := prod.ChildText(".ui-pdp-buybox__quantity p")             // caso seja o ultimo em estoque, essa string sera "Último disponível!", e caso não seja, vai estar em branco
	ultimoEstoque := "Último disponível!"

	if len(estoqueUltimo) != 0 && estoqueUltimo == ultimoEstoque {
		a.estoque = "1"
		return
	}
	if estoqueNaoUltimo == "" && estoqueUltimo == "" {
		a.estoque = "0"
		return
	} else {
		// ta "funcionado", mas precisa de mais testes
		// utils.LogarErroFunc("estoque", map[string]any{
		// 	"estoqueNaoUltimo": estoqueNaoUltimo,
		// 	"estoqueUltimo":    estoqueUltimo,
		// 	"anuncio":          a.link,
		// }, nil)
		estoqueNaoUltimo = strings.Replace(estoqueNaoUltimo, "+", "", -1)
		s := strings.Split(estoqueNaoUltimo, " ")
		a.estoque = s[0][1:]
	}
}

// Pega descrição do anuncio. Vem com algumas formatações html simplificada mas por enquanto é relevada
func (a *anuncioML) Descricao(prod colly.HTMLElement) {
	a.descricao = prod.ChildText(".ui-pdp-description__content")
}

// Cria o csv em uma pasta "extracoes" com os dados extraidos dos anuncios da pesquisa
func ExportarCSV(buscaML string, anuncios []anuncioML) error {
	/*
	   caso for converter para binário e quiser travar onde a pasta será criada, modificar para o seguinte código

	   // Caminho relativo à pasta do executável
	   	execPath, err := os.Executable()
	   	var pastaExtracoes string
	   	if err != nil {
	   		// Fallback: usa o diretório atual
	   		fmt.Println("⚠️  Aviso: não foi possível detectar o caminho do executável. Usando diretório atual.")
	   		pastaExtracoes = "extracoes"
	   	} else {
	   		execDir := filepath.Dir(execPath)
	   		pastaExtracoes = filepath.Join(execDir, "extracoes")
	   	}

	   	// Garante que a pasta existe
	   	if err := os.MkdirAll(pastaExtracoes, os.ModePerm); err != nil {
	   		return fmt.Errorf("erro ao criar pasta extracoes: %v", err)
	   	}

	   	// Caminho completo para o arquivo
	   	caminhoCompleto := filepath.Join(pastaExtracoes, nomeArquivo)

	   	// Cria o arquivo CSV
	   	arquivo, err := os.Create(caminhoCompleto)
	*/
	const pastaDestino = "extracoes"
	if err := os.MkdirAll(pastaDestino, os.ModePerm); err != nil {
		return fmt.Errorf("não foi possível criar a pasta [%s]. Erro %v", pastaDestino, err)
	}

	// define nome do arquivo. Por padrão: o-que-foi-pesquisado_no_ML_data-hoje_hora-hoje
	nomePesquisa := strings.Replace(buscaML, " ", "-", -1)
	dataExecucao := time.Now().Format("02-01-2006_15-04")

	nomeArquivo := fmt.Sprintf("%s_%s.csv", nomePesquisa, dataExecucao)
	caminhoArquivo := fmt.Sprintf("%s/%s", pastaDestino, nomeArquivo)

	// criar arquivo
	arquivo, err := os.Create(caminhoArquivo)
	if err != nil {
		return fmt.Errorf("não foi possível criar o arquivo [%s]. Erro: %v", nomeArquivo, err)
	}
	defer arquivo.Close()

	writer := csv.NewWriter(arquivo)
	writer.Comma = ';' // separador é ; para evitar problema com formatação de numeros e algum texto que possivelmente possa ter ',' nele
	defer writer.Flush()
	cabecalho := []string{"titulo", "preco_base", "preco_atual", "quantidade_vendas",
		"estoque", "patrocinado", "tem_full", "nota", "quantidade_reviews", "link_anuncio",
		"vendedor", "vendedor_link"}
	if err := writer.Write(cabecalho); err != nil {
		return fmt.Errorf("erro ao adicionar cabeçalho ao csv: [%v]", err)
	}
	for _, anuncio := range anuncios {
		linha := []string{
			anuncio.titulo,
			strconv.FormatFloat(anuncio.precoBase, 'f', 2, 64),
			strconv.FormatFloat(anuncio.precoAtual, 'f', 2, 64),
			anuncio.quantidadeVendas,
			anuncio.estoque,
			anuncio.patrocinado,
			anuncio.full,
			anuncio.nota,
			anuncio.quantidadeReviews,
			anuncio.link,
			anuncio.vendedorML.nome,
			anuncio.vendedorML.linkVendedor,
		}
		if err := writer.Write(linha); err != nil {
			return fmt.Errorf("erro ao escrever linha para o csv [%v]", err)
		}
	}
	return nil
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

	resultadoScrapper := resultadoInicialPesquisaML{
		anuncios: []anuncioML{},
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
		s, err := TratarQtdResultados(qtdResultados)
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
		anuncio := anuncioML{
			titulo: prod.ChildText(".ui-pdp-title"),
			link:   prod.Request.URL.String(),
		}

		anuncio.Patrocinado()
		anuncio.GetFrete(*prod)
		anuncio.Full(*prod)
		anuncio.NotaAvaliacao(*prod)
		anuncio.GetQtdReview(*prod)
		anuncio.PrecoAtual(*prod)
		anuncio.PrecoBase(*prod)
		anuncio.QtdVendas(*prod)
		anuncio.GetVendedor(*prod)
		anuncio.VendedorLink(*prod)
		anuncio.Estoque(*prod)
		anuncio.Descricao(*prod)

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

	err := ExportarCSV(termoBuscaML, resultadoScrapper.anuncios)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("CSV com os dados dos anuncios criado com sucesso")
	}
}
