docker pull postgres:14
# docker run --name local-postgres-angularandgo -e POSTGRES_PASSWORD=sven1 -e POSTGRES_USER=sven1 -e POSTGRES_DB=angularandgo -p 5432:5432 -d postgres:14
docker run --name local-postgres-angularandgo -e POSTGRES_PASSWORD=sven1 -e POSTGRES_USER=sven1 -e POSTGRES_DB=angularandgo --network=host -d postgres:14

# docker start local-postgres-angularandgo
# docker stop local-postgres-angularandgo
# docker exec -it local-postgres-angularandgo bash
# pg_dump -h localhost -U sven1 -d angularandgo -c > angularandgo.sql
# psql -h localhost -U sven1 -d angularandgo < angularandgo.sql