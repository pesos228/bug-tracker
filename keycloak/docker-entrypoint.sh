#!/bin/sh
set -e

TEMPLATE_FILE="/tmp/realm-template.json"
OUTPUT_FILE="/opt/keycloak/data/import/myrealm-realm.json"

echo "Generating Keycloak realm configuration..."
mkdir -p /opt/keycloak/data/import

envsubst < ${TEMPLATE_FILE} > ${OUTPUT_FILE}

echo "Realm config generated at ${OUTPUT_FILE}"
echo "Starting Keycloak..."

exec /opt/keycloak/bin/kc.sh "$@"