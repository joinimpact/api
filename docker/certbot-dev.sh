#!/bin/sh

docker run -ti --rm --hostname=certbot --name=certbot \
              -v ${DATA_PATH}/joinimpact/certbot/www:/var/www/certbot:rw \
              -v ${DATA_PATH}/joinimpact/certbot/letsencrypt:/etc/letsencrypt:rw \
              -v ${DATA_PATH}/joinimpact/certbot/log:/var/log/letsencrypt:rw \
              certbot/certbot \
              certonly --webroot -w /var/www/certbot \
              --rsa-key-size 2048 --agree-tos --force-renewal \
              -d dev.joinimpact.org
