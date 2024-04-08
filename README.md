# [Avito Backend Trainee Project](https://github.com/avito-tech/backend-trainee-assignment-2024?tab=readme-ov-file)
## Some thoughts...
> Found out that it's possible to use codegenerators to parse `openapi.yml` file and generated handlers. But generated code was too strange and I decided to do it by myself....

> Decided to use two tables: banners and features; this allows us to have multiple features for single banner. Scheme is shown in a picture.  

<img src="./assets/db_scheme.png" width=50%>

> How to get 403 status for /user_banner...? Under these conditions probably when token is incorrect.

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
- [ ] Implement SQL queries execution for:
    - [x] GET /user_banner
    - [x] GET /banner
    - [ ] POST /banner
    - [ ] PATCH /banner/{id}
    - [ ] DELETE /banner/{id}
- [ ] Implement users & admins auth via JWT
    - [ ] Users auth
    - [ ] Admins auth
- [ ] Add cache (HashMap/Redis)
- [ ] Add some data to be loaded into db for load tests and optimization
- [ ] Add load tests
- [ ] Add functional tests
- [x] Pack app into image and make docker compose file