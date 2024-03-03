# sql-transactions

This version only support bun.Tx for now.

## Dependencies
1. `go 1.21.5`

## Usage
```go
dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", "dbUser", "dbPassword", "dbHost", "dbPort", "dbName")

sqlDB, err := sql.Open("mysql", dsn)

db := bun.NewDB(sqlDB, mysqldialect.New())

trans, err := NewTransaction(db)
if err != nil {
    log.Fatal(err)
}

tx, err := trans.Init()
if err != nil {
    log.Fatal(err)
}

err = trans.Exec(tx, func(ctx context.Context, tx *bun.Tx) error {
if err := test1(ctx, tx); err != nil {
    return err
}
if err := test2(ctx, tx); err != nil {
    return err
}
return test3(ctx, tx)
})
```