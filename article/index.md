Initialize the project

```
go mod init github.com/<your_username>/golang_internet_clipboard
```

Install the required modules

```sh
go get github.com/gin-gonic/gin
go get github.com/jmoiron/sqlx
go get github.com/lib/pq
go get github.com/pressly/goose/v3/cmd/goose@latest
```

```sh
goose -dir db/migrations create create_clips_table sql
```

Add the delete_after field to the schema

```sh
goose create add_delete_after_column_to_clips sql
```
