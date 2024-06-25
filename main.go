package main

import (
	"fmt"
)

func main() {
    scheduler := NewJobScheduler()

	// Add jobs Section
    scheduler.AddJob("Job1", []string{})
    scheduler.AddJob("Job2", []string{"Job1"})
	scheduler.AddJob("Job3", []string{"Job2"})
	scheduler.AddJob("Job4", []string{"Job3", "Job2"})

	fmt.Println("\nInitial jobs and their dependencies:")
    scheduler.DisplayJobs()
    fmt.Println("\nInitial job queue:")
    scheduler.DisplayJobQueue()

	// Add dependency Section
	fmt.Println("\nAdding Dependency Job4 -> Job1")
	scheduler.AddDependency("Job4", "Job1")
    fmt.Println("\nJobs after adding dependency:")
    scheduler.DisplayJobs()
    fmt.Println("\nJob queue after adding dependency:")
    scheduler.DisplayJobQueue()

	// Remove Dependency Section
	fmt.Println("\nRemoving Dependency Job4 -> Job1")
	scheduler.RemoveDependency("Job4", "Job1")
	fmt.Println("\nJobs after removing dependency:")
	scheduler.DisplayJobs()
    fmt.Println("\nJob queue after removing dependency:")
    scheduler.DisplayJobQueue()

	// Remove jobs Section
	scheduler.RemoveJob("Job2")

	fmt.Println("\nAfter removing job Job2:")
	fmt.Println("Jobs after removing job:")
	scheduler.DisplayJobs()
	fmt.Println("\nJob queue after removing job:")
	scheduler.DisplayJobQueue()

	// Process job Section
	scheduler.ProcessJob()
	fmt.Println("\nAfter processing a job:")
	fmt.Println("Jobs after processing:")
	scheduler.DisplayJobs()
	fmt.Println("\nJob queue after processing:")
	scheduler.DisplayJobQueue()
}
