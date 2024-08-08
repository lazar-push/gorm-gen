// Code generated by github.com/lazar-push/gorm-gen. DO NOT EDIT.
// Code generated by github.com/lazar-push/gorm-gen. DO NOT EDIT.
// Code generated by github.com/lazar-push/gorm-gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"github.com/lazar-push/gorm-gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q          = new(Query)
	Bank       *bank
	CreditCard *creditCard
	Customer   *customer
	Person     *person
	User       *user
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	Bank = &Q.Bank
	CreditCard = &Q.CreditCard
	Customer = &Q.Customer
	Person = &Q.Person
	User = &Q.User
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:         db,
		Bank:       newBank(db, opts...),
		CreditCard: newCreditCard(db, opts...),
		Customer:   newCustomer(db, opts...),
		Person:     newPerson(db, opts...),
		User:       newUser(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	Bank       bank
	CreditCard creditCard
	Customer   customer
	Person     person
	User       user
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:         db,
		Bank:       q.Bank.clone(db),
		CreditCard: q.CreditCard.clone(db),
		Customer:   q.Customer.clone(db),
		Person:     q.Person.clone(db),
		User:       q.User.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:         db,
		Bank:       q.Bank.replaceDB(db),
		CreditCard: q.CreditCard.replaceDB(db),
		Customer:   q.Customer.replaceDB(db),
		Person:     q.Person.replaceDB(db),
		User:       q.User.replaceDB(db),
	}
}

type queryCtx struct {
	Bank       IBankDo
	CreditCard ICreditCardDo
	Customer   ICustomerDo
	Person     IPersonDo
	User       IUserDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		Bank:       q.Bank.WithContext(ctx),
		CreditCard: q.CreditCard.WithContext(ctx),
		Customer:   q.Customer.WithContext(ctx),
		Person:     q.Person.WithContext(ctx),
		User:       q.User.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
