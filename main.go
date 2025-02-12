// # Implementation:
// - Implement a task tracking system using maps and slices
// - Tasks should have priority, status, and deadlines
// - Implement sorting by different criteria
// - Implement custom JSON marshaling

package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// Task represents a task in the task tracking system
type Task struct {
	ID        int
	Name      string
	Priority  int
	Status    string
	Deadline  time.Time
	CreatedAt time.Time
}

// TaskList represents a list of tasks
type TaskList struct {
	Tasks []Task
}

// AddTask adds a task to the task list
func (tl *TaskList) AddTask(task Task) {
	tl.Tasks = append(tl.Tasks, task)
}

// SortByPriority sorts tasks by priority
func (tl *TaskList) SortByPriority() {
	sort.Slice(tl.Tasks, func(i, j int) bool {
		return tl.Tasks[i].Priority > tl.Tasks[j].Priority
	})
}

// SortByDeadline sorts tasks by deadline
func (tl *TaskList) SortByDeadline() {
	sort.Slice(tl.Tasks, func(i, j int) bool {
		return tl.Tasks[i].Deadline.Before(tl.Tasks[j].Deadline)
	})
}

// SortByCreatedAt sorts tasks by creation time
func (tl *TaskList) SortByCreatedAt() {
	sort.Slice(tl.Tasks, func(i, j int) bool {
		return tl.Tasks[i].CreatedAt.Before(tl.Tasks[j].CreatedAt)
	})
}

// MarshalJSON custom JSON marshaling for TaskList
func (tl TaskList) MarshalJSON() ([]byte, error) {
	type Alias TaskList
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&tl),
	})
}

func main() {
	tasks := TaskList{}

	tasks.AddTask(Task{ID: 1, Name: "Task 1", Priority: 3, Status: "Pending", Deadline: time.Now().AddDate(0, 0, 5), CreatedAt: time.Now()})
	tasks.AddTask(Task{ID: 2, Name: "Task 2", Priority: 1, Status: "In Progress", Deadline: time.Now().AddDate(0, 0, 3), CreatedAt: time.Now()})
	tasks.AddTask(Task{ID: 3, Name: "Task 3", Priority: 2, Status: "Completed", Deadline: time.Now().AddDate(0, 0, 7), CreatedAt: time.Now()})

	fmt.Println("Tasks:")
	for _, task := range tasks.Tasks {
		fmt.Printf("ID: %d, Name: %s, Priority: %d, Status: %s, Deadline: %s, CreatedAt: %s\n", task.ID, task.Name, task.Priority, task.Status, task.Deadline, task.CreatedAt)
	}

	fmt.Println("\nSorting by priority:")
	tasks.SortByPriority()
	for _, task := range tasks.Tasks {
		fmt.Printf("ID: %d, Name: %s, Priority: %d, Status: %s, Deadline: %s, CreatedAt: %s\n", task.ID, task.Name, task.Priority, task.Status, task.Deadline, task.CreatedAt)
	}

	fmt.Println("\nSorting by deadline:")
	tasks.SortByDeadline()
	for _, task := range tasks.Tasks {
		fmt.Printf("ID: %d, Name: %s, Priority: %d, Status: %s, Deadline: %s, CreatedAt: %s\n", task.ID, task.Name, task.Priority, task.Status, task.Deadline, task.CreatedAt)
	}

	fmt.Println("\nSorting by creation time:")
	tasks.SortByCreatedAt()
	for _, task := range tasks.Tasks {
		fmt.Printf("ID: %d, Name: %s, Priority: %d, Status: %s, Deadline: %s, CreatedAt: %s\n", task.ID, task.Name, task.Priority, task.Status, task.Deadline, task.CreatedAt)
	}

	jsonData, err := json.Marshal(tasks)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	} else {
		fmt.Println("\nJSON representation:")
		fmt.Println(string(jsonData))
	}

}
