package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/robfig/cron/v3"
	"gitlab.com/koshkaj/macler/backend/core"
	"gitlab.com/koshkaj/macler/backend/db"
	"gitlab.com/koshkaj/macler/backend/types"
	"gitlab.com/koshkaj/macler/backend/util"
	"gitlab.com/koshkaj/macler/backend/validators"
	"log"
	"net/http"
	"sync"
)

type Server struct {
	*echo.Echo
	*db.Mongo
	mq *core.RabbitMQ

	mu sync.Mutex
	Cr *cron.Cron
}

func NewServer(e *echo.Echo, mq *core.RabbitMQ, cr *cron.Cron, mongo *db.Mongo) *Server {
	return &Server{Echo: e, mq: mq, Cr: cr, Mongo: mongo}
}

func InitServer() *Server {
	e := echo.New()
	mq := core.InitMQ()
	mongo := db.InitDb(context.Background())
	loc := util.LoadLocalTime()
	cr := cron.New(cron.WithLocation(loc))
	cr.Start()
	e.Validator = validators.NewValidator()
	e.Use(middleware.Logger())

	server := NewServer(e, mq, cr, mongo)
	indexGroup := server.Group("/")
	{
		indexGroup.Add("GET", "", server.handleIndex)
		indexGroup.Add("GET", "healthz", server.handleHealthz)
		indexGroup.Add("GET", "readiness", server.handleReadiness)
		indexGroup.Add("POST", "schedule", server.handleSchedule)
		indexGroup.Add("GET", "cronList", server.handleScheduleList)
	}
	server.listenToMQ()
	return server
}

func (s *Server) listenToMQ() {
	messages, err := s.mq.Channel.Consume(s.mq.ListenQueue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for m := range messages {
			var incomingData []interface{}
			err = json.Unmarshal(m.Body, &incomingData)
			if err != nil {
				log.Fatal("Something wrong happened while unmarshaling")
			}
			fmt.Println(incomingData)
		}
	}()
}

func (s *Server) scheduleCron(schedule string, data *types.ScheduleInput) (cron.EntryID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Cr.AddFunc(schedule, func() {
		ctx := context.Background()
		payloadJson, err := json.Marshal(&data)
		if err != nil {
			log.Fatal(err)
		}
		payload := amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         payloadJson,
		}
		err = s.mq.Channel.PublishWithContext(ctx,
			"",
			s.mq.DeliverQueue.Name,
			false,
			false,
			payload,
		)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("sent : %s -> parser\n", string(payloadJson))
	})
}

func (s *Server) handleIndex(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func (s *Server) handleHealthz(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"health": "ok",
	})
}

func (s *Server) handleReadiness(c echo.Context) error {
	if s.mq.Conn.IsClosed() {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"status": "unavailable",
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ready",
	})
}
