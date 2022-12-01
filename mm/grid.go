package grid

import (
	"strconv"
	"strings"
	"sync"
	"time"
)

type GridTrading struct {
	Base *strategy.StrategyBase
}

func (g GridTrading) Run(params []byte) error {
	panic("implement me")
}

func Grid(c exchange.Exchange, log log.Logger, db *mongo.MongoDB, symbol string, maxU, timeLimit, BalanceMin, BalanceMax float64) {
	exit := false
	token := strings.Split(symbol, "_")
	var wg sync.WaitGroup
	var Bid1, Ask1 float64
	var LastRound []string
	PairInfo, err := c.GetPairInfo(symbol)
	if err != nil {
		log.Error("GetPair", err)
	}
	decimal, _ := strconv.ParseFloat(PairInfo.MinBaseAmount, 64)
	wg.Add(1)
	go func() {
		for {
			if exit == true {
				wg.Done()
			}
		OrderLoop:
			depth, err := c.Depth(symbol, "5")
			if err != nil {
				log.Error("Depth", err)
			}

			bid := depth.Bids[0].Price
			ask := depth.Asks[0].Price
			Bid1, _ = strconv.ParseFloat(bid, 64)
			Ask1, _ = strconv.ParseFloat(ask, 64)

			ol := []models.OrderList{}
			ol = append(ol, models.OrderList{
				Side:  base.BID,
				Price: strconv.FormatFloat(Bid1+decimal, 'f', 7, 64),
				Size:  strconv.FormatFloat(tool.RandFloat(5, maxU)/(Bid1+decimal), 'f', 7, 64),
			}, models.OrderList{
				Side:  base.ASK,
				Price: strconv.FormatFloat(Ask1-decimal, 'f', 7, 64),
				Size:  strconv.FormatFloat(tool.RandFloat(5, maxU)/(Ask1-decimal), 'f', 7, 64),
			})
			orders, err := c.LimitOrders(symbol, ol)
			if err != nil {
				log.Error("Orders", err)
			}
			if LastRound != nil {
				cancel, err := c.CancelOrder(symbol, LastRound[0])
				if err != nil {
					log.Error("Cancel Order", err)
				}
				log.Info("Cancel", cancel)
				cancel, err = c.CancelOrder(symbol, LastRound[1])
				if err != nil {
					log.Error("Cancel Order", err)
				}
				log.Info("Cancel", cancel)
			}
			time.Sleep(time.Duration(timeLimit) * time.Second)

			info1, err := c.GetOrder(symbol, orders[0])
			if err != nil {
				log.Error("GetOrder", err)
			}
			info2, err := c.GetOrder(symbol, orders[1])
			if err != nil {
				log.Error("GetOrder", err)
			}
			if info1.Status == base.FILLED && info2.Status == base.FILLED {
				LastRound = nil
				goto OrderLoop

			} else if info1.Status == base.OPEN && info2.Status == base.OPEN {
				LastRound = orders
				goto OrderLoop
			}
			time.Sleep(60 * time.Second)
			LastRound = orders
			goto OrderLoop

		}
	}()

	wg.Add(1)
	go func() {
		for {
			if exit == true {
				wg.Done()
			}
			time.Sleep(300 * time.Millisecond)

			currentUSDT, err := c.GetAccountBalance(token[1])
			if err != nil {
				log.Error("[GetAccountBalance] err is: ", err)
			}

			USDTExit, _ := strconv.ParseFloat(currentUSDT, 64)

			if USDTExit < float64(BalanceMin) || USDTExit > float64(BalanceMax) {
				exit = true
			}
		}
	}()
	wg.Wait()

}
