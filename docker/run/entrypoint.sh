#!/usr/bin/env sh

set -e

if [[ -z "${SKIP_WAIT_FOR}" ]];
then
  /wait-for ${IMPACT_DATABASE_HOST}:${IMPACT_DATABASE_PORT:-5432} --timeout=10 -- echo "database is up"
fi

if [[ -n "${APP_START_TIMEOUT}" ]];
then
  sleep ${APP_START_TIMEOUT}
fi

echo
echo "Startig app"
echo

exec $@
