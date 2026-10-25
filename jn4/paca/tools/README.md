# Run Paca on a PC with no internet

Copy this whole `paca` folder onto the other PC (USB disk or zip). Do not use `git clone`. Git does not store `tools` (except this file), `node_modules`, or the `.env` files, and those are required.

This bundle is Windows 64-bit. It will not run on Linux.

The copied programs are about 1.3 GB, plus the `node_modules` folders already in the project.

## What is in `tools`

| Folder | What it is |
|---|---|
| `go` | Go 1.27.3, used to run the API |
| `gopath` | Go libraries for the API. Already downloaded. |
| `bun` | Bun, used to run the web app and the realtime service |
| `node` | Node.js 24.21.0. The web dev server runs on this. |
| `rustfs` | RustFS, the file store. The API will not start without it. |
| `caddy` | Caddy and the site config. The browser uses this on port 3000. |
| `database` | A dump of the `paca` database from PostgreSQL 18.6 |
| `installers` | Memurai (Redis) installer |

Two programs are not in this folder, because they are Windows services and cannot be copied as files:

- **PostgreSQL 18.** Install it on the offline PC before you start. The installer was not on the machine that built this bundle. During setup, set the `postgres` user password to `postgres`.
- **Memurai.** Install `tools\installers\Memurai-for-Redis-v8.2.0.msi`. Leave the port at **6379** and do not set a password. The free edition stops itself after 10 days; start the Memurai service again after that.

`minio.exe` is intentionally not included. That build refuses file uploads unless a license is installed.

## 1. Database

Install PostgreSQL 18, then in PowerShell:

```powershell
$Root = "C:\path\to\paca"   # change this to the copied folder
$env:PGPASSWORD = "postgres"
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -h localhost -U postgres -d postgres -c "CREATE DATABASE paca OWNER postgres;"
& "C:\Program Files\PostgreSQL\18\bin\psql.exe" -h localhost -U postgres -d paca -f "$Root\tools\database\paca.sql"
```

If `CREATE DATABASE` says the database already exists, skip that line and run only the restore.

`services\api\.env` expects this URL:

```
postgres://postgres:postgres@localhost:5432/paca?sslmode=disable
```

If the offline Postgres password is not `postgres`, change `DATABASE_URL` in that file.

## 2. Prepare each PowerShell window

The API, the web app, realtime, RustFS, and Caddy each need their own window. Paste this into every window before the start command. Change `$Root` to the copied folder.

```powershell
$Root = "C:\path\to\paca"
$env:Path = "$Root\tools\go\bin;$Root\tools\bun;$Root\tools\node;" + $env:Path
$env:GOROOT = "$Root\tools\go"
$env:GOPATH = "$Root\tools\gopath"
$env:GOMODCACHE = "$Root\tools\gopath\pkg\mod"
$env:GOPROXY = "off"
$env:GOSUMDB = "off"
$env:GOTOOLCHAIN = "local"
$env:CGO_ENABLED = "0"
```

`GOPROXY=off` stops Go from trying to download libraries. Do not run `bun install` or `go mod download`.

## 3. Start RustFS

```powershell
New-Item -ItemType Directory -Force -Path "$Root\tools\rustfs-data" | Out-Null
$env:RUSTFS_ACCESS_KEY = "rustfsadmin"
$env:RUSTFS_SECRET_KEY = "rustfsadmin"
$env:RUSTFS_VOLUMES = "$Root\tools\rustfs-data"
$env:RUSTFS_ADDRESS = "0.0.0.0:9000"
$env:RUSTFS_CONSOLE_ADDRESS = "0.0.0.0:9001"
$env:RUSTFS_CONSOLE_ENABLE = "true"
& "$Root\tools\rustfs\rustfs.exe"
```

Leave the window open. The console is http://localhost:9001 (`rustfsadmin` / `rustfsadmin`).

## 4. Start the API

```powershell
Set-Location "$Root\services\api"
go run ./cmd/api
```

Leave the window open. Wait until the log says `starting server` on port 8080. The first start compiles the API from `tools\gopath` and can take a few minutes.

## 5. Start realtime

```powershell
Set-Location "$Root\services\realtime"
bun run dev
```

Leave the window open. It listens on port 3001 and reads `services\realtime\.env`.

## 6. Start the web app

```powershell
Set-Location "$Root\apps\web"
bun run dev --host 127.0.0.1 --port 5173
```

Leave the window open.

## 7. Start Caddy

```powershell
& "$Root\tools\caddy\caddy.exe" run --config "$Root\tools\caddy\Caddyfile"
```

Leave the window open. Then open **http://localhost:3000**.

Sign in as **admin** / **adminpassword**.

Caddy sends `/api` to the API, `/ws` to realtime, `/storage` to RustFS, and the rest to the web app. Use port 3000 in the browser, not 5173.

## What stays empty without a network

- Plugin marketplace and the update banner. Both call GitHub.
- AI agents. They were not included. Their libraries were never downloaded, and the sandboxes need Docker images.
