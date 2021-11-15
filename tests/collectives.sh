#!/bin/bash

AUTH=https://auth.prod.monitalks.io
API=https://int-api.prod.monitalks.io/function
FAAS=https://openfaas.prod.monitalks.io/function

COGNITO_URL=${AUTH}/cognito

API_PROFILE_URL=${API}/profile/my
FAAS_PROFILE_URL=${FAAS}/profile/my

ACCESS=$(curl -s -X POST --data $CREDENTIALS $COGNITO_URL | jq -r '.Authorization')

echo $ACCESS

# Random Collectives
BETA=a7aaabc9-4053-4596-9e51-37a2295fb6c1
OCEAN=a7aaabc9-4053-4596-9e51-37a2295fb6a9
POST=d29efbcb-10e9-4e3d-9ee3-dddab7e0fddd

# Sample URLs
# COLL_URL=${API}/posts/collectives/$BETA/
COLL_URL=${FAAS}/posts/collectives/?collectiveId=$BETA
# POST_URL=${FAAS}/posts/$POST
# POSTS_URL=${FAAS}/posts?search=\&page=1\&limit=10

# Retrievals
curl -H "Authorization: $ACCESS" $COLL_URL | jq '.[] | {objectId,collectiveId,ownerDisplayName,body}'
# curl -H "Authorization: $ACCESS" $POST_URL | jq '{objectId,collectiveId,ownerDisplayName}'
# curl -H "Authorization: $ACCESS" $POSTS_URL | jq '.'
