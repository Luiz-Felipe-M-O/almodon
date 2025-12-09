# Almodon - Manual do Desenvolvedor

## Índice

1. [Visão Geral do Projeto](#visão-geral-do-projeto)
2. [Arquitetura](#arquitetura)
3. [Stack Tecnológica](#stack-tecnológica)
4. [Estrutura do Projeto](#estrutura-do-projeto)
5. [Configuração de Desenvolvimento](#configuração-de-desenvolvimento)
6. [Modelo de Domínio](#modelo-de-domínio)
7. [Design da API](#design-da-api)
8. [Arquitetura Frontend](#arquitetura-frontend)
9. [Design do Banco de Dados](#design-do-banco-de-dados)
10. [Autenticação e Autorização](#autenticação-e-autorização)
11. [Diretrizes de Estilo de Código](#diretrizes-de-estilo-de-código)
12. [Fluxo de Desenvolvimento](#fluxo-de-desenvolvimento)
13. [Testes](#testes)
14. [Deploy](#deploy)
15. [Contribuindo](#contribuindo)

## Visão Geral do Projeto

**Almodon** é um sistema de gestão de estoque desenvolvido para auxiliar as atividades do almoxarifado interno da Odontologia da UFVJM (Universidade Federal dos Vales do Jequitinhonha e Mucuri).

### Propósito
- Gerenciar inventário de materiais e suprimentos odontológicos
- Rastrear lotes de itens com datas de validade
- Gerenciar requisições de clínicas e laboratórios
- Fornecer gestão de usuários com controle de acesso baseado em funções
- Gerar relatórios e gerenciar fornecedores

### Contribuidores da Equipe
- Alan Barbosa Lima ([@alan-b-lima](https://github.com/alan-b-lima))
- Breno Augusto Braga Oliveira ([@bragabreno](https://github.com/bragabreno))
- Lucas Rocha Oliveira ([@Lucas-Rocha-Oliveira](https://github.com/Lucas-Rocha-Oliveira))
- Luiz Felipe Melo Oliveira ([@Luiz-Felipe-M-O](https://github.com/Luiz-Felipe-M-O))
- Otávio Gomes Calazans ([@otaviogomes03](https://github.com/otaviogomes03))
- Rafael Gomes Silva ([@rafleGomes](https://github.com/rafleGomes))

## Arquitetura

O Almodon segue uma abordagem de **Arquitetura Limpa** com princípios de **Domain-Driven Design (DDD)**, implementado como uma aplicação monolítica com clara separação de responsabilidades.

### Arquitetura de Alto Nível

```
┌────────────────────────────────────────────────────┐
│               Camada de Apresentação               │
│ ┌───────────────────────┐┌───────────────────────┐ │
│ │  Web UI (TypeScript)  ││  REST API (Go/HTTP)   │ │
│ └───────────────────────┘└───────────────────────┘ │
└────────────────────────────────────────────────────┘
┌────────────────────────────────────────────────────┐
│                Camada de Aplicação                 │
│ ┌────────────────────────────────────────────────┐ │
│ │           Serviços & Portões de Auth           │ │
│ └────────────────────────────────────────────────┘ │
└────────────────────────────────────────────────────┘
┌────────────────────────────────────────────────────┐
│                 Camada de Domínio                  │
│ ┌───────┐┌───────────┐┌────────┐┌───────────┐      │
│ │ User  ││ Material  ││  Item  ││  Session  │ ...  │
│ └───────┘└───────────┘└────────┘└───────────┘      │
└────────────────────────────────────────────────────┘
┌────────────────────────────────────────────────────┐
│              Camada de Infraestrutura              │
│ ┌────────────────┐┌──────────────────────────────┐ │
│ │ Mapa (Efêmero) ││ Banco de Dados (Persistente) │ │
│ └────────────────┘└──────────────────────────────┘ │
└────────────────────────────────────────────────────┘
```

### Padrões Arquiteturais

1. **Padrão Repositório**: Abstração de acesso a dados
2. **Padrão Serviço**: Lógica de negócio encapsulada
3. **Portões de Autorização**: Controle de acesso
4. **Objetos de Transporte**: Comunicação entre camadas
5. **Entidate-Repositório-Serviço-Transporte**: Organização do domínio

## Tecnologias Utilizadas

### Backend

- **Linguagem**: Go 1.25.1
- **Banco de Dados**: PostgreSQL (com driver pgx/v5)
- **Geração SQL**: sqlc para acesso com tipos seguros à SQL
- **Servidor HTTP**: Servidor HTTP nativo do Go
- **Autenticação**: Sistema de sessão customizado similar ao JWT
- **Hash de Senhas**: bcrypt

### Frontend

- **Linguagem**: TypeScript (target ESNext)
- **Framework UI**: Framework customizado leve (jsxmm)
- **Ferramenta de Build**: Compilador TypeScript nativo

### Ferramentas de Desenvolvimento

- **Geração de Código**: sqlc para consultas de banco de dados
- **Controle de Versão**: Git
- **Licença**: GPLv3

## Estrutura do Projeto

```
almodon/
├── cmd/                 // Pontos de entrada da aplicação
│   ├── main.go          // Servidor principal da aplicação
│   └── test/main.go     // Servidor/utilitários de teste
├── internal/            // Código privado da aplicação
│   ├── api/v1/          // Implementação da API versão 1
│   ├── auth/            // Autenticação e autorização
│   ├── domain/          // Entidades de domínio e lógica de negócio
│   │   ├── item/        // Gestão de itens do estoque
│   │   ├── material/    // Gestão do catálogo de materiais
│   │   ├── promotion/   // Sistema de promoções
│   │   ├── requisition/ // Gestão de requisições
│   │   ├── session/     // Gestão de sessões de usuário
│   │   └── user/        // Gestão de usuários
│   ├── support/         // Utilitários compartilhados e infraestrutura
│   └── xerrors/         // Erros específicos do domínio
├── pkg/                 // Pacotes públicos reutilizáveis
│   ├── closer/          // Utilitários de limpeza de recursos
│   ├── errors/          // Framework de tratamento de erros
│   ├── hash/            // Hash de senhas
│   ├── heap/            // Utilitários de estruturas de dados
│   ├── opt/             // Tipos opcionais
│   └── uuid/            // Utilitários UUID
├── ui/web/              // Aplicação frontend
│   ├── src/             // Código fonte TypeScript
│   │   ├── internal/    // Módulos internos
│   │   └── module/      // Módulos reutilizáveis
│   ├── index.html       // Template HTML principal
│   └── tsconfig.json    // Configuração TypeScript
├── db/                  // Arquivos do banco de dados
│   ├── schema.sql       // Schema do banco de dados
│   ├── query.sql        // Consultas SQL para sqlc
│   └── connection.go    // Lógica de conexão do banco
└── bin/                 // Binários compilados
```

### Organização de Domínio

Cada domínio segue uma estrutura consistente:

```
domain/{entity}/
├── entity.go     // Entidade de domínio com regras de negócio
├── repository.go // Interface do repositório
├── service.go    // Interface do serviço
├── transport.go  // Objetos de transferência de dados
├── repository/
│   ├── db.go     // Implementação PostgreSQL
│   └── map.go    // Implementação em memória
├── resource/
│   └── http.go   // Handlers HTTP
└── service/
    ├── auth.go   // Portões de autorização
    └── core.go   // Lógica de negócio principal
```

## Configuração de Desenvolvimento

### Pré-requisitos

1. **Go 1.25.1+**
2. **PostgreSQL 13+** (para produção)
3. **npm** (opcional, para ferramentas adicionais)
4. **tsc** (Compilador do TypeScript)
5. **Git**

### Configuração de Desenvolvimento Local

1. **Clone o repositório**:

    ```bash
    git clone https://github.com/alan-b-lima/almodon.git
    cd almodon
    ```

2. **Instale as dependências Go**:

    ```bash
    go mod download
    ```

3. **Configure o banco de dados** (opcional para desenvolvimento):

    ```bash
    # Criar banco de dados
    createdb almodon_dev

    # Aplicar schema
    psql almodon_dev < db/schema.sql
    ```

4. **Compile a aplicação**:

    ```bash
    go build -o bin/almodon cmd/main.go
    ```

5. **Execute o servidor**:

    ```bash
    ./bin/almodon
    ```

6. **Acesse a aplicação**:

    - Servidor roda em `http://localhost:4545`
    - Interface web é servida no caminho raiz

### Modo de Desenvolvimento

A aplicação suporta armazenamento baseado em arquivos para desenvolvimento:
- Dados de usuário: `.data/users.json`
- Dados de material: `.data/materials.json`
- Dados de item: `.data/items.json`

Esses arquivos são criados e gerenciados automaticamente pela aplicação.

### Frontend TypeScript

O frontend é construído com TypeScript e não requer etapa de build para desenvolvimento:

1. **Configurar TypeScript**:

```bash
cd ./ui/web/
tsc
# Configuração TypeScript já fornecida
```

2. **Desenvolvimento**:

   - Arquivos são servidos diretamente pelo servidor Go
   - Compilação TypeScript acontece no navegador durante desenvolvimento
   - Use extensões `.ts` para importações

## Modelo de Domínio

### Entidades Principais

#### User (`internal/domain/user/`)

- **Campos Principais**:
  - `uuid`: Identificador único
  - `siape`: ID do funcionário
  - `name`: Nome completo
  - `email`: Endereço de e-mail
  - `password`: Senha com hash
  - `role`: Função de autorização

#### Material (`internal/domain/material/`)

- **Campos Principais**:
  - `uuid`: Identificador único
  - `name`: Nome do material
  - `siads`: Código SIADS
  - `catmat`: Código CATMAT
  - `ecampus`: Código eCampus
  - `description`: Descrição detalhada
  - `unit`: Unidade de medida
  - `minQuantity`: Nível mínimo de estoque

#### Item (`internal/domain/item/`)

- **Campos Principais**:
  - `uuid`: Identificador único
  - `material`: Referência ao material
  - `supplier`: Referência ao fornecedor
  - `quantity`: Quantidade atual
  - `unitCost`: Custo por unidade
  - `arrival`: Data de chegada
  - `expiration`: Data de validade
  - `invoice`: Número da nota fiscal
  - `lot`: Número do lote
  - `notes`: Observações adicionais

#### Session (`internal/domain/session/`)

- **Propósito**: Gerencia sessões de autenticação de usuários
- **Campos Principais**:
  - `uuid`: ID da sessão
  - `user`: Referência do usuário
  - `expires`: Tempo de expiração da sessão

#### Requisition (`internal/domain/requisition/`)

- **Propósito**: Gerencia solicitações de materiais de clínicas/laboratórios
- **Campos Principais**:
  - `uuid`: Identificador único
  - `author`: Usuário solicitante
  - `status`: Status da solicitação
  - `destination`: Departamento solicitante
  - `entries`: Lista de itens solicitados

## Design da API

### Endpoints REST

A API segue convenções RESTful com padrões consistentes em todos os domínios:

#### Gestão de Usuários
```
GET    /api/v1/users        # Listar usuários
POST   /api/v1/users        # Criar usuário
GET    /api/v1/users/{uuid} # Obter usuário por UUID
PATCH  /api/v1/users/{uuid} # Atualizar usuário
DELETE /api/v1/users/{uuid} # Deletar usuário
POST   /api/v1/users/auth   # Autenticar usuário
GET    /api/v1/users/me     # Obter usuário atual
```

#### Gestão de Materiais
```
GET    /api/v1/materials                # Listar materiais
POST   /api/v1/materials                # Criar material
GET    /api/v1/materials/{uuid}         # Obter material
PATCH  /api/v1/materials/{uuid}         # Atualizar material
DELETE /api/v1/materials/{uuid}         # Deletar material
GET    /api/v1/materials/siads/{code}   # Listar por SIADS
GET    /api/v1/materials/catmat/{code}  # Listar por CATMAT
GET    /api/v1/materials/ecampus/{code} # Listar por eCampus
```

#### Gestão de Itens
```
GET    /api/v1/items                 # Listar itens
POST   /api/v1/items                 # Criar item
GET    /api/v1/items/{uuid}          # Obter item
PATCH  /api/v1/items/{uuid}          # Atualizar item
DELETE /api/v1/items/{uuid}          # Deletar item
GET    /api/v1/items/material/{uuid} # Listar por material
GET    /api/v1/items/supplier/{uuid} # Listar por fornecedor
```

### Formato de Requisição/Resposta

#### Resposta Padrão de Lista

```json
{
  "offset": 0,
  "length": 10,
  "records": [...],
  "total_records": 17
}
```

#### Resposta Padrão de Erro

```json
{
  "kind": "entrada inválida",
  "title": "email-inválido",
  "message": "endereço de email é obrigatório",
  "cause": null
}
```

#### Paginação

Todos os endpoints de lista suportam:

- `offset`: Registro inicial (padrão: 0)
- `limit`: Máximo de registros (padrão: 10, máx: sem limite definido)

### Objetos de Transporte

Cada domínio define objetos de transporte para comunicação da API:

- **Create**: Para requisições POST (criação de entidade)
- **Patch**: Para requisições PATCH (atualizações parciais)
- **Result**: Para dados de resposta
- **ListResult**: Para respostas paginadas
- **CreateResult**: Para respostas de criação

Exemplo:

```go
type Create struct {
    Name        string `json:"name"`
    SIADS       string `json:"siads"`
    CATMAT      string `json:"catmat"`
    Description string `json:"description"`
}

type Patch struct {
    Name        opt.Opt[string] `json:"name"`
    Description opt.Opt[string] `json:"description"`
}
```

## Arquitetura Frontend

### Framework: jsxmm

O frontend usa um framework customizado leve chamado **jsxmm**, (JavaScript XML Minimal Markup, ou JSX--, se tiver de mal humor):

#### Conceitos Principais

1. **Reatividade baseada em Sinais**: Usa um sistema de sinais para atualizações reativas
2. **Baseado em Componentes**: Componentes UI reutilizáveis
3. **TypeScript-first**: Suporte completo ao TypeScript
4. **Sem Etapa de Build**: Executa diretamente no navegador durante o desenvolvimento

#### Arquitetura

```
src/
├── main.ts           // Ponto de entrada da aplicação
├── internal/
│   ├── api.ts        // Construção do cliente API
│   ├── auth/         // Lógica de autenticação
│   ├── component/    // Componentes UI reutilizáveis
│   ├── context/      // Contexto da aplicação e roteamento
│   ├── domain/       // Views específicas do domínio
│   │   ├── item/     // UI de gestão de itens
│   │   ├── material/ // UI de gestão de materiais
│   │   └── user/     // UI de gestão de usuários
│   ├── pages/        // Páginas estáticas
│   └── support/      // Utilitários e helpers
└── module/           // Módulos reutilizáveis
    ├── errors/       // Tratamento de erros
    ├── jsxmm/        // Núcleo do framework
    └── uuid/         // Utilitários UUID
```

### Padrões Frontend Principais

#### 1. Padrão Gateway

```typescript
interface Gateway {
    List(offset: number, limit: number): Promise<ListResponse>
    Get(uuid: UUID): Promise<Response>
    Create(req: Entity): Promise<UUID>
    Patch(uuid: UUID, req: PartialEntity): Promise<void>
    Delete(uuid: UUID): Promise<void>
}
```

#### 2. Componentes View

```typescript
class MaterialView implements Context {
    constructor(gateway: material.Gateway) {}
    
    Final(): boolean {
        // Assim que verdadeiro, não será chamado HTML novamente
    }

    HTML(): HTMLElement {
        // Lógica de renderização do componente
    }
}
```

#### 3. Estado baseado em Sinais
```typescript
const user = new Signal.State<user.Response | null>(null)

new Signal.Effect(() => {
    const current_user = user.value
    // Reagir a mudanças do usuário
})
```

### Navegação e Roteamento

A aplicação usa um sistema de roteamento simples baseado em hash:
- `#items`: Gestão de itens
- `#materials`: Gestão de materiais  
- `#users`: Gestão de usuários
- `#profile`: Perfil do usuário
- `#about`: Página sobre

## Design do Banco de Dados

### Visão Geral do Schema

O banco de dados segue um design relacional otimizado para gestão de inventário:

#### Tabelas Principais

```sql
-- Tabela de usuários
usuarios (
    siape VARCHAR(20) PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    senha_hash VARCHAR(255) NOT NULL,
    perfil VARCHAR(50) NOT NULL,
    data_criacao TIMESTAMP DEFAULT NOW()
)

-- Catálogo de materiais
produtos (
    id SERIAL PRIMARY KEY,
    nome VARCHAR(255) NOT NULL,
    descricao TEXT,
    codigo_ecampus VARCHAR(50),
    siads VARCHAR(50),
    catmat VARCHAR(50),
    estoque_minimo INTEGER DEFAULT 0,
    unidade VARCHAR(20) NOT NULL
)

-- Lotes do inventário
lotes (
    id SERIAL PRIMARY KEY,
    id_produto INTEGER REFERENCES produtos(id),
    id_fornecedor INTEGER REFERENCES fornecedores(id),
    codigo_lote VARCHAR(100) NOT NULL,
    data_validade DATE NOT NULL,
    quantidade_atual NUMERIC(10,3),
    preco_unitario NUMERIC(10,2)
)
```

### Camada de Acesso a Dados

#### Integração SQLC

O projeto usa **sqlc** para gerar código de acesso ao banco type-safe:

```yaml
# sqlc.yaml
version: "2"
sql:
  - schema: "db/schema.sql"
    queries: "db/query.sql"  
    engine: "postgresql"
    gen:
      go:
        package: "db"
        out: "db"
        sql_package: "pgx/v5"
```

#### Implementação do Repository

Atualmente usando repositórios em memória com arquivos JSON persistentes para desenvolvimento:

```go
type Map struct {
    uuid repo.Index[uuid.UUID, int]

    repo []material.Entity
    mu   sync.RWMutex

    datapath string
}
```

O deploy de produção usará repositórios PostgreSQL implementando as mesmas interfaces.

## Autenticação e Autorização

### Sistema de Autenticação

#### Autenticação baseada em Sessão

- Usuários se autenticam com SIAPE (ID do funcionário) e senha
- Autenticação bem-sucedida cria uma sessão com expiração
- UUID da sessão é retornado ao cliente para requisições subsequentes

#### Segurança de Senha

- Senhas são hash com bcrypt com fator de custo 12
- Validação de senha inclui requisitos de comprimento

### Sistema de Autorização

#### Controle de Acesso baseado em Funções (RBAC)

```go
type Role string

const (
    User  Role = "user"  // Acesso básico
    Admin Role = "admin" // Acesso administrativo
    Chief Role = "chief" // Acesso completo ao sistema
)
```

#### Sistema de Permissões

```go
// Definir permissões
var permAdmin = auth.Permit(auth.Admin)
var permUser  = auth.Permit(auth.User)

// Verificar autorização
if err := service.Authorize(permAdmin, actor); err != nil {
    return xerrors.ErrUnauthorized
}
```

#### Portões de Autorização

Cada serviço implementa portões de autorização:

```go
type Gate struct {
    material.Service
    actor auth.Actor
}

func (s *Gate) Create(req material.Create) (uuid.UUID, error) {
    if err := service.Authorize(permAdmin, s.actor); err != nil {
        return uuid.UUID{}, err
    }
    return s.Service.Create(req)
}
```

### Recursos de Segurança

1. **Hash de Senhas**: bcrypt com alto fator de custo
2. **Gestão de Sessão**: Expiração baseada em tempo
3. **Hierarquia de Funções**: Usuários, Admins, Chefes
4. **Portões de Permissão**: Autorização a nível de método
5. **Proteção CSRF**: Validação baseada em sessão

## Diretrizes de Estilo de Código

O projeto segue padrões rígidos de codificação conforme definido em `CONTRIBUTING.md`:

### Padrões de Código Go

1. **Idioma**: Todos os identificadores devem estar em inglês
2. **Convenções de Nomeação**:
   - Identificadores exportados: PascalCase
   - Identificadores não exportados: camelCase  
   - Constantes: PascalCase (com prefixo `_` opcional para não exportados)
   - Rótulos: PascalCase

3. **Nomeação de Funções**: Baseado na ação e valor de retorno
   - Prefira o padrão `<ação> [<objeto>]`
   - Métodos: Nomes curtos de receptor (1-2 letras)
   - Construtores: padrão `New<Tipo>`

4. **Padrões de Entidade**:

```go
type User struct {
    name string
}

func (u *User) Name() string         // Getter (sem prefixo "Get")
func (u *User) SetName(string) error // Setter (sempre retorna erro)
```

5. **Formatação**: Todo código deve ser formatado com `go fmt`

Mais detalhes em [Diretrizes de Estilo Go](./CONTRIBUTING.md).

### Padrões TypeScript

1. **Sistema de Módulos**: ESModules com extensões `.ts`
2. **Segurança de Tipos**: Configuração TypeScript rígida
3. **Design de Interface**: Padrões de gateway consistentes
4. **Estrutura de Componentes**: Componentes baseados em classes

### Tratamento de Erros

Framework de erro customizado com tipos de erro estruturados:

```go
type Error struct {
    Kind    Kind      // Categoria do erro
    Title   string    // Título legível por humanos  
    Message string    // Mensagem detalhada
    Cause   error     // Causa subjacente
}
```

## Fluxo de Desenvolvimento

### Fluxo Git

1. **Branching**: Branches de feature a partir do `master`
2. **Commits**: Mensagens de commit descritivas em inglês
3. **Pull Requests**: Obrigatório para todas as mudanças
4. **Code Review**: Todos os PRs devem ser revisados

### Processo de Build

1. **Go Build**:

```bash
go build -o bin/almodon cmd/main.go
```

2. **TypeScript**: Nenhuma etapa de build necessária (execução direta)

3. **Banco de Dados**:

```bash
sqlc generate # Gerar código do banco de dados
```

## Deploy

### Deploy de Produção

1. **Configuração do Banco de Dados**:

```sql
CREATE DATABASE almodon_prod;
\c almodon_prod
\i db/schema.sql
```

2. **Configuração de Ambiente**:

```bash
export DATABASE_URL="postgres://user:pass@localhost/almodon_prod"
export PORT=8080
```

3. **Build e Execução**:

```bash
go build -ldflags="-s -w" -o almodon cmd/main.go
./almodon
```

### Monitoramento e Logging

A aplicação inclui logging estruturado:
- Logging de requisição/resposta HTTP
- Rastreamento de erros
- Monitoramento de performance

## Contribuindo

### Primeiros Passos

1. **Fork** do repositório
2. **Clone** seu fork localmente
3. **Crie** uma branch de feature
4. **Faça** suas mudanças seguindo as diretrizes de estilo
5. **Teste** suas mudanças completamente
6. **Envie** um pull request

### Diretrizes

- Todo código deve estar em inglês
- Siga os padrões arquiteturais estabelecidos
- Inclua tratamento de erro apropriado
- Adicione testes para novas funcionalidades
- Atualize a documentação conforme necessário

### Processo de Pull Request

1. **Descrição**: Descrição clara das mudanças
2. **Testes**: Evidência de testes
3. **Documentação**: Atualizada se necessário
4. **Revisão**: Abordar todos os comentários de revisão
5. **Merge**: Squash e merge quando aprovado

---

## Apêndice

### Comandos Comuns

```bash
# Executar o servidor
go run ./cmd/

# Build para produção
go build -o ./bin/almodon.exe ./cmd/main.go

# Gerar código do banco de dados
sqlc generate

# Formatar código Go
go fmt ./...

# Executar testes
go test ./...

# Verificar dependências
go mod tidy
```

### Recursos Úteis

- [Documentação Go](https://golang.org/doc/)
- [Documentação PostgreSQL](https://www.postgresql.org/docs/)
- [Documentação sqlc](https://sqlc.dev/)
- [Documentação TypeScript](https://www.typescriptlang.org/docs/)

### Links do Projeto

- **Repositório**: https://github.com/alan-b-lima/almodon
- **Licença**: GPLv3 - veja [LICENSE](./LICENSE)
- **Issues**: https://github.com/alan-b-lima/almodon/issues

---

_Última atualização: 3 de Dezembro de 2025 - Versão: 1.0.0_