# Configuração e tipos

## Cliente e transporte

```go
client, err := raims.NewClient(apiKey,
    raims.WithTimeout(80*time.Second),
    raims.WithMaxNetworkRetries(2),
)
```

O timeout padrão é 80 segundos por tentativa, incluindo a leitura da resposta.
O transporte padrão limita a conexão a dez segundos e reutiliza conexões.
Passe um contexto com prazo para limitar toda a chamada, incluindo retries e
esperas. `WithRequestTimeout` ajusta o prazo de uma chamada.

`WithHTTPClient` aceita um `*http.Client` da aplicação. O SDK copia sua
configuração, recusa redirects e não usa cookies. O transporte fornecido
precisa respeitar o contexto e preservar a política de envio único de escritas.
`CloseIdleConnections` encerra conexões ociosas do transporte usado pelo cliente.
Quando compartilhar um transporte, essa chamada também afeta suas conexões ociosas.

`WithBaseURL` aceita uma origem HTTPS confiável, sem caminho, query ou credenciais.
HTTP é aceito somente em IP de loopback para validação local. A origem recebe
a credencial; use somente um destino sob seu controle.

## API key e OAuth

`NewClient` usa uma chave de API. A chave determina o ambiente; o cliente não
oferece seleção de ambiente nem envia `X-Raims-Environment` nesse modo.

Para um aplicativo autorizado por OAuth:

```go
client, err := raims.NewOAuthClient(accessToken, raims.Sandbox,
    []raims.Environment{raims.Sandbox, raims.Production},
)
```

Informe somente os ambientes efetivamente concedidos à conexão. O SDK recebe
um access token obtido pelo fluxo da aplicação; autorização, armazenamento e
renovação do token continuam sob responsabilidade dela. Crie outro cliente
quando renovar o token ou mudar de ambiente.

O cliente envia o ambiente escolhido e recusa endpoints incompatíveis antes
do HTTP. Simulações aceitam somente Sandbox. Não há troca automática de ambiente.
A API continua responsável por conferir a credencial, o recurso e as permissões.

## Omissão, null e valores

Campos de entrada usam `Optional[T]`. Seu valor zero omite o campo.
`Set(valor)` envia o valor, inclusive `false`, zero ou texto vazio.
`Null[T]()` envia null quando o campo permitir.

```go
params := &raims.UpdateCustomerParams{
    FullName: raims.Null[string](),
    Metadata: raims.Set(map[string]*string{}),
}
```

Esse exemplo remove o nome e preserva metadata pelo envio de `{}`. Para apagar
todo o mapa, use `raims.Null[map[string]*string]()`. Para remover uma chave,
envie seu valor como nil. Use um mapa ou slice vazio inicializado para enviar
`{}` ou `[]`; um mapa ou slice nil é serializado como null.

As constantes de enums são geradas do contrato. Inteiros usam `int64`, inclusive
nas constantes numéricas. Em literais, use `Set[int64](15000)` para corresponder
ao campo sem inferir `int`.

## Respostas

Respostas têm campos tipados, `LastResponse` e `RawJSON`. `HasField` distingue
um campo ausente; `IsNull` identifica null explícito. `ExtraFields` preserva
campos futuros e `RawJSON` guarda a resposta original. Use `RawJSON` quando
precisar conservar exatamente a representação recebida.

`json.Marshal` preserva campos futuros, omissão e null das respostas recebidas.
Ele também inclui alterações feitas nos campos tipados e em `ExtraFields`.
`HasField` e `IsNull` descrevem a resposta recebida; não acompanham edições
posteriores na struct. Os nomes JSON são associados exatamente, para que um
campo desconhecido não substitua uma propriedade conhecida.

Os bytes de `ExtraFields` são independentes de `RawJSON`. Você pode editar um
valor extra ou acrescentar bytes a ele sem alterar o snapshot recebido nem os
indicadores de presença e null.

A validação estrutural confere campos obrigatórios, tipos e nullability.
Centavos e datas não passam por ponto flutuante. Novos valores de enums em
respostas continuam legíveis; decida como tratá-los na aplicação.

`LastResponse` mantém status, `RequestID`, `RetryAfter` e
`IdempotencyReplayed`. O status `202` permanece pendente, sem polling automático.
`Charges.Revoke` usa `RevokeChargeResult` para representar tanto `200` quanto `202`.

Métodos de documentos e preparação de upload retornam os dados fornecidos pela
API. O SDK não baixa as URLs retornadas nem executa o upload com a chave de API.

## Diagnóstico

`WithRequestObserver` recebe método, caminho sem IDs, tentativa, duração, status
e Request-Id após cada tentativa. O callback pode ser chamado por goroutines
diferentes e precisa ser seguro para uso concorrente. Uma falha no callback
não altera o resultado da chamada.

O callback executa de forma síncrona e deve retornar rapidamente. Seu tempo
de execução não descarta uma resposta conclusiva já recebida, inclusive em
escritas. Os prazos do contexto continuam controlando o transporte e novas
tentativas elegíveis.

Agregue erros de conexão, HTTP e protocolo para detectar indisponibilidade ou
mudança incompatível de resposta. Confirme recuperação por uma nova leitura
válida. Evite registrar objetos completos, payloads e URLs temporárias.
