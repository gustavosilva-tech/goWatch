# GoWatch

Monitor de uptime HTTP simples escrito em Go, sem dependências externas.
Verifica periodicamente uma lista de URLs em paralelo e expõe o estado mais
recente através de uma API JSON.

## Como correr

```bash
go run . -targets "https://example.com,https://go.dev"
```

Ou com Docker:

```bash
docker build -t gowatch .
docker run -p 8080:8080 -e GOWATCH_TARGETS="https://example.com" gowatch
```

## Configuração

Cada opção pode ser passada como flag ou como variável de ambiente
(as flags têm prioridade).

| Flag        | Variável           | Default          | Descrição                              |
|-------------|--------------------|------------------|----------------------------------------|
| `-targets`  | `GOWATCH_TARGETS`  | *(obrigatório)*  | URLs http(s) separadas por vírgulas     |
| `-addr`     | `GOWATCH_ADDR`     | `127.0.0.1:8080` | Endereço onde o servidor escuta         |
| `-interval` | `GOWATCH_INTERVAL` | `30s`            | Intervalo entre rondas de verificação   |
| `-workers`  | `GOWATCH_WORKERS`  | `4`              | Número de verificações em simultâneo    |
| `-timeout`  | `GOWATCH_TIMEOUT`  | `5s`             | Timeout de cada pedido                  |

Um URL é considerado **online** se a resposta final (depois de seguir
redirects) tiver um status `2xx`.

## API

### `GET /health`

Liveness do próprio GoWatch.

```json
{"service": "gowatch", "status": "ok"}
```

### `GET /status`

Resultado da última ronda de verificações. `updated_at` é `null` até a
primeira ronda terminar.

```json
{
  "updated_at": "2026-10-26T19:40:50.549Z",
  "summary": {"total": 2, "online": 1, "offline": 1},
  "results": [
    {
      "url": "https://example.com",
      "online": true,
      "status_code": 200,
      "response_time_ms": 84,
      "checked_at": "2026-10-26T19:40:50.548Z"
    },
    {
      "url": "http://127.0.0.1:1",
      "online": false,
      "status_code": 0,
      "response_time_ms": 0,
      "checked_at": "2026-10-26T19:40:50.548Z",
      "error": "dial tcp 127.0.0.1:1: connect: connection refused"
    }
  ]
}
```

## Estrutura

| Ficheiro       | Responsabilidade                                       |
|----------------|--------------------------------------------------------|
| `main.go`      | Arranque, rotas HTTP e graceful shutdown               |
| `config.go`    | Leitura e validação de flags / variáveis de ambiente   |
| `checker.go`   | Verificação de um único URL                            |
| `monitor.go`   | Worker pool que verifica vários URLs em paralelo       |
| `scheduler.go` | Corre as verificações a cada `interval`                |
| `store.go`     | Guarda o último resultado de forma thread-safe         |

## Testes

```bash
go test -race ./...
```
