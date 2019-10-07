package persistence

import (
	"encoding/binary"
	"encoding/json"
	"log"

	"github.com/AsynkronIT/protoactor-go/actor"
	bolt "go.etcd.io/bbolt"

	"github.com/alittlebrighter/saltshakers/configuration"
	"github.com/alittlebrighter/saltshakers/utils"
)

func BoltDBProducer() actor.Actor {
	return &BoltDBActor{BaseActor: utils.NewBaseActor("persistence.boltdb")}
}

type BoltDBActor struct {
	*utils.BaseActor

	dbPath  string
	db      *bolt.DB
	buckets map[string]bool
}

func (state *BoltDBActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case Create:
		err := state.db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(msg.EntityType))
			if err != nil {
				return err
			}

			if !msg.Upsert || len(msg.Entity.GetId()) == 0 {
				id, err := b.NextSequence()
				if err != nil {
					return err
				}
				msg.Entity.SetId(itob(id))
			}

			marshaled, err := json.Marshal(msg.Entity)
			if err != nil {
				return err
			}
			return b.Put(msg.Entity.GetId(), marshaled)
		})
		if err != nil {
			log.Println(state.Name(), err)
		}
		context.Respond(msg)
	case GetOne:
		err := state.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(msg.EntityType))

			marshaled := b.Get(msg.Id)
			json.Unmarshal(marshaled, msg.Entity)
			log.Println("got household", msg.Entity)
			context.Respond(msg)
			return nil
		})
		if err != nil {
			log.Println(state.Name(), err)
		}
	case Query: // only gets all for now
		log.Println(state.Name(), ": querying", msg.EntityType)
		err := state.db.View(func(tx *bolt.Tx) error {
			msg.Entities = []interface{}{}

			b := tx.Bucket([]byte(msg.EntityType))
			if b == nil {
				context.Respond(msg)
				return nil
			}

			c := b.Cursor()

			empty := msg.Model

			for k, v := c.First(); k != nil; k, v = c.Next() {
				next := empty

				// TODO apply filters

				json.Unmarshal(v, &next)
				msg.Entities = append(msg.Entities, next)
			}

			context.Respond(msg)

			return nil
		})
		if err != nil {
			log.Println(state.Name(), err)
		}
	case []configuration.PersistenceConfig:
		for _, config := range msg {
			if config.Kind() != configuration.Bolt {
				continue
			}
			state.init(config)
		}
	case *actor.Started:
		state.buckets = make(map[string]bool)
		if len(state.dbPath) == 0 {
			return
		}

		var err error
		state.db, err = bolt.Open(state.dbPath, 0666, nil)
		if err != nil {
			log.Println(state.Name(), err)
		}
	case *actor.Stopped:
		state.db.Close()
	case *actor.Restarting:
		state.Restarting(context, msg)
	default:
		context.Forward(state.Children())
	}
}

func (state *BoltDBActor) init(config configuration.PersistenceConfig) {
	state.dbPath = config.Params["dbPath"].(string)

	var err error
	state.db, err = bolt.Open(state.dbPath, 0666, nil)
	if err != nil {
		log.Println(state.Name(), err)
	}
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
