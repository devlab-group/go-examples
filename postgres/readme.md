# PostgreSQL implementation

Basic PostgreSQL operations

Before use install necessary libraries:

```
go get github.com/go-pg/pg
go get github.com/go-pg/migrations
go get gopkg.in/ini.v1
# Or
dep ensure
```

Or with bake

```
bake go get github.com/go-pg/pg
bake go get github.com/go-pg/migrations
bake go get gopkg.in/ini.v1
# Or
bake dep ensure
```

Create `db.ini` file from the sample and set your own configs:

```
cp db.ini.sample db.ini
```

Run database migrations

```
go run main.go migrations up
```

Or with bake

```
bake run main.go migrations up
```

## License

MIT.
