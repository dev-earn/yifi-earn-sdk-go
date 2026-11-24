package main

import "os"

type cycleConfig struct {
	strategyID  string
	inputToken  string
	inputAmount string
	engine      string
	chain       string
	chainID     int64
}

func cycleConfigForRun() cycleConfig {
	cfg := cycleConfig{
		strategyID:  cycleStrategyID,
		inputToken:  cycleInputToken,
		inputAmount: cycleInputAmount,
		engine:      "lifi",
		chain:       "bsc",
		chainID:     56,
	}
	if os.Getenv("YIFI_CYCLE_ENGINE") == "portals" {
		cfg.strategyID = "bsc:0xa9251ca9de909cb71783723713b21e4233fbf1b1"
		cfg.inputAmount = "5000000000000000"
		cfg.engine = "portals"
	}
	if value := os.Getenv("YIFI_CYCLE_STRATEGY_ID"); value != "" {
		cfg.strategyID = value
	}
	if value := os.Getenv("YIFI_CYCLE_INPUT_TOKEN"); value != "" {
		cfg.inputToken = value
	}
	if value := os.Getenv("YIFI_CYCLE_CHAIN"); value != "" {
		cfg.chain = value
		if value == "ethereum" {
			cfg.chainID = 1
		}
	}
	if value := os.Getenv("YIFI_CYCLE_INPUT_AMOUNT"); value != "" {
		cfg.inputAmount = value
	}
	return cfg
}
