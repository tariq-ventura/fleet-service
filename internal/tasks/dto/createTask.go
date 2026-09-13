package tasks_dto

import "time"

type CreateTaskRequest struct {
	TaskID        string    `json:"taskId" binding:"required,max=100"`
	Title         string    `json:"title" binding:"required,max=200"`
	Description   string    `json:"description" binding:"omitempty,max=2000"`
	Type          string    `json:"type" binding:"required,max=100"`
	Status        string    `json:"status" binding:"required,max=100"`
	ScheduledDate time.Time `json:"scheduledDate" binding:"required"`
	Origin        string    `json:"origin" binding:"required,max=250"`
	Destination   string    `json:"destination" binding:"required,max=250"`
	Latitude      int64     `json:"latitude" binding:"required"`
	Longitude     int64     `json:"longitude" binding:"required"`
	Assignee      string    `json:"assignee" binding:"required,max=200"`
}
