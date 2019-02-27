package federation

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/blocksafe/go/support/db/dbtest"
	"github.com/blocksafe/go/support/http/httptest"
)

func TestHandler(t *testing.T) {
	db := dbtest.Postgres(t).Load(`
    CREATE TABLE people (id character varying, name character varying, domain character varying);
    INSERT INTO people (id, name, domain) VALUES 
      ('GD2GJPL3UOK5LX7TWXOACK2ZPWPFSLBNKL3GTGH6BLBNISK4BGWMFBBG', 'scott', 'blocksafe.org'),
      ('GCYMGWPZ6NC2U7SO6SMXOP5ZLXOEC5SYPKITDMVEONLCHFSCCQR2J4S3', 'bartek', 'blocksafe.org');
  `)
	defer db.Close()

	driver := &ReverseSQLDriver{
		SQLDriver: SQLDriver{
			DB:                db.Open().DB,
			Dialect:           db.Dialect,
			LookupRecordQuery: "SELECT id FROM people WHERE name = ? AND domain = ?",
		},
		LookupReverseRecordQuery: "SELECT name, domain FROM people WHERE id = ?",
	}

	defer driver.DB.Close()

	handler := &Handler{driver}
	server := httptest.NewServer(t, handler)
	defer server.Close()

	// Good name request
	server.GET("/federation").
		WithQuery("type", "name").
		WithQuery("q", "scott*blocksafe.org").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("account_id").
		ValueEqual("account_id", "GD2GJPL3UOK5LX7TWXOACK2ZPWPFSLBNKL3GTGH6BLBNISK4BGWMFBBG")

	// No record in DB
	server.GET("/federation").
		WithQuery("type", "name").
		WithQuery("q", "jed*blocksafe.org").
		Expect().
		Status(http.StatusNotFound).
		JSON().Object().
		ContainsKey("code").
		ValueEqual("code", "not_found")

	// Invalid addresses
	server.GET("/federation").
		WithQuery("type", "name").
		WithQuery("q", "scott**blocksafe.org").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ContainsKey("code").
		ValueEqual("code", "invalid_query")

	server.GET("/federation").
		WithQuery("type", "name").
		WithQuery("q", "scott").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ContainsKey("code").
		ValueEqual("code", "invalid_query")

	// Missing query
	server.GET("/federation").
		WithQuery("type", "name").
		WithQuery("q", "").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ContainsKey("code").
		ValueEqual("code", "invalid_request").
		ContainsKey("message").
		ValueEqual("message", "q parameter is blank")

		// Different domain
	server.GET("/federation").
		WithQuery("type", "name").
		WithQuery("q", "scott*example.com").
		Expect().
		Status(http.StatusNotFound).
		JSON().Object().
		ContainsKey("code").
		ValueEqual("code", "not_found")

	// Reverse federation questions

	// Good request
	server.GET("/federation").
		WithQuery("type", "id").
		WithQuery("q", "GD2GJPL3UOK5LX7TWXOACK2ZPWPFSLBNKL3GTGH6BLBNISK4BGWMFBBG").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("blocksafe_address").
		ValueEqual("blocksafe_address", "scott*blocksafe.org")

	// No record in DB
	server.GET("/federation").
		WithQuery("type", "id").
		WithQuery("q", "GA3R753JKGXU6ETHNY3U6PYIY7D6UUCXXDYBRF4XURNAGXW3CVGQH2ZA").
		Expect().
		Status(http.StatusNotFound).
		JSON().Object().
		ContainsKey("code").
		ValueEqual("code", "not_found")

	// TXID request
	server.GET("/federation").
		WithQuery("type", "txid").
		WithQuery("q", "hello").
		Expect().
		Status(http.StatusNotImplemented).
		JSON().Object().
		ContainsKey("code").
		ValueEqual("code", "not_implemented")

	// Invalid type
	server.GET("/federation").
		WithQuery("type", "foo").
		WithQuery("q", "hello").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ContainsKey("code").
		ValueEqual("code", "invalid_request")

}

