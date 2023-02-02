package cores

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

type (
	FnTrx[T any]  func(sessCtx context.Context, sessDB *mongo.Database) (T, error)
	TrxMng[T any] struct{}
)

// Transaction 몽고 디비에는 lock and wait 기능이 없다.
func Transaction[T any](
	ctx context.Context,
	mdb *Client,
	fnTrx FnTrx[T],
) (T, error) {
	mng := &TrxMng[T]{}
	return mng.trx(ctx, mdb, fnTrx)
}

func (x *TrxMng[T]) trx(
	ctx context.Context,
	mdb *Client,
	fnTrx FnTrx[T],
) (T, error) {
	if mdb.canTrx {
		return x.doTrx(ctx, mdb, fnTrx)
	} else {
		return x.doNotTrx(ctx, mdb, fnTrx)
	}
}

func (x *TrxMng[T]) doTrx(
	ctx context.Context,
	mdb *Client,
	fnTrx FnTrx[T],
) (T, error) {
	emptyRes := *new(T)

	session, err := mdb.StartSession()
	if err != nil {
		return emptyRes, err
	}
	defer session.EndSession(ctx)

	res, err := session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		return fnTrx(sessCtx, sessCtx.Client().Database(mdb.dbNm))
	}, options.
		Transaction().
		SetWriteConcern(writeconcern.New(
			writeconcern.WMajority(),
			writeconcern.WTimeout(1*time.Second),
		)),
	)

	if err != nil {
		return emptyRes, err
	}

	parsed, isOk := res.(T)
	if !isOk {
		return emptyRes, fmt.Errorf("res is not T value")
	}

	return parsed, nil
}

func (x *TrxMng[T]) doNotTrx(
	ctx context.Context,
	mdb *Client,
	fnTrx FnTrx[T],
) (T, error) {
	return fnTrx(ctx, mdb.UsedDatabase())
}
