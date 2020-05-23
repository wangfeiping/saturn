package cli

import (
	"bufio"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/wangfeiping/saturn/x/ace/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	aceTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	aceTxCmd.AddCommand(flags.PostCommands(
		GetCmdTxPlay(cdc),
	)...)

	return aceTxCmd
}

// GetCmdTxPlay is the CLI command for sending play Tx
func GetCmdTxPlay(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "play [func] [arg1,arg2,...] [seed] [address]",
		Short: "send a Tx(request for a one-step-play command)",
		Args:  cobra.ExactArgs(4), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			function := args[0]
			argsStr := args[1]
			seed := &types.Seed{Seed: []byte(args[2])}
			address := args[3]
			fmt.Printf("play: %s %s\n", function, argsStr)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).
				WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.
				NewCLIContextWithInputAndFrom(
					inBuf, address).WithCodec(cdc)

			// msg, err := types.NewMsgAce(cliCtx.GetFromAddress())
			msg, err := types.NewMsgPlay(
				"LuckyAce", "LuckyAce-30", "",
				seed, "draw", argsStr,
				cliCtx.GetFromAddress())

			if err != nil {
				fmt.Printf("new msg error: %v\n", err)
				return err
			}
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
