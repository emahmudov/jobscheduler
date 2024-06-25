package main

import (
	"container/list"
	"fmt"
)

type Set map[string]struct{}

func NewSet(items ...string) Set{
	s := Set{}
	for _, item := range items {
        s[item] = struct{}{}
    }
	return s
}

func (s Set) Add(item string) {
    s[item] = struct{}{}
}

func (s Set) Remove(item string) {
    delete(s, item)
}

func (s Set) Contains(item string) bool {
    _, exists := s[item]
    return exists
}

func (s Set) List() []string {
	items := []string{}
	for item := range s {
        items = append(items, item)
    }
	return items
}

func (s Set) Size() int {
	return len(s)
}




// Job represents a job with an ID and a set of dependencies
type Job struct {
	ID           string
	Dependencies Set
}

// JobScheduler handles scheduling and processing of jobs
type JobScheduler struct {
	jobs map[string]Job    				// ishleri saxlayacayiq 
	dependencyMap map[string]Set		// her ishin asililiqlarini izlemek 
	dependentJobs map[string]Set		// her ishin diger ishlerden asili olduqunu izlemek
	completeJob map[string]bool			// tamamlanmish ishlerin izlenmesi
	jobQueue *list.List		// ishlerin novbesi
}

// NewJobScheduler initializes a new JobScheduler
func NewJobScheduler() *JobScheduler {
	return &JobScheduler{
		jobs:   make(map[string]Job),
		dependencyMap: make(map[string]Set),
		dependentJobs: make(map[string]Set),
		completeJob:   make(map[string]bool),
		jobQueue:   list.New(),
	}
}

// AddJob adds a new job to the scheduler
// Parameters:
// - id: the unique identifier for the job
// - dependencies: a list of job IDs that this job depends on
func (js *JobScheduler) AddJob(id string, dependencies []string) {
	if _, exists := js.jobs[id]; exists {
		return
	}

    depSet := NewSet(dependencies...)
	newJob := Job{ID: id, Dependencies: depSet}
	js.jobs[id] = newJob

    js.dependencyMap[id] = depSet

    if len(dependencies) == 0 {
		js.jobQueue.PushBack(id)
	}

    for _, dep := range dependencies {
		if js.dependentJobs[dep] == nil {
			js.dependentJobs[dep] = NewSet()
		}
		js.dependentJobs[dep].Add(id)
	}
}



// RemoveJob removes a job from the scheduler
// Parameters:
// - id: the unique identifier of the job to remove
func (js *JobScheduler) RemoveJob(id string) {
	if _, exists := js.jobs[id]; !exists {
		return
	}

    // İş `jobs` xəritəsindən silinir
    delete(js.jobs, id)

    // İş `dependencyMap` xəritəsindən silinir
    dependencies := js.dependencyMap[id]
	delete(js.dependencyMap, id)

    // İşin asılı olduğu işlərdən asılılığı silinir
    for dep := range dependencies {
		if js.dependentJobs[dep] != nil {
			js.dependentJobs[dep].Remove(id)
			if js.dependentJobs[dep].Size() == 0 {
				delete(js.dependentJobs, dep)
			}
		}
	}

    // İşi `dependentJobs` xəritəsindən silirik
    dependents := js.dependentJobs[id]
	for dep := range dependents {
		js.dependencyMap[dep].Remove(id)
		if js.dependencyMap[dep].Size() == 0 {
			js.jobQueue.PushBack(dep)
		}
	}
	delete(js.dependentJobs, id)

    // İş `jobQueue` növbəsindən silinir
    for e := js.jobQueue.Front(); e != nil; e = e.Next() {
		if e.Value == id {
			js.jobQueue.Remove(e)
			break
		}
	}

    // İş `completeJob` xəritəsindən silinir
    delete(js.completeJob, id)
}



// AddDependency adds a dependency to a job
// Parameters:
// - jobID: the ID of the job to add a dependency to
// - dependencyID: the ID of the job that is a dependency
func (js *JobScheduler) AddDependency(jobID, dependencyID string) {
	if _, jobExists := js.jobs[jobID]; !jobExists {
		return
	}
	if _, depExists := js.jobs[dependencyID]; !depExists {
		return
	}
	if js.dependencyMap[jobID].Contains(dependencyID) {
		return
	}

	js.dependencyMap[jobID].Add(dependencyID)
	if js.dependentJobs[dependencyID] == nil {
		js.dependentJobs[dependencyID] = NewSet()
	}
	js.dependentJobs[dependencyID].Add(jobID)

	for e := js.jobQueue.Front(); e != nil; e = e.Next() {
		if e.Value == dependencyID {
			js.jobQueue.Remove(e)
			break
		}
	}
}


// RemoveDependency removes a dependency from a job
// Parameters:
// - jobID: the ID of the job to remove a dependency from
// - dependencyID: the ID of the dependency to remove
func (js *JobScheduler) RemoveDependency(jobID, dependencyID string) {
	if _, jobExists := js.jobs[jobID]; !jobExists {
		return
	}
	if _, depExists := js.jobs[dependencyID]; !depExists {
		return
	}

	js.dependencyMap[jobID].Remove(dependencyID)
	if js.dependencyMap[jobID].Size() == 0 {
		js.jobQueue.PushBack(jobID)
	}

	js.dependentJobs[dependencyID].Remove(jobID)
	if js.dependentJobs[dependencyID].Size() == 0 {
		delete(js.dependentJobs, dependencyID)
	}
}



// GetNextJob retrieves the next job to be processed based on dependencies
// Returns:
// - *Job: the next job to be processed or nil if no jobs are available
func (js *JobScheduler) GetNextJob() *Job {
	if js.jobQueue.Len() == 0 {
		return nil
	}

	element := js.jobQueue.Front()
	jobID := element.Value.(string)
	js.jobQueue.Remove(element)

	job, exists := js.jobs[jobID]
	if !exists {
		return nil
	}
	return &job
}



// ProcessJob processes the next job in the queue
func (js *JobScheduler) ProcessJob() {
	job := js.GetNextJob()
	if job == nil {
		fmt.Println("Dont have any jobs")
		return
	}

	fmt.Printf("Processing job: %v\n", job.ID)
	js.completeJob[job.ID] = true

	for dep := range js.dependentJobs[job.ID] {
		js.dependencyMap[dep].Remove(job.ID)
		if js.dependencyMap[dep].Size() == 0 {
			js.jobQueue.PushBack(dep)
		}
	}
	delete(js.dependentJobs, job.ID)
}


// DisplayJobQueue displays the current job queue
func (js *JobScheduler) DisplayJobQueue() {
	fmt.Println("Current Job Queue:")
	for e := js.jobQueue.Front(); e != nil; e = e.Next() {
		fmt.Printf("Job ID: %v\n", e.Value)
	}
}

// DisplayJobs displays all jobs with their details
func (js *JobScheduler) DisplayJobs() {
	fmt.Println("All jobs")
	for id, job := range js.jobs {
		fmt.Printf("Job ID: %s, Dependencies: %v\n", id, job.Dependencies.List())
	}
}

