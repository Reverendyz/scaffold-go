package hello

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "hello",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(strings.Join(args, " "))
	},
}
