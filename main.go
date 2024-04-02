package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/raydium-io/raydium-rpc/raydium"
	"github.com/raydium-io/raydium-rpc/raydium/types"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "solana-sniper",
		Usage: "Solana Sniper Bot for Raydium",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "rpc",
				Aliases:  []string{"r"},
				Required: true,
				Usage:    "RPC URL",
			},
			&cli.StringFlag{
				Name:    "wallet",
				Aliases: []string{"w"},
				Value:   os.Getenv("RAYDIUM_WALLET"),
				Usage:   "Wallet private key",
			},
			&cli.Int64Flag{
				Name:    "amount",
				Aliases: []string{"a"},
				Value:   int64(10000000000),
				Usage:   "Amount of SOL or USDC to buy",
			},
		},
		Action: func(c *cli.Context) error {
			rpc := c.String("rpc")
			wallet := c.String("wallet")
			amount := c.Int64("amount")

			client, err := raydium.NewClient(rpc)
			if err != nil {
				return fmt.Errorf("failed to create client: %v", err)
			}

			payer, err := types.Base58ToPublicKey(wallet)
			if err != nil {
				return fmt.Errorf("failed to parse payer wallet: %v", err)
			}

			for {
				newPools, err := getNewPools(client, payer)
				if err != nil {
					log.Printf("failed to get new pools: %v", err)
					continue
				}

				for _, pool := range newPools {
					amountBig := big.NewInt(amount)
					_, err = buy(client, payer, pool, amountBig)
					if err != nil {
						log.Printf("failed to buy tokens: %v", err)
						continue
					}

					log.Printf("bought %s tokens from pool %s", pool.Token.Mint, pool.Address)
				}

				time.Sleep(1 * time.Second)
			}
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getNewPools(client *raydium.Client, payer types.PublicKey) ([]types.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newPools, err := client.GetNewPools(ctx, payer)
	if err != nil {
		return nil, fmt.Errorf("failed to get new pools: %v", err)
	}

	return newPools, nil
}

func buy(client *raydium.Client,Client, payer types.PublicKey, pool types.Pool, amount *big.Int) (types.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := client.Buy(ctx, payer, pool.Address, amount, 1)
 payer types.PublicKey, pool types.Pool, amount *big.Int) (*types.TokenAccount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to buy tokens: %v", err)
	}

	return tx, nil
}