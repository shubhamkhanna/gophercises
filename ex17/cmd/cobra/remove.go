package cobra

import (
	"fmt"
	"gophercises/ex17/cipher"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := cipher.File(encodingKey, secretsPath())
		key := args[0]
		_, err := v.Get(key)
		if err == nil {
			err = v.Remove(key)
			if err != nil {
				return
			}
		} else {
			fmt.Printf("%v is not present!\n", key)
			return
		}
		fmt.Printf("Key : %v removed successfully!\n", key)
	},
}

func init() {
	RootCmd.AddCommand(removeCmd)
}
