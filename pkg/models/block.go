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

package models

const BlockTableName = "blocks"

type Block struct {
	BlockHash         string `pg:"blockHash,pk" json:"blockHash"`
	Network           string `pg:"network" json:"network"`
	BlockNumber       uint64 `pg:"blockNumber" json:"blockNumber"`
	PrevioudBlockHash string `pg:"preBlockHash" json:"preBlockHash"`
	DataHash          string `pg:"dataHash" json:"dataHash"`
	CreatedAt         int64  `pg:"createdAt" json:"createdAt"`
	BlockSize         int    `pg:"blockSize" json:"blockSize"`
	TxCount           int    `pg:"txCount" json:"txCount"`
}
