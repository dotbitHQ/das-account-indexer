# Das-Account-Indexer
This repo introduces a simple server, which provided some APIs for search DAS account's records or reverse records
## Prerequisites
* Ubuntu >= 18.04
* MYSQL >= 8.0
* go version >= 1.15.0 
* Redis >= 5.0 (for cache, not necessary)

## Install

```shell
# get the code
git clone https://github.com/DeAccountSystems/das-account-indexer.git

# edit conf/config.yaml for your own convenient
vi conf/config.yaml

# init mysql database
mysql -uroot -p
> source das-account-indexer/tables/sql.sql
> quit;

# compile and run
cd das-account-indexer
make default
./das_account_indexer_server --config=conf/config.yaml
# it will take about 3 hours to synchronize to the latest data(Dec 15, 2021)
```


## API List
### Get DAS Account-Indexer Info
Shows the current status of das-account-indexer-server
#### Request
* path: `/v1/indexer/info`
* param: None
#### Response

```json
{
  "errno": 0,
  "errmsg": "",
  "data": {
    "is_latest_block_number": true,
    "current_block_number": 0
  }
}
```
#### Usage
```shell
curl -X POST http://127.0.0.1:8121/v1/indexer/info
```
or json rpc style:
```shell
curl -X POST http://127.0.0.1:8121 -d'{"jsonrpc": "2.0","id": 1,"method": "das_indexerInfo","params": []}'
```

### Get Account's Basic Info
#### Request
* path: `/v1/account/info`
* param: 
```json
{"account": "phone.bit"}
```
#### Response
```JavaScript
{
  "errno": 0,
  "errmsg": "",
  "data": {
    "out_point": {
      "tx_hash": "0xabb6b2f502e9d992d00737a260e6cde53ad3f402894b078f60a52e0392a17ec8",
      "index": 0
    },
    "account_info": {
      "account": "phone.bit",
      "account_id_hex": "0x5f560ec1edc638d7dab7c7a1ca8c3b0f6ed1848b",
      "next_account_id_hex": "0x5f5c20f6cd95388378771ca957ce665f084fe23b",
      "create_at_unix": 1626955542,
      "expired_at_unix": 1658491542,
      "status": 1, // 0: normal 1: on sale
      "das_lock_arg_hex": "0x0559724739940777947c56c4f2f2c9211cd5130fef0559724739940777947c56c4f2f2c9211cd5130fef",
      "owner_algorithm_id": 5, // 3: eth personal sign 4: tron sign 5: eip-712
      "owner_das_type": 1, // 1: evm chain 3: tron
      "owner_address": "0x59724739940777947c56c4f2f2c9211cd5130fef",
      "manager_algorithm_id": 5,
      "manager_das_type": 1,
      "manager_address": "0x59724739940777947c56c4f2f2c9211cd5130fef"
    }
  }
}
```
#### Usage
```shell
curl -X POST http://127.0.0.1:8121/v1/account/info -d'{"account":"phone.bit"}'
```
or json rpc style:
```shell
curl -X POST http://127.0.0.1:8121 -d'{"jsonrpc": "2.0","id": 1,"method": "das_accountInfo","params": [{"account":"phone.bit"}]}'
```

### Get Account Records Info
#### Request
* path: `/v1/account/records`
* param:
```json
{"account": "phone.bit"}
```
#### Response
```JavaScript
{
  "errno": 0,
  "errmsg": "",
  "data": {
    "account": "phone.bit",
    "records": [
      {
        "key": "address.btc",
        "label": "Personal account",
        "value": "3EbtqPeAZbX6wmP6idySu4jc2URT8LG2aa",
        "ttl": "300"
      },
      {
        "key": "address.eth",
        "label": "Personal account",
        "value": "0x59724739940777947c56C4f2f2C9211cd5130FEf",
        "ttl": "300"
      }
      // ...
    ]
  }
}
```
#### Usage
```shell
curl -X POST http://127.0.0.1:8121/v1/account/records -d'{"account":"phone.bit"}'
```
or json rpc style:
```shell
curl -X POST http://127.0.0.1:8121 -d'{"jsonrpc": "2.0","id": 1,"method": "das_accountRecords","params": [{"account":"phone.bit"}]}'
```

### Get Address Reverse Record Info
#### Request
* path: `/v1/reverse/record`
* param:
```JavaScript
{
    "das_type": 1, // 1: evm chain 3: tron
    "address": "0xc9f53b1d85356B60453F867610888D89a0B667Ad"
}
```
#### Response
```json
{
    "errno": 0,
    "errmsg": "",
    "data": {
        "account": ""
    }
}
```
#### Usage

