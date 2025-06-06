#!/bin/bash
set -e  # Exit immediately if any command fails
ssh -t deployer@172.17.0.2 <<EOF
    REPO_DIR="financas_v2"

    if [ -d "financas_v2" ]; then
        if [ -d "$REPO_DIR/.git" ]; then
            echo "O repositório '$REPO_DIR' existe e é um repositório Git."
        else
            echo "Repositorio existe mas nao tem a pasta .git"
            rm -rf ./*
            git clone git@github.com:kaduWorkspace/financas_v2.git
        fi
    else
        echo "O repositório '$REPO_DIR' nao existe."
        git clone git@github.com:kaduWorkspace/financas_v2.git
    fi
    cd financas_v2
    git checkout master
    git pull
EOF
