#!/bin/bash
set -e  # Exit immediately if any command fails
ssh -t deployer@172.17.0.2 <<EOF
    cd financas_v2
    echo "Building image..."
    docker build \
            --network=host \
            --build-arg BUILDKIT_INLINE_CACHE=1 \
            -t kaduhod/fin .
EOF

