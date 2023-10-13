# URL Shortener

Simple and minimalism URL Shortener

## Made out of

- Golang for backend
- XORM as ORM
- Echo as HTTP Framework
- go-playground/validator as Validator
- Cobra as cli and configuration manager
- TS for frontend
- Vue3 as Frontend framework
- Ant Design Vue as UI Components
- ...

## Screenshots

<img src="./.docs/main-dark.png" width="45%"></img> <img src="./.docs/main-light.png" width="45%"></img> <img src="./.docs/url-dark.png" width="45%"></img> <img src="./.docs/url-light.png" width="45%"></img> <img src="./.docs/url-dark.png" width="45%"></img> <img src="./.docs/url-light.png" width="45%"></img>

## Running project

just run project with `--help` to see options and their env equivalent

### Locally

**Unix Users Only** use run script, it will install dependencies, format the code and run the project

### Docker

Available as `ghcr.io/mhkarimi1383/url-shortener:main`

simply run that with the options/envs as you want

## TODO

- [ ] Add more deploy/build options
- [ ] API Documentation (swaggo is buggy and generating errors with some datatypes)
- [ ] Vue3 Frontend
- [ ] Webhook and Reporting Support
