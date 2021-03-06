package main

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/thrasher-corp/gocryptotrader/common"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/engine"
	exchange "github.com/thrasher-corp/gocryptotrader/exchanges"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
)

const (
	totalWrappers = 20
)

func main() {
	var err error
	engine.Bot, err = engine.New()
	if err != nil {
		log.Fatalf("Failed to initialise engine. Err: %s", err)
	}

	engine.Bot.Settings = engine.Settings{
		DisableExchangeAutoPairUpdates: true,
	}

	log.Printf("Loading exchanges..")
	var wg sync.WaitGroup
	for x := range exchange.Exchanges {
		err := engine.Bot.LoadExchange(exchange.Exchanges[x], true, &wg)
		if err != nil {
			log.Printf("Failed to load exchange %s. Err: %s",
				exchange.Exchanges[x],
				err)
			continue
		}
	}
	wg.Wait()
	log.Println("Done.")

	log.Printf("Testing exchange wrappers..")
	results := make(map[string][]string)
	wg = sync.WaitGroup{}
	exchanges := engine.Bot.GetExchanges()
	for x := range exchanges {
		exch := exchanges[x]
		wg.Add(1)
		go func(e exchange.IBotExchange) {
			results[e.GetName()] = testWrappers(e)
			wg.Done()
		}(exch)
	}
	wg.Wait()
	log.Println("Done.")

	log.Println()
	for name, funcs := range results {
		pct := float64(totalWrappers-len(funcs)) / float64(totalWrappers) * 100
		log.Printf("Exchange %s wrapper coverage [%d/%d - %.2f%%] | Total missing: %d", name, totalWrappers-len(funcs), totalWrappers, pct, len(funcs))
		log.Printf("\t Wrappers not implemented:")

		for x := range funcs {
			log.Printf("\t - %s", funcs[x])
		}
		log.Println()
	}
}

func testWrappers(e exchange.IBotExchange) []string {
	p := currency.NewPair(currency.BTC, currency.USD)
	assetType := asset.Spot
	if !e.SupportsAsset(assetType) {
		assets := e.GetAssetTypes()
		rand.Seed(time.Now().Unix())
		assetType = assets[rand.Intn(len(assets))] // nolint:gosec // basic number generation required, no need for crypo/rand
	}

	var funcs []string

	_, err := e.FetchTicker(p, assetType)
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "FetchTicker")
	}

	_, err = e.UpdateTicker(p, assetType)
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "UpdateTicker")
	}

	_, err = e.FetchOrderbook(p, assetType)
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "FetchOrderbook")
	}

	_, err = e.UpdateOrderbook(p, assetType)
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "UpdateOrderbook")
	}

	_, err = e.FetchTradablePairs(asset.Spot)
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "FetchTradablePairs")
	}

	err = e.UpdateTradablePairs(false)
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "UpdateTradablePairs")
	}

	_, err = e.FetchAccountInfo()
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "GetAccountInfo")
	}

	_, err = e.GetExchangeHistory(p, assetType, time.Time{}, time.Time{})
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "GetExchangeHistory")
	}

	_, err = e.GetFundingHistory()
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "GetFundingHistory")
	}

	_, err = e.SubmitOrder(nil)
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "SubmitOrder")
	}

	_, err = e.ModifyOrder(nil)
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "ModifyOrder")
	}

	err = e.CancelOrder(nil)
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "CancelOrder")
	}

	_, err = e.CancelAllOrders(nil)
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "CancelAllOrders")
	}

	_, err = e.GetOrderInfo("1", p, assetType)
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "GetOrderInfo")
	}

	_, err = e.GetOrderHistory(nil)
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "GetOrderHistory")
	}

	_, err = e.GetActiveOrders(nil)
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "GetActiveOrders")
	}

	_, err = e.GetDepositAddress(currency.BTC, "")
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "GetDepositAddress")
	}

	_, err = e.WithdrawCryptocurrencyFunds(nil)
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "WithdrawCryptocurrencyFunds")
	}

	_, err = e.WithdrawFiatFunds(nil)
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "WithdrawFiatFunds")
	}
	_, err = e.WithdrawFiatFundsToInternationalBank(nil)
	if err == common.ErrNotYetImplemented {
		funcs = append(funcs, "WithdrawFiatFundsToInternationalBank")
	}

	return funcs
}
