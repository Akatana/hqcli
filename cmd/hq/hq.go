package main

import (
	"log"
	"os"

	"github.com/Akatana/hqcli/internal"
	"github.com/Akatana/hqcli/pkg/api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("HQ_TOKEN")
	if token == "" {
		log.Fatal("Missing HQ_TOKEN in .env file")
	}

	client := api.NewClient(token)

	projects, err := client.GetProjects()
	if err != nil {
		log.Fatalf("Error fetching projects: %v", err)
	}

	selectedProject, err := internal.SelectProject(projects)
	if err != nil {
		log.Fatalf("Project selection cancelled: %v", err)
	}

	tasks, err := client.GetTasksForProject(selectedProject.ID)
	if err != nil {
		log.Fatalf("Error fetching tasks: %v", err)
	}

	selectedTask, err := internal.SelectTask(tasks)
	if err != nil {
		log.Fatalf("Task selection cancelled: %v", err)
	}
	println(selectedTask.ID)
}
