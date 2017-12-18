# PostgreSQL implementation

Basic PostgreSQL operations

Before use install necessary libraries:

```shell
go get github.com/go-pg/pg
go get gopkg.in/ini.v1
# Or
dep ensure
```

Or with bake

```shell
bake go get github.com/go-pg/pg
bake go get gopkg.in/ini.v1
# Or
bake dep ensure
```

Create `db.ini` file from the sample and set your own configs:

```
cp db.ini.sample db.ini
```

## License

MIT.
