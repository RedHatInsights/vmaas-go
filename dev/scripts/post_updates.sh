#!/bin/bash

curl -kv -XPOST http://localhost:8080/api/patch/v1/updates \
                -H "Content-Type: application/json" \
                -d '{"package_list":["pkg-sec-errata1-1:1.1-1.noarch"], "repository_list": ["content-set-name-1"]}'
