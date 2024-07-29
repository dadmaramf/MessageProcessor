#!/bin/sh

set -e

sleep 5

./migrator --migrations-path=./migrations --storage-path="user:password@postgres:5432/messages"

./message_processor &

./consumer