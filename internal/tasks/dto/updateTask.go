package tasks_dto

import "time"

type UpdateTaskRequest struct {
	TaskID        *string    `json:"taskId" binding:"omitempty,max=100"`
	Title         *string    `json:"title" binding:"omitempty,max=200"`
	Description   *string    `json:"description" binding:"omitempty,max=2000"`
	Type          *string    `json:"type" binding:"omitempty,max=100"`
	Status        *string    `json:"status" binding:"omitempty,max=100"`
	ScheduledDate *time.Time `json:"scheduledDate"`
	Origin        *string    `json:"origin" binding:"omitempty,max=250"`
	Destination   *string    `json:"destination" binding:"omitempty,max=250"`
	Latitude      *int64     `json:"latitude"`
	Longitude     *int64     `json:"longitude"`
	Assignee      *string    `json:"assignee" binding:"omitempty,max=200"`
}