func TestNameHandler(t *testing.T) {
	db := dbtest.Postgres(t).Load(`
    CREATE TABLE people (id character varying, name character varying, domain character varying);
    INSERT INTO people (id, name, domain) VALUES 
      ('GD2GJPL3UOK5LX7TWXOACK2ZPWPFSLBNKL3GTGH6BLBNISK4BGWMFBBG', 'scott', 'blocksafe.org'),
      ('GCYMGWPZ6NC2U7SO6SMXOP5ZLXOEC5SYPKITDMVEONLCHFSCCQR2J4S3', 'bartek', 'blocksafe.org');
  `)
	defer db.Close()

	driver := &SQLDriver{
		DB:                db.Open().DB,
		Dialect:           db.Dialect,
		LookupRecordQuery: "SELECT id FROM people WHERE name = ? AND domain = ?",
	}

	defer driver.DB.Close()

	handler := &Handler{driver}
	server := httptest.NewServer(t, handler)
	defer server.Close()

	// Good name request
	server.GET("/federation").
		WithQuery("type", "name").
		WithQuery("q", "scott*blocksafe.org").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("account_id").
		ValueEqual("account_id", "GD2GJPL3UOK5LX7TWXOACK2ZPWPFSLBNKL3GTGH6BLBNISK4BGWMFBBG")

	// Reverse request
	server.GET("/federation").
		WithQuery("type", "id").
		WithQuery("q", "GA3R753JKGXU6ETHNY3U6PYIY7D6UUCXXDYBRF4XURNAGXW3CVGQH2ZA").
		Expect().
		Status(http.StatusNotImplemented).
		JSON().Object().
		ContainsKey("code").
		ValueEqual("code", "not_implemented")
}

type ForwardTestDriver struct{}

func (fd ForwardTestDriver) LookupForwardingRecord(query url.Values) (*Record, error) {
	if query.Get("acct") == "1234" {
		return &Record{
			AccountID: "GD2GJPL3UOK5LX7TWXOACK2ZPWPFSLBNKL3GTGH6BLBNISK4BGWMFBBG",
			MemoType:  "id",
			Memo:      "1",
		}, nil
	} else if query.Get("acct") == "4321" {
		return &Record{
			AccountID: "GD2GJPL3UOK5LX7TWXOACK2ZPWPFSLBNKL3GTGH6BLBNISK4BGWMFBBG",
			MemoType:  "text",
			Memo:      "test",
		}, nil
	} else {
		return nil, nil
	}
}

func (fd ForwardTestDriver) LookupRecord(name string, domain string) (*Record, error) {
	return nil, nil
}

func TestForwardHandler(t *testing.T) {
	handler := &Handler{ForwardTestDriver{}}
	server := httptest.NewServer(t, handler)
	defer server.Close()

	// Good forward request
	server.GET("/federation").
		WithQuery("type", "forward").
		WithQuery("forward_type", "bank_account").
		WithQuery("acct", "1234").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("account_id").
		ValueEqual("account_id", "GD2GJPL3UOK5LX7TWXOACK2ZPWPFSLBNKL3GTGH6BLBNISK4BGWMFBBG").
		ContainsKey("memo_type").
		ValueEqual("memo_type", "id").
		ContainsKey("memo").
		ValueEqual("memo", "1")

		// Good forward request
	server.GET("/federation").
		WithQuery("type", "forward").
		WithQuery("forward_type", "bank_account").
		WithQuery("acct", "4321").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("account_id").
		ValueEqual("account_id", "GD2GJPL3UOK5LX7TWXOACK2ZPWPFSLBNKL3GTGH6BLBNISK4BGWMFBBG").
		ContainsKey("memo_type").
		ValueEqual("memo_type", "text").
		ContainsKey("memo").
		ValueEqual("memo", "test")

	// Not Found forward request
	server.GET("/federation").
		WithQuery("type", "forward").
		WithQuery("forward_type", "bank_account").
		WithQuery("acct", "8888").
		Expect().
		Status(http.StatusNotFound).
		JSON().Object().
		ContainsKey("code").
		ValueEqual("code", "not_found")
}