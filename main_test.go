package main

import (
	"testing"

	"github.com/caiowirthmann/market2csv/scraper/mercadolivre"
)

func TestTratarQtdResultados(t *testing.T) {
	tests := []string{
		"15.585 resultados", "785 reusltados", "1 resultado",
	}
	// entrada := "15.585 resultados"

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			resultado, err := mercadolivre.TratarQtdResultados(tt)
			if err != nil {
				t.Fatalf("Erro ao converter string [%v]", err)
			}
			t.Logf("[%s [%T]] convertido para [%v [%T]] com sucesso", tt, tt, resultado, resultado)

		})
	}

}
