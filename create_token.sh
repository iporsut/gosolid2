#!/usr/bin/env bash

curl -XPOST -v -H 'Content-Type: application/json' localhost:8080/auth/login -d '{"username":"admin","password":"password"}'
