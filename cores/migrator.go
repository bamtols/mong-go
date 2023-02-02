package cores

import (
	"context"
	"fmt"
	"github.com/flowd-cores/fn-go/independent/fnReflect"
	"github.com/flowd-cores/fn-go/independent/fnString"
	"github.com/flowd-cores/fn-go/independent/lbGqlgen/scalars"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type (
	IFMigrator interface {
		Migrate(models ...IFDocument) error
	}

	IFDocument interface {
		GetColNm() string
		GetMigrateList() MigrateList
	}

	Migrate struct {
		MigNm string
		Fn    func(ctx context.Context, col *mongo.Collection) (err error)
	}

	MigrateList []Migrate

	MigV1 struct {
		Client   *mongo.Client
		Database *mongo.Database
	}
)

const (
	MigrateHistoryCollectionNm = "migrations"
)

func NewMigrator(m *mongo.Client, dbNm string) IFMigrator {
	return &MigV1{
		Client:   m,
		Database: m.Database(dbNm),
	}
}

func (x *MigV1) Migrate(models ...IFDocument) (err error) {
	ctx := context.TODO()
	var migList DocMigrateList
	if migList, err = x.loadMigData(ctx); err != nil {
		return
	}

	for _, model := range append([]IFDocument{
		&DocMigrate{},
	}, models...) {
		colNm := model.GetColNm()
		migData := migList.FindCol(colNm)
		migList := model.GetMigrateList()
		migCol := x.Database.Collection(colNm)

		for _, migrate := range migList {
			migNm := fnString.NewChain(migrate.MigNm).RemoveSpace().String()

			if migData.IsMigrated(migNm) {
				continue
			}

			if err = migrate.Fn(ctx, migCol); err != nil {
				return
			}

			if err = x.updateMigData(ctx, colNm, migNm); err != nil {
				return
			}
		}
	}

	return
}

func (x *MigV1) loadMigData(ctx context.Context) (res DocMigrateList, err error) {
	res = make([]DocMigrate, 0)

	var cur *mongo.Cursor
	if cur, err = x.Database.
		Collection(MigrateHistoryCollectionNm).
		Find(ctx, &bson.M{}); err != nil {
		return
	}

	if err = cur.Err(); err != nil {
		return
	}

	if err = cur.All(ctx, &res); err != nil {
		return
	}

	return
}

func (x *MigV1) updateMigData(ctx context.Context, colNm, migNm string) (err error) {
	_, err = x.Database.
		Collection(MigrateHistoryCollectionNm).
		UpdateOne(ctx, &bson.M{
			"colNm": colNm,
		}, &bson.M{
			"$set": &bson.M{
				fmt.Sprintf("migrated.%s", migNm): true,
				"updatedAt":                       time.Now(),
			},
		}, &options.UpdateOptions{
			Upsert: fnReflect.ToPointer(true),
		})
	return
}

type (
	MigData    map[string]bool
	DocMigrate struct {
		Id        scalars.ObjectID `bson:"_id"`
		ColNm     string           `bson:"colNm"`
		Migrated  MigData          `bson:"migrated"`
		CreatedAt time.Time        `bson:"createdAt"`
		UpdatedAt time.Time        `bson:"updatedAt"`
	}
	DocMigrateList []DocMigrate
)

func (x *DocMigrate) GetColNm() string {
	return MigrateHistoryCollectionNm
}

func (x *DocMigrate) GetMigrateList() MigrateList {
	return MigrateList{
		{
			MigNm: "uniqueColNm",
			Fn: func(ctx context.Context, col *mongo.Collection) (err error) {
				_, err = col.Indexes().CreateOne(ctx, mongo.IndexModel{
					Keys: &bson.D{
						{"colNm", 1},
					},
					Options: &options.IndexOptions{
						Unique: fnReflect.ToPointer(true),
					},
				})
				return
			},
		},
	}
}

func (x *DocMigrate) IsMigrated(migNm string) bool {
	v, isOk := x.Migrated[migNm]
	if !isOk {
		return false
	}
	return v
}

func (x *DocMigrateList) FindCol(colNm string) *DocMigrate {
	for _, migrate := range *x {
		if migrate.ColNm == colNm {
			return &migrate
		}
	}

	return &DocMigrate{
		Id:        scalars.NewObjectID(),
		ColNm:     colNm,
		Migrated:  make(MigData),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
