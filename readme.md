# market2csv

🎯 Uma ferramenta de linha de comando (CLI) para extrair dados de anúncios de marketplaces como Mercado Livre, Shopee, Shein (em breve) e exportá-los para um arquivo `.csv`.

---

## 🛍️ O que o `market2csv` faz?

- Aceita um termo de busca diretamente no terminal
- Faz scraping dos principais marketplaces (inicialmente Mercado Livre)
- Coleta dados detalhados dos anúncios:
  - Título
  - Preço atual e original
  - Frete
  - Descrição
  - Estoque
  - Quantidade de reviews
  - Nota
  - Vendedor + link do vendedor
  - Link do anúncio
- Exporta todos os dados para um arquivo `market2csv.csv`

---

## 🚀 Como usar

### Pré-requisitos
- Go instalado ([instale aqui](https://golang.org/doc/install))
- Conexão com a internet

### Clone o repositório

```bash
git clone https://github.com/seunome/market2csv.git
cd market2csv
```

### Execute o programa

```bash
go run main.go
```

Você será guiado no terminal para digitar:
1. O termo da busca (ex: `"smartphone xiaomi"`)
2. A quantidade de anúncios a serem coletados

---

## 📁 Exemplo de saída (.csv)

| Título                        | Preço   | Frete | Estoque | Nota | Link do Anúncio          |
|------------------------------|---------|-------|---------|------|---------------------------|
| Smartphone Xiaomi Redmi 12   | R$ 1.199| Grátis| 25      | 4.8  | https://...               |

---

## 🛠️ Roadmap

- [x] Mercado Livre scraper
- [ ] Shopee scraper
- [ ] Shein scraper
- [ ] Suporte a múltiplos idiomas
- [ ] Exportar para JSON

---

## ⚠️ Aviso legal

Esta ferramenta é feita para fins educacionais. O uso de scraping deve sempre respeitar os [termos de uso](https://www.mercadolivre.com.br/ajuda/Termos-e-condicoes-gerais-de-uso_1403) dos sites acessados.

---

## 📄 Licença

MIT © [Seu Nome](https://github.com/seunome)