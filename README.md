# feedex
demo for querying trx fee of uniswap

# Requirements
- Realtime data recording of UniswapV3 USDC/ETH pool transaction fee
- Historical data batch data recording of UniswapV3 USDC/ETH pool transaction fee
- Convert fee in ETH to USDT using historical/live price for ETH/USDT 
- Batch query by given time period
- Single query by given transaction hash

## Tech Requirements
- RESTful API follwoing Swagger UI standards
- Read config from file
- Gracefully start and stop
- Rate limit for external api calls
- Dockerize your applications withdocker-compose

# Tech Design
## Database Design
We can use mysql to store all the required data of UniswapV3 USDC/ETH transactions.
refer to `scripts/mysql/init.sql`

## API Design
### Query trsanction fee of single transaction
input: 
- trx_hash, string

output:
- trx_fee, string, decimal number

### Batch query transaction fees given time period
input:
- symbol, string, WETH/USDC by default
- start_time, int, unix timestamp in seconds
- end_time, int, unix timestamp in seconds
- limit, int, limit of results size, 20 by default
- page, int, starting from 0

output:
- result, array of json struct
  - TrxHash, string
  - TrxFeeUsdt, string, decimal number
  - TrxTime, int, unix timestamp in seconds

## Architecture
![alt text](image.png)

# Build & Run

## build
run `make build`

## run
befor running, need to set up the mysql database
1. start mysql server locally
2. execute `scripts/mysql/init.sql`
3. modify the `config.yml` accordingly
4. run `make run`


## Swagger docs

visit `http://localhost:8080/swagger/index.html`

# Test
Unit tests are mainly generated by ChatGPT.