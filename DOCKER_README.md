

## Preparing
* copy `.env.example` as `.env` and set variables which will be used by docker-compose, e.g. `DATA_PATH`
* if you need to change some container run options e.g. host bind address - just copy `docker-compose.override.yml.example` as `docker-compose.override.yml` and add custom options into this one
* if you need to use custom application variables for debugging or something - just copy `envvars.conf` as `envvars.conf.something` and set `APP_ENV_FILE=./envvars.conf.something` in `.env` file

## Build and start
* checkout `api` and `impact-frontend` repositories into same subdirectory and run:

```shell
~$ docker-compose build
~$ docker-compose up -d
```
