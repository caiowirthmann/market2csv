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

Entre no site para mais informações e atualizações sobre a ferramenta:
| -----> https://caiowirthmann.github.io/market2csv-site/

	- Como usar
	- Instalação e Download
	- Exemplos
	- Limitações
	- FAQ e Suporte

Github com o código-fonte aberto e gratuitp da ferramenta:
| -----> github.com/caiowirthmann/market2csv

contato:
| ---> telegram: @caiowp
| ---> discord: musaaaa4

Se curtir a ferramenta, ajude a divulga-lá. Uma ⭐ no projeto no github ajuda bastante também
==============================================================================================
`
)

func PrintTelaTerminal() {
	fmt.Print(asciiBanner)
}
