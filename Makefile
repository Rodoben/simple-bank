migrateup:
	migrate -path db/migration -database "postgresql://username1:strongpassword@localhost:5432/simplebank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://username1:strongpassword@localhost:5432/simplebank?sslmode=disable" -verbose down
up:
	docker compose up -d

