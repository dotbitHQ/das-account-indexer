* [API List](#api-list)
    * [Indexer Api](#indexer-api)
        * [Get Server Info](#get-server-info)
            * [Request](#request)
            * [Response](#response)
            * [Usage](#usage)
        * [Get Account's Basic Info](#get-accounts-basic-info)
            * [Request](#request-1)
            * [Response](#response-1)
            * [Usage](#usage-1)
        * [Get Account Records Info](#get-account-records-info)
            * [Request](#request-2)
            * [Response](#response-2)
            * [Usage](#usage-2)
        * [Get Address Reverse Record Info](#get-address-reverse-record-info)
            * [Request](#request-3)
            * [Response](#response-3)
            * [Usage](#usage-3)
    * [Reverse Api](#reverse-api)
        * [Get Server Info](#get-server-info)
            * [Request](#request)
            * [Response](#response)
            * [Usage](#usage)
        * [Get Address Reverse Record Info](#get-address-reverse-record-info)
            * [Request](#request-3)
            * [Response](#response-3)
            * [Usage](#usage-3)
* [<em>Deprecated API List</em>](#deprecated-api-list)
    * [Get Server Info](#get-server-info)
        * [Request](#request)
        * [Response](#response)
        * [Usage](#usage)
    * [<em>Get Account's Basic Info And Records</em>](#get-accounts-basic-info-and-records)
        * [<em>Request</em>](#request-4)
        * [<em>Response</em>](#response-4)
        * [<em>Usage</em>](#usage-4)
    * [<em>Get Related Accounts By Owner Address</em>](#get-related-accounts-by-owner-address)
        * [<em>Request</em>](#request-5)
        * [<em>Response</em>](#response-5)
        * [<em>Usage</em>](#usage-5)
* [Error Code](#error-code)

## API List

### Indexer Api

#### Get Server Info

##### Request

* path: `/v1/server/info`
* param: none

##### Response

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

##### Usage

```shell
curl -X POST http://127.0.0.1:8122/v1/server/info
```

or json rpc style:

```shell
curl -X POST http://127.0.0.1:8122 -d'{"jsonrpc": "2.0","id": 1,"method": "das_serverInfo","params": []}'
```

#### Get Account's Basic Info

##### Request

* path: `/v1/account/info`
* param:

```json
{
  "account": "phone.bit"
}
```

##### Response

```json
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
      "status": 1,
      "das_lock_arg_hex": "0x0559724739940777947c56c4f2f2c9211cd5130fef0559724739940777947c56c4f2f2c9211cd5130fef",
      "owner_algorithm_id": 5,
      // 3: eth personal sign 4: tron sign 5: eip-712
      "owner_address": "0x59724739940777947c56c4f2f2c9211cd5130fef",
      "manager_algorithm_id": 5,
      "manager_address": "0x59724739940777947c56c4f2f2c9211cd5130fef"
    }
  }
}
```

```json
{
  "errno": 20007,
  "errmsg": "account not exist",
  "data": null
}
```

##### Usage

```shell
curl -X POST http://127.0.0.1:8122/v1/account/info -d'{"account":"phone.bit"}'
```

or json rpc style:

```shell
curl -X POST http://127.0.0.1:8122 -d'{"jsonrpc": "2.0","id": 1,"method": "das_accountInfo","params": [{"account":"phone.bit"}]}'
```

#### Get Account Records Info

##### Request

* path: `/v1/account/records`
* param:

```json
{
  "account": "phone.bit"
}
```

##### Response

```json
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

```json
{
  "errno": 20007,
  "errmsg": "account not exist",
  "data": null
}
```

##### Usage

```shell
curl -X POST http://127.0.0.1:8122/v1/account/records -d'{"account":"phone.bit"}'
```

or json rpc style:

```shell
curl -X POST http://127.0.0.1:8122 -d'{"jsonrpc": "2.0","id": 1,"method": "das_accountRecords","params": [{"account":"phone.bit"}]}'
```

#### Get Address Reverse Record Info

##### Request

* path: `/v1/reverse/record`
* param:

```json
{
  "type": "blockchain",
  "key_info": {
    "coin_type": "",
    // ETH-60 195-TRX 714-BNB 966-Matic
    "chain_id": "",
    // ETH-1 BSC-56 Polygon-137
    "key": ""
    // address
  }
}
```

##### Response

```json
{
  "errno": 0,
  "errmsg": "",
  "data": {
    "account": ""
  }
}
```

```json
{
  "errno": 10000,
  "errmsg": "coin_type [601] and chain_id [1] is invalid",
  "data": null
}
```

##### Usage

```shell
curl -X POST http://127.0.0.1:8122/v1/reverse/record -d'{"type": "blockchain","key_info":{"coin_type": "60","chain_id": "1","key": "0xc9f53b1d85356B60453F867610888D89a0B667Ad"}}'
```

or json rpc style:

```shell
curl -X POST http://127.0.0.1:8122 -d'{"jsonrpc": "2.0","id": 1,"method": "das_reverseRecord","params": [{"das_type":1,"address":"0xc9f53b1d85356B60453F867610888D89a0B667Ad"}]}'
```

### Reverse Api

### Deprecated API List

#### Get Account's Basic Info And Records

##### _Request_

* path: `/v1/search/account`
* param:

```json
{
  "account": "phone.bit"
}
```

##### _Response_

```json
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

```json
{
  "errno": 20007,
  "errmsg": "account not exist",
  "data": null
}
```

##### _Usage_

```shell
curl -X POST http://127.0.0.1:8121/v1/search/account -d'{"account":"phone.bit"}'
```

or json rpc style:

```shell
curl -X POST http://127.0.0.1:8121 -d'{"jsonrpc": "2.0","id": 1,"method": "das_searchAccount","params": ["phone.bit"]}'
```

#### _Get Related Accounts By Owner Address_ `Deprecated`

##### _Request_

* path: `/v1/address/account`
* param:

```json
{
  "address": "0x773BCCE3B8b41a37CE59FD95F7CBccbff2cfd2D0"
}
```

##### _Response_

```json
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

##### _Usage_

```shell
curl -X POST http://127.0.0.1:8121/v1/address/account -d'{"address":"0x773BCCE3B8b41a37CE59FD95F7CBccbff2cfd2D0"}'
```

or json rpc style:

```shell
curl -X POST http://127.0.0.1:8121 -d'{"jsonrpc": "2.0","id": 1,"method": "das_getAddressAccount","params": ["0x773BCCE3B8b41a37CE59FD95F7CBccbff2cfd2D0"]}'
```

## Error Code

```go

const (
ApiCodeSuccess        Code = 0
ApiCodeError500       Code = 500
ApiCodeParamsInvalid  Code = 10000
ApiCodeMethodNotExist Code = 10001
ApiCodeDbError        Code = 10002

ApiCodeAccountFormatInvalid Code = 20006
ApiCodeAccountNotExist      Code = 20007
)

```
    