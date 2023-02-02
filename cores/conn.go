package cores

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type (
	IConnect struct {
		Host         string
		Username     string
		Password     string
		DbNm         string
		Direct       *bool
		ReplicaSetNm *string
	}

	Client struct {
		*mongo.Client
		canTrx bool
		dbNm   string
	}
)

func Connect(i *IConnect) (res *Client, err error) {
	ctx := context.TODO()

	op := options.Client().
		SetRegistry(
			bson.NewRegistryBuilder().
				RegisterCodec(DecimalCodecRegister()).
				RegisterCodec(ObjectIDCodecRegister()).
				Build(),
		).
		SetReadConcern(readconcern.Majority()).
		SetWriteConcern(writeconcern.New(
			writeconcern.WMajority(),
		)).
		ApplyURI(fmt.Sprintf("mongodb://%s", i.Host)).
		SetAuth(options.Credential{
			Username: i.Username,
			Password: i.Password,
		})

	if i.Direct != nil {
		op = op.SetDirect(*i.Direct)
	}

	if i.ReplicaSetNm != nil {
		op = op.SetReplicaSet(*i.ReplicaSetNm)
	}

	var cli *mongo.Client
	if cli, err = mongo.Connect(ctx, op); err != nil {
		return
	}

	if err = cli.Ping(ctx, nil); err != nil {
		return
	}

	return &Client{
		Client: cli,
		canTrx: i.ReplicaSetNm != nil,
		dbNm:   i.DbNm,
	}, nil
}

func (x *Client) UsedDatabase(opts ...*options.DatabaseOptions) *mongo.Database {
	return x.Database(x.dbNm, opts...)
}
