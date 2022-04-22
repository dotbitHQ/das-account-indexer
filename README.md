 * [Prerequisites](#prerequisites)
 * [Install](#install)
 * [API Usage](#usage)
 * [Others](#others)
    

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

## Docker Install
```bash
# if you already have a mysql database installed, just run
docker run -dp 8121-8123:8121-8123 -v $PWD/config/config.yaml:/app/config/config.yaml --name bit-indexer-server slagga/bit-indexer

# if not, you need docker-compose to automate the installation
docker-compose up -d
```

## API Usage
[Here](https://github.com/DeAccountSystems/das-account-indexer/blob/main/API.md) are the APIs details.

* If you are a newcomer, just read [API List](https://github.com/DeAccountSystems/das-account-indexer/blob/main/API.md) 
* If you are come from [das_account_indexer](https://github.com/DeAccountSystems/das_account_indexer), you probably need do nothing, the new APIs are compatible with the old ones. More details see [deprecated-api-list](https://github.com/DeAccountSystems/das-account-indexer/blob/main/API.md#deprecated-api-list), but we still suggest you replace with the corresponding new APIs




## Others
* [What is DAS](https://github.com/DeAccountSystems/das-contracts/blob/master/docs/en/Overview-of-DAS.md)
* [What is a DAS transaction on CKB](https://github.com/DeAccountSystems/das-contracts/blob/master/docs/en/Data-Structure-and-Protocol/Transaction-Structure.md)
* [How to install MySQL8.0](https://github.com/DeAccountSystems/das-database/wiki/How-To-Install-MySQL-8.0)
