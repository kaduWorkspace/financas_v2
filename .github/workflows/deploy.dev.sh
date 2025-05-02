#!/bin/bash
ssh -t deployer@172.17.0.2 <<EOF
    docker pull kaduhod/fin
    docker stop cdb
    docker rm cdb
    docker run --name cdb -d -p 3000:3000 kaduhod/fin
EOF
