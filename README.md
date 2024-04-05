# avito-backed
## Some thoughts...
> Found out that it's possible to use codegenerators to parse `openapi.yml` file and generated handlers. But generated code was too strange and I decided to do it by myself....

> Decided to use two tables: banners and features; this allows us to have multiple features for single banner. Scheme is shown in a picture.  
![Database Scheme](./assets/db_scheme.png)
## Setup
### PostgreSQL migration
```shell
make install-migrate
make migrate
```
### Server
```shell
make run
```

## Load testing
Load testing can be performed via Apache Bench. Results can be viewied inside `tests` directory.
```shell
make load-test
```

## TODO
- [x] REST API
- [x] PostgreSQL Tables
- [ ] Implement SQL queries execution
- [ ] Implement users & admins auth via JWT
- [ ] Add cache (HashMap/Redis)
- [ ] Add functional tests
- [ ] Add load tests
- [ ] Pack app into image and make docker compose file