### Run locally

```shell
docker-compose up -d psql
time sleep 10
docker-compose up api
```

### Useful

- `docker-compose up tests`

> Using the same database for development and testing, because don't want to run another container, small project.

- Swagger available [here](./docs/swagger.yaml).

