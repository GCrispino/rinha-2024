# rinha-de-backend-2024-q1-participacao-go
Participação na Rinha de Backend 2024 

## Tecnologias utilizadas:
- `nginx` como load balancer;
- `postgres` como banco de dados;
- `go` para api com as seguintes bibliotecas: 
    - [`echo`](https://echo.labstack.com/) para o servidor;
    - [`go-json`](https://github.com/goccy/go-json/) para serialização de JSON no servidor;
    - [`pgx`](https://github.com/jackc/pgx) como driver de Postgres;
    - [`automaxprocs`](https://github.com/uber-go/automaxprocs) para configurar automaticamente a variável de ambiente [`GOMAXPROCS`](https://dave.cheney.net/tag/gomaxprocs).