```shell
curl -X POST http://127.0.0.1:8121/v1/reverse/record -d'{"chain_type":1,"address":"0xc9f53b1d85356B60453F867610888D89a0B667Ad"}'
```
or json rpc style:
```shell
curl -X POST http://127.0.0.1:8121 -d'{"jsonrpc": "2.0","id": 1,"method": "das_reverseRecord","params": [{"das_type":1,"address":"0xc9f53b1d85356B60453F867610888D89a0B667Ad"}]}'
```
## _Deprecated API List_
Deprecated APIs will be removed in the future, if you rely on these APIs, please see [das-database](https://github.com/DeAccountSystems/das-database) for more help or do some secondary developments based on this repo
### _Get Account's Basic Info And Records_
#### _Request_
* path: `/v1/search/account`
* param:
```json
{"account": "phone.bit"}
```
#### _Response_
```JavaScript
{
  "errno": 0,
  "errmsg": "",
  "data": {
    "out_point": {
      "tx_hash": "0xabb6b2f502e9d992d00737a260e6cde53ad3f402894b078f60a52e0392a17ec8",
      "index": 0
    },
    "account_data": {
      "account": "phone.bit",
      "account_id_hex": "0x5f560ec1edc638d7dab7c7a1ca8c3b0f6ed1848b",
      "next_account_id_hex": "0x5f5c20f6cd95388378771ca957ce665f084fe23b",
      "create_at_unix": 1626955542,
      "expired_at_unix": 1658491542,
      "status": 1,
      "das_lock_arg_hex": "0x0559724739940777947c56c4f2f2c9211cd5130fef0559724739940777947c56c4f2f2c9211cd5130fef",
      "owner_address_chain": "ETH",
      "owner_lock_args_hex": "0x0559724739940777947c56c4f2f2c9211cd5130fef",
      "owner_address": "0x59724739940777947c56c4f2f2c9211cd5130fef",
      "manager_address_chain": "ETH",
      "manager_address": "0x59724739940777947c56c4f2f2c9211cd5130fef",
      "manager_lock_args_hex": "0x0559724739940777947c56c4f2f2c9211cd5130fef",
      "records": [
        {
          "key": "address.btc",
          "label": "Personal account",
          "value": "3EbtqPeAZbX6wmP6idySu4jc2URT8LG2aa",
          "ttl": "300"
        },
        {
          "key": "address.eth",
          "label": "Personal account",
          "value": "0x59724739940777947c56C4f2f2C9211cd5130FEf",
          "ttl": "300"
        }
        // ...
      ]
    }
  }
}
```
#### _Usage_
```shell
curl -X POST http://127.0.0.1:8121/v1/search/account -d'{"account":"phone.bit"}'
```
or json rpc style:
```shell
curl -X POST http://127.0.0.1:8121 -d'{"jsonrpc": "2.0","id": 1,"method": "das_searchAccount","params": ["phone.bit"]}'
```
### _Get Related Accounts By Owner Address_
#### _Request_
* path: `/v1/address/account`
* param:
```json
{"address": "0x773BCCE3B8b41a37CE59FD95F7CBccbff2cfd2D0"}
```
#### _Response_
```JavaScript
{
  "errno": 0,
  "errmsg": "",
  "data": [
    {
      "out_point": {
        "tx_hash": "0xdad77b108e447f4ddd905214021594d69ef50a5b06baf84686031a0d9b45265c",
        "index": 0
      },
      "account_data": {
        "account": "werwefdsft3.bit",
        "account_id_hex": "0xb97565e427dca668f9989c6a2149d8ab3ef37a29",
        "next_account_id_hex": "0xb97577b49a2f5889627d1baa5af5129c4c1ebf9d",
        "create_at_unix": 1631618255,
        "expired_at_unix": 1664968655,
        "status": 1,
        "das_lock_arg_hex": "0x05773bcce3b8b41a37ce59fd95f7cbccbff2cfd2d005773bcce3b8b41a37ce59fd95f7cbccbff2cfd2d0",
        "owner_address_chain": "ETH",
        "owner_lock_args_hex": "0x05773bcce3b8b41a37ce59fd95f7cbccbff2cfd2d0",
        "owner_address": "0x773bcce3b8b41a37ce59fd95f7cbccbff2cfd2d0",
        "manager_address_chain": "ETH",
        "manager_address": "0x773bcce3b8b41a37ce59fd95f7cbccbff2cfd2d0",
        "manager_lock_args_hex": "0x05773bcce3b8b41a37ce59fd95f7cbccbff2cfd2d0",
        "records": [
          {
            "key": "profile.twitter",
            "label": "",
            "value": "egtfghdfhfg",
            "ttl": "300"
          },
          {
            "key": "profile.facebook",
            "label": "",
            "value": "沃尔特图和",
            "ttl": "300"
          }
        ]
      }
    }
    // ...
  ]
}
```
#### _Usage_
```shell
curl -X POST http://127.0.0.1:8121/v1/address/account -d'{"address":"0x773BCCE3B8b41a37CE59FD95F7CBccbff2cfd2D0"}'
```
or json rpc style:
```shell
curl -X POST http://127.0.0.1:8121 -d'{"jsonrpc": "2.0","id": 1,"method": "das_getAddressAccount","params": ["0x773BCCE3B8b41a37CE59FD95F7CBccbff2cfd2D0"]}'
```

## Others
* [What is DAS](https://github.com/DeAccountSystems/das-contracts/blob/master/docs/en/Overview-of-DAS.md)
* [What is a DAS transaction on CKB](https://github.com/DeAccountSystems/das-contracts/blob/master/docs/en/Data-Structure-and-Protocol/Transaction-Structure.md)
* [How to install MySQL8.0](https://github.com/DeAccountSystems/das-database/wiki/How-To-Install-MySQL-8.0)
