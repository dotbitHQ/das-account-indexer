* [Prerequisites](#prerequisites)
* [Install &amp; Run](#install--run)
    * [Source Compile](#source-compile)
    * [Docker](#docker)
* [API Usage](#api-usage)
* [Others](#others)

# Das-Account-Indexer

This repo introduces a simple server, which provided some APIs for search DAS account's records or reverse records

## Prerequisites

* Ubuntu >= 18.04
* MYSQL >= 8.0
* go version >= 1.21.3
* Redis >= 5.0 (for cache, not necessary)
* [ckb-node](https://github.com/nervosnetwork/ckb) (Must be synced to latest height and add `Indexer` module to ckb.toml)
* If the version of the dependency package is too low, please install `gcc-multilib` (apt install gcc-multilib)
* Machine configuration: 4c8g200G
## Install & Run

### Source Compile

```shell
# get the code
git clone https://github.com/dotbitHQ/das-account-indexer.git

# init config/config.yaml for your own convenient
cp config/config.example.yaml config/config.yaml

# create mysql database
mysql -uroot -p
> create database das_account_indexer;
> quit;

# compile and run
cd das-account-indexer
make default
./das_account_indexer_server --config=config/config.yaml
# it will take about 3 hours to synchronize to the latest data(Dec 15, 2021)
```

### Docker

* docker >= 20.10
* docker-compose >= 2.2.2

```bash
sudo curl -L "https://github.com/docker/compose/releases/download/v2.2.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
docker-compose up -d
```

_if you already have a mysql installed, just run_

```bash
docker run -dp 8122:8122 -v $PWD/config/config.yaml:/app/config/config.yaml --name das-indexer-server admindid/das-account-indexer:latest
```

## API Usage

[Here](https://github.com/dotbitHQ/das-account-indexer/blob/main/API.md) are the APIs details.

* If you are a newcomer, just read [API List](https://github.com/dotbitHQ/das-account-indexer/blob/main/API.md)
* If you are come from [das_account_indexer](https://github.com/dotbitHQ/das_account_indexer), you probably need do
  nothing, the new APIs are compatible with the old ones. More details
  see [deprecated-api-list](https://github.com/dotbitHQ/das-account-indexer/blob/main/API.md#deprecated-api-list), but
  we still suggest you replace with the corresponding new APIs

## Others

* [What is DAS](https://github.com/dotbitHQ/did-contracts/blob/docs/docs/en/design/Overview-of-DAS.md)
* [What is a DAS transaction on CKB](https://github.com/dotbitHQ/did-contracts/blob/docs/docs/en/developer/Transaction-Structure.md)
* [How to install MySQL8.0](https://github.com/dotbitHQ/das-database/wiki/How-To-Install-MySQL-8.0)
