#!/bin/bash

AUTH=https://auth.prod.monitalks.io
FAAS=https://openfaas.prod.monitalks.io/function

COGNITO_URL=${AUTH}/cognito

FAAS_SIGNUP_URL=${FAAS}/auth/signup2

ACCESS=$(curl -s -X POST -H "Cache-Control: no-cache" -H "Content-Type: application/x-www-form-urlencoded" -d $CREDENTIAL_EMAIL -d $CREDENTIAL_PASS $COGNITO_URL | jq -r '.Authorization')

echo $ACCESS

NEW_EMAIL=username=project.scapa%2B${RANDOM}%40gmail.com
NEW_PASS=password=In%21ferno666
RESIDENCY=residency=IE
BIRTHDATE=birthdate=1/1/1970

curl -H "Authorization: ${ACCESS}" -X POST -d $NEW_EMAIL -d $NEW_PASS -d $RESIDENCY -d $BIRTHDATE ${FAAS_SIGNUP_URL} 

echo
