#!/bin/sh
set -e

if  [[ -z "$ ADMIN_PASS" ]]; then
    echo " ADMIN_PASS not defined"
    exit 1
fi


if [[ -z "$TELEGRAM_TOKEN" ]]; then
    echo "TELEGRAM_TOKEN not defined"
    exit 1
else
    if [[ -z "$TELEGRAM_CHAT_ID" ]]; then
        echo "TELEGRAM_CHAT_ID not defined"
        exit 1
    else
        ./app
    fi
fi




