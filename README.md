# avito-backed
## Some thoughts...
> Found out that it's possible to use codegenerators to parse `openapi.yml` file and generated handlers. But generated code was too strange and I decided to do it by myself....

> Decided to use two tables: banners and features; this allows us to have multiple features for single banner. Scheme is shown in a picture.
![Database Scheme](./assets/db_scheme.png)
## Setup
### PostgreSQL migration
```
make install-migrate
make migrate
```
### Server
```
