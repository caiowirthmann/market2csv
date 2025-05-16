package mercadolivre_test

import (
	"testing"

	"github.com/caiowirthmann/market2csv/scraper/mercadolivre"
)

func TestConverterPrecoFloat(t *testing.T) {
	teste := []struct {
		caso     string
		s1       string
		s2       string
		esperado float64
	}{
		{caso: "real e centavos",
			s1:       "135",
			s2:       "45",
			esperado: 135.45,
		},
		{caso: "só real sem centavos",
			s1:       "39",
			s2:       "",
			esperado: 39.00},
		{caso: "valor acima de 999, pontuação milhar",
			s1:       "2.549",
			s2:       "99",
			esperado: 2549.99},
		{caso: "só centavos. NÃO DEVE ACONTECER NO MeLi, R$ minimo é 7,00; mas...",
			s1:       "",
			s2:       "99",
			esperado: 0.99},
		{caso: "tudo vazio",
			s1:       "",
			s2:       "",
			esperado: 0.0},
	}

	for _, tt := range teste {
		t.Run(tt.caso, func(t *testing.T) {
			resultado, err := mercadolivre.ConverterPrecoFloat(tt.s1, tt.s2)
			if err != nil && tt.esperado == 0.0 {
				t.Logf("Erro esperado na conversão e logado. OK --> Valor setado para 0.0")
			}
			if tt.esperado == resultado {
				t.Log("OK")
			} else {
				t.Fatalf("valor não convetido corretamente [%s e %s => %v] não bate com"+
					"[%v]", tt.s1, tt.s2, tt.esperado, resultado)
			}
		})

	}

}

func TestTratarQtdResultados(t *testing.T) {
	teste := []struct {
		caso         string
		valor        string
		esperado     int64
		erroEsperado bool
	}{{
		caso:         "string normal",
		valor:        "153 resultados",
		esperado:     153,
		erroEsperado: false,
	}, {
		caso:         "string vazia",
		valor:        "",
		esperado:     0,
		erroEsperado: false,
	}, {
		caso:         "numero muito grande --> chance quase 0 de rolar",
		valor:        "19399348",
		esperado:     19399348,
		erroEsperado: false,
	}, {
		caso:         "string sem espaço",
		valor:        "208resultados",
		esperado:     0,
		erroEsperado: true,
	}}

	for _, tt := range teste {
		t.Run(tt.caso, func(t *testing.T) {
			resultado, err := mercadolivre.TratarQtdResultados(tt.valor)
			if err != nil && tt.erroEsperado {
				t.Logf("elemento [%s] com erro no MeLi. NÃO DEVE ACONTECER", tt.valor)
			}
			if err != nil && !tt.erroEsperado {
				t.Logf("erro ao pegar string. Valor setado para 0")
			}
			if tt.esperado == resultado {
				t.Log("OK")
			} else {
				t.Fatalf("[%s] não bate com [%v]", tt.valor, resultado)
			}
		})
	}
}

func TestTratarQtdVendas(t *testing.T) {
	teste := []struct {
		caso         string
		texto        string
		esperado     string
		erroEsperado bool
	}{{
		caso:         "prod novo com vendas",
		texto:        "NOVO | +25 vendidos",
		esperado:     "25",
		erroEsperado: false,
	}, {
		caso:         "prod novo SEM vendas",
		texto:        "Novo",
		esperado:     "0",
		erroEsperado: false,
	}, {
		caso:         "string vazia (casos de catalogo. Não tem informação de qtd vendas e do vendedor",
		texto:        "",
		esperado:     "0",
		erroEsperado: true,
	}, {
		caso:         "só qtd vendas (NÃO ACONTECE), mas MeLi pode mudar isso",
		texto:        "+50 vendidos",
		esperado:     "0",
		erroEsperado: true,
	}}

	for _, tt := range teste {
		t.Run(tt.caso, func(t *testing.T) {
			resultado, err := mercadolivre.TratarQtdVendas(tt.texto)
			if tt.erroEsperado && err == nil {
				t.Logf("caso anormal do MeLi [%s]. Setado para 0", tt.texto)
			}
			if tt.esperado == resultado {
				t.Logf("OK")
			} else {
				t.Fatalf("INPUT [%s] === ESPERADO [%s] | RECEBIDO [%s]", tt.texto, tt.esperado, resultado)
			}
			if err != nil && tt.erroEsperado {
				t.Logf("OK com erro [%s] esperado e tratado", err)
			}
		})
	}
}
