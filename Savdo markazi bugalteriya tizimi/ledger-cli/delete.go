package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/darcys22/godbledger/proto/transaction"

	"github.com/darcys22/godbledger/godbledger/cmd"

	"google.golang.org/grpc"

	"github.com/urfave/cli/v2"
)

var commandDeleteTransaction = &cli.Command{
	Name:      "delete",
	Usage:     "ledger-cli delete <transaction_id>",
	ArgsUsage: "[]",
	Description: `
	Deletes a transaction from the database
`,
	Flags: []cli.Flag{},
	Action: func(ctx *cli.Context) error {
		err, cfg := cmd.MakeConfig(ctx)
		if err != nil {
			return fmt.Errorf("Could not make config (%v)", err)
		}

		if ctx.NArg() > 0 {
			address := fmt.Sprintf("%s:%s", cfg.Host, cfg.RPCPort)
			log.WithField("address", address).Info("GRPC Dialing on port")
			opts := []grpc.DialOption{}

			if cfg.CACert != "" && cfg.Cert != "" && cfg.Key != "" {
				tlsCredentials, err := loadTLSCredentials(cfg)
				if err != nil {
					return fmt.Errorf("Could not load TLS credentials (%v)", err)
				}
				opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
			} else {
				opts = append(opts, grpc.WithInsecure())
			}

			// Set up a connection to the server.
			conn, err := grpc.Dial(address, opts...)
			if err != nil {
				return fmt.Errorf("Could not connect to GRPC (%v)", err)
			}
			defer conn.Close()
			client := transaction.NewTransactorClient(conn)

			ctxtimeout, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			req := &transaction.DeleteRequest{
				Identifier: ctx.Args().Get(0),
			}
			r, err := client.DeleteTransaction(ctxtimeout, req)
			if err != nil {
				return fmt.Errorf("Could not call Delete Transaction Method (%v)", err)
			}
			log.Infof("Delete Transaction Response: %s", r.GetMessage())
		} else {
			return errors.New("This command requires an argument")
		}

		return nil
	},
}

var commandVoidTransaction = &cli.Command{
	Name:      "void",
	Usage:     "ledger-cli void <transaction_id>",
	ArgsUsage: "[]",
	Description: `
	Reverses a transaction by creating an identical inverse and tags both transactions as void 
`,
	Flags: []cli.Flag{},
	Action: func(ctx *cli.Context) error {
		err, cfg := cmd.MakeConfig(ctx)
		if err != nil {
			return fmt.Errorf("Could not make config (%v)", err)
		}

		if ctx.NArg() > 0 {
			address := fmt.Sprintf("%s:%s", cfg.Host, cfg.RPCPort)
			log.WithField("address", address).Info("GRPC Dialing on port")
			opts := []grpc.DialOption{}

			if cfg.CACert != "" && cfg.Cert != "" && cfg.Key != "" {
				tlsCredentials, err := loadTLSCredentials(cfg)
				if err != nil {
					return fmt.Errorf("Could not load TLS credentials (%v)", err)
				}
				opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
			} else {
				opts = append(opts, grpc.WithInsecure())
			}

			// Set up a connection to the server.
			conn, err := grpc.Dial(address, opts...)
			if err != nil {
				return fmt.Errorf("Could not connect to GRPC (%v)", err)
			}
			defer conn.Close()
			client := transaction.NewTransactorClient(conn)

			ctxtimeout, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			req := &transaction.DeleteRequest{
				Identifier: ctx.Args().Get(0),
			}
			r, err := client.VoidTransaction(ctxtimeout, req)
			if err != nil {
				return fmt.Errorf("Could not call Void Transaction Method (%v)", err)
			}
			log.Infof("Void Transaction Response: %s", r.GetMessage())
		} else {
			return errors.New("This command requires an argument")
		}

		return nil
	},
}
