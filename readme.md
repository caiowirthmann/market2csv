# Market2csv

Market2csv é uma ferramenta Open Souce Gratuíta para extração de dados de anúncios de diversos marketplaces e exportação para um `.csv` com base em um termo de busca

Fácil de usar e simples, não precisa de login na conta do marketplace, pagamento ou qualquer outra informação

Compatível com sistemas Windows quanto Unix

<br>

## 🛍️ O que o `market2csv` faz?

- Aceita um termo de busca diretamente no terminal
- Permite que o usuário defina quantos anuncios serão analisados
- Coleta os seguintes dados de cada anúncios:
    - Titulo
    - Condição (Novo, usado, recondicionado)
    - Preço (Preço Base, Preço com desconto)
    - Quantidade de Vendas
    - Estoque
    - Anuncio é patrocinado?
    - Anuncio está no Full?
    - Nota de avaliação
    - Quantidade de reviews
    - Link do anuncio
    - Nome do vendedor
    - Link da loja do vendedor
    - Tipo de loja do vendedor (Padrão, Eshop, Loja Oficial)
    - Descrição completa
    - Ficha técnica completa (exportado para um arquivo separado)
  
- Exporta todos os dados para um arquivo `.csv` com o nome da pesquisa e data para fácil acesso e indexação
- Ficha técnica é exportada para um arquivo `.json`
- Gera um log de qualquer erro que aconteça durante a extração para facilitar a correção de bugs

<br>


## 🚀 Como usar
### Demo do funcionamento da ferramenta


![Como usar a ferramenta](assets/como-usar.gif)

Ao abrir a ferramenta, será solicitado que você digite o que quer pesquisar nos marketplaces para realizar a extração dos dados:


![pesquisa](assets/pesquisa.png)

Após a ferramenta identificar e mostrar a quantidade de anúncios encontrados, será solicitado quantos anúncios você quer analisar:
- Você pode digitar qualquer valor entre 1 ou a quantidade máxima de anúncios encontrados
    - Os anúncios são analisados por **ordem de relevância**. Exemplo:
    - Caso digite `10`, serão analisados os *10 primeiros resultados*
    - Anuncios patrocinados são considerados nesse ranking de relevância
- E caso queira **TODOS OS ANÚNCIOS ENCONTRADOS**, basta digitar `0 (zero)`:
    - Isso **levará um tempo** caso existam muitos resultados para a sua pesquisa, ou se for uma pesquisa mais genérica, como `camiseta branca` por exemplo, que retora mais de *60 mil anuncios*

![quantidade](assets/quantidade.png)

A ferramenta irá mostrar um texto informativo do progresso da extração dos dados dos anúncio, e caso aconteça algum problema na extração de algum dado, será mostrado também uma mensagem informativa.

E ao final, uma mensagem informando sobre a criação do arquivo `.csv` com os dados coletados e com o arquivo `.json` da ficha técnica completa de cada anúncio:

![final](assets/final.png)

---
<br>


### Executando a ferramenta no Windows

1. Baixe o arquivo `.exe` ou extraia a ferramenta da pasta compactada `.zip` e salve na pasta onde quer que os arquivos das extrações sejam salvas

2. Execute a ferramenta (pode ser que você tenha que dar permissão para executar como administrador, dependendo de como esteja configurado seu computador)

3. As pastas `extracoes` e `logs` onde são gerados os arquivos são criados na mesma pasta onde o executável do `market2csv` está salvo

*Dependendo de como esteja configurada algumas opções no seu computador, pode ser necessário executar a ferramenta como **administrador**. Para dar permissão de administrador, basta clicar com o botão direito no arquivo e clicar na opção **executar como administrador**.<br><br>Ou se quiser tornar essa opção padrão para esse arquivo, clique em **Propriedades >> Atalho >> Avançado >> Executar como administrador** e depois clicar em OK*
<br>

### Executando a ferramenta em sistemas Unix (Linux, MacOs) pelo binário compilado

1. Baixe o arquivo `binário` ou extraia a ferramenta da pasta compactada `.tar.gz` e salve na pasta onde quer que os arquivos das extrações sejam salvas

```bash
# criando a pasta onde irá salvar a ferramenta
mkdir -p market2csv

# extraindo o .tar.gz para a pasta market2csv
tar -xzvf market2csv-linux.v1.0.0.tar.gz -C market2csv
```

2. De permissão e execute o arquivo

```bash
sudo chmod +x market2csv-linux
./market2csv-linux
```

