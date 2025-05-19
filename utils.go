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
			t.Description,
			strconv.FormatBool(t.Status),
			humanize.Time(t.CreatedAt), // e.g. “a minute ago”
		)
	}
	w.Flush()
}
