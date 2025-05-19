package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

func handleCommand(db *sql.DB) {
	// check for the command length
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: tasks <add | list | complete> [option]")
		os.Exit(1)
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
		addCmd.Parse(os.Args[2:])
		if addCmd.NArg() < 1 {
			fmt.Fprintln(os.Stderr, "usage: tasks add <description>")
			os.Exit(1)
		}
		description := addCmd.Arg(0)
		task := Task{
			Name:        description,
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
	case "list":
		listCmd.Parse(os.Args[2:])
		tasks, err := getTasks(db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting tasks: %v\n", err)
			os.Exit(1)
		}
		if *listAllCmd {
			printTasks(tasks)
		} else {
			printTasks(tasks)
		}
	case "complete":
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
	case "delete":
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
	default:
		fmt.Fprintln(os.Stderr, "Usage: tasks <add | list | complete> [option]")
		os.Exit(1)
	}
}
