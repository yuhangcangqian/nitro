// Copyright 2025, Offchain Labs, Inc.
// For license information, see:
// https://github.com/OffchainLabs/nitro/blob/master/LICENSE.md

package arbtest

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/core/rawdb"

	"github.com/offchainlabs/nitro/util/arbmath"
	"github.com/offchainlabs/nitro/validator/valnode"
)

var withL1 = true

func TestStorageTrie(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	builder := NewNodeBuilder(ctx).DefaultConfig(t, withL1)

	// retryableSetup is being called by tests that validate blocks.
	// For now validation only works with HashScheme set.
	builder.execConfig.Caching.StateScheme = rawdb.HashScheme
	builder.nodeConfig.BlockValidator.Enable = false
	builder.nodeConfig.Staker.Enable = true
	builder.nodeConfig.BatchPoster.Enable = true
	builder.nodeConfig.ParentChainReader.Enable = true
	builder.nodeConfig.ParentChainReader.OldHeaderTimeout = 10 * time.Minute

	valConf := valnode.TestValidationConfig
	valConf.UseJit = true
	_, valStack := createTestValidationNode(t, ctx, &valConf)
	configByValidationNode(builder.nodeConfig, valStack)

	cleanup := builder.Build(t)
	defer cleanup()

	ownerTxOpts := builder.L2Info.GetDefaultTransactOpts("Owner", ctx)
	_, bigMap := builder.L2.DeployBigMap(t, ownerTxOpts)

	// Store enough values to use just over 32M gas
	values := big.NewInt(1431)

	userTxOpts := builder.L2Info.GetDefaultTransactOpts("Faucet", ctx)
	tx, err := bigMap.StoreValues(&userTxOpts, values)
	Require(t, err)

	receipt, err := builder.L2.EnsureTxSucceeded(tx)
	Require(t, err)

	if receipt.GasUsed != 32_002_907 {
		t.Errorf("Want GasUsed: 32002907: got: %d", receipt.GasUsed)
	}

	// Clear about 75% of them, and add another 10%
	toClear := arbmath.BigDiv(arbmath.BigMul(values, big.NewInt(75)), big.NewInt(100))
	toAdd := arbmath.BigDiv(arbmath.BigMul(values, big.NewInt(10)), big.NewInt(100))

	tx, err = bigMap.ClearAndAddValues(&userTxOpts, toClear, toAdd)
	Require(t, err)

	receipt, err = builder.L2.EnsureTxSucceeded(tx)
	Require(t, err)

	// Ensures that the validator gets the same results as the executor
	validateBlockRange(t, []uint64{receipt.BlockNumber.Uint64()}, true, builder)
}
