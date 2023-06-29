package api

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"gitlab.com/koshkaj/macler/backend/types"
	"net/http"
)

func (s *Server) handleSchedule(c echo.Context) error {
	var inputData types.ScheduleInput
	if err := c.Bind(&inputData); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(inputData); err != nil {
		return err
	}
	schedule := fmt.Sprintf("0/1 %s * * *", inputData.Interval)
	_, err := s.scheduleCron(schedule, &inputData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":  "could not schedule a task",
			"detail": err.Error()})
	}
	go func() {
		mongoInputData := &types.CronMongoInput{
			Schedule: schedule,
			Data:     inputData,
		}
		s.Mongo.InsertCronJob(context.Background(), mongoInputData)
	}()

	return c.JSON(http.StatusOK, map[string]string{"detail": "succesfully scheduled a task"})
}

func (s *Server) handleScheduleList(c echo.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	crontabs := s.Cr.Entries()
	for _, cront := range crontabs {
		fmt.Println(cront.Schedule)
	}
	return c.JSON(http.StatusOK, "")
}
