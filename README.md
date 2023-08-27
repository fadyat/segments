### Run locally

```shell
docker-compose up -d psql
time sleep 10
docker-compose up api
```

### Useful

- `docker-compose up tests`

> Using the same database for development and testing, because don't want to run another container, small project.
>
> Testcases are done for all database queries, but not for the API endpoints, because i'm lazy.

- Swagger available [here](./docs/swagger.yaml).

### Decisions

- Migrations are [here](./migrations/postgres)

> Users are inserted using a raw query because it is not clear from the task how to work with this
