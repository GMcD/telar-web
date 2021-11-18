#!/bin/bash

AUTH=https://auth.prod.monitalks.io
FAAS=https://openfaas.prod.monitalks.io/function

COGNITO_URL=${AUTH}/cognito

FAAS_PROFILE_URL=${FAAS}/profile/my

# curl -s -X POST -H "Cache-Control: no-cache" -H "Content-Type: application/x-www-form-urlencoded" -d $CREDENTIAL_EMAIL -d $CREDENTIAL_PASS $COGNITO_URL

ACCESS=$(curl -s -X POST -H "Cache-Control: no-cache" -H "Content-Type: application/x-www-form-urlencoded" -d $CREDENTIAL_EMAIL -d $CREDENTIAL_PASS $COGNITO_URL | jq -r '.Authorization')

echo $ACCESS

curl -H "Authorization: ${ACCESS}" ${FAAS_PROFILE_URL}

echo
