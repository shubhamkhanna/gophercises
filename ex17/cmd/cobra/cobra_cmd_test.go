package cobra

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestGetCmd(t *testing.T) {
	var myCmd *cobra.Command
	// Cover error conditon when key is not set.
	getCmd.Run(myCmd, []string{"key100"})

	// Happy flow when key is set.
	setCmd.Run(myCmd, []string{"key1", "value"})
	getCmd.Run(myCmd, []string{"key1"})
	assert.Equal(t, getCmd.Use, "get")
	assert.Equal(t, getCmd.Short, "Gets a secret in your secret storage")
}

func TestSetCmd(t *testing.T) {
	var myCmd *cobra.Command
	// Happy flow when key is set.
	setCmd.Run(myCmd, []string{"key1", "value"})
	// Cover error conditon when key is not set.
	tmpKey := encodingKey
	encodingKey = "/invalid"
	setCmd.Run(myCmd, []string{"k", "v"})
	assert.Equal(t, setCmd.Use, "set")
	assert.Equal(t, setCmd.Short, "Sets a secret in your secret storage")
	defer func() {
		encodingKey = tmpKey
	}()
}

func TestRemoveCmd(t *testing.T) {
	var myCmd *cobra.Command
	// Happy flow when key is set.
	removeCmd.Run(myCmd, []string{"key1"})
	// Cover error conditon when key is not set.
	removeCmd.Run(myCmd, []string{"k500"})
	assert.Equal(t, removeCmd.Use, "remove")
	assert.Equal(t, removeCmd.Short, "Removes a secret in your secret storage")
}
