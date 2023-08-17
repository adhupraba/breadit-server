cd internal/migrations
goose postgres postgresql://<user>:<password>@<host>:<port>/<dbname> up
cd ../..