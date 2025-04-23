package assets

import (
	"fmt"
)

const (
	asciiBanner = `
Ferramenta extratora de dados de anúncios de marketplaces


███╗   ███╗ █████╗ ██████╗ ██╗  ██╗███████╗████████╗   ██████╗        ██████╗███████╗██╗   ██╗
████╗ ████║██╔══██╗██╔══██╗██║ ██╔╝██╔════╝╚══██╔══╝   ╚════██╗      ██╔════╝██╔════╝██║   ██║
██╔████╔██║███████║██████╔╝█████╔╝ █████╗     ██║█████╗ █████╔╝█████╗██║     ███████╗██║   ██║
██║╚██╔╝██║██╔══██║██╔══██╗██╔═██╗ ██╔══╝     ██║╚════╝██╔═══╝ ╚════╝██║     ╚════██║╚██╗ ██╔╝
██║ ╚═╝ ██║██║  ██║██║  ██║██║  ██╗███████╗   ██║      ███████╗      ╚██████╗███████║ ╚████╔╝ 
╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝   ╚═╝      ╚══════╝       ╚═════╝╚══════╝  ╚═══╝  
                                                                                              


Dúvidas, sugestões, bugs, mais informações:
| -----> contato: github.com/caiowirthmann/market2csv


Está é uma ferramenta gratuita e open souce.
O código está disponível no link do github acima e contribuições são bem-vindas!

Se você gostar da ferramenta, uma ⭐ no repositório ajudaria bastante para divulgar o projeto!
===============================================================================================

`
)

func PrintTelaTerminal() {
	fmt.Print(asciiBanner)
}
