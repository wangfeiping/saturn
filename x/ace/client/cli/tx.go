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

	aceTxCmd.AddCommand(flags.PostCommands(GetCmdTxPlay(cdc))...)

	return aceTxCmd
}

// GetCmdTxPlay is the CLI command for sending play Tx
/*
 * $ ./saturncli tx ace play LuckyAce draw "a seed" \
  cosmos1ah0pfkn6kwuj24ks2uqqu8m387ffwnm7fvsugm --chain-id saturn-testnet-0
*/
func GetCmdTxPlay(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "play [game-ace-id] [func] [seed] [address]",
		Short: "send a Tx(command of one-step-play for a game)",
		Args:  cobra.ExactArgs(4), // Does your request require arguments
		RunE: func(cmd *cobra.Command, args []string) error {
			aceID := args[0]
			function := args[1]
			seed := types.Seed{Hash: []byte(args[2])}
			address := args[3]
			fmt.Printf("play: %s %s\n", aceID, function)

			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).
				WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.
				NewCLIContextWithInputAndFrom(
					inBuf, address).WithCodec(cdc)

			// query game id
			// res, _, err := cliCtx.QueryWithData(
			// 	fmt.Sprintf("custom/%s/%s",
			// 		types.QuerierRoute, types.QueryGames), []byte(aceID))
			// if err != nil {
			// 	fmt.Printf("query game error: %v\n", err.Error())
			// 	return nil
			// }
			// var out types.Game
			// cdc.MustUnmarshalJSON(res, &out)
			// fmt.Println(out.Info)
			// fmt.Printf("Game       : %s\t%s\t%d\n", out.AceID, out.Type, out.GameID)
			// fmt.Printf("IsGroupGame: %t\n", out.IsGroupGame)

			var err error
			var out types.Game
			out.GameID = 100
			// create, sign and send play Tx
			// msg, err := types.NewMsgAce(cliCtx.GetFromAddress())
			msg := types.NewMsgPlay(
				aceID, out.GameID, 0,
				seed, "draw", "arg1,arg2,arg3",
				cliCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			err = utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			if err != nil {
				fmt.Printf("send play Tx error: %v\n", err)
				return err
			}
			fmt.Printf("Ok!")
			return err
		},
	}
}