Caso queira, você pode criar um `alias` para executar a ferramenta de qualquer lugar, basta entrar no seu `.bashrc` e adicionar a seguinte linha

*Substitua o nome do alias e o caminho onde a ferramenta está salva para o seu caso*
```bash
alias market2csv='~/home/market2csv/market2csv-linux'
```

---
<br>

### Pré-requisitos para executar diretamente pelo código-fonte


- Go `+1.24.0` instalado ([Caso não tenha siga as instruções aqui](https://golang.org/doc/install))
- Git instalado ([Caso não tenha siga as intruções aqui](https://git-scm.com/))

### Clone o repositório

```bash
git clone https://github.com/caiowirthmann/market2csv.git
cd market2csv
```
e execute a ferramenta

```bash
go run main.go
```

---
<br>

## 🧾 Exemplo do arquivo gerado (.csv)

> Mercado Livre

<!-- |titulo|condição|preco_base|preco_atual|quantidade_vendas|estoque|patrocinado|tem_full|nota|quantidade_reviews|link_anuncio|vendedor|vendedor_link|tipo_loja|descricao|
|------|--------|-------|-------|------|-------|-----|-----|----|-----|------------|--------|-------------|---------|---------|
|Produto 1|novo|65.99|60.99|50|4|não|sim|4.8|75|link_anuncio_marketplace|vendedor x|link_vendedor_marketplace|Loja oficial|descricao do produto completa| -->

<style>
  .styled-table-wrapper {
    overflow-x: auto;
  }

  .styled-table {
    width: 100%;
    border-collapse: collapse;
    margin: 1rem 0;
    font-size: 0.9rem;
    font-family: sans-serif;
    color:rgb(21, 194, 125);
  }

  .styled-table thead tr {
    background-color: #222;
    text-align: left;
  }

  .styled-table th,
  .styled-table td {
    padding: 10px 12px;
    border: 1px solid #444;
    white-space: nowrap;
  }

  .styled-table td {
    overflow-wrap: anywhere;
  }

  @media (max-width: 768px) {
    .styled-table {
      font-size: 0.75rem;
    }

    .styled-table th,
    .styled-table td {
      padding: 6px 8px;
    }
  }
</style>

<div class="styled-table-wrapper">
  <table class="styled-table">
    <thead>
      <tr>
        <th>titulo</th>
        <th>condição</th>
        <th>preco_base</th>
        <th>preco_atual</th>
        <th>quantidade_vendas</th>
        <th>estoque</th>
        <th>patrocinado</th>
        <th>tem_full</th>
        <th>nota</th>
        <th>quantidade_reviews</th>
        <th>link_anuncio</th>
        <th>vendedor</th>
        <th>vendedor_link</th>
        <th>tipo_loja</th>
        <th>descricao</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>Xiaomi Redmi 13 Dual Sim 256-gb 8-gb Ram Global</td>
        <td>Novo</td>
        <td>1399.00</td>
        <td>1287.08</td>
        <td>50</td>
        <td>10</td>
        <td>não</td>
        <td>Não</td>
        <td>5.0</td>
        <td>3</td>
        <td>https://produto.mercadolivre.com.br/MLB-5349435042-xiaomi-redmi-13-dual-sim-256-gb-8-gb-ram-global-_JM?searchVariation=187511295831#is_advertising=true&searchVariation=187511295831&position=1&search_layout=stack&type=pad&tracking_id=03312ff9-b89b-4f6e-945b-a855779b7a86&is_advertising=true&ad_domain=VQCATCORE_LST&ad_position=1&ad_click_id=OTYzNWFlYmEtOGZkMi00ZGRmLThjMjYtN2I2YjlmMjViOGFh
        </td>
        <td>PECSHOPP</td>
        <td>https://lista.mercadolivre.com.br/_CustId_1797698755?item_id=MLB5349435042&category_id=MLB1055&seller_id=1797698755&client=recoview-selleritems&recos_listing=true#origin=vip&component=sellerData&typeSeller=classic
        </td>
        <td>Padrão</td>
        <td>APARELHO VERSÃO GLOBAL LACRADO + NOTA FISCAL + BRINDE</td>
      </tr>
      <tr>
        <td>Celular Xiaomi Redmi 14c 256gb 8gb Ram Dual Sim</td>
        <td>Novo</td>
        <td>1359.00</td>
        <td>1168.74</td>
        <td>50</td>
        <td>50</td>
        <td>não</td>
        <td>Não</td>
        <td>5.0</td>
        <td>3</td>
        <td>https://produto.mercadolivre.com.br/MLB-5341837680-celular-xiaomi-redmi-14c-256gb-8gb-ram-dual-sim-_JM?searchVariation=187431799701#is_advertising=true&searchVariation=187431799701&position=2&search_layout=stack&type=pad&tracking_id=03312ff9-b89b-4f6e-945b-a855779b7a86&is_advertising=true&ad_domain=VQCATCORE_LST&ad_position=2&ad_click_id=MTg5MWNiMzgtMTVkNy00OTM2LWFkMWYtYWY1MTM0NGE4NDc2
        </td>
        <td>XIAMIELTRO</td>
        <td>https://lista.mercadolivre.com.br/_CustId_1827271932?item_id=MLB5341837680&category_id=MLB1055&seller_id=1827271932&client=recoview-selleritems&recos_listing=true#origin=vip&component=sellerData&typeSeller=classic
        </td>
        <td>Padrão</td>
        <td>FRETE GRÁTIS - Produto Original com Garantia - Loja 100% confiável</td>
      </tr>
      <tr>
        <td>Redmi Note 13 Pro 5G 256 GB Preto 8 GB RAM</td>
        <td>Novo</td>
        <td>1727.99</td>
        <td>1727.99</td>
        <td>10000</td>
        <td>25</td>
        <td>não</td>
        <td>Não</td>
        <td>4.7</td>
        <td>5090</td>
        <td>https://www.mercadolivre.com.br/redmi-note-13-pro-5g-256-gb-preto-8-gb-ram/p/MLB29739167#polycard_client=search-nordic&searchVariation=MLB29739167&wid=MLB5334363160&position=5&search_layout=stack&type=product&tracking_id=03312ff9-b89b-4f6e-945b-a855779b7a86&sid=search</td>
        <td>VENDASBESTSHOP</td>
        <td>https://lista.mercadolivre.com.br/_CustId_771360381?item_id=MLB5334363160&category_id=MLB1055&seller_id=771360381&client=recoview-selleritems&recos_listing=true#origin=pdp&component=sellerData&typeSeller=classic</td>
        <td>Padrão</td>
        <td>Fotografia profissional no seu bolso  -  Descubra infinitas possibilidades para suas fotos com as 3 câmeras principais da sua equipe. Teste sua criatividade e brinque com iluminação, diferentes planos e efeitos para obter ótimos resultados.</td>
      </tr>
      <tr>
        <td>Xiaomi Poco C75 Dual Sim 256 Gb Verde 8 Gb RAM</td>
        <td>Novo</td>
        <td>858.00</td>
        <td>858.00</td>
        <td>5000</td>
        <td>25</td>
        <td>não</td>
        <td>Não</td>
        <td>4.7</td>
        <td>2122</td>
        <td>https://www.mercadolivre.com.br/xiaomi-poco-c75-dual-sim-256-gb-verde-8-gb-ram/p/MLB42312701#polycard_client=search-nordic&searchVariation=MLB42312701&wid=MLB3980403061&position=8&search_layout=stack&type=product&tracking_id=03312ff9-b89b-4f6e-945b-a855779b7a86&sid=search</td>
        <td>CHAVESANGELO70</td>
        <td>https://lista.mercadolivre.com.br/_CustId_141417082?item_id=MLB3980403061&category_id=MLB1055&seller_id=141417082&client=recoview-selleritems&recos_listing=true#origin=pdp&component=sellerData&typeSeller=classic</td>
        <td>Padrão</td>
        <td>Fotografia profissional no seu bolso  -  Descubra infinitas possibilidades para suas fotos com as 3 câmeras principais da sua equipe. Teste sua criatividade e brinque com iluminação, diferentes planos e efeitos para obter ótimos resultados.</td>
      </tr>
    </tbody>
  </table>
</div>



*arquivo `.json` resumido apenas para exemplificar e não ficar extenso. A ferramenta extrai **TODOS** os campos que foram preenchidos na ficha técnica do anúncio*


```json
[
 {
  "titulo": "Xiaomi Redmi 13 Dual Sim 256-gb 8-gb Ram Global",
  "link": "https://produto.mercadolivre.com.br/MLB-5349435042-xiaomi-redmi-13-dual-sim-256-gb-8-gb-ram-global-_JM?searchVariation=187511295831#is_advertising=true\u0026searchVariation=187511295831\u0026position=1\u0026search_layout=stack\u0026type=pad\u0026tracking_id=03312ff9-b89b-4f6e-945b-a855779b7a86\u0026is_advertising=true\u0026ad_domain=VQCATCORE_LST\u0026ad_position=1\u0026ad_click_id=OTYzNWFlYmEtOGZkMi00ZGRmLThjMjYtN2I2YjlmMjViOGFh",
  "ficha_tecnica": {
   "Abertura do diafragma da câmera frontal": "f 2.2",
   "Abertura do diafragma da câmera traseira": "f 1.69/f 2.4/f 2.4/f 2.4",
   "Acessórios incluídos": "1 cabo USB",
   "Altura fechado": "0 mm",
   "Altura x Largura x Profundidade": "16.26 cm x 7.48 cm x 8.96 mm",
   "Ano de lançamento": "2021",
   "Autonomia da batería em standby": "0 h",
   "Autonomia de conversação": "0 h",
   "Brilho máximo da tela": "1.120 cd/m²",
   "Brilho máximo da tela secundária": "0 cd/m²",
   "Camada original de personalização do sistema operacional": "MIUI 12",
   "Capacidade da bateria": "4,78 Ah",
   "Capacidade máxima do cartão de memória": "0 KB",
   "Características principais das câmeras": "Foto HDR, Modo beleza IA, Modo profissional, Panorama mode, Super Night View 2.0",
   "Classificação IP": "IP53",
   "Com Bluetooth": "Sim",
   "Com GPS": "Sim",
   "Com IMEI": "Sim",
   "Com NFC": "Sim",
   "Com Wi-Fi": "Sim",
   "Com acelerômetro": "Sim",
   "Com barômetro": "Sim",
   "Com bateria removível": "Não",
   "Com bússola": "Sim",
   "Com carregamento rápido": "Sim",
   "Com carregamento sem fio": "Sim",
   "Com conector USB": "Sim",
   "Com conector jack 3.5 mm": "Não",
   "Com câmera": "Sim",
   "Com flash na câmara frontal": "Não",
   "Com giroscópio": "Sim",
   "Com leitor de impressão digital": "Sim",
   "Com radio": "Não",
   "Com ranhura para cartão de memória": "Não",
   "Com reconhecimento de mão": "Sim",
   "Com reconhecimento de íris": "Não",
   "Com reconhecimento facial": "Sim",
   "Com sensor de proximidade": "Sim",
   "Com sintonizador de TV": "Não",
   "Com teclado QWERTY físico": "Não",
   "Com tela dobrável": "Não",
   "Com tela tátil": "Sim",
   "Duração do produto": "0 meses",
   "Fabricante": "Xiaomi Communications Co., Ltd.",
   "Inclui lápis": "Não",
   "Largura fechada": "0 mm",
   "Linha": "Mi",
   "Marca": "Xiaomi",
   "Memória RAM": "8 GB",
   "Memória interna": "256 GB",
   "Modelo": "Xiaomi 13 lite",
   "Modelo alfanumérico": "24049RN28L",
   "Modelo de GPU": "Adreno 650",
   "Modelo detalhado": "8 GB" 
  }
 }
]
```
<br>
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

*Apesar do mercado livre não disponibilizar uma forma de acessar a quantidade de estoque disponível no anúncio, existem meios de descobrir o número exato. É necessário que o vendedor não tenha definido um limite de unidades por compra no anúncio (na opção de quantidade aparece um menu com a opção `mais de x unidades`) e que você tenha tempo de sobra. Basta clicar na opção e ir trocando para um valor que não aparece a mensagem de "sem estoque" quando você digitar o valor. Tentativa e erro, bateção de lata*

**A FERRAMENTA NÃO FAZ ISSO**

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

🔴 **Alta prioridade** || 🟡 **Média prioridade** || ⚪ **Baixa prioridade**


- [ ] 🔴 Incluir funcionalidade para Shopee
- [ ] 🔴 Incluir funcionalidade para Amazon
- [ ] 🔴 Incluir funcionalidade para Magazine Luiza
- [ ] 🔴 Seleção de quais marketplaces serão analisados
- [ ] Configuração para permitir personalização do arquivo de exportação:
    - [ ] ⚪ Incluir/Não incluir campo no `.csv`
    - [ ] ⚪ Ordem das colunas
    - [ ] 🟡 Configuração de quais campos são incluidos no arquivo da ficha técnica para identificação do anúncio (por padrão, são incluidos título e link do anúncio)
    - [ ] 🟡 Pasta de exportação (nome e local)
- [ ] Novas opções de exportação:
    - [ ] ⚪ Exportar dados do anúncio para planilhas:
        - [ ] Excel: `.xls .xlsx`
        - [ ] LibreOffice/OpenOffice: `.ods`
---