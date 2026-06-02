# Contribuindo ao Almodon

Contribuições ao Almodon são limitadas, atualmente, apenas ao integrantes da sua equipe de desenvolvimento, reporte de bugs, críticas e sugestões são bem-vindas através do sistema de issues do GitHub.

## Diretrizes de Contribuição de Código

Essa seção lista as diretrizes para a contribuição de código, ou seja, o processo de desenvolvimento e o formato de commits, issues e pull requests.

### 1. Começar a Contribuir

O desenvolvimento é feito por meio the _one-off branches_, ou seja, temos a master, quando qualquer alteração deseja ser feita, vai-se efetuar um _fork_ do repositório, criar uma nova branch, desenvolver o desejado, e então fazer o seu pull request.

Não serão aceitos pull requests direto na master. Talvez em algum momento haja mais de uma branch ativa, mas isso será comunicado, contudo continuaria sendo o mesmo modelo, _one-off_ destas novas branchs.

### 2. Formato de Commits

Com exceção do primeiro commit "ready, set, Go!", todo commit deve seguir a estrutura:

```bnf
<tipo> [ "(" <escopo> ")" ] [ "!" ] ": " <mensagem>
```

O `<tipo>` pode ser um dos valores a seguir:

- `feat`: para a adição de uma nova funcionalidade
- `doc`: para a documentação de funcionalidades
- `fix`: para a correção de bugs
- `test`: para a adição ou modificação de testes
- `refactor`: para mudanças no código que não adicionam funcionalidades ou consertam bugs
- `style`: para mudanças que não afetam o significado do código (espaços em branco, formatação, etc)
- `merge`: para commits de merge

O `<escopo>` refere-se ao pacote ou unidade semântica afetada pelo commit, mudanças que envolvem mais de um pacote ou unidade semântica em um único commit devem ser evitados.

O `!` opcional descreve `BREAKING CHANGE`, que é presente quando uma API pública é alterada, que pode fazer outros pacotes, em cascata, terem que ser alterados.

A `<mensagem>` é uma descrição curta, normalmente uma única oração no infinitivo, comumente, mas não obrigatoriamente em inglês, que descreve as alterações feita nesse commit. A mensagem não deve terminar com ponto final nem conter letras maiúsculas, exceto para identificadores e siglas.

### 3. Formato de issues

Issues são usadas como bilhetes de desenvolvimento. Tarefa devem ser postadas, principalmente, por meio dessas. O título deve começar com um verbo no infinitivo, em português, e ser uma descrição curta da issue. A descrição, também em português, deve conter uma explicação detalhada do que a issue propõe.

Uma vez postada, a issue vai para o `Backlog` do projeto, onde é possível associar uma data limite e, se acreditar passível, passá-las para o estado `Ready`. Quando tomar responsabilidade sobre uma issue, comente-a brevemente e a mova para `In Progress`, assim que pronto, com um pull request associado, move-la para `In Review`, onde irá aguardar revisão.

### 4. Formato das branches

Com exceção da branch `master` e possíveis auxiliares futuras, ou seja, toda _one-off_ branch deve ser nomeadas na forma `<tipo> / <unidade>`.

O `<tipo>` deve ser um dos a seguir:

- `feat`: para a adição de uma nova funcionalidade
- `doc`: para a documentação de funcionalidades
- `fix`: para a correção de bugs
- `test`: para a adição ou modificação de testes
- `refactor`: para mudanças no código que não adicionam funcionalidades ou consertam bugs
- `style`: para mudanças que não afetam o significado do código (espaços em branco, formatação, etc)

E a `<unidade>` é algo que descreve o que está sendo feito, como: `e2e-user` para testes ponta-a-ponta (end-to-end) do recurso de usuário; ou `stock` para a implementação do estoque.

### 5. Formato de pull requests

Pull requests devem possuir um título e uma descrição. O título deve seguir o mesmo formato dos títulos das issues, ou seja, começar com um verbo no infinitivo, em português, e ser uma descrição curta do pull request.

A descrição, também em português, deve conter uma explicação detalhada das alterações feitas no pull request, incluindo o motivo dessas alterações e quaisquer informações relevantes para a revisão do código.

