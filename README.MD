# Desafio FullCycle rate limiter

Este projeto implementa um middleware de rate limiting baseado em IP e tokens de acesso, utilizando o _Redis_ como armazenamento.

## Tecnologias

- [Go](https://go.dev/)
- [Chi](https://go-chi.io/#/)
- [Redis](http://redis.io/insight/)
- [Docker](https://www.docker.com/)

## Como executar
- Execute o comando `docker compose up -d` para iniciar o Redis e o servidor.
- Caso deseje modificar os parâmetros de configuração, edite o arquivo `cmd/api/env.json`.


## Configuração Rate Limiter
### IP
- O arquivo `cmd/api/env.json` define os parâmetros relacionados ao controle de acesso por IP.


```json
  "rate_limiter": {
    "by_ip": {
      "time_window": 5,
      "max_requests": 3,
      "blocked_duration": 10
    },
    ....
  }
```

- time_window: Tempo em segundos que o rate limiter irá contar as requisições.
- max_requests: Quantidade máxima de requisições permitidas no período de tempo definido em time_window.
- blocked_duration: Tempo em segundos que o IP será bloqueado após atingir o limite de requisições.

### Token
- O arquivo `cmd/api/env.json` define os parâmetros relacionados ao controle de acesso por token.

```json
  "rate_limiter": {
    ...
    "by_api_key": {
      "time_window": 10,
      "max_requests": 5,
      "blocked_duration": 10
    }
  }
```

- time_window: Tempo em segundos que o rate limiter irá contar as requisições.
- max_requests: Quantidade máxima de requisições permitidas no período de tempo definido em time_window.
- blocked_duration: Tempo em segundos que o IP será bloqueado após atingir o limite de requisições.


Esses são os valores defaults, mas você pode personalizar esses valores para um token ou IP específico.

## Rotas

### Criar o rate limiter do IP ou Token
- Para customizar o rate limit de um IP ou Token, basta acessar a rota /rate-limiter e enviar um JSON com os valores desejados. 

#### Para customizar o rate limiter do IP
- Para customizar o rate limit de um IP, basta executar o comando curl abaixo. Não é necessário enviar o IP no corpo da requisição, pois o IP é obtido automaticamente pelo middleware. 
Os valores que estão no arquivo `env.json` deixarão de ser usados para esse IP que está sendo customizado.

```c
curl --location 'http://localhost:8080/rate-limiter' \
--header 'Content-Type: application/json' \
--data '{
    
   "max_requests" : 3,
    "time_window" : 5,
    "block_duration": 5
}'
```

Resposta:

```json
{
    "key": "$[SEU_IP_AQUI]",
    "max_requests": 3,
    "time_window": 5,
    "block_duration": 5
}
```

#### Para customizar o rate limiter do Token
- Para customizar o rate limit do Token, basta executar o curl abaixo. Um novo token será criado com os valores customizados e os valores que estão no arquivo `env.json` deixarão de ser usados para esse token.

```c
curl --location 'http://localhost:8080/rate-limiter' \
--header 'Content-Type: application/json' \
--data '{
    "is_token" : true,
   "max_requests" : 3,
    "time_window" : 5,
    "block_duration": 5
}'
```
Resposta:

```json
{
    "key":  "$[SEU_TOKEN_AQUI]",
    "max_requests": 3,
    "time_window": 5,
    "block_duration": 5
}
```

### Hello World
- Essa rota serve para testar o rate limiter, basta acessar a rota /hello-world. 

```c
curl --location 'http://localhost:8080/hello-world' \
--header 'API_KEY: $[SEU_TOKEN_AQUI]''
```

>Vale ressaltar que o header API_KEY é opcional. Caso não seja enviado, a regra do rate limiter será aplicado apenas ao IP do cliente.


## Como testar?
Para aprender mais sobre o CLI, você pode executar o seguinte comando:

```docker
 docker compose run --rm go-cli-test -h 
```

Resposta:

```
Usage of ./cli-test:
  -k string
        API Key for the request
  -m string
        HTTP method to use (default "GET")
  -r int
        Maximum amount of requests to send (default 100)
  -t int
        Time in seconds of each request (default 1)
  -url string
        URL to test
```
### Testando com IP:
- Após [executar a aplicação](#como-executar), você pode executar o seguinte comando para testar o rate limiter:

```docker
docker compose run --rm go-cli-test -url http://go-app:8080/hello-world -m GET -t 1 -r 10
```

### Testing with API key

- Primeiro é necessário [criar um token](#para-customizar-o-rate-limiter-do-token)
- Com o token criado, é necessário executar o seguinte comando com a **flag -k** com o valor do token.

```docker
docker compose run --rm go-cli-test -url http://go-app:8080/hello-world -m GET -t 1 -r 10  -k $[SEU_TOKEN_AQUI]
```

## Como funciona?
- A struct de controle de rate limit é essa:

```go
// RateLimiter is a struct that implements rate limiting logic.
// It's designed to limit the number of Requests allowed within a specified time window.
type RateLimiter struct {
    // Requests is a slice of time.Time that holds the timestamps of the incoming Requests.
    // It is used to keep track of the Requests that have been made and whether a new request
    // should be allowed or not based on the timing of previous Requests.
    Requests []time.Time

    // TimeWindowSec specifies the duration in SECONDS of the time window for which the Requests are counted.
    // For instance, if TimeWindow is set to 60, the RateLimiter will only consider
    // the number of Requests in the last minute.
    TimeWindowSec int64

    // MaxRequests defines the maximum number of Requests that are allowed within the TimeWindow.
    // For example, if we have TimeWindowSec = 1 and MaxRequests = 100, we obtain the max request limit 100 req/s
    // Other example, if we have TimeWindowSec = 60 and MaxRequests = 100, we obtain the max request limit 100 req/min
    MaxRequests int

    // lock is a mutex that ensures that access to the Requests slice is synchronized across multiple goroutines.
    lock sync.Mutex
}

```

- **Start Request:** O processo começa quando um cliente faz uma request a qualquer [endpoint](#rotas).

- **Get rate limit config:** O sistema verifica o endereço IP do cliente e verifica se há alguma [customização](#para-customizar-o-rate-limiter-do-ip) de rate limit para esse endereço IP. Se não houver customização, ele usa os valores padrão definidos no arquivo de configuração `env.json`. Caso haja uma customização, ele usa os valores customizados para esse endereço IP. 
O mesmo ocorre para o token. Quando o header **API_KEY** é enviado, o sistema verifica se há uma [customização](#para-customizar-o-rate-limiter-do-token) de rate limit para esse token. Se não houver customização, ele usa os valores padrão definidos no arquivo de configuração `env.json`. Caso haja uma customização, ele usa os valores customizados para esse token.

- **Save Request:** A request é salva no Redis. O Redis armazena as informações da request associadas ao endereço IP do cliente ou Token.

- **API key or IP exist?:** O sistema consulta o Redis para verificar se já existe algum registro de request para o endereço IP ou Token. Se não houver uma entrada existente, ele cria uma nova instância no Redis para rastreamento.

- **Add request in the request list:** Se existir uma entrada para o endereço IP ou token, a nova request é adicionada à lista associada a esse IP ou Token no Redis. Essa lista mantém o controle de todas as requests dentro de uma determinada time window.

- **Verify Blocked Duration:** Antes de prosseguir, o sistema verifica se o cliente foi bloqueado devido ao excesso de requests. Se o cliente estiver atualmente bloqueado, o processo é encerrado e falha.

- **Verify Limit:** Se o cliente não estiver bloqueado, o sistema verifica se a request atual excede o rate limit estabelecido (o número máximo de requests permitidas na time window).

- **Continue to the endpoint/End:** Se a request não exceder o rate limit, o processo continua para o endpoint, ou seja, a API tratará a request normalmente. Se o rate limit for excedido, o processo é encerrado, indicando que a request não será processada, e normalmente o cliente receberá uma resposta de erro de rate limiting.

## 📄 Licença

Este projeto está licenciado sob a Licença MIT. Veja o arquivo [LICENSE](./LICENSE) para mais detalhes.