#!/usr/bin/env bash

sleep 1

echo "$1" > messages.txt

if [ "$1" = "000" ]; then
  exit 1
fi
