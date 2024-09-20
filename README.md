
docker run --name dbconnect -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=databack -p 5432:5432 -d postgres:16


goose postgres "postgres://postgres:password@localhost:5432/databack" up
