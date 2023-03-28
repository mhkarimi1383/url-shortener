# URL Shortener

Simple and minimalism URL Shortener

The Goal is to be minimalist but scalable and powerful

> Also I want to Learn more about some libs/stacks that I used in this project ;)

## Project Status

> I Will try to keep this section of readme up to date

Connecting Login page to API

## Powered By

### Backend

- Golang
- Echo
- XORM
- Cobra

### Frontend (in [client/url-shortener](client/url-shortener))

- TypeScript
- VueJS 3
- Nuxt 3
- Vuestic UI

## Running Project

### Frontend

Just like any other NuxtJS applications

You can take a look at [client/url-shortener/README.md](client/url-shortener/README.md)

### Backend

Like any other go application (Make sure CGO is enabled and `gcc` is present) run

```bash
go mod download
```

to get dependencies

run

```bash
go run .
```

to start application

or build a binary by running

```bash
go build .
```

### Backend and Frontend toghether

You can build the project by running `./build.sh` script and run it by `./run.sh`

or simply use `./start.sh` to do both

it will start both frontend and backend, you can see frondend at the same port as backend (in `/admin` path) checkout [./api/server.go](./api/server.go) to see more info
