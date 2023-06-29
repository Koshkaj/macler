package types

type CronMongoInput struct {
	Schedule string        `json:"schedule"`
	Data     ScheduleInput `json:"data"`
}
