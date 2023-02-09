package hdModel

import (
	"context"
	"github.com/bamtols/mong-go/extends/scalars"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type (
	IUpdate[T any] struct {
		Filter           bson.M
		UpdateValues     IFData[T]
		FindAndUpdateOpt *options.FindOneAndUpdateOptions
		UpdateOpt        *options.UpdateOptions
	}
)

type (
	toolUpdate[T any] struct {
	}
)

func newToolUpdate[T any]() *toolUpdate[T] {
	return &toolUpdate[T]{}
}

func (x *toolUpdate[T]) Update(ctx context.Context, mdb *mongo.Database, params *IUpdate[T]) {

}

func (x *toolUpdate[T]) UpdateWithTrx(ctx context.Context, mdb *mongo.Database, params *IUpdate[T]) (md *Model[T], err error) {
	col := mdb.Collection(params.UpdateValues.GetColNm())

	filter := params.Filter
	filter[mFieldProcState] = ProcStateCommit

	now := time.Now()
	procInfo := &ProcInfo{
		HistoryId: scalars.NewObjectID(),
		ProcAt:    now,
	}

	var res = col.FindOneAndUpdate(
		ctx,
		&filter,
		&bson.M{
			"$set": &bson.M{
				mFieldProcState: ProcStateProcessing,
				mFiledProcInfo:  procInfo,
				mFieldData:      params.UpdateValues.Data(),
				mFiledUpdatedAt: now,
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

	// 히스토리 생성
	// 업데이트 전 데이터로 히스토리 작성
	md.ProcInfo = procInfo
	md.colNm = params.UpdateValues.GetColNm()
	if err = x.createHistory(ctx, mdb, md); err != nil {
		return
	}

	// 새로운 모델로 데이터 업데이트
	md.isValid = true
	md.isLoaded = true
	md.migList = params.UpdateValues.GetMigrateList()
	md.ProcState = ProcStateProcessing
	md.Data = params.UpdateValues.Data()

	return
}

func (x *toolUpdate[T]) createHistory(ctx context.Context, mdb *mongo.Database, prevMd *Model[T]) (err error) {
	col := mdb.Collection(prevMd.historyColNm())
	// todo 이곳부터 시작
}
