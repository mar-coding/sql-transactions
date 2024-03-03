package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"log"
)

func ExampleNewTransaction() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", "dbUser", "dbPassword", "dbHost", "dbPort", "dbName")

	sqlDB, err := sql.Open("mysql", dsn)

	db := bun.NewDB(sqlDB, mysqldialect.New())

	// Start a transaction
	trans, err := NewTransaction(db)
	if err != nil {
		log.Fatal(err)
	}

	err = trans.Exec(func(ctx context.Context, tx *bun.Tx) error {
		if err := test1(ctx, tx); err != nil {
			return err
		}
		if err := test2(ctx, tx); err != nil {
			return err
		}
		return test3(ctx, tx)
	})
}

func test1(ctx context.Context, tx *bun.Tx) error {
	// first test
	return nil
}

func test2(ctx context.Context, tx *bun.Tx) error {
	// second test
	return nil
}

func test3(ctx context.Context, tx *bun.Tx) error {
	// third test
	return nil
}
