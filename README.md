# Saturn-Ace Demo

## Build

$ git clone https://github.com/wangfeiping/saturn.git

$ cd $GOPATH/src/github.com/wangfeiping/saturn

$ make all

Get help info

$ saturnd -h

$ saturncli -h

## Start

$ saturnd start

## Play

$ saturncli tx ace play LuckyAce draw "@random" \
    $(saturncli keys show -a test --keyring-backend=test) \
    --chain-id saturn-testnet-0 \
    --keyring-backend=test

## Test

``` sh

$ saturncli keys add test --keyring-backend=test

- name: test
  type: local
  address: cosmos1dt2mf4cjvnpthl83wcp3zvypexfa2lqya6ptfv
  pubkey: cosmospub1addwnpepqvfeshxpp0j2uz5cskyj3658tge4jxyn0g9mre9w98e59z3qqks227kjqp7
  mnemonic: ""
  threshold: 0
  pubkeys: []

**Important** write this mnemonic phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

carpet kid rice blush cement room start noise bread boring what mutual decide return edit switch electric reveal bracket cup nephew scene measure govern

$ saturncli keys show -a test --keyring-backend=test
cosmos1dt2mf4cjvnpthl83wcp3zvypexfa2lqya6ptfv

$ saturncli tx send \
    cosmos12rpact6fczpp2v9n3sxk2l5efwqjy0c77e4lm5 \
    cosmos1dt2mf4cjvnpthl83wcp3zvypexfa2lqya6ptfv \
    100000000chip --chain-id saturn-testnet-0

$ saturncli query account \
    $(saturncli keys show -a test --keyring-backend=test) \
    --chain-id saturn-testnet-0

$ saturncli query account cosmos1dt2mf4cjvnpthl83wcp3zvypexfa2lqya6ptfv
  address: cosmos1dt2mf4cjvnpthl83wcp3zvypexfa2lqya6ptfv
  coins:
  - denom: chip
    amount: "100000000"
  public_key: ""
  account_number: 9
  sequence: 0

$ saturncli tx sign tx_ace_play_.json \
      --from $(saturncli keys show -a test --keyring-backend=test) \
      --offline --chain-id saturn-testnet-0 \
      --sequence 0 --account-number 9 \
      --keyring-backend=test \
    > signed_tx_ace_play_.json

#!/bin/bash

for i in 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15
do
  echo "Running $i times"

  saturncli tx ace play LuckyAce draw "@random" \
      $(saturncli keys show -a test --keyring-backend=test) \
      --chain-id saturn-testnet-0 \
      --keyring-backend=test \
      --generate-only > tx_ace_play_$i.json

  saturncli tx sign tx_ace_play_.json \
      --from $(saturncli keys show -a test --keyring-backend=test) \
      --offline --chain-id saturn-testnet-0 \
      --sequence 1 --account-number 45 \
      --keyring-backend=test \
    > signed_tx_ace_play_$i.json

  saturncli tx broadcast signed_tx_ace_play_$i.json
done
```

## keys export&import

``` sh

saturncli keys export test

saturncli keys import test test.pri

```
