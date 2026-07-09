# 🏛️ RioZelo — Painel 1746 Transparente

> Plataforma full-stack de zeladoria urbana para o cidadão carioca registrar ocorrências e acompanhar indicadores em tempo real.

![Go](https://img.shields.io/badge/Back--end-Go-00ADD8?style=flat-square&logo=go&logoColor=white)
![Next.js](https://img.shields.io/badge/Front--end-Next.js%20App%20Router-000000?style=flat-square&logo=nextdotjs&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Estilo-Tailwind%20CSS-38BDF8?style=flat-square&logo=tailwindcss&logoColor=white)
![SSE](https://img.shields.io/badge/Tempo%20Real-SSE-FF6B35?style=flat-square)
![MVP](https://img.shields.io/badge/Status-MVP%201-22C55E?style=flat-square)

---

## 📋 Índice

- [Sobre o Projeto](#-sobre-o-projeto)
- [Fluxo do Usuário](#-fluxo-do-usuário)
- [Dashboard e Tempo Real](#-dashboard-e-tempo-real)
- [Arquitetura Técnica](#️-arquitetura-técnica)
- [Requisitos do Sistema](#-requisitos-do-sistema-mvp-1)
- [Histórias de Usuário](#-histórias-de-usuário)
- [Estrutura do Projeto](#-estrutura-do-projeto)
- [Como Rodar Localmente](#-como-rodar-localmente)

---

## 📌 Sobre o Projeto

Cidades inteligentes precisam de canais de comunicação eficientes e dados transparentes. O **RioZelo** é uma plataforma full-stack que permite ao cidadão carioca registrar ocorrências urbanas de forma guiada e acompanhar, em tempo real, os indicadores de zeladoria da cidade — como bairros e categorias com maior índice de problemas.

**Duas telas, uma missão:**

| Tela | Persona | Função |
|---|---|---|
| 🧑‍💻 Tela do Cidadão | Morador do Rio | Formulário dinâmico para registrar ocorrências |
| 📊 Tela do Operador | COR / Central 1746 | Dashboard reativo com ranking, gráfico e feed ao vivo |

---

## 🗺️ Fluxo do Usuário

O sistema foi desenhado para **mitigar o erro humano** no preenchimento através de um fluxo dinâmico em árvore:

```
1. Macro Categoria
   └── 2. Subcategoria (filtrada dinamicamente)
           ├── Opção específica → Prossegue
           └── "Outro" → Campo de texto livre (renderizado condicionalmente)
                   └── 3. Bairro (lista pré-definida) + Rua (input de texto)
                           └── ✅ Botão Enviar (liberado apenas com todos os campos válidos)
```

**Exemplo de fluxo:**
> Usuário seleciona `Iluminação Pública` → aparece `Lâmpada Apagada | Poste Danificado | Outro` → seleciona `Poste Danificado` → preenche `Bairro: Tijuca` + `Rua: Rua Conde de Bonfim` → envia.

---

## 📊 Dashboard e Tempo Real

Após o envio, a plataforma processa os dados instantaneamente via **SSE (Server-Sent Events)** e atualiza três visões na tela do operador:

- **Ranking de Ocorrências** — lista ordenada das macro categorias mais afetadas (ex: Asfalto lidera com 45% das reclamações)
- **Gráfico de Barras Vertical** — comparativo visual do volume por categoria ou bairro
- **Feed de Denúncias Recentes** — scroll das últimas 5–6 ocorrências, atualizado em < 1 segundo, sem refresh de página

---

## 🏗️ Arquitetura Técnica

### Back-end — Go (Golang)

```
POST /api/ocorrencias  →  Handler Go
                               └── Goroutine computa novo ranking (async)
                                       └── Canal SSE empurra update para o front-end
GET  /api/stream       →  Conexão SSE aberta por cliente
```

| Decisão | Justificativa |
|---|---|
| **Go** | Baixíssimo consumo de memória com milhares de conexões SSE abertas simultaneamente |
| **Goroutines** | Ranking recomputado em background sem bloquear a resposta HTTP |
| **SSE** | Push server → client sem overhead de WebSocket para este caso de uso |
| **Estado em memória** | Sem dependência de banco de dados no MVP 1 — ranking vive em struct Go |

### Front-end — Next.js (App Router) + Tailwind CSS

| Decisão | Justificativa |
|---|---|
| **App Router** | Server Components para dados estáticos, Client Components para interatividade |
| **`"use client"`** | Gerenciamento de estado do formulário dinâmico no cliente (sem recarregamento) |
| **Tailwind CSS** | Responsividade Mobile-First — cidadão acessa pelo celular na rua |
| **SPA** | Toda a experiência (formulário + dashboard) em uma única tela reativa |

### Diagrama de Fluxo de Dados

```
Cidadão
  │
  ├─[Formulário Next.js]──► POST /api/ocorrencias ──► Handler Go
  │                                                        │
  │                                               Goroutine computa ranking
  │                                                        │
  └─[Dashboard Next.js] ◄── SSE /api/stream ◄────── Canal Go
         │
         ├── Ranking atualizado
         ├── Gráfico atualizado
         └── Feed atualizado
```

---

## 📋 Requisitos do Sistema (MVP 1)

### ✅ Requisitos Funcionais

| ID | Nome | Descrição |
|---|---|---|
| RF01 | Fluxo Dinâmico | Seleção de Macro Categoria filtra automaticamente as opções de Subcategoria |
| RF02 | Campo Condicional | Ao selecionar "Outro", renderizar obrigatoriamente um campo de texto livre |
| RF03 | Localização Declarativa | Coletar Bairro (lista pré-definida) e Rua (input de texto) da ocorrência |
| RF04 | Streaming (SSE) | Exibir novas denúncias no feed em tempo real, sem atualizar a página |
| RF05 | Agregação de Dados | Calcular e exibir ranking das categorias mais afetadas e gráfico de barras comparativo |

### ⚙️ Requisitos Não-Funcionais

| ID | Nome | Descrição |
|---|---|---|
| RNF01 | Alta Concorrência | Servidor Go com baixo consumo de memória para múltiplas conexões SSE simultâneas |
| RNF02 | Interface Responsiva | Next.js + Tailwind CSS com Mobile-First e carregamento rápido |
| RNF03 | Monopágina (SPA) | Toda a experiência de inserção e visualização em uma única tela reativa |

---

## 📖 Histórias de Usuário

### 👤 Persona 1: O Cidadão Carioca

**US01 — Fluxo Guiado de Denúncia**

> **Como** um cidadão morador do Rio de Janeiro,
> **Eu quero** registrar um problema de zeladoria através de um formulário inteligente e dinâmico,
> **Para que** eu não precise adivinhar termos técnicos e consiga enviar o endereço exato pelo celular rapidamente.

**Critérios de Aceitação:**
- [ ] O segundo select só aparece após o primeiro ser preenchido
- [ ] Ao selecionar "Outro", um campo de texto surge imediatamente
- [ ] O botão de envio só é liberado com todos os campos (Bairro e Rua) validados

---

### 🖥️ Persona 2: O Operador da Central (COR / 1746)

**US02 — Monitoramento Visual e Reativo**

> **Como** um gestor ou operador da central de operações,
> **Eu quero** visualizar um feed de denúncias que atualiza sozinho e ver gráficos dinâmicos na mesma tela,
> **Para que** eu identifique picos de problemas e bairros críticos instantaneamente, sem atualizar o sistema manualmente.

**Critérios de Aceitação:**
- [ ] Novas denúncias aparecem no topo do feed em menos de 1 segundo
- [ ] O feed retém no máximo as últimas 5–6 ocorrências para manter a interface limpa
- [ ] O ranking e o gráfico recalculam automaticamente a cada nova denúncia recebida

---

## 📁 Estrutura do Projeto

```
riozelo/
├── backend/                  # Servidor Go
│   ├── main.go               # Entry point, rotas e SSE handler
│   ├── handlers/
│   │   ├── ocorrencia.go     # POST /api/ocorrencias
│   │   └── stream.go         # GET /api/stream (SSE)
│   ├── ranking/
│   │   └── ranking.go        # Lógica de agregação e ranking em memória
│   └── models/
│       └── ocorrencia.go     # Structs de dados
│
├── frontend/                 # Next.js App Router
│   ├── app/
│   │   ├── page.tsx          # Página principal (layout SPA)
│   │   └── layout.tsx
│   ├── components/
│   │   ├── FormularioCidadao.tsx   # RF01, RF02, RF03 — formulário dinâmico
│   │   ├── FeedDenuncias.tsx       # RF04 — feed SSE em tempo real
│   │   ├── RankingCategorias.tsx   # RF05 — ranking das ocorrências
│   │   └── GraficoBarras.tsx       # RF05 — gráfico de barras
│   └── lib/
│       └── categorias.ts     # Mapeamento Macro Categoria → Subcategorias
│
└── README.md
```

---

## 🚀 Como Rodar Localmente

### Pré-requisitos

- [Go 1.22+](https://golang.org/dl/)
- [Node.js 18+](https://nodejs.org/)

### 1. Clone o repositório

```bash
git clone https://github.com/seu-usuario/riozelo.git
cd riozelo
```

### 2. Suba o back-end (Go)

```bash
cd backend
go run main.go
# Servidor rodando em http://localhost:8080
```

### 3. Suba o front-end (Next.js)

```bash
cd frontend
npm install
npm run dev
# Aplicação disponível em http://localhost:3000
```

### 4. Acesse

| URL | O quê |
|---|---|
| `http://localhost:3000` | Aplicação completa (formulário + dashboard) |
| `http://localhost:8080/api/ocorrencias` | Endpoint POST para criar ocorrências |
| `http://localhost:8080/api/stream` | Endpoint SSE para o feed em tempo real |

---

## 🗺️ Roadmap

- [x] MVP 1 — Formulário dinâmico + SSE + Ranking + Gráfico (tudo em uma tela)
- [ ] MVP 2 — Persistência com banco de dados (PostgreSQL)
- [ ] MVP 3 — Autenticação de operadores e painel administrativo
- [ ] MVP 4 — Mapa do Rio com geolocalização das ocorrências por bairro

---

<div align="center">

Feito com ☕ e Golang para o Rio de Janeiro.

</div>
