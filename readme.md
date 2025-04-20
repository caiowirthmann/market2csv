# market2csv

🎯 Market2csv é uma ferramenta Open Souce gratuíta para extração de dados de anúncios de diversos marketplaces e exportação para um `.csv` com base em um termo de busca

Fácil de usar e simples, não precisa de login na conta do marketplace, pagamento ou qualquer outra informação

Compatível com sistemas Windows quanto Unix

---
<br>

## 🛍️ O que o `market2csv` faz?

- Aceita um termo de busca diretamente no terminal
- Permite que o usuário defina quantos anuncios serão analisados
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
  
- Exporta todos os dados para um arquivo `.csv` com o nome da pesquisa e data para fácil acesso e indexação
- Gera um log de qualquer erro que aconteça durante a extração para facilitar a correção de bugs
---
<br>

## 🚀 Como usar
......

### Pré-requisitos
- Go `+1.24.0` instalado ([Caso não tenha siga as instruções aqui](https://golang.org/doc/install))


### Clone o repositório

```bash
git clone https://github.com/caiowirthmann/market2csv.git
cd market2csv
```

### Execute o programa

```bash
go run main.go
```

Você será guiado no terminal para digitar:
1. O termo da busca (ex: `"smartphone xiaomi"`)
2. A quantidade de anúncios a serem extraídos

>Os arquivos CSV gerados serão salvos automaticamente em uma pasta chamada extracoes, criada no mesmo local do executável — isso funciona tanto em Windows quanto em sistemas Unix.

---

## 🧾 Exemplo do arquivo gerado (.csv)

|titulo|preço base|preço atual|quantidade vendas|estoque|patrocinado|tem Full?|nota|quantidade reviews|link anuncio|descricao|nome vendedor|link vendedor|
|------------|-----|------|-----------|----------|-------|-----|-----|-------|------|-----|------|------|
|Produto1| 65,99|53,50|30|5|nao|sim|5.0|55|linkdoanuncio.com|descricao do produto|vendedor x|linkvendedor.com

---


## ℹ️ Notas sobre limitações e comportamento da ferramenta no **Mercado Livre**

Algumas informações exibidas nos anúncios do Mercado Livre são disponibilizadas de forma limitada ou não-exatas (principalmente por questões de privacidade implementadas pelo Mercado Livre) Abaixo seguem as explicações de cada campo em que isso ocorre:

### 📊 Quantidade de vendas

O número de vendas exibido nos anúncios do Mercado Livre segue um padrão de **faixas** após 5 unidades vendidas. Então:

- Anuncios com até 5 unidades, o número exato é mostrado.
- Após isso (5+ vendas), o site exibe apenas pelo **PRIMEIRO VALOR** de um **INTERVALO APROXIMADO**

```code
5 a 9 - Exibido como +5 Vendidos
10 a 24 - Exibido como +10 Vendidos
25 a 49 - Exibido como +25 Vendidos
50 a 99 - Exibido como +50 Vendidos
100 a 499 - Exibido como +100 Vendidos
500 a 999 - Exibido como +500 Vendidos

A partir de 1000 unidades vendidas começa a aparece "mil" de forma literal ao invés de "000"

1000 a 4.999 - Exibido como +1000 Vendidos
5mil a 9.999 - Exibido como +5mil Vendidos
10mil a 49.999 - Exibido como +10mil Vendidos
50mil a 99.999 - Exibido como +50mil Vendidos
+100mil - Exibido como +100mil Vendidos
```
### 📦 Estoque

Mesmo caso da quantidade de vendas, o numero **EXATO** é disponibilizado somente até **5 unidades**, após isso é disponibilizado em **faixas**

> caso o anúncio não tenha estoque (o unidades), o anúncio fica pausado e só é possível acessá-lo se tiver o link dele. Ele não aparece na busca

> 1 unidade - Último disponível!

```code
5 a 9 - +5 disponíveis
10 a 24 - +10 disponíveis
25 a 49 - +25 disponíveis
50 a ... - +50 disponíveis
```

*Apesar do mercado livre não disponibilizar uma forma de acessar a quantidade de estoque disponível no anúncio, existem meios de descobrir o número exato. É necessário que o vendedor não tenha definido um limite de unidades por compra no anúncio (na opção de quantidade aparece um menu com a opção `mais de x unidades`) e que você tenha tempo de sobra. Basta clicar na opção e ir mudando para um valor que não aparece a mensagem de "sem estoque" quando você digitar o valor. Tentativa e erro, bateção de lata*

### 🚚 Frete / Full

Frete *por enquanto* não é uma opção disponibilizada na ferramenta:

- Calculo do frete é multifatorial:
    - Categoria do anúncio
    - Reputação do vendedor
    - Valor de venda do anúncio
    - Dimensões do produto
    - Forma de entrega
    - Região despacho x Região entrega
    > Esse valor é calculado e mostrado no anúncio quando se está **logado** na sua conta do mercado livre, e selecionadA a forma de entrega <mark>Mercado envios, Full, Mercado Envios Flex, Frete a combinar com o vendedor</mark> e o endereço de entrega (é possível ter mais de um endereço de entrega, mas na conta você seleciona um "padrão" que será usado para esses calculos e mostrado primeiro)

A coluna de `Full` da ferramenta mostra apenas se o anúncio tem Mercado Envios Full ou não. O valor em sí segue a mesma questão citada acima

## 🛠️ Roadmap de funcionalidades e melhorias

- [x] Mercado Livre
    - [ ] Ficha técnica
- [ ] Shopee
- [ ] Amazon
- [ ] Shein
- [ ] Exportação do arquivo para `JSON`
- [ ] Personalização do arquivo de exportação:
    - [ ] Incluir/Não incluir campo
    - [ ] Ordem
    - [ ] 


---