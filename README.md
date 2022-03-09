# vmaas-go
Go reimplementation of original [vmaas app](https://github.com/RedHatInsights/vmaas)

## Run app
~~~bash
docker-compose up --build -d
./dev/scripts/post_updates.sh # post testing /updates request
PGPASSWORD=passwd psql -d testdb -h localhost -U admin -p 5433 # connect to the database
~~~

## Run tests
~~~
docker-compose -f docker-compose.test.yml up --build test
~~~
