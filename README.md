# Payment gateway handler service
### Overview
Payment gateway handler service (IPG,MPG)  

**All the things PG does :**
1. Inits MYSQL connection
2. Reads all psp gateway configs
3. Inits mux router with path prefix
4. Inits request route
5. Inits psp confirm route
6. Inits psp form redirect route (form-acceptable psp gates)
7. Starts listening with router

### Ports and configs
see `pg/cmd/configs`

- Make commands:
    * `make deps` get dependencies
    * `make build` build locally
    * `make test` run tests
    * `make clean` clean binaries
    * `make docker-build` build and run docker
    * `make db-schema` generate db schema# Krakend-ModifiedErrorResponse
# pg-master
