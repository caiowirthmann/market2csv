package assets

import (
	"fmt"
)

const (
	asciiBanner = `
███╗   ███╗ █████╗ ██████╗ ██╗  ██╗███████╗████████╗   ██████╗        ██████╗███████╗██╗   ██╗
████╗ ████║██╔══██╗██╔══██╗██║ ██╔╝██╔════╝╚══██╔══╝   ╚════██╗      ██╔════╝██╔════╝██║   ██║
██╔████╔██║███████║██████╔╝█████╔╝ █████╗     ██║█████╗ █████╔╝█████╗██║     ███████╗██║   ██║
██║╚██╔╝██║██╔══██║██╔══██╗██╔═██╗ ██╔══╝     ██║╚════╝██╔═══╝ ╚════╝██║     ╚════██║╚██╗ ██╔╝
██║ ╚═╝ ██║██║  ██║██║  ██║██║  ██╗███████╗   ██║      ███████╗      ╚██████╗███████║ ╚████╔╝ 
╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝   ╚═╝      ╚══════╝       ╚═════╝╚══════╝  ╚═══╝  

Ferramenta extratora de dados de anúncios de marketplaces                                                                           

Entre no site para mais informações sobre a ferramenta:
| -----> https://caiowirthmann.github.io/market2csv-site/

	- Instalação
	- Como usar
	- Exemplos
	- Limitações
	- Link para download


Github com o código aberto da ferramenta:
| -----> github.com/caiowirthmann/market2csv

Está é uma ferramenta gratuita e open souce.
O código está disponível no link do github acima e contribuições são bem-vindas!

Se você gostou da ferramenta, de uma ⭐ no repositório ajudaria bastante para divulgar o projeto!
===================================================================================================

`
)

func PrintTelaTerminal() {
	fmt.Print(asciiBanner)
}
