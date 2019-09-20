package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

func main() {
	conn, err := pgx.Connect(context.Background(), "")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "drop type if exists foobar;")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = conn.Exec(context.Background(), "create type foobar as (foo text, bar int8);")
	if err != nil {
		log.Fatalln(err)
	}

	var buf []byte
	err = conn.QueryRow(context.Background(), "select '(hey,42)'::foobar").Scan(&buf)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(buf, string(buf))

	err = conn.QueryRow(context.Background(), "select '(hey,42)'::foobar", pgx.QueryResultFormats{pgx.BinaryFormatCode}).Scan(&buf)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(buf)

	var record pgtype.Record
	err = conn.QueryRow(context.Background(), "select '(hey,42)'::foobar", pgx.QueryResultFormats{pgx.BinaryFormatCode}).Scan(&record)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v\n", record)

	var foobar Foobar
	err = conn.QueryRow(context.Background(), "select '(hey,42)'::foobar", pgx.QueryResultFormats{pgx.BinaryFormatCode}).Scan(&foobar)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v\n", foobar)
}

type Foobar struct {
	foo string
	bar int64
}

func (fb *Foobar) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	var record pgtype.Record
	err := record.DecodeBinary(ci, src)
	if err != nil {
		return err
	}

	err = record.Fields[0].AssignTo(&fb.foo)
	if err != nil {
		return err
	}

	err = record.Fields[1].AssignTo(&fb.bar)
	if err != nil {
		return err
	}

	return nil
}
