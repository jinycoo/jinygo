/**------------------------------------------------------------**
 * @filename commands/commands.go
 * @author   jinycoo - caojingyin@jinycoo.com
 * @version  1.0.0
 * @date     2018/08/10 15:27
 * @desc     commands - add commands
 **------------------------------------------------------------**/
package commands

import "github.com/spf13/cobra"

func AddCommands(cmd *cobra.Command) {
	addVersion(cmd)
	addCreate(cmd)
}
