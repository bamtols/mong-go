package models

import (
	"github.com/bamtols/mong-go/extends/scalars"
	"time"
)

type (
	// HDModel T 는 데이터 모델
	HDModel[T any] struct {
		Id        scalars.ObjectID `bson:"_id"`
		Data      T                `bson:"data"`
		ProcState ProcState        `bson:"procState"`
		ProcInfo  *HDProcInfo      `bson:"procInfo"`
		CreatedAt time.Time        `bson:"createdAt"`
		UpdatedAt time.Time        `bson:"updatedAt"`
	}

	HDProcInfo struct {
		PrevHistoryId scalars.ObjectID `bson:"prevHistoryId"`
		ProcAt        time.Time        `bson:"procAt"`
	}

	HDModelHistory[T any] struct {
		Id        scalars.ObjectID `bson:"_id"`
		ModelId   scalars.ObjectID `bson:"modelId"`
		Data      T                `bson:"data"`
		CreatedAt time.Time        `bson:"createdAt"`
	}
)

type ProcState string

const (
	ProcStateProcessing ProcState = "Processing"
	ProcStateCommit     ProcState = "Commit"
	ProcStateRollback   ProcState = "Rollback"
)
