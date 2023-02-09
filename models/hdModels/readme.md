# History driven model management
* MongoDB -> RDBMS 처럼 Lock and wait 이 없다 => 동시에 쓰기작업이 들어오면 나중것을 탈락시킨다. (더 나은 방법이 있는지는 아직 찾지 못했다. 큐도 적절하지 않다.)
* 2 way commit, rollback 기능 지원을 위하여 몽고 데이터 모델은 특정 struct 를  extend 하여 만든다.

* Update 할때만 적용된다. 2 way commit 이 되지 않으면 Pending 상태가 유지되고, 1분이 지나도록 커밋이 되지 않으면 (스케쥴러)가 자동으로 롤백되게 만든다.
* 데이터 업데이트 할때는, 히스토리 생성을 위해서라도 반드시 데이터를 한번은 불러오게 되어있다.

~~~mermaid
sequenceDiagram
    Request ->> MongoModel: Update 
    MongoModel ->> MongoModelHistory: Copy
   
    activate MongoModelHistory
        MongoModelHistory ->> MongoDB : Insert Prev data
        MongoModelHistory ->> MongoModel : HistoryID
    deactivate MongoModelHistory
    
    activate MongoModel
        MongoModel ->> MongoDB : Update(with Pending state and HistoryID)
        MongoModel ->> Request : HistoryID 
    deactivate MongoModel
    
    
    alt Commit
        Request ->> MongoModel : HistoryID
        activate MongoModel
            MongoModel ->> MongoDB : Update(with deactivate PendingState)
            MongoModel ->> Request : Commit done
        deactivate MongoModel
    else Rollback
        Request ->> MongoModel : HistoryID
        activate MongoModel
            MongoModel ->> MongoModelHistory : Rollback
            MongoModelHistory ->> MongoDB : Rollback(with deactivate PendingState)
            MongoModel ->> Request : Rollback done
        deactivate MongoModel
    end
~~~

~~~mermaid
stateDiagram-v2

state MongoModel {
    state MngData {
        Id
        HistoryID
        note right of  HistoryID
            롤백, 커밋을 위한 데이터
        end note
    }
    

    
    Data
}

state MongoModelHistory {
    state HistoryMngData {
        HistoryId
    }
    
    HistoryData
}




~~~

# Migrate
* migNm 은 중복을 허락하지 않는다.


# JetBrains Go live template
### lbMongo.MigrateList
~~~go 
lbMongo.MigrateList{
    {
        MigNm: "",
        Fn: func(ctx context.Context, col *mongo.Collection) (err error) {
            _, err = col.Indexes().CreateOne(ctx, mongo.IndexModel{
                Keys: bson.D{
                    {"", -1},
                },
                Options: &options.IndexOptions{},
            })
            return
        },
    },
}

~~~

### lbMongo.Transaction
~~~go
lbMongo.Transaction[$RES$]($CTX$, $MDB$, func(sessCtx mongo.SessionContext, sessDB *mongo.Database) (res $RES$, err error) {

	return
})
~~~