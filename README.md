## Desafio #06 - Multithread Auction

**Objetivo:** Adicionar ao projeto a funcionalidade de fechamento automático com uso de go routines, baseado em tempo definido em variável de ambiente.

### Requisitos:

- Adicionar uma função que irá calcular o tempo do leilão, baseado em parâmetros previamente definidos em variáveis de ambiente
- Uma nova go routine que validará a existência de um leilão (auction) vencido (que o tempo já se esgotou) e que deverá realizar o update, fechando o leilão (auction)
- Um teste para validar se o fechamento acontece de forma automatizada

#### 🧭 Parametrização

```dotenv
##> cmd/auction/.env

GO_ENV=prd

BATCH_INSERT_INTERVAL=60s
MAX_BATCH_SIZE=4
AUCTION_INTERVAL=60s

MONGO_INITDB_ROOT_USERNAME: admin
MONGO_INITDB_ROOT_PASSWORD: admin
MONGODB_URL=mongodb://admin:admin@mongodb:27017/auctions?authSource=admin
MONGODB_DB=auction
```

#### 🚀 Execução:
Para executar a aplicação em ambiente local, basta utilizar o docker-compose disponível na raiz do projeto. Para isso, execute o comando abaixo:
```bash
$ docker-compose up # -d (para executar em background)
```

> 💡 **Portas necessárias:**
> - Aplicação: 8080
> - MongoDB: 27017

#### 🧪 Teste:

Foi incluído, conforme pedido no escopo da atividade, um teste de integração para validar o fechamento automático do leilão. Para executar o teste, basta executar o comando abaixo:

```bash
$ go test -v github.com/tiagoncardoso/fc-pge-auction-multithread/internal/infra/database/auction # -count=1 (para evitar cache de execução)
```
