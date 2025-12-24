# SlugShrink

SlugShrink is a self-hosted URL shortener with a dashboard for managing domains and redirects, plus a public redirect server for resolving short links.

## Tech stack

- Backend: Go
- API: ConnectRPC (Connect over HTTP/1.1 + h2c)
- Frontend: SvelteKit (Svelte)

## Auth modes

SlugShrink supports several different modes for authenticating users. They can be configured in the frontend.

- Authless: No login; all requests are treated as admin.
- Basic: Password login with a JWT cookie (`auth_token`) signed by `JWT_SECRET`. Tokens expire after 1 hour.
- Proxy: Trusts upstream headers for user/group info; groups are matched for permissions. Make sure your application is only available behind the proxy. Requests that do not have the headers set will be allowed as normal user.

If you have locked yourself out (e.g. set Basic Auth but forgot password) you can remove the first record in the `auth_config` database table to reset it to authless.

## User permissions

- User: Create and view redirects; edit and delete only redirects they created.
- Superuser: Create, view, update, and delete all redirects.
- Admin: Everything a superuser can do, plus manage domains and configure auth mode.
- Proxy mode note: To let users edit/delete only their own redirects, the `user_id_header` must contain a unique user identifier. You are responsible for ensuring its correctness.

## Development helpers

The `dev/` folder contains a Traefik setup that forwards to the local frontend (`http://localhost:5173`) and injects proxy-auth headers. It exposes multiple entrypoints (ports 8091–8094) to simulate different users by setting `X-Proxy-Groups` and `X-Proxy-User`.

## Deployment

- Docker Compose: see `compose.example.yaml` for a Postgres-backed setup.

The dashboard is served on port 8080, and the server doing the actual redirects on 8081.

## Environment variables

| Variable      | Options / Default                                         | Description                                             |
| ------------- | --------------------------------------------------------- | ------------------------------------------------------- |
| `DB_TYPE`     | `sqlite` (default), `postgres`, `mysql`, `mariadb`        | Database driver selection.                              |
| `DB_PATH`     | `local.db` (default); Docker image sets `/data/sqlite.db` | SQLite database file path (used when `DB_TYPE=sqlite`). |
| `DB_HOST`     | required for `postgres`, `mysql`, `mariadb`               | Database host.                                          |
| `DB_PORT`     | required for `postgres`, `mysql`, `mariadb`               | Database port.                                          |
| `DB_USER`     | required for `postgres`, `mysql`, `mariadb`               | Database user.                                          |
| `DB_PASSWORD` | required for `postgres`, `mysql`, `mariadb`               | Database password.                                      |
| `DB_NAME`     | required for `postgres`, `mysql`, `mariadb`               | Database name.                                          |
| `DB_SSLMODE`  | `disable` (default)                                       | Postgres SSL mode (used when `DB_TYPE=postgres`).       |
| `JWT_SECRET`  | required for Basic auth                                   | Secret key used to sign login JWTs.                     |
| `LOG_LEVEL`   | `info` (default), `debug`, `warn`, `error`                | Logging verbosity.                                      |
