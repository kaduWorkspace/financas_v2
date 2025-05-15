#!/bin/bash
docker run -d -p 5432:5432 --name findb -e POSTGRES_PASSWORD=123456 -e POSTGRES_USER=root -e POSTGRES_DB=fin postgres

