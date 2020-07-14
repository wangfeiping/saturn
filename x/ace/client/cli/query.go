package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/wangfeiping/saturn/x/ace/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group ace queries under a subcommand
	aceQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	aceQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdQuerySecret(cdc),
			GetCmdQueryGames(cdc),
			GetCmdQueryRounds(cdc),
			GetCmdQueryWinners(cdc),
		)...,
	)

	return aceQueryCmd
}

// GetCmdQuerySecret queries secret public key
func GetCmdQuerySecret(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "secret",
		Short: "secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s",
					types.QuerierRoute, types.QuerySecret), nil)
			if err != nil {
				fmt.Printf("query secret error: %v\n", err.Error())
				return nil
			}

			var out types.Secret
			cdc.MustUnmarshalJSON(res, &out)

			home := viper.GetString(flags.FlagHome)
			if err := os.MkdirAll(home, os.ModePerm); err != nil {
				fmt.Printf("mkdir error: %v\n", err)
				return err
			}
			secretFile := path.Join(home, "ace_secret.json")
			if err = ioutil.WriteFile(secretFile, res, 0644); err != nil {
				fmt.Printf("save secret error: %v\n", err)
				return err
			}
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdQueryGames queries infomations of game
func GetCmdQueryGames(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "games",
		Short: "games",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			var data []byte
			if len(args) > 0 {
				data = []byte(args[0])
			}
			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s",
					types.QuerierRoute, types.QueryGames), data)
			if err != nil {
				fmt.Printf("query game error: %v\n", err.Error())
				return nil
			}
			// fmt.Println("result: " + string(res))
			if data == nil {
				var out []types.Game
				cdc.MustUnmarshalJSON(res, &out)
				// return cliCtx.PrintOutput(out)
				for _, g := range out {
					fmt.Printf("Game       : %s\t%s\t%s\n", g.AceID, g.Type, g.GameID)
					fmt.Printf("IsGroupGame: %t\n", g.IsGroupGame)
				}
				return nil
			}
			var out types.Game
			cdc.MustUnmarshalJSON(res, &out)
			fmt.Println(out.Info)
			fmt.Printf("Game       : %s\t%s\t%s\n", out.AceID, out.Type, out.GameID)
			fmt.Printf("IsGroupGame: %t\n", out.IsGroupGame)
			// return cliCtx.PrintOutput(out)
			return nil
		},
	}
}

// GetCmdQueryRounds queries secret public key
func GetCmdQueryRounds(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "rounds",
		Short: "rounds",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			var data []byte
			if len(args) > 0 {
				data = []byte(args[0])
			}
			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s",
					types.QuerierRoute, types.QueryRounds), data)
			if err != nil {
				fmt.Printf("query secret error: %v\n", err.Error())
				return nil
			}

			var out []types.Play
			cdc.MustUnmarshalJSON(res, &out)

			home := viper.GetString(flags.FlagHome)
			if err := os.MkdirAll(home, os.ModePerm); err != nil {
				fmt.Printf("mkdir error: %v\n", err)
				return err
			}
			secretFile := path.Join(home, "ace_secret.json")
			if err = ioutil.WriteFile(secretFile, res, 0644); err != nil {
				fmt.Printf("save secret error: %v\n", err)
				return err
			}
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdQueryWinners queries all winners about a game
func GetCmdQueryWinners(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "winners",
		Short: "winners",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s",
					types.QuerierRoute, types.QueryWinners), nil)
			if err != nil {
				fmt.Printf("query winners error: %v\n", err.Error())
				return nil
			}
			fmt.Printf("query winners: %s", string(res))

			// var out types.Round
			// cdc.MustUnmarshalJSON(res, &out)
			// return cliCtx.PrintOutput(out)
			return nil
		},
	}
}
