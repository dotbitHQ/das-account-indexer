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
      "owner_address": "0x59724739940777947c56c4f2f2c9211cd5130fef",
      "manager_algorithm_id": 5,
      "manager_address": "0x59724739940777947c56c4f2f2c9211cd5130fef"
    }
  }
}
```

#### Usage

```shell
curl -X POST http://127.0.0.1:8122/v1/account/info -d'{"account":"phone.bit"}'
```

or json rpc style:

```shell
curl -X POST http://127.0.0.1:8122 -d'{"jsonrpc": "2.0","id": 1,"method": "das_accountInfo","params": [{"account":"phone.bit"}]}'
```

### Reverse Api