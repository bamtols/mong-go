package hdModel

import (
	"context"
	"github.com/bamtols/fn-go/fn/fnParams"
	"github.com/bamtols/mong-go/cores"
	"github.com/bamtols/mong-go/extends/scalars"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type (
	// Model T 는 데이터 모델
	Model[T any] struct {
		Id        scalars.ObjectID `bson:"_id"`
		Data      T                `bson:"data"`
		ProcState ProcState        `bson:"procState"`
		ProcInfo  *ProcInfo        `bson:"procInfo"`
		CreatedAt time.Time        `bson:"createdAt"`
		UpdatedAt time.Time        `bson:"updatedAt"`

		// 관리용 값
		isValid  bool              `bson:"-"`
		isLoaded bool              `bson:"-"`
		colNm    string            `bson:"-"`
		migList  cores.MigrateList `bson:"-"`
	}

	ModelGroup[T any] interface {
		Model[T]
	}

	ProcInfo struct {
		HistoryId scalars.ObjectID `bson:"historyId"`
		ProcAt    time.Time        `bson:"procAt"`
	}

	ModelHistory[T any] struct {
		Id        scalars.ObjectID `bson:"_id"`
		ModelId   scalars.ObjectID `bson:"modelId"`
		Data      T                `bson:"data"`
		CreatedAt time.Time        `bson:"createdAt"`
	}

	IFData[T any] interface {
		cores.IFDocument
		Data() T
	}
)

const (
	mFieldProcState = "procState"
	mFiledProcInfo  = "procInfo"
	mFiledUpdatedAt = "updatedAt"
	mFieldData      = "data"
)

type (
	ProcState string
)

const (
	ProcStateProcessing ProcState = "Processing"
	ProcStateCommit     ProcState = "Commit"
	ProcStateRollback   ProcState = "Rollback"
)

func NewModel[T any](data IFData[T]) *Model[T] {
	now := time.Now()
	return &Model[T]{
		Id:        scalars.NewObjectID(),
		Data:      data.Data(),
		ProcState: ProcStateCommit,
		ProcInfo:  nil,
		CreatedAt: now,
		UpdatedAt: now,

		// 관리용
		isValid:  true,
		isLoaded: false,
		colNm:    data.GetColNm(),
		migList:  data.GetMigrateList(),
	}
}

func UpdateOne[T any](ctx context.Context, mdb *mongo.Database, params *IUpdate[T], withTrx ...bool) (md *Model[T], err error) {
	if fnParams.Pick(withTrx) {
		col := mdb.Collection(params.UpdateValues.GetColNm())

		filter := params.Filter
		filter[mFieldProcState] = ProcStateCommit

		historyId := scalars.NewObjectID()
		now := time.Now()

		var res = col.FindOneAndUpdate(
			ctx,
			&filter,
			&bson.M{
				"$set": &bson.M{
					mFieldProcState: ProcStateProcessing,
					mFiledProcInfo: &ProcInfo{
						HistoryId: historyId,
						ProcAt:    now,
					},
				},
			},
			params.FindAndUpdateOpt,
		)

		if err = res.Err(); err != nil {
			return
		}

		md = &Model[T]{}
		if err = res.Decode(md); err != nil {
			return
		}

	} else {

	}

	return
}

func UpdateAll[T any](filter *bson.M, data IFData[T], opt *options.UpdateOptions, withTrx ...bool) ([]*Model[T], error) {
	if fnParams.Pick(withTrx) {

	} else {

	}

	panic("notImpl")
}

func Commit[T any](historyId scalars.ObjectID) (*Model[T], error) {
	panic("notImpl")
}

func (x *Model[T]) Update(data T, withTrx ...bool) (*Model[T], error) {
	panic("notImpl")
}

func (x *Model[T]) historyColNm() string {
	panic("notImpl")
}
