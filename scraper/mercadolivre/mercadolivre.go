package mercadolivre

// package para scrape do mercado livre
// contem todas as funções exclusivas ao tratamento dos dados e estrutura do scrape

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/caiowirthmann/market2csv/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

// dados dos anuncios
type Anuncio struct {
	precoBase, precoAtual float64
	descricao             string
	titulo                string
	link                  string
	nota                  string
	quantidadeReviews     string
	patrocinado           string
	full                  string
	estoque               string
	condicao              string
	quantidadeVendas      string

	vendedor
	FichaTecnica

	// - Ficha tecnica ==> entender como implementar esse campo, mas por enquanto é só uma idéia

}

// dados vendedor
type vendedor struct {
	nome         string
	tipoLoja     string
	linkVendedor string
}

type FichaTecnica struct {
	caracteristicas map[string]any
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
// Melhor um dado qualitativo preciso do que um quantitativo impreciso
func TratarQtdVendas(textoQtdVendas string) (qtdVendas string, err error) {
	// quando não tem vendas, fica no formato {CONDICAO}, não tem |
	if strings.Contains(textoQtdVendas, "|") {
		s := strings.Split(textoQtdVendas, "|")
		s[1] = strings.Replace(s[1], "+", "", -1)
		corte := strings.Index(s[1], "v") //indice da palavra "vendidas"
		s[1] = strings.Replace(s[1][:corte], " ", "", -1)
		s[1] = strings.Replace(s[1], "mil", "000", -1)

		return s[1], nil
	}
	return "0", fmt.Errorf("anuncio sem vendas, ou produto de catálogo(pag. sem info vendedor) ou layout alterado")
}

// check se é patrocinado pela URL. Na query tem a tag is_advertising=true, indicando que teve impulsionamento pelo mercado ADS
// só é possível checar pela URL da pagina de resultados
func (a *Anuncio) temPatrocinado() {
	if strings.Contains(a.link, "is_advertising=true") {
		a.patrocinado = "sim"
	} else {
		a.patrocinado = "não"
	}
}

// tipo da loja do vendedor é possível pegar pelo link do vendedor e as tags da query do link
// Na função construtora do Anuncio, ela tem que vir APÓS a func vendedorLink()
func (a *Anuncio) tipoLojaVendedor() {

	tipos := map[string]string{
		"typeSeller=official_store": "Loja oficial",
		"typeSeller=eshop":          "Eshop",
		"typeSeller=classic":        "Padrão",
	}

	for chave, tipo := range tipos {
		if strings.Contains(a.linkVendedor, chave) {
			a.tipoLoja = tipo
			return
		}
	}
	utils.LogarErroFunc("tipoLojaVendedor - DESCONHECIDO", map[string]any{
		"vendedor": a.vendedor.nome,
		"linkLoja": a.linkVendedor,
	}, fmt.Errorf("não foi possível determinar o tipo de loja do vendedor. Por padrão o campo ficará em branco"))
}

// check se tem full pela existencia do texto "enviado pelo", já texto full é um .svg que precede o texto
func (a *Anuncio) temFull(prod *colly.HTMLElement) {
	if temFull := prod.ChildText(".ui-pdp-promotions-pill-label__text"); len(temFull) != 0 {
		a.full = "sim"
	} else {
		a.full = "não"
	}
}

// check para nota do anuncio, já que anuncio pode não ter rating disponivel
func (a *Anuncio) notaAvaliacao(prod *colly.HTMLElement) {
	if rating := prod.ChildText(".ui-pdp-review__rating"); len(rating) == 0 {
		a.nota = "Sem nota"
	} else {
		a.nota = rating
	}
}

// check para quantidade reviews, mesmo caso da rating
func (a *Anuncio) qtdAvaliacoes(prod *colly.HTMLElement) {
	if qtdReviews := prod.ChildText(".ui-pdp-review__amount"); len(qtdReviews) == 0 {
		a.quantidadeReviews = "Sem reviews"
	} else {
		a.quantidadeReviews = qtdReviews[1 : len(qtdReviews)-1] // para remover () da string da qtd de reviews
	}
}

// Check inicial de PREÇO COM DESCONTO ou PREÇO ATUAL DO ANUNCIO (anuncio que NÃO tem desconto).
// Valor é "construido" na pagina do ML por 2 elementos: MONEY-AMOUNT_FRACTION e MONEY-AMOUNT_CENTS. Se for um preço "cheio", não tem o cents.
// Por isso a função constrói o valor primeiro pegando o FRACTION e depois checando a existencia do cents, criando a string e tratando ela com a função ConverterPrecoFloat()
func (a *Anuncio) montarPrecoAtual(prod *colly.HTMLElement) {
	precoAtual := prod.DOM.Find(".ui-pdp-price__second-line span.andes-money-amount__fraction").First().Text()
	if checkprecoAtualCentavos := prod.DOM.Find(".ui-pdp-price__second-line span.andes-money-amount__cents").First().Text(); len(checkprecoAtualCentavos) != 0 { //caso tenha centavos
		precoAtualConvertido, err := ConverterPrecoFloat(precoAtual, checkprecoAtualCentavos)
		if err != nil {
			utils.LogarErroFunc("precoAtual", map[string]any{
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
		x, err := ConverterPrecoFloat(precoAtual, "")
		if err != nil {
			fmt.Println(err)
		}
		a.precoAtual = x
	}
}

// Check inicial de PREÇO ORIGINAL DO ANUNCIO (anuncio que tem desconto). Por "original", entende-se o preço base do anuncio ANTES dos descontos.
// O Mercado Livre só mostra essa linha quando existe algum desconto no anuncio. Por isso em casos no qual NÃO existe desconto, esse valor é igual ao precoAtual
// Para fins analíticos, faz mais sentido manter os dois valores iguais do que colocar 0(zero). Ex: caso for calculado um percentual de desconto, o calculo seria feito errado se não rolasse tratamento na função de desconto.
// Valor é "construido" na pagina do ML por 2 elementos: MONEY-AMOUNT_FRACTION e MONEY-AMOUNT_CENTS. Se for um preço "cheio", não tem o cents.
// Por isso a função constrói o valor primeiro pegando o FRACTION e depois checando a existencia do cents, criando a string e tratando ela com a função ConverterPrecoFloat()
func (a *Anuncio) montarPrecoBase(prod *colly.HTMLElement) {
	if checkPrecoBase := prod.DOM.Find(".ui-pdp-price__original-value span.andes-money-amount__fraction").First().Text(); len(checkPrecoBase) != 0 {
		precoBaseConvertido, err := ConverterPrecoFloat(checkPrecoBase, prod.DOM.Find(".ui-pdp-price__original-value span.andes-money-amount__cents").First().Text())
		if err != nil {
			utils.LogarErroFunc("precoBase", map[string]any{
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

// É mantido como string porque o Mercado Livre só disponibiliza a quantidade de vendas por uma range.
// Olhar função TratarQtdVendas() para explicação das ranges
func (a *Anuncio) qtdVendas(prod *colly.HTMLElement) {
	s, err := TratarQtdVendas(prod.ChildText("span.ui-pdp-subtitle"))
	if err != nil {
		utils.LogarErroFunc("qtdVendas", map[string]any{
			"texto":       s,
			"anuncio":     a.titulo,
			"linkAnuncio": a.link,
		}, err)
		fmt.Println("anuncio sem vendas OU com formato do texto não reconhecido. Cheque o log gerado para mais detalhes")
	}
	a.quantidadeVendas = s
}

// Trata a string {CONDIÇÃO | xx vendidos} que aparece nos anúncios
// e pega somente a condição
func (a *Anuncio) condAnuncio(prod *colly.HTMLElement) {
	c := prod.ChildText("span.ui-pdp-subtitle")

	// log caso não exista esse elemento na página
	if c == "" {
		utils.LogarErroFunc("condAnuncio", map[string]any{
			"link": a.link,
		}, fmt.Errorf("não foi possível extrair a condição do item vendido. Cheque o log gerado para mais detalhes"))
	}
	s := strings.Split(c, "|")
	s[0] = strings.TrimSpace(s[0])
	a.condicao = s[0]
}

// Trata string e remove "Vendido por "
func (a *Anuncio) vendedorNome(prod *colly.HTMLElement) {
	prefixo := "Vendido por "
	vendedor := prod.ChildText(".ui-seller-data-header__title-container")

	if vendedor == "" || len(vendedor) == 0 {
		utils.LogarErroFunc("vendedorNome", map[string]any{
			"vendedor":     vendedor,
			"link_anuncio": a.link,
		}, fmt.Errorf("erro ao extrair nome do vendedor. Verificar anuncio e vendedor"))
	}
	if strings.Contains(vendedor, prefixo) {
		vendedor = strings.Replace(vendedor, prefixo, "", -1)
		vendedor = strings.TrimSpace(vendedor)
		a.vendedor.nome = vendedor
		return
	}
	a.vendedor.nome = vendedor
}

// Pega link do vendedor do produto no ML
func (a *Anuncio) vendedorLink(prod *colly.HTMLElement) {

	link := prod.Request.AbsoluteURL(prod.ChildAttr("div.ui-seller-data-footer__container a", "href"))
	if link == "" || len(link) == 0 {
		utils.LogarErroFunc("vendedorLink", map[string]any{
			"linkVendedor": link,
			"linkAnuncio":  a.link,
		}, fmt.Errorf("elemento html não tem link do vendedor. Examinar pagina"))
		fmt.Println("Erro ao pegar o link do vendedor no mercado livre. Cheque o log para mais detalhes")
	}
	a.vendedor.linkVendedor = link
}

// Alem do numero, pode aparecer "Ultimo disponível" --> Nesse caso irá ser transformado para 1
// Como o texto está envolvido por ( ), é removido por filtrar o 1º e ultimo caracter da string. E isso só acontece se não for o ultimo disponível
// Função busca as duas tags já que ML traz em lugares diferentes a informação caso seja o ultimo em estoque (genial isso kkkk)
func (a *Anuncio) montarEstoque(prod *colly.HTMLElement) {
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
		estoqueNaoUltimo = strings.Replace(estoqueNaoUltimo, "+", "", -1)
		s := strings.Split(estoqueNaoUltimo, " ")
		a.estoque = s[0][1:]
		return
	}
}

// Pega descrição do anuncio. Vem com algumas formatações html simplificada mas por enquanto é relevada
func (a *Anuncio) extrairDescricao(prod *colly.HTMLElement) {
	textoDesc, err := prod.DOM.Find(".ui-pdp-description__content").Html()
	if err != nil {
		utils.LogarErroFunc("extrairDescricao", map[string]any{
			"textoDesc": textoDesc,
			"anuncio":   a.link,
		}, fmt.Errorf("conteudo html do DOM da descrição com problemas. Checar manualmente"))
		fmt.Println("Erro ao extrair descrição do anúncio. Por padrão o campo ficará em branco. Cheque o log para mais detalhes")
		a.descricao = ""
		return
	}

	// troca <br> e <br/> por \n
	reBR := regexp.MustCompile(`(?i)<br\s*/?>`)
	textoDesc = reBR.ReplaceAllString(textoDesc, " - ")

	// remove tag html da desc caso tenha
	reTags := regexp.MustCompile(`<[^>]+>`)
	texto := reTags.ReplaceAllString(textoDesc, "")

	texto = strings.TrimSpace(texto)

	// duplica "" para não quebrar csv em importação
	// texto = strings.ReplaceAll(texto, `"`, `""`)

	a.descricao = `"` + texto + `"`
}

// pega link do anuncio completo da request
func (a *Anuncio) extrairLinkAnuncio(prod *colly.HTMLElement) {
	a.link = prod.Request.URL.String()
}

// pega titulo do anúncio da página do anuncio
func (a *Anuncio) tituloAnuncio(prod *colly.HTMLElement) {
	a.titulo = prod.ChildText(".ui-pdp-title")
}

// percorre cada tabela das seções da ficha tecnica do anuncio
func (a *Anuncio) montarFichaTecnica(prod *colly.HTMLElement) {
	a.FichaTecnica.caracteristicas = make(map[string]any)
	prod.DOM.Find("tr.ui-vpp-striped-specs__row").Each(func(i int, ficha *goquery.Selection) {
		key := ficha.Find("th div.andes-table__header__container").Text()
		val := ficha.Find("td span.andes-table__column--value").Text()

		if key != "" && val != "" {
			a.FichaTecnica.caracteristicas[key] = val
		}
	})
}

// TODO: metodo fallback de gerar na mesma pasta de execução caso não consiga criar pasta
// Cria o csv em uma pasta "extracoes" com os dados extraidos dos anuncios da pesquisa
func ExportarCSV(buscaML string, anuncios []Anuncio) error {
	var pastaDestino = "extracoes"

	caminhoExecutavel, err := os.Executable()
	if err != nil {
		fmt.Println("Erro ao identificar o caminho do executável:", err)
	}
	caminhoReal, err := filepath.EvalSymlinks(caminhoExecutavel)
	if err != nil {
		fmt.Println("Erro ao identificar o caminho do executável:", err)
	}

	var dirBase string
	if strings.HasPrefix(caminhoReal, os.TempDir()) {
		// caso de execução da ferramenta NÃO compilada
		dirBase = "."
	} else {
		dirBase = filepath.Dir(caminhoReal)
	}

	dirExtracao := filepath.Join(dirBase, pastaDestino)
	if err := os.MkdirAll(dirExtracao, os.ModePerm); err != nil {
		return fmt.Errorf("⚠️ não foi possível criar a pasta [%s]. Erro %v", pastaDestino, err)
	}

	// define nome do arquivo. Por padrão: o-que-foi-pesquisado_no_ML_data-hoje_hora-hoje
	nomePesquisa := strings.ReplaceAll(buscaML, " ", "-")
	dataExecucao := time.Now().Format("02-01-2006_15-04")
	nomeArquivo := fmt.Sprintf("%s_%s.csv", nomePesquisa, dataExecucao)

	// gera caminhoabsoluto  do arquivo cross-plataform
	caminhoArquivo := filepath.Join(dirExtracao, nomeArquivo)

	// criar arquivo
	arquivo, err := os.Create(caminhoArquivo)
	if err != nil {
		return fmt.Errorf("⚠️ não foi possível criar o arquivo [%s]. Erro: %v", nomeArquivo, err)
	}
	defer arquivo.Close()

	writer := csv.NewWriter(arquivo)
	writer.Comma = ';' // separador é ; para evitar problema com formatação de numeros e algum texto que possivelmente possa ter ',' nele
	defer writer.Flush()
	cabecalho := []string{"titulo", "condição", "preco_base", "preco_atual", "quantidade_vendas",
		"estoque", "patrocinado", "tem_full", "nota", "quantidade_reviews", "link_anuncio",
		"vendedor", "vendedor_link", "tipo_loja", "descricao"}
	if err := writer.Write(cabecalho); err != nil {
		return fmt.Errorf("erro ao adicionar cabeçalho ao csv: [%v]", err)
	}
	for _, anuncio := range anuncios {
		linha := []string{
			anuncio.titulo,
			anuncio.condicao,
			strconv.FormatFloat(anuncio.precoBase, 'f', 2, 64),
			strconv.FormatFloat(anuncio.precoAtual, 'f', 2, 64),
			anuncio.quantidadeVendas,
			anuncio.estoque,
			anuncio.patrocinado,
			anuncio.full,
			anuncio.nota,
			anuncio.quantidadeReviews,
			anuncio.link,
			anuncio.vendedor.nome,
			anuncio.vendedor.linkVendedor,
			anuncio.vendedor.tipoLoja,
			anuncio.descricao,
		}
		if err := writer.Write(linha); err != nil {
			return fmt.Errorf("erro ao escrever linha para o csv [%v]", err)
		}
	}
	return err
}

// TODO: metodo fallback de gerar na mesma pasta de execução caso não consiga criar pasta
// TODO: extrair MLB do link do anuncio
// cria json com a ficha tecnica de cada anuncio {titulo, link, ficha_tecnica}
// e salva na pagina de extracoes também, mas com sufixo "_fichatecnica"
func ExportarFichaTecnica(buscaML string, anuncios []Anuncio) error {

	var pastaDestino = "extracoes"

	caminhoExecutavel, err := os.Executable()
	if err != nil {
		fmt.Println("Erro ao identificar o caminho do executável:", err)
	}
	caminhoReal, err := filepath.EvalSymlinks(caminhoExecutavel)
	if err != nil {
		fmt.Println("Erro ao identificar o caminho do executável:", err)
	}

	var dirBase string
	if strings.HasPrefix(caminhoReal, os.TempDir()) {
		// caso de execução da ferramenta NÃO compilada
		dirBase = "."
	} else {
		dirBase = filepath.Dir(caminhoReal)
	}

	dirExtracao := filepath.Join(dirBase, pastaDestino)
	if err := os.MkdirAll(dirExtracao, os.ModePerm); err != nil {
		return fmt.Errorf("⚠️ não foi possível criar a pasta [%s]. Erro %v", pastaDestino, err)
	}

	// define nome do arquivo. Por padrão: o-que-foi-pesquisado_no_ML_data-hoje_hora-hoje
	nomePesquisa := strings.ReplaceAll(buscaML, " ", "-")
	dataExecucao := time.Now().Format("02-01-2006_15-04")
	nomeArquivo := fmt.Sprintf("%s_%s_fichatecnica.json", nomePesquisa, dataExecucao)

	// gera caminho absoluto do arquivo cross-plataform
	caminhoArquivo := filepath.Join(dirExtracao, nomeArquivo)

	// criar arquivo
	arquivo, err := os.Create(caminhoArquivo)
	if err != nil {
		return fmt.Errorf("⚠️ não foi possível criar o arquivo com as fichas técnicas [%s]. Erro: %v", nomeArquivo, err)
	}
	defer arquivo.Close()

	type ExportJson struct {
		Titulo       string         `json:"titulo"`
		Link         string         `json:"link"`
		FichaTecnica map[string]any `json:"ficha_tecnica"`
	}
	var exportados []ExportJson

	for _, anuncio := range anuncios {
		exportados = append(exportados, ExportJson{
			Titulo:       anuncio.titulo,
			Link:         anuncio.link,
			FichaTecnica: anuncio.FichaTecnica.caracteristicas,
		})
	}
	encoder := json.NewEncoder(arquivo)
	encoder.SetIndent("", " ")

	return encoder.Encode(exportados)

}

// Encapsula funções do scrape das info do anuncio
func NovoAnuncio(prod *colly.HTMLElement) Anuncio {
	var a Anuncio
	a.tituloAnuncio(prod)
	a.extrairLinkAnuncio(prod)
	a.extrairDescricao(prod)
	a.montarEstoque(prod)
	a.montarPrecoAtual(prod)
	a.montarPrecoBase(prod)
	a.notaAvaliacao(prod)
	a.qtdAvaliacoes(prod)
	a.qtdVendas(prod)
	a.condAnuncio(prod)
	a.temFull(prod)
	a.temPatrocinado()
	a.vendedorLink(prod)
	a.montarFichaTecnica(prod)
	a.tipoLojaVendedor()
	a.vendedorNome(prod)
	return a
}
