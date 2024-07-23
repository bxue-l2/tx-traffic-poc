# tx-traffic-poc

# how to run
```./run.sh```

# configuration
```cat .env```

Important parameters
- TRAFFIC_GENERATOR_NUM_INSTANCES  : number of concurrent wallet sending the transaction
- TRAFFIC_GENERATOR_REQUEST_INTERVAL  : how frequent each wallet is sending the transaction
- TRAFFIC_GENERATOR_PAD_SIZE : how many random bytes to be padded into a transaction

# demo
Create a chain
```anvil --block-time 3```
Create transactions
