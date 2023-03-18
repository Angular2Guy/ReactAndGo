FROM ubuntu:jammy
VOLUME /tmp
ARG APP_FILE
ADD backend/$APP_FILE /application
ADD backend/config/properties.env /config/properies.env
ENTRYPOINT exec /application