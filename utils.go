package main

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/dustin/go-humanize"
)

func printTasks(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTask\tCompleted\tCreated")
	for _, t := range tasks {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
			t.ID,
			t.Name,
			strconv.FormatBool(t.Status),
			humanize.Time(t.CreatedAt), // e.g. “a minute ago”
		)
	}
	w.Flush()
}

func printTask(task Task) {
	var zeroTask Task
	if task == zeroTask {
		fmt.Println("No tasks found")
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTask")
	fmt.Fprintf(w, "%d\t%s\t\n",
		task.ID,
		task.Name, // e.g. “a minute ago”
	)
	fmt.Fprintf(w, task.Description)
	w.Flush()
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("tasks add <task_name> : Add a new task")
	fmt.Println("tasks list [<id>] : Show task")
	fmt.Println("tasks list -a : List all tasks")
	fmt.Println("tasks complete <id> : Mark task as complete")
	fmt.Println("tasks delete <id> : Delete task")
	fmt.Println("tasks -help")
}
