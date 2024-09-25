migrateup:
	migrate -path db/migration -database "postgresql://username1:strongpassword@localhost:5432/simplebank?sslmode=disable" -verbose up

migrateupaws:
	migrate -path db/migration -database "postgresql://root:QPwoeiruty12345@simplebank.cfuwaoceu4bc.ap-south-1.rds.amazonaws.com:5432/simple_bank?sslmode=require" -verbose up	

migratedown:
	migrate -path db/migration -database "postgresql://username1:strongpassword@localhost:5432/simplebank?sslmode=disable" -verbose down
up:
	docker compose up -d

test:
	go test -v -cover -short ./...

