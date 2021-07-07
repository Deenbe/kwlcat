# gopi

> Set of core kwlcat services written in go

All services in this solution are executed via same cli.

## Build

```shell
make install
```

## Run

### Launch dependencies

```shell
./launch_collector # runs otel daemon and a local jager instance
```

### idgen 
API to generate an always increasing integer. IDs are stored in DynamoDB.

```shell
gopi idgen --port 8080 --idgen-table-name <dynamodb table name>
```

### namegen
API to generate a random name. Depends on `idgen` api to generate a unique id for generate name.

```shell
gopi namegen --port 8081 --idgen-api-base-url http://localhost:8080 
```

## Try

### Generate a name

```shell
curl localhost:8081/names/next
```

### See traces generated

```shell
open http://localhost:16686 
```

