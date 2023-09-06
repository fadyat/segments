## Segments

Backend service for dynamic user segmentation with core functionality:

- **Dynamic Segment Management**
  > provides the ability to create, delete, and manage user segments, enabling flexible experimentation and targeted
  user grouping.
- **Real-time History Tracking**
  > keeps a comprehensive record of user segment membership changes, allowing users to track precisely when users are
  > added or removed from segments.
- **Time-Based Membership Management**
  > can specify a time duration for which a user should remain in a segment, ensuring seamless control over segment
  > membership based on defined time frames.
- **Automated User Distribution**
  > automatic distribution feature, enabling the assignment of users to segments based
  > on predefined percentages, streamlining user segmentation and testing processes.

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

- Created [`Transactor`](./internal/repository/transactor.go) interface for storing/getting transaction
  from context

> Allows you to hide the implementation of the transaction and use it in the repository layer
> without passing the transaction as an argument to each method.