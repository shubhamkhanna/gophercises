package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {

	t.Run("it checks valid File", func(t *testing.T) {
		err := Initialize("problems.csv", 10)
		assert.Equal(t, err, nil)
	})

	t.Run("it checks invalid File", func(t *testing.T) {
		err := Initialize("problems1.csv", 10)
		assert.NotEqual(t, err, nil)
	})

	t.Run("it checks invalid File content", func(t *testing.T) {
		err := Initialize("bad_file.csv", 10)
		assert.NotEqual(t, err, nil)
	})

	t.Run("it tests Correct result", func(t *testing.T) {
		problems := []problem{
			{"5 + 5", "10"},
		}

		content := []byte("10")
		tmpfile, _ := ioutil.TempFile("", "example")

		defer os.Remove(tmpfile.Name())

		if _, err := tmpfile.Write(content); err != nil {
			log.Fatal(err)
		}

		if _, err := tmpfile.Seek(0, 0); err != nil {
			log.Fatal(err)
		}
		os.Stdin = tmpfile
		Initialize("problems.csv", 10)
		if correct != 1 {
			t.Errorf("You have scored %d out of %d\n", correct, len(problems))
		}
	})

	t.Run("it checks time up case", func(t *testing.T) {
		correct = 0
		defer Initialize("problems.csv", 0)
		if correct != 0 {
			t.Errorf("Time up test failed!")
		}
	})

}
