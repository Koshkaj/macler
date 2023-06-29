package core

import (
	"context"
	"encoding/json"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/koshkaj/macler/parser/db"
	"gitlab.com/koshkaj/macler/parser/providers"
	"gitlab.com/koshkaj/macler/parser/types"
	"gitlab.com/koshkaj/macler/parser/util"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"sync"
)

type ParserFunc func(context.Context, map[string]struct{}, chan<- types.DataNormalized, types.IncommingEventMessage)

// rac xdeba while loopamde
// loopis shignit
// 1. Gavamzadot cvladebi mongodb contexti, channeli da rabbitmq +
// 2. gavushvat event listen while loopi romelic mousmens events
//     - validacia gavuketot events romelic movida, tu araa validuti ar gavushvat parseri

// 3. rogorcki eventi mova mag queue dan yoveli provaideristvis davcallot parseri calcalke rutinebad
// 4. DAVELODOT yvela rutinas sanam morcheba ro mere mtliani data chavwerot db shi
// 5. datas chawerastan ertad gavgzavnot es data output queue shi

func InitParser() {
	parserFuncs := map[string]ParserFunc{
		"myhome": providers.RunMyhomeParser,
		//"ss":     providers.RunParseSS,
		//"livo":   providers.RunParseLivo,
		//"area":   providers.RunParseArea,
	}
	ctx := context.Background()
	mongoDB := db.InitDb(ctx)
	finalDatach := make(chan types.DataNormalized)
	var finalDataSlice []interface{}
	mq := InitMQ()
	defer mq.Conn.Close()
	msgs, err := mq.Channel.Consume(mq.ListenQueue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case msg, ok := <-msgs: // MOVIDA EVENTI
			if !ok {
				log.Fatal(errors.New("error while reading from the messages channel"))
			}
			jsonMessage, err := util.ValidateEventMessage(msg.Body)
			if err != nil {
				log.Printf("Message : %s is not valid", string(msg.Body))
				continue
			}

			go func() {
				for data := range finalDatach {
					finalDataSlice = append(finalDataSlice, data)
				}
			}()
			var wg sync.WaitGroup
			wg.Add(len(parserFuncs))
			for provider, function := range parserFuncs {
				siteSeenIDs := mongoDB.QuerySeenIDs(ctx, bson.M{"site": provider, "phone": jsonMessage.Phone})
				go func(prov string, fun ParserFunc) {
					defer wg.Done()
					util.NormalizeData(prov, &jsonMessage)
					fun(ctx, siteSeenIDs, finalDatach, jsonMessage)
				}(provider, function)
			}
			wg.Wait()

			go func() {
				mongoDB.InsertSeenIDs(ctx, finalDataSlice)
			}()

			finalDataJson, err := json.Marshal(finalDataSlice)
			if err != nil {
				log.Fatal(errors.New("json marshalling error"))
			}
			payload := amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "application/json",
				Body:         finalDataJson,
			}
			go func() {
				err := mq.Channel.PublishWithContext(ctx,
					"",
					mq.DeliverQueue.Name,
					false,
					false,
					payload,
				)
				if err != nil {
					log.Print("Error while publishing a message")
				}
			}()
			log.Printf("Message : %s finished processing", msg.MessageId)
		case <-ctx.Done():
			break
		}
	}

}
