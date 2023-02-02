package cores

import (
	"context"
	"github.com/flowd-cores/fn-go/independent/fnParams"
	"github.com/flowd-cores/fn-go/independent/lbGqlgen/scalars"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type (
	Model[T any] struct {
		Id        scalars.ObjectID `bson:"_id"`
		Data      T                `bson:"data"`
		Pending   *ModelPending    `bson:"pending"`
		CreatedAt time.Time        `bson:"createdAt"`
		UpdatedAt time.Time        `bson:"updatedAt"`
	}

	IFModelParams[T any] interface {
		Model[T]
	}

	ModelHistory[T any] struct {
		Id        scalars.ObjectID `bson:"_id"`
		TrxState  TrxState         `bson:"trxState"`
		Data      T                `bson:"data"`
		CreatedAt time.Time        `bson:"createdAt"`
	}

	TrxState string

	ModelPending struct {
		HistoryId scalars.ObjectID `bson:"historyId"`
		PendingAt time.Time        `bson:"pendingAt"`
	}
)

const (
	TrxStatePending  TrxState = "Pending"
	TrxStateCommit   TrxState = "Commit"
	TrxStateRollback TrxState = "Rollback"
)

func (x *Model[T]) Update(ctx context.Context, data T, client *Client, withTrx ...bool) (*Model[T], error) {
	if fnParams.Pick(withTrx) {
		return x.updateWithTrx(ctx, data, client)
	} else {
		return x.update(ctx, data, client)
	}
}

func (x *Model[T]) updateWithTrx(ctx context.Context, data T, client *Client) (*Model[T], error) {
	return Transaction(ctx, client, func(sessCtx context.Context, sessDB *mongo.Database) (m *Model[T], err error) {

		panic("notImpl")
		return
	})
}

func (x *Model[T]) update(ctx context.Context, data T, client *Client) (m *Model[T], err error) {
	panic("notImpl")
}

func (x *Model[T]) NewHistory(state ...TrxState) *ModelHistory[T] {
	return &ModelHistory[T]{
		Id:        scalars.NewObjectID(),
		TrxState:  fnParams.PickWithDefault(state, TrxStatePending),
		Data:      x.Data,
		CreatedAt: time.Now(),
	}
}
