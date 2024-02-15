# remix

The **remix** service allows us to register short video clips captured from tapes, then
later queue those clips to be played on demand.

- **OpenAPI specification:** https://golden-vcr.github.io/remix/

## Development Guide

On a Linux or WSL system:

1. Install [Go 1.21](https://go.dev/doc/install)
2. Clone the [**terraform**](https://github.com/golden-vcr/terraform) repo alongside
   this one, and from the root of that repo:
    - Ensure that the module is initialized (via `terraform init`)
    - Ensure that valid terraform state is present
    - Run `terraform output -raw env_remix_local > ../remix/.env` to populate an `.env`
      file.
    - Run `./local-db.sh up` to ensure that a local Postgres server is running in
      Docker.
3. Ensure that the [**auth**](https://github.com/golden-vcr/auth?tab=readme-ov-file#development-guide)
   server is running locally.
4. From the root of this repository:
    - Run [`./db-migrate.sh`](./db-migrate.sh) to ensure that migrations have been
      applied to the local database.
    - Run [`go run ./cmd/server`](./cmd/server/main.go) to start up the server.

Once done, the remix server will be running at http://localhost:5010.
