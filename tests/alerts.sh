#!/bin/bash

JWT=eyJraWQiOiI3WHV5NTdsc0I3TnJ2MUo4TkM5T0FlaEhhQ1pBVjJGRVwveCs2YzhlXC9CYXM9IiwiYWxnIjoiUlMyNTYifQ.eyJzdWIiOiI1ZDM1ZjNhMy01MmVjLTQxN2QtYTAzMy0yNzEzMGQ0ZTNjZmMiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiaXNzIjoiaHR0cHM6XC9cL2NvZ25pdG8taWRwLmV1LXdlc3QtMi5hbWF6b25hd3MuY29tXC9ldS13ZXN0LTJfWldtZko2dkpzIiwiY29nbml0bzp1c2VybmFtZSI6IjVkMzVmM2EzLTUyZWMtNDE3ZC1hMDMzLTI3MTMwZDRlM2NmYyIsIm9yaWdpbl9qdGkiOiJkZWQ0NmQ2Zi1kZDgxLTQ5ZjAtYTQ5Ni0zYjVhYWU2MTgyMDEiLCJhdWQiOiJnY2xnN2k1ajQxY3VjYWl2a3RxMjc2bTE1IiwiZXZlbnRfaWQiOiIzMjQyOTljOC0xZWVmLTQwYzYtYmRlMy02MjAxYmU4NjhiZTIiLCJ0b2tlbl91c2UiOiJpZCIsImF1dGhfdGltZSI6MTYzMzg1MTMyMiwibmFtZSI6IkdhcnkgTWFjRG9uYWxkIiwiZXhwIjoxNjMzOTE4MDM2LCJpYXQiOjE2MzM5MTA4MzYsImp0aSI6ImUzZDZkZmIzLWM5YjAtNGUxOS1iMDI2LTk0ODM3YTE0NWNiNCIsImVtYWlsIjoicHJvamVjdC5zY2FwYUBnbWFpbC5jb20ifQ.CbA_CMZgPKdlEZPgFEPfgZEEUaqNEL-pirB4yoNOM5odOA0k8q6ooSS5htnwVFIOOGu3pjt5K4-6jDQbKi-WVfqKYpEbd6p4eFyzb2zhZn3zjYvAyAwFYVDWSU489uuytRxoDrxM6SDpzkcnJVFltUPVqZGIkxcTNtPq6td_GcQmg73AUehu9Uswg9FPotXYcDvZwtMqwGV4AhH72LT9EVCPBbpxLJHzkpFbSz-2JMB780CKhntDJCluN9Dw-QEo_7yjZWu1clUo01bg2ygLhpdHFh6kBb71s_DKMQGZEFIj3NOP6978Md8laneKVsNM3yRmvojQuIzqg024IB1NGw

ALERT_URL=https://openfaas.prod.monitalks.io/function/notifications?page=1

curl -H "Authorization: $JWT" $ALERT_URL | jq '.[] | {ownerDisplayName,ownerAvatar}'
