package tasks

// Task is made to run in the background. It holds a function called Update
type Task struct {
	Update func()
}
