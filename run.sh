#!/bin/bash

go build

set -a
source .env
set +a
echo $TRAFFIC_GENERATOR_HOSTNAME
./tx-traffic-poc

