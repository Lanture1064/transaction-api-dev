/*
Copyright 2023 The Bestchains Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"

	"github.com/bestchains/bc-explorer/pkg/errorsq"
	"github.com/bestchains/bc-explorer/pkg/viewer"
	"github.com/go-pg/pg/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"k8s.io/klog/v2"
)

var (
	db   = flag.String("db", "pg", "which database to use, default is pg(postgresql)")
	dsn  = flag.String("dsn", "postgres://bestchains:Passw0rd!@127.0.0.1:5432/bc-explorer?sslmode=disable", "database connection string")
	addr = flag.String("addr", ":9998", "used to listen and serve http requests")
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	if err := run(); err != nil {
		klog.Error(err)
	}
}

func run() error {
	pctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errq := errorsq.NewErrorsq(pctx, func(err error) {
		klog.Errorln("bc-explorer", err)
	})

	klog.Infoln("Starting a blockchain explorer viewer server")

	klog.Infoln("init db")
	block := viewer.NewBlockLoggerHandler()
	var transaction viewer.Transaction
	if *db == "pg" {
		klog.Infoln("Using postgreSQL")
		opts, err := pg.ParseURL(*dsn)
		if err != nil {
			return err
		}
		pgDB := pg.Connect(opts)
		defer pgDB.Close()
		if err := pgDB.Ping(pctx); err != nil {
			panic(err)
		}

		block = viewer.NewBlockHandler(pgDB)
		transaction = viewer.NewTxHandler(pgDB)
	}
	klog.Infoln("Creating http server")
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		Immutable:     true,
		AppName:       "bc-explorer-viewer",
	})

	viewerHandler := viewer.NewViewHandler(transaction, block)
	app.Use(cors.New(cors.ConfigDefault))
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	// TODO: register handlers
	// app.Get("/blocks", handler.List)
	app.Get("/networks/:network/blocks", viewerHandler.ListBlocks)
	app.Get("/networks/:network/blocks/:blockHash", viewerHandler.GetBlock)

	app.Get("/networks/:network/transactions", viewerHandler.ListTransactions)
	app.Get("/networks/:network/transactions/:txHash", viewerHandler.GetTransactionByTxHash)
	app.Get("/networks/:network/transactionsCount", viewerHandler.CountTransactionsCreatedByOrg)

	if err := app.Listen(*addr); err != nil {
		errq.Send(err)
	}
	return nil
}