Caso o pull request esteja relacionado a uma issue, deve-se referenciar a issue na descrição, caso ele feche alguma issue, deve-se usar a sintaxe `resolve #<número da issue>`, pois "resolve" é a palavra reservada para fechar issues automaticamente no GitHub, e essa é a que mais se aproxima do português com um sentido natural.

## Diretrizes para a Padronização de Código

Essa seção lista todas as condições consideradas essenciais para o cultivo de uma base de código sem surpresas. A violação de qualquer uma dessas diretrizes resultará em rejeição de pull requests.

### 1. Todos os identificadores devem estar em inglês

Como a linguagem Go tem sua estrutura de controle e biblioteca padrão escritas em inglês, é natural que todo o código seja em inglês.

### 2. Padrão de nomeação de identificadores

- Identificadores exportados, naturalmente, seguem PascalCase, até mesmo constantes e variáveis exportadas. Isto é, não usar SCREAMING_SNAKE_CASE.

- Identificadores não constantes não exportados e parametros e variáveis em escopos seguem snake_case.

- (Talvez seja revogada) Identificadores constantes não exportados seguem _PascalCase, que é o mesmo que PascalCase, mas com um underline (`_`), prefixado.

- Identificadores de rótulos seguem PascalCase.

### 3. Nomeação de funções e métodos, tipos e variáveis

Funções e método são nomeadas a partir do que fazem e retornam, e caso necessário, para suprir a ausência de sobrecarga, o que recebem. É preferível `<ação> [ <objeto de ação> ]`. Para métodos, é comum que o receptor tenha apenas uma ou duas letras, pois seu tipo já apresenta informações suficiente. Também não é incomum para funções e baixo nível terem nomes com apenas uma ou duas letras, porém, para as políticas de alto nível, essa prática é desencorajada.

Para tipos, a menos que esteja referindo-se à uma entidade (no contexto de modelagem), ela não precisa carregar o nome da entidade. Por exemplo, em Java, é comum que exista `UserController`, já em Go, como o pacote é sempre referenciado por todos os seus identificadores exportados, `user.UserController` é exagerado, portando conforma-se apenas a `Controller`.

Para entidades, tipos que são um subconjunto de sua representação real obrigada por _setters_, deve-se seguir o seguinte modelo:

```go
type User struct {
    name string
}

func (u *User) Name() string
func (u *User) SetName(string) error
```

Os _getters_ não possuem "get" no seu nome e os _setters_ sempre possibilitam o retorno de um error, mesmo que, no momento, aceitem todo valor.

Construtores seguem o padrão `New <tipo>`, normalmente retornado um tipo (ou uma referência para esse tipo) e, opcionalmente, um erro. Como:

```go
func NewUser(string) (User, error)
```

Se um pacote compreender apenas um tipo expressivo, como o pacote `errors` da biblioteca padrão, usar apenas `New` é encorajado.

### 4. Código deve ser formatado pelo `go fmt`

A ferramenta `fmt` da cadeia de ferramentas do compilador fornece um estilo bem definido, um formato canônico, que faz esse documento ser mais simples. O comando `go fmt` é integrado na grande maioria de editores de texto juntamente a outras ferramentas da linguagem.

### 5. Escopos devem não apresentar linhas em branco no começo ou no final

Escopos, blocos de código delimitados por `{` e `}`, devem sempre possuir uma declaração no primeira e última linha, ou seja, o código a seguir é encorajado:

```go
func main() {
    var num int
    fmt.Print("Enter a number: ")
    fmt.Scanf("%d\n", &num)

    if num <= 0 {
        fmt.Printf("The number %d is non-positive")
    } else {
        fmt.Printf("The number %d is positive")
    }
}
```

Enquanto o código a seguir será rejeitado:

```go
func main() {
                                                \\ <- linha em branco no início de escopo
    var num int
    fmt.Print("Enter a number: ")
    fmt.Scanf("%d\n", &num)

    if num <= 0 {
        fmt.Printf("The number %d is non-positive")
                                                \\ <- linha em branco no final de escopo
    } else {
        fmt.Printf("The number %d is positive")
    }
}
```
