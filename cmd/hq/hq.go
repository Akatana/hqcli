package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Akatana/hqcli/internal"
	"github.com/Akatana/hqcli/pkg/api"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hq",
	Short: "HelloHQ CLI Tool",
	Long:  `A CLI to work with the HQ v2 API`,
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Create a user report",
	Run:   reportHandler,
}

var (
	duration float64
	breakHrs float64
	startStr string
	endStr   string
	name     string
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	reportCmd.Flags().Float64VarP(&duration, "duration", "d", 0, "Duration in hours")
	reportCmd.Flags().Float64VarP(&breakHrs, "break", "b", 0, "Break duration in hours")
	reportCmd.Flags().StringVarP(&startStr, "start", "s", "", "Start time (HH:MM)")
	reportCmd.Flags().StringVarP(&endStr, "end", "e", "", "End time (HH:MM)")
	reportCmd.Flags().StringVarP(&name, "name", "n", "", "Report name")

	reportCmd.MarkFlagRequired("name")
	reportCmd.MarkFlagsOneRequired("duration", "start")

	rootCmd.AddCommand(reportCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func reportHandler(cmd *cobra.Command, args []string) {
	token := os.Getenv("HQ_TOKEN")
	user := os.Getenv("HQ_USER_ID")
	if token == "" || user == "" {
		log.Fatal("Missing HQ_TOKEN or HQ_USER_ID in .env file")
	}
	userId, err := strconv.Atoi(user)
	if err != nil {
		log.Fatal("invalid HQ_USER_ID in .env file, must be an integer")
	}
	client := api.NewClient(token)

	var startTime, endTime time.Time
	now := time.Now()
	loc := now.Location()
	if startStr != "" && endStr != "" {
		startTime, err = parseTime(startStr, loc)
		if err != nil {
			log.Fatalf("Invalid start time: %v", err)
		}
		endTime, err = parseTime(endStr, loc)
		if err != nil {
			log.Fatalf("Invalid end time: %v", err)
		}
	} else if duration > 0 {
		endTime = now.Truncate(time.Minute)
		workDuration := time.Duration((duration - breakHrs) * float64(time.Hour))
		startTime = endTime.Add(-workDuration)
	} else {
		log.Fatal("Either --start and --end or --d must be provided")
	}

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

	report := api.UserReporting{
		UserID:  userId,
		TaskID:  selectedTask.ID,
		Name:    name,
		StartOn: startTime.UTC(),
		EndOn:   endTime.UTC(),
	}

	if err := client.CreateUserReporting(report); err != nil {
		log.Fatalf("Failed to create report: %v", err)
	}

	fmt.Println("✅ Report submitted.")
}

func parseTime(hhmm string, loc *time.Location) (time.Time, error) {
	today := time.Now().In(loc).Format("2006-01-02")
	layout := "2006-01-02T15:04:05"
	return time.ParseInLocation(layout, today+"T"+hhmm+":00", loc)
}
