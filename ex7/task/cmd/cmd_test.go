package cmd

import (
	"gophercises/ex7/task/db"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var home, _ = homedir.Dir()
var dbPath = filepath.Join(home, "tasks.db")

func TestAddCmd(t *testing.T) {
	db.Init(dbPath)
	var myCmd *cobra.Command

	addCmd.Run(myCmd, []string{"task1"})
	// Covers error conditon when DB is closed.
	db.CloseDB()
	addCmd.Run(myCmd, []string{""})
	assert.Equal(t, addCmd.Use, "add")
	assert.Equal(t, addCmd.Short, "Adds a task to your task list.")

	defer func() {
		db.CloseDB()
	}()
}

func TestListCmd(t *testing.T) {
	db.Init(dbPath)
	var myCmd *cobra.Command
	listCmd.Run(myCmd, []string{"task1"})

	t.Run("It checks no tasks condition", func(t *testing.T) {
		doCmd.Run(myCmd, []string{"all"})
		listCmd.Run(myCmd, []string{""})
	})

	t.Run("Checks error conditon when DB is closed.", func(t *testing.T) {
		db.CloseDB()
		listCmd.Run(myCmd, []string{""})
	})
	
	assert.Equal(t, listCmd.Use, "list")
	assert.Equal(t, listCmd.Short, "Lists all of your tasks.")

	defer func() {
		db.CloseDB()
	}()
}

func TestDoCmd(t *testing.T) {
	db.CloseDB()
	db.Init(dbPath)
	var myCmd *cobra.Command
	args := []string{"task 1  task 2  task 3"}
	task := strings.Join(args, " ")
	db.CreateTask(task)

	t.Run("It marks specific task as done", func(t *testing.T) {
		doCmd.Run(myCmd, []string{"1"})
	})

	t.Run("It throws an error for in valid task", func(t *testing.T) {
		doCmd.Run(myCmd, []string{"1-2"})
	})

	t.Run("It marks all tasks as done", func(t *testing.T) {
		doCmd.Run(myCmd, []string{"all"})
	})

	t.Run("It checks condition id should be > 0", func(t *testing.T) {
		doCmd.Run(myCmd, []string{"0"})
	})

	t.Run("Checks error conditon when DB is closed.", func(t *testing.T) {
		db.CloseDB()
		doCmd.Run(myCmd, []string{"1"})
	})

	assert.Equal(t, doCmd.Use, "do")
	assert.Equal(t, doCmd.Short, "Marks a task as complete")

	defer func() {
		db.CloseDB()
	}()
}
