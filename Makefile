DB_DRIVER=mysql
DB_DSN=poker:poker@tcp(localhost:3307)/poker?parseTime=true

migrate-up:
	goose -dir migrations $(DB_DRIVER) "$(DB_DSN)" up

migrate-down:
	goose -dir migrations $(DB_DRIVER) "$(DB_DSN)" down

migrate-status:
	goose -dir migrations $(DB_DRIVER) "$(DB_DSN)" status

migration:
	goose -dir migrations create $(name) sql