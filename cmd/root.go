// This file is part of Compact-Chain.
// Copyright (c) 2023 0xSharma
// Licensed under the MIT License (MIT).
// See the LICENSE file in the project root for license information.
package cmd // root.go

// Package cmd implements the root command for the Compact-Chain CLI.
// It provides commands to start the node, send transactions, and display the version.

import (
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	//	"github.com/0xsharma/compact-chain/cmd/sendtx"

	"github.com/0xsharma/compact-chain/config"
	"github.com/0xsharma/compact-chain/core"
	"github.com/0xsharma/compact-chain/types"
	"github.com/0xsharma/compact-chain/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "v1.1.0" // Version of the Compact-Chain CLI
	// rootCmd represents the base command when called without any subcommands.
	// It is used to define the root command for the CLI application.
	// It can be used to add global flags or configuration settings.
	rootCmd = &cobra.Command{} // rootCmd represents the base command when called without any subcommands.
	// versionCmd represents the command to print the version number of Compact-Chain.

	versionCmd = &cobra.Command{ // versionCmd represents the command to print the version number of Compact-Chain.
		// Use:   "version",` // Use is the command name used in the CLI to invoke this command.`
		Use: "version", // Use is the command name used in the CLI to invoke this command.
		// Short: "Print the version number of Compact-Chain", // Short is a brief description of the command.
		// Run: func(cmd *cobra.Command, args []string) { // Run is the function that gets executed when this command is called.
		Short: "Print the version number of Compact-Chain", //	 Short is a brief description of the command.
		Run: func(cmd *cobra.Command, args []string) { // Run is the function that gets executed when this command is called.
			fmt.Println(version) // Print the version number of Compact-Chain
		},
	}

	startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the Compact-Chain node",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Starting Compact-Chain node\n\n")
			nodeID, _ := strconv.ParseInt(args[0], 10, 0)
			startBlockchainNode(nodeID)
		},
	}

	sendTxCmd = &cobra.Command{
		Use:   "send-tx",
		Short: "Send a transaction to the Compact-Chain node",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Sending transaction to Compact-Chain node\n\n")

			flags := cmd.Flags()

			to, _ := flags.GetString("to")
			value, _ := flags.GetInt64("value")
			privateKey, _ := flags.GetString("privatekey")
			nonce, _ := flags.GetInt64("nonce")
			rpcAddr, _ := flags.GetString("rpc")

			sendTxCfg := &sendTxConfig{
				To:         to,
				Value:      value,
				PrivateKey: privateKey,
				Nonce:      nonce,
				RPCAddr:    rpcAddr,
			}

			SendTx(sendTxCfg)
		},
	}

	demoCmd = &cobra.Command{
		Use:   "demo",
		Short: "Demo the Compact-Chain node",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Starting Compact-Chain node\n\n")
			demoBlockchain()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

// nolint : errcheck
func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(demoCmd)
	rootCmd.AddCommand(sendTxCmd)

	sendTxCmd.PersistentFlags().String("to", "", "To Address")
	viper.BindPFlag("to", sendTxCmd.PersistentFlags().Lookup("to"))
	cobra.MarkFlagRequired(sendTxCmd.PersistentFlags(), "to")

	sendTxCmd.PersistentFlags().Int64("value", 0, "Value to send")
	viper.BindPFlag("value", sendTxCmd.PersistentFlags().Lookup("value"))
	cobra.MarkFlagRequired(sendTxCmd.PersistentFlags(), "value")

	sendTxCmd.PersistentFlags().String("privatekey", "", "Private key to sign transaction")
	viper.BindPFlag("privatekey", sendTxCmd.PersistentFlags().Lookup("privatekey"))
	cobra.MarkFlagRequired(sendTxCmd.PersistentFlags(), "privatekey")

	sendTxCmd.PersistentFlags().Int64("nonce", 0, "Nonce of transaction")
	viper.BindPFlag("nonce", sendTxCmd.PersistentFlags().Lookup("nonce"))
	cobra.MarkFlagRequired(sendTxCmd.PersistentFlags(), "nonce")

	sendTxCmd.PersistentFlags().String("rpc", "", "RPC endpoint of node")
	viper.BindPFlag("rpc", sendTxCmd.PersistentFlags().Lookup("rpc"))
	cobra.MarkFlagRequired(sendTxCmd.PersistentFlags(), "rpc")
}

var (
	homePath, _ = os.UserHomeDir()
	dbPath      = homePath + "/.compact-chain/db"
	stateDbPath = homePath + "/.compact-chain/statedb"
)

func demoBlockchain() {
	config := &config.Config{
		ConsensusDifficulty: 16,
		ConsensusName:       "pow",
		DBDir:               dbPath + "demo",
		StateDBDir:          stateDbPath + "demo",
		MinFee:              big.NewInt(100),
		RPCPort:             ":1711",
		BalanceAlloc:        map[string]*big.Int{},
		P2PPort:             ":6060",
		Peers:               []string{"localhost:6061"},
		BlockTime:           2,
		SignerPrivateKey:    util.HexToPrivateKey("c3fc038a9abc0f483e2e1f8a0b4db676bce3eaebd7d9afc68e1e7e28ca8738a1"),
	}

	chain := core.NewBlockchain(config)
	if chain.LastBlock.Number.Int64() == 0 {
		fmt.Println("Number : ", chain.LastBlock.Number, "Hash : ", chain.LastBlock.DeriveHash().String())
	} else {
		fmt.Println("LastNumber : ", chain.LastBlock.Number, "LastHash : ", chain.LastBlock.DeriveHash().String())
	}

	lastNumber := chain.LastBlock.Number

	for i := lastNumber.Int64() + 1; i <= lastNumber.Int64()+10; i++ {
		time.Sleep(2 * time.Second)

		err := chain.AddBlock([]byte(fmt.Sprintf("Block %d", i)), []*types.Transaction{}, make(chan bool), config.SignerPrivateKey)
		if err != nil {
			fmt.Println("Error Adding Block", err)
		}

		fmt.Println("Number : ", chain.LastBlock.Number, "Hash : ", chain.LastBlock.DeriveHash().String())
	}
}

func startBlockchainNode(nodeId int64) {
	fmt.Println("Starting node", nodeId)

	config := &config.Config{
		ConsensusDifficulty: 20,
		ConsensusName:       "pow",
		DBDir:               dbPath + fmt.Sprint(nodeId),
		StateDBDir:          stateDbPath + fmt.Sprint(nodeId),
		MinFee:              big.NewInt(100),
		RPCPort:             ":1711" + fmt.Sprint(nodeId),
		BalanceAlloc: map[string]*big.Int{
			"0xa52c981eee8687b5e4afd69aa5006548c24d7685": big.NewInt(1000000000000000000), // Allocating funds to 0xa52c981eee8687b5e4afd69aa5006548c24d7685
		},
		P2PPort:          ":6060" + fmt.Sprint(nodeId),
		Peers:            []string{"localhost:60601", "localhost:60602", "localhost:60603"},
		BlockTime:        4,
		SignerPrivateKey: util.HexToPrivateKey("c3fc038a9abc0f483e2e1f8a0b4db676bce3eaebd7d9afc68e1e7e28ca8738a" + fmt.Sprint(nodeId)),
		Mine:             true,
	}

	core.StartBlockchain(config)
}
