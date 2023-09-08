set -a            
source .env
set +a

cd internal/migrations
goose postgres $DB_URL up
cd ../..