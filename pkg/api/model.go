package api

import "time"

type Client struct {
	token string
}

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Task struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	ParentID *int   `json:"parentId,omitempty"`
	Parent   *Task  `json:"parent,omitempty"`
}

type UserReporting struct {
	Name          string    `json:"name"`
	StartOn       time.Time `json:"startOn"`
	EndOn         time.Time `json:"endOn"`
	BreakDuration float32   `json:"breakDuration"`
	UserID        int       `json:"userId"`
	TaskID        int       `json:"taskId"`
}
