# distribution-service
This service is responsible for the distribution of 10% of the DAOs daily received funds between all validators. This ensures that the validators always have sufficient funds to operate on the network.

## Mechanics
```mermaid
sequenceDiagram
    loop Cronjob
        Service->>LiquidWallet: ListReceivedByAddress
        LiquidWallet-->>Service: txDetails{amount: n}
        Service-->>Service: Compare amount to last occurrence

        Service->>Planetmint: GetValidators
        Planetmint-->>Service: validators{operatorAddress: "valoper"}
        loop for every active validator
        Service->>Planetmint: GetValidatorDelegations
        Planetmint-->>Service: delegations{delegatorAddress: "valaddress"}
        end

        loop for every delegatorAddress
        Service->>R2P: GetRecieveAddress
        R2P-->>Service: liquidAddress
        end

        Service-->>Service: Calculate share per delegator address

        loop for every liquidAddress
        Service->>ShamirCoordinator: IssueTransaction
        end
    end
```

## Execution
The service can be executed via the following go command without having it previously built:
```
go run cmd/rddl-claim-service/main.go
```

## Configuration
The service needs to be configured via the ```./app.toml``` file or environment variables. The defaults are
```
asset = '7add40beb27df701e02ee85089c5bc0021bc813823fedb5f1dcb5debda7f3da9' //  liquid asset
confirmations = 10 // min number of confirmations on liquid
cron = '* * * * * *' // cron expression for periodic exectution
fund-address = '' // liquid address that receives funds to distribute
planetmint-rpc-host = '127.0.0.1:9090' // plmnt rpc host
r2p-host = 'planetmint-go-testnet-3.rddl.io' // r2p host
rpc-host = 'planetmint-go-testnet-3.rddl.io:18884' // elements rpc host
rpc-pass = 'password' // elements rpc password
rpc-user = 'user' // elements rpc user
shamir-host = '127.0.0.1:9091' // shamir coordinator
wallet = 'dao' // liquid wallet
```

## Data storage
This service stores every occurrence of a distribution in a LevelDB. It needs this to compare how much tokens have been received by the DAO wallet since the last occurrence.