# Desafio FullCycle rate limiter

Este projeto implementa um middleware de rate limiting baseado em IP e tokens de acesso, utilizando Redis como armazenamento.

## ⚙️ Tecnologias

- [Go](https://go.dev/)
- [Chi](https://go-chi.io/#/)
- [Redis](http://redis.io/insight/)
- [Docker](https://www.docker.com/)

## Como executar
- Basta executar o comando `docker-compose up -d` para subir o Redis e o servidor. Caso queira editar os valores de configuração, basta editar o arquivo `cmd/api/env.json`.


## ⚙️ Configuração Rate Limiter
### IP
- No arquivo `cmd/api/env.json` possui os valores relacioados ao rate limiter por IP.


```
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
- No arquivo `cmd/api/env.json` possui os valores relacioados ao rate limiter por token.

```
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

## 📝 Rotas

### Criar o rate limiter do IP ou Token
- Para customizar o rate limit de um IP ou Token, basta acessar a rota /rate-limiter/ e enviar um JSON com os valores desejados. 

#### Para customizar o rate limiter do IP
- Para customizar o rate limit do IP, basta executar o curl abaixo. Não é necessário enviar o IP no body da requisição pois o IP é obtido automaticamente pelo middleware.
Os valores que estão no arquivo env.json deixarão de ser usados para esse IP que está sendo customizado.

```
curl --location 'http://localhost:8080/rate-limiter' \
--header 'Content-Type: application/json' \
--data '{
    
   "max_requests" : 3,
    "time_window" : 5,
    "block_duration": 5
}'
```

Resposta:

```
{
    "key": "$[SEU_IP_AQUI]",
    "max_requests": 3,
    "time_window": 5,
    "block_duration": 5
}
```

#### Para customizar o rate limiter do Token
- Para customizar o rate limit do Token, basta executar o curl abaixo. Um novo token será criado com os valores customizados e os valores que estão no arquivo env.json deixarão de ser usados para esse token.

```
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

```
{
    "key":  "$[SEU_TOKEN_AQUI]",
    "max_requests": 3,
    "time_window": 5,
    "block_duration": 5
}
```

### Hello World
- Para testar o rate limiter, basta acessar a rota /hello-world. 

```
curl --location 'http://localhost:8080/hello-world' \
--header 'API_KEY: e38f28e492c34d7a95a208476720b63510c52663c05be73e7e905bdf2d3dccec'
```
>Vale ressaltar que o header API_KEY é opcional, caso não seja enviado, o rate limiter será aplicado apenas para o IP do cliente.



## How can I test?
To learn more about the CLI using for test, you can run this:

``` docker compose run --rm go-cli-test -h ```

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
### Testing with IP:
- Apos [executar a aplicação](#como-executar), você pode executar o seguinte comando para testar o rate limiter:

```
docker compose run --rm go-cli-test -url http://go-app:8080/hello -m GET -t 1 -r 10
```

### Testing with API key

- Primeiro é necessario [criar um token](#para-customizar-o-rate-limiter-do-token)
- Com o token criado, é necessario executar o seguinte comando com a flag -k com o valor do token.

```
docker compose run --rm go-cli-test -url http://go-app:8080/hello -m GET -t 1 -r 10  -k $[SEU_TOKEN_AQUI]
```

## 🏭 Como funciona?


## 📄 Licença

Este projeto está licenciado sob a Licença MIT. Veja o arquivo [LICENSE](./LICENSE) para mais detalhes.