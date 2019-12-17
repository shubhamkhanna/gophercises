package db

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mitchellh/go-homedir"
)

var home, _ = homedir.Dir()
var dbPath = filepath.Join(home, "tasks.db")

func TestInit(t *testing.T) {
	t.Run("It checks valid scenario.", func(t *testing.T) {
		Init(dbPath)
		assert.NotEqual(t, db.Path(), nil)
		defer func() {
			CloseDB()
		}()
	})

	t.Run("It throws an error.", func(t *testing.T) {
		err := Init("")
		assert.NotEqual(t, err, nil)
	})
}

func TestCreateTask(t *testing.T) {
	Init(dbPath)
	t.Run("It creates new task.", func(t *testing.T) {
		CreateTask("my new task")
		tasks, _ := AllTasks()
		task := tasks[len(tasks)-1]
		assert.Equal(t, task.Value, "my new task")
		DeleteTask(task.Key)
	})

	t.Run("It checks an error condition.", func(t *testing.T) {
		CloseDB()
		_, err := CreateTask("my new task")
		assert.NotEqual(t, err, nil)
	})
	defer func() {
		CloseDB()
	}()
}

func TestAllTask(t *testing.T) {
	Init(dbPath)
	t.Run("It returns all tasks.", func(t *testing.T) {
		CreateTask("my new task")
		tasks, _ := AllTasks()
		assert.NotEqual(t, len(tasks), 0)
		DeleteTask(tasks[len(tasks)-1].Key)
	})

	t.Run("It checks an error condition.", func(t *testing.T) {
		CloseDB()
		_, err := AllTasks()
		assert.NotEqual(t, err, nil)
	})
	defer func() {
		CloseDB()
	}()
}

func TestDelteTask(t *testing.T) {
	Init(dbPath)
	t.Run("It returns all tasks.", func(t *testing.T) {
		CreateTask("my new task")
		tasks, _ := AllTasks()
		beforeTaskCounts := len(tasks)
		DeleteTask(tasks[len(tasks)-1].Key)
		tasks, _ = AllTasks()
		afterTaskCounts := len(tasks)
		assert.NotEqual(t, beforeTaskCounts, afterTaskCounts)
	})

	t.Run("It checks an error condition.", func(t *testing.T) {
		CloseDB()
		err := DeleteTask(1)
		assert.NotEqual(t, err, nil)
	})
	defer func() {
		CloseDB()
	}()
}
