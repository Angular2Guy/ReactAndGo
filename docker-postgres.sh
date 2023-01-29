docker pull postgres:14
# docker run --name local-postgres-reactandgo -e POSTGRES_PASSWORD=sven1 -e POSTGRES_USER=sven1 -e POSTGRES_DB=reactandgo -p 5432:5432 -d postgres:14
docker run --name local-postgres-reactandgo -e POSTGRES_PASSWORD=sven1 -e POSTGRES_USER=sven1 -e POSTGRES_DB=reactandgo --network=host -d postgres:14

# docker start local-postgres-reactandgo
# docker stop local-postgres-reactandgo
# docker exec -it local-postgres-reactandgo bash
# pg_dump -h localhost -U sven1 -d reactandgo -c > reactandgo.sql
# psql -h localhost -U sven1 -d reactandgo < reactandgo.sql