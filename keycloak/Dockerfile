FROM registry.access.redhat.com/ubi9/ubi:latest AS builder

RUN dnf install -y gettext --setopt=tsflags=nodocs && dnf clean all

FROM quay.io/keycloak/keycloak:26.3.0

USER root

COPY --from=builder /usr/bin/envsubst /usr/bin/

COPY ./myrealm-realm-template.json /tmp/realm-template.json
COPY ./docker-entrypoint.sh /opt/keycloak/bin/docker-entrypoint.sh

USER keycloak

ENTRYPOINT ["/opt/keycloak/bin/docker-entrypoint.sh"]

CMD ["start", "--import-realm"]