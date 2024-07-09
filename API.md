* [API List](#api-list)
    * [Get Server Info](#get-server-info)
    * [Get Account Info](#get-account-info)
    * [Get Account List](#get-account-list)
    * [Get Account Records Info](#get-account-records-info)
    * [Get Valid Reverse Addresses](#Get-Valid-Reverse-Addresses)
    * [Get Reverse Record Info](#get-reverse-record-info)
    * [Get Sub-Account List](#get-sub-account-list)
    * [Verify Sub-Account](#verify-sub-account)
    * [Get Batch Account Records Info](#get-batch-account-records-info)
    * [Get Batch Reverse Record Info](#Get-Batch-Reverse-Record-Info)
    * [Get Batch register Info](#get-batch-register-info)
    * [Get Did Cell List](#get-did-cell-list)
    * [Get Account Records Info V2](#get-account-records-info-v2)

* [<em>Deprecated API List</em>](#deprecated-api-list)
    * [<em>Get Account Basic Info And Records</em>](#get-account-basic-info-and-records-deprecated)
    * [<em>Get Related Accounts By Owner Address</em>](#get-related-accounts-by-owner-address-deprecated)
* [Error](#error)
    * [Error Example](#error-example)
    * [Error Code](#error-code)

  
## API List

Please familiarize yourself with the meaning of some common parameters before reading the API list:

| param                                                                        | description                                        |
|:-----------------------------------------------------------------------------|:---------------------------------------------------|
| type                                                                         | Filled with "blockchain" for now                   |
| [coin_type](https://github.com/satoshilabs/slips/blob/master/slip-0044.md)   | 60: eth, 195: trx, 9006: bsc, 966: matic, 3: doge  |
| account                                                                      | Contains the suffix `.bit` in it                   |
| key                                                                          | Generally refers to the blockchain address for now |


#### Full Functional Indexer

```json
https://indexer-v1.did.id
```

This service can query all data, but it is recommended that developers setup their own Indexer for the sake of decentralization.


### Get Server Info

**Request**
* host: `indexer-v1.did.id`
* path: `/v1/server/info`
* param: none

**Response**

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "is_latest_block_number": true,
    "current_block_number": 6088191,
    "chain": "testnet" 
  }
}
```

**Usage**

```shell
curl -X POST https://indexer-v1.did.id/v1/server/info
```

or json rpc style:

```shell
curl -X POST https://indexer-v1.did.id -d'{"jsonrpc": "2.0","id": 1,"method": "das_serverInfo","params": []}'
```

### Get Account Info

**Request**
* host: `indexer-v1.did.id`
* path: `/v1/account/info`
* param: none

**Response**

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "account": "",
    "account_id": ""
  }
}
```

**Usage**

```shell
curl -X POST https://indexer-v1.did.id/v1/account/info -d'{"account":"","account_id":""}'
```


### Get Reverse Record Info
* You need to set an alias for it to take effect.
* [How to set an alias](https://app.did.id/alias)

**Request**
* host: `indexer-v1.did.id`
* path: `/v1/reverse/record`
* param:

```json
{
  "type": "blockchain",
  "key_info": {
    "coin_type": "",
    "key": ""
  }
}
```

**Response**

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "account": "",
    "account_alias": "",
    "display_name": ""
  }
}
```


**Usage**

```shell
curl -X POST https://indexer-v1.did.id/v1/reverse/record -d'{"type": "blockchain","key_info":{"coin_type": "60","key": "0x9176acd39a3a9ae99dcb3922757f8af4f94cdf3c"}}'
```

or json rpc style:

```shell
curl -X POST https://indexer-v1.did.id -d'{"jsonrpc": "2.0","id": 1,"method": "das_reverseRecord","params": [{"type": "blockchain","key_info":{"coin_type": "60","key": "0x9176acd39a3a9ae99dcb3922757f8af4f94cdf3c"}}]}'
```

### Get Batch Reverse Record Info
* You need to set an alias for it to take effect.
* [How to set an alias](https://app.did.id/alias)

**Request**
* host: `indexer-v1.did.id`
* path: `/v1/batch/reverse/record`
* param:
  * support up to 100 addresses
```json
{
  "batch_key_info":[
    {
      "type": "blockchain",
      "key_info": {
        "coin_type": "", // 60: ETH, 195: TRX, 9006: BNB, 966: Matic, 3: doge
        "key": "" // address
      }
    }
  ]
}
```

**Response**

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "list": [
      {
        "account": "",
        "account_alias": "",
        "display_name": "",
        "err_msg": ""
      }
      //...
    ]
  }
}
```

**Usage**

```shell
curl -X POST https://indexer-v1.did.id/v1/batch/reverse/record -d'{"batch_key_info":[{"type": "blockchain","key_info":{"coin_type": "60","key": "0x9176acd39a3a9ae99dcb3922757f8af4f94cdf3c"}}]}'
```

or json rpc style:

```shell
curl -X POST https://indexer-v1.did.id -d'{"jsonrpc": "2.0","id": 1,"method": "das_batchReverseRecord","params": [{"batch_key_info":[{"type": "blockchain","key_info":{"coin_type": "60","key": "0x9176acd39a3a9ae99dcb3922757f8af4f94cdf3c"}}]}]}'
```

### Get Batch register Info

batch get account register info, currently can only check whether the account can be registered

**Request**
* host: `indexer-v1.did.id`
* path: `/v1/batch/register/info`
* param:
  * support up to max 50 account
```json
{
  "batch_account": [
    "xxxx",
    "test1",
    "20230906"
  ]
}
```

**Response**

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "list": [
      {
        "account": "xxxx",
        "can_register": false
      },
      {
        "account": "test1",
        "can_register": true
      },
      {
        "account": "20230906",
        "can_register": false
      }
    ]
  }
}
```

**Usage**

```shell
curl -X POST https://indexer-v1.did.id/v1/batch/register/info -d '{"batch_account": ["xxxxx", "test1.bit", "20230906.bit"]}'
```

or json rpc style:

```shell
curl -X POST https://indexer-v1.did.id -d '{"jsonrpc": "2.0","id": 1,"method": "das_batchRegisterInfo","params": [{"batch_account": ["xxxxx", "test1.bit", "20230906.bit"]}]}'
```

### Get Did Cell List

batch get .bit DOBs 

**Request**
* host: `indexer-v1.did.id`
* path: `/v1/did/list`
* param:
  * did_type: 
```json
{
  "type": "blockchain",
  "key_info": {
    "coin_type": "", 
    "key": ""
  },
  "page": 1,
  "size": 10,
  "did_type": 1
}
```

**Response**

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "list": [
      {
        "outpoint": "",
        "account_id": "",
        "account": "",
        "args": "",
        "expired_at": 0,
        "did_cell_status": 1
      }
    ]
  }
}
```

**Usage**

```shell
curl -X POST https://indexer-v1.did.id/v1/did/list -d '{"type": "blockchain","key_info": {"coin_type": "", "key": ""},"page": 1,"size": 10,"did_type": 1}'
```

### Get Record List

batch get account register info, currently can only check whether the account can be registered

**Request**
* host: `indexer-v1.did.id`
* path: `/v1/record/list`
* param:
  * support up to max 50 account
```json
{
  "account": ""
}
```

**Response**

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "account": "",
    "records": [
      {
        "key":"",
        "label": "",
        "value": "",
        "ttl": ""
      }
    ]
  }
}
```

**Usage**

```shell
curl -X POST https://indexer-v1.did.id/v1/record/list -d '{"account": ""}'
```



### Verify Sub-Account

**Request**

* host: `indexer-v1.did.id`
* path: `/v1/sub/account/verify`
* param:
  * account: main account (Choose one with sub_account)
  * sub_account: sub account (Choose one with account)
  * address: address
  * verify_type: 0 (verify owner, default), 1 (verify manager)

For account and sub_account There are two verification methods:
1. Verify if you are the owner or manager of a arbitrary sub account of a main account
```json

{
  "account": "phone.bit",
  "address": "0xb77067fd217a8215953380bcb1cae0a1be2def31",
  "verify_type": 0
}
```
2. Verify if you are the owner or manager of a specific sub account
```json
{
  "sub_account":"test01.phone.bit",
  "address": "0xb77067fd217a8215953380bcb1cae0a1be2def31",
  "verify_type": 0
}
```
**Response**
```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "is_subdid": false
  }
}
```
```shell
curl --location 'https://indexer-v1.did.id/v1/sub/account/verify' \
--header 'Content-Type: application/json' \
--data '{
    "account":"phone.bit",
    "sub_account":"a.phone.bit",
    "address":"0xb77067fd217a8215953380bcb1cae0a1be2def31",
    "verify_type":0

}'
```
or json rpc style:

```shell
curl -X POST https://indexer-v1.did.id -d'{"jsonrpc": "2.0","id": 1,"method": "das_subAccountVerify","params": [{"account":"phone.bit","sub_account":"a.phone.bit","address":"0xb77067fd217a8215953380bcb1cae0a1be2def31","verify_type":0}]}'
```
### Get Account Basic Info

**Request**

* host: `indexer-v1.did.id`
* path: `/v1/account/info`
* param:
  * You can provide either `account` or `account_id`. The `account_id` will be used, if you provide both.

```json
{
  "account": "phone.bit",
  "account_id": ""
}
```

**Response**
  * status: 0-normal, 1-on sale, 3-cross-chain, 4-approval-enable

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "out_point": {
      "tx_hash": "0xabb6b2f502e9d992d00737a260e6cde53ad3f402894b078f60a52e0392a17ec8",
      "index": 0
    },
    "account_info": {
      "account": "phone.bit",
      "account_alias":"",        
      "account_id_hex": "0x5f560ec1edc638d7dab7c7a1ca8c3b0f6ed1848b",
      "next_account_id_hex": "0x5f5c20f6cd95388378771ca957ce665f084fe23b",
      "create_at_unix": 1626955542,
      "expired_at_unix": 1658491542,
      "status": 1,
      "das_lock_arg_hex": "0x0559724739940777947c56c4f2f2c9211cd5130fef0559724739940777947c56c4f2f2c9211cd5130fef",
      "owner_algorithm_id": 5, // 3: eth personal sign, 4: tron sign, 5: eip-712
      "owner_key": "0x59724739940777947c56c4f2f2c9211cd5130fef",
      "manager_algorithm_id": 5,
      "manager_key": "0x59724739940777947c56c4f2f2c9211cd5130fef", 
      "enable_sub_account": 0, // 0-disable 1-enable
      "display_name":""
    }
  }
}
```

**Usage**

```shell
curl -X POST https://indexer-v1.did.id/v1/account/info -d'{"account":"phone.bit"}'
```

or json rpc style:

```shell
curl -X POST https://indexer-v1.did.id -d'{"jsonrpc": "2.0","id": 1,"method": "das_accountInfo","params": [{"account":"phone.bit"}]}'
```

### Get Account List

**Request**

* host: `indexer-v1.did.id`
* path: `/v1/account/list`
* param:
  
```json
{
  "type": "blockchain",
  "key_info": {
    "coin_type": "", // 60: ETH, 195: TRX, 9006: BNB, 966: Matic, 3: doge
    "key": "" // address
  },
  "role": "owner" // owner,manager
}
```

**Response**

```json
{
  "err_no":0,
  "err_msg":"",
  "data":{
    "account_list": [
      {
        "account":"",
        "account_alias":"",
        "display_name":"",
        "registered_at": 1666268687,
        "expired_at": 1729340687
      }
    ],
    "total":1
  }
}
```

**Usage**

```shell
curl -X POST https://indexer-v1.did.id/v1/account/list -d'{"type": "blockchain","key_info":{"coin_type": "60","key": "0x3a6cab3323833f53754db4202f5741756c436ede"}}'
```

or json rpc style:

```shell
curl -X POST https://indexer-v1.did.id -d'{"jsonrpc": "2.0","id": 1,"method": "das_accountList","params": [{"type": "blockchain","key_info":{"coin_type": "60","key": "0x3a6cab3323833f53754db4202f5741756c436ede"}}]}'
```

### Get Account Records Info

> This is deprecated, please use [Records Info V2 below](https://github.com/dotbitHQ/das-account-indexer/blob/main/API.md#get-account-records-info-v2) instead.

**Request**

* host: `http://127.0.0.1:8122`
* path: `/v1/account/records`
* param:

```json
{
  "account": "phone.bit"
}
```

**Response**
* key: https://github.com/dotbitHQ/cell-data-generator/blob/master/data/record_key_namespace.txt

```json
{
  "err_no": 0,
  "err_msg": "",
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


**Usage**

```shell
curl -X POST https://indexer-v1.did.id/v1/account/records -d'{"account":"phone.bit"}'
```

or json rpc style:

```shell
curl -X POST https://indexer-v1.did.id -d'{"jsonrpc": "2.0","id": 1,"method": "das_accountRecords","params": [{"account":"phone.bit"}]}'
```

### Get Account Records Info V2

The return field [key] from [SLIP-0044](https://github.com/satoshilabs/slips/blob/master/slip-0044.md).

**Request**

* host: `http://127.0.0.1:8122`
* path: `/v2/account/records`
* param:

```json
{
  "account": "phone.bit"
}
```

**Response**
* key: https://github.com/satoshilabs/slips/blob/master/slip-0044.md

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "account": "phone.bit",
    "records": [
      {
        "key": "address.0",
        "label": "Personal account",
        "value": "3EbtqPeAZbX6wmP6idySu4jc2URT8LG2aa",
        "ttl": "300"
      },
      {
        "key": "address.60",
        "label": "Personal account",
        "value": "0x59724739940777947c56C4f2f2C9211cd5130FEf",
        "ttl": "300"
      }
      // ...
    ]
  }
}
```


**Usage**

```shell
curl -X POST https://indexer-v1.did.id/v2/account/records -d'{"account":"phone.bit"}'
```

or json rpc style:

```shell
curl -X POST https://indexer-v1.did.id -d'{"jsonrpc": "2.0","id": 1,"method": "das_accountRecordsV2","params": [{"account":"phone.bit"}]}'
```

### Get Batch Account Records Info

The return field [key] from [SLIP-0044](https://github.com/satoshilabs/slips/blob/master/slip-0044.md).

**Request**

* host: `http://127.0.0.1:8122`
* path: `/v1/batch/account/records`
* param:
  * Support up to 100 accounts

```json
{
  "accounts": ["phone.bit","..."] 
}
```

**Response**
* key: https://github.com/dotbitHQ/cell-data-generator/blob/master/data/record_key_namespace.txt

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
      "list":[{
        "account": "phone.bit",
        "account_id": "",
        "records": [
          {
            "key": "address.eth",
            "label": "Personal account",
            "value": "0x59724739940777947c56C4f2f2C9211cd5130FEf",
            "ttl": "300"
          }
      ]
    }
    ]
  }
}
```


**Usage**

```shell
curl -X POST https://indexer-v1.did.id/v1/batch/account/records -d'{"accounts":["phone.bit"]}'
```

or json rpc style:

```shell
curl -X POST https://indexer-v1.did.id -d'{"jsonrpc": "2.0","id": 1,"method": "das_batchAccountRecords","params": [{"accounts":["phone.bit"]}]}'
```

### Get Valid Reverse Addresses
**Request**

* host: `http://127.0.0.1:8122`
* path: `/v1/account/reverse/address`
* param:

```json
{
  "account": "20230725.bit"
}
```

**Response**
  * coin_type: 60-evm, 195-tron, 3-doge, 309-ckb

```json
{
  "err_no":0, 
  "err_msg":"",
  "data":{
  "list":[
    {
      "type":"blockchain",
      "key_info":{
        "coin_type":"60",
        "key":"0x15a33588908cf8edb27d1abe3852bf287abd3891"
      }
    },
    {
      "type":"blockchain",
      "key_info":{
        "coin_type":"195",
        "key":"TQoLh9evwUmZKxpD1uhFttsZk3EBs8BksV"
      }
    }
  ]
  }
}
```

**Usage**

```shell
curl -X POST https://indexer-v1.did.id/v1/account/reverse/address -d'{"account":"20230725.bit"}'
```

or json rpc style:

```shell
curl -X POST https://indexer-v1.did.id -d'{"jsonrpc": "2.0","id": 1,"method": "das_accountReverseAddress","params": [{"account":"20230725.bit"}]}'
```



### Get Sub-Account List

**Request**

* host: `http://127.0.0.1:8122`
* path: `/v1/sub/account/list`
* param:

```json
{
  "account": "0x.bit",
  "page": 1,
  "size": 20
}
```

**Response**

* enable_sub_account: 0-unenabled, 1-enabled

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "account": "0x.bit",
    "account_id_hex": "0x35612d221d6c02564c36935f81ec8568b07a39f3",
    "enable_sub_account": 1,
    "sub_account_total": 300,
    "sub_account_list": [
      {
        "account": "1234.0x.bit",
        "account_id_hex": "0x673cfe5216652c3b401d904d776dedf4a2ce9f41",
        "create_at_unix": 0,
        "expired_at_unix": 0,
        "owner_algorithm_id": 5,// 3: eth personal sign, 4: tron sign, 5: eip-712
        "owner_sub_aid": 0,
        "owner_key": "0x...",
        "manager_algorithm_id": 5,
        "manager_sub_aid": 0,
        "manager_key": "0x...",
        "display_name":""
      }
    ]
  }
}
```

**Usage**

```shell
curl -X POST https://indexer-v1.did.id/v1/sub/account/list -d'{"account":"0x.bit","page":1,"size":20}'
```

or json rpc style:

```shell
curl -X POST https://indexer-v1.did.id -d'{"jsonrpc": "2.0","id": 1,"method": "das_subAccountList","params": [{"account":"0x.bit","page":1,"size":20}]}'
```

## _Deprecated API List_

### _Get Account Basic Info And Records_ `Deprecated`

 _**Request**_

* path: `/v1/search/account`
* param:

```json
{
  "account": "phone.bit"
}
```

 _**Response**_

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": {
    "out_point": {
      "tx_hash": "0xabb6b2f502e9d992d00737a260e6cde53ad3f402894b078f60a52e0392a17ec8",
      "index": 0
    },
    "account_data": {
      "account": "phone.bit",
      "account_alias":"",        
      "display_name":"",
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

 _**Usage**_

```shell
curl -X POST https://indexer-v1.did.id/v1/search/account -d'{"account":"phone.bit"}'
```

or json rpc style:

```shell
curl -X POST https://indexer-v1.did.id -d'{"jsonrpc": "2.0","id": 1,"method": "das_searchAccount","params": ["phone.bit"]}'
```

### _Get Related Accounts By Owner Address_ `Deprecated`

 _**Request**_

* path: `/v1/address/account`
* param:

```json
{
  "address": "0x773BCCE3B8b41a37CE59FD95F7CBccbff2cfd2D0"
}
```

 _**Response**_

```json
{
  "err_no": 0,
  "err_msg": "",
  "data": [
    {
      "out_point": {
        "tx_hash": "0xdad77b108e447f4ddd905214021594d69ef50a5b06baf84686031a0d9b45265c",
        "index": 0
      },
      "account_data": {
        "account": "werwefdsft3.bit",
        "account_alias":"",
        "display_name":"",
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

 _**Usage**_

```shell
curl -X POST https://indexer-v1.did.id/v1/address/account -d'{"address":"0x773BCCE3B8b41a37CE59FD95F7CBccbff2cfd2D0"}'
```

or json rpc style:

```shell
curl -X POST https://indexer-v1.did.id -d'{"jsonrpc": "2.0","id": 1,"method": "das_getAddressAccount","params": ["0x773BCCE3B8b41a37CE59FD95F7CBccbff2cfd2D0"]}'
```

## Error
### Error Example
```json
{
  "err_no": 20007,
  "err_msg": "account not exist",
  "data": null
}
```
### Error Code
```go

const (
  ApiCodeSuccess              Code = 0
  ApiCodeError500             Code = 500
  ApiCodeParamsInvalid        Code = 10000
  ApiCodeMethodNotExist       Code = 10001
  ApiCodeDbError              Code = 10002
  
  ApiCodeAccountFormatInvalid Code = 20006
  ApiCodeAccountNotExist      Code = 20007
)

```
    
