package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func handleCommand(db *sql.DB) {
	// check for the command length
	if len(os.Args) < 2 {
		handleError()
	}

	if os.Args[1] == "-help" {
		handleError()
	}

	// add command
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listAllCmd := listCmd.Bool("a", false, "list all tasks")

	completeCmd := flag.NewFlagSet("complete", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	// Use swich on Args[1] which holds the command name
	switch os.Args[1] {
	case "add":
		addTaskHandler(addCmd, db)
	case "list":
		listTaskHandler(listCmd, listAllCmd, db)
	case "complete":
		completeTaskHandler(completeCmd, db)
	case "delete":
		deleteTaskHandler(deleteCmd, db)
	default:
		handleError()
	}

}
func handleError() {
	if os.Args[1] == "-help" {
		printHelp()
		os.Exit(0)
	}
	fmt.Fprintln(os.Stderr, "Use [tasks -help] for more information.")
	os.Exit(1)
}
func addTaskHandler(addCmd *flag.FlagSet, db *sql.DB) {
	addCmd.Parse(os.Args[2:])
	if addCmd.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: tasks add <taskName>")
		os.Exit(1)
	}
	taskName := addCmd.Arg(0)
	fmt.Fprintln(os.Stdout, "Enter task description separated by '.'")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	description := strings.ReplaceAll(input, ". ", ".\n")
	task := Task{
		Name:        taskName,
		Status:      false,
		CreatedAt:   time.Now(),
		Description: description,
	}
	result, err := addTask(db, task)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error adding task: %v\n", err)
		os.Exit(1)
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting last insert id: %v\n", err)
		os.Exit(1)
	}
	task, err = getTask(db, int(id))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting task: %v\n", err)
		os.Exit(1)
	}
	printTasks([]Task{task})
}
func listTaskHandler(listCmd *flag.FlagSet, listAllCmd *bool, db *sql.DB) {
	listCmd.Parse(os.Args[2:])

	if *listAllCmd {
		tasks, err := getTasks(db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting tasks: %v\n", err)
			os.Exit(1)
		}
		printTasks(tasks)
	} else if len(os.Args) > 2 && os.Args[2] != "" {
		idStr := os.Args[2]
		taskID, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "task id should be integer value.")
			os.Exit(1)
		}
		task, err := getTask(db, taskID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting task: %v\n", err)
			os.Exit(1)
		}
		printTask(task)
	} else {
		fmt.Fprintln(os.Stderr, "Please refer to the help section for usage.")
		os.Exit(1)
	}
}

func completeTaskHandler(completeCmd *flag.FlagSet, db *sql.DB) {
	completeCmd.Parse(os.Args[2:])
	if completeCmd.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: tasks complete <id>")
		os.Exit(1)
	}
	id := completeCmd.Arg(0)
	taskID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting id to int: %v\n", err)
		os.Exit(1)
	}
	_, err = completeTask(db, taskID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error completing task: %v\n", err)
		os.Exit(1)
	}
	task, err := getTask(db, taskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting task: %v\n", err)
		os.Exit(1)
	}
	printTasks([]Task{task})
}

func deleteTaskHandler(deleteCmd *flag.FlagSet, db *sql.DB) {
	deleteCmd.Parse(os.Args[2:])
	if deleteCmd.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: tasks delete <id>")
		os.Exit(1)
	}

	id := deleteCmd.Arg(0)
	taskID, err := strconv.Atoi(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting id to int: %v\n", err)
		os.Exit(1)
	}
	_, err = deleteTask(db, taskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error completing task: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Task deleted successfully")

}
