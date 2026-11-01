FROM ubuntu:26.04
VOLUME /tmp
ENV GIN_MODE=release
ARG APP_FILE
ADD backend/$APP_FILE /application
ADD backend/config/properties.env /config/properies.env
ENTRYPOINT exec /application