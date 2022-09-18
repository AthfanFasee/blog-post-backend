database :
	migrate -path ./migrations -database "postgres://blogpost:secret@localhost:542/blogpost?sslmode=disable" -verbose up
migrateDown :
    migrate -path ./migrations -database "postgres://blogpost:secret@localhost:542/blogpost?sslmode=disable" down
