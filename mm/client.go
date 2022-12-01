package binance

import (
	"context"
	"errors"
	"github.com/adshao/go-binance/v2"
	"github.com/bitly/go-simplejson"
	"strconv"
	"time"
)

type Client struct {
	client *binance.Client
}

func (c *Client) New(params []byte) error {
	if c.client == nil {

		sj, err := simplejson.NewJson(params)
		if err != nil {
			return err
		}
		_ = sj.Get("url").MustString()
		apiKey := sj.Get("apiKey").MustString()
		secretKey := sj.Get("secretKey").MustString()
		_ = sj.Get("password").MustString()

		// init client by config
		binance.UseTestnet = false
		c.client = binance.NewClient(apiKey, secretKey)

		return nil
	}

	return errors.New("binance client has not been initialized")
}

func (c *Client) GetAccountBalance(currency string) (string, error) {

	account, err := c.client.
		NewGetAccountService().
		Do(context.Background())
	if err != nil {
		return "", err
	}

	for _, balance := range account.Balances {
		if balance.Asset == currency {
			a, _ := strconv.ParseFloat(balance.Free, 64)
			b, _ := strconv.ParseFloat(balance.Locked, 64)
			c := strconv.FormatFloat(a+b, 'f', 5, 64)
			return c, nil
		}
	}
	return "", nil
}

func (c *Client) MarketOrder(symbol, side, size string) (string, error) {
	var s binance.SideType
	if side == base.BID {
		s = binance.SideTypeBuy
	} else if side == base.ASK {
		s = binance.SideTypeSell
	}

	order, err := c.client.NewCreateOrderService().
		Symbol(symbol).
		Side(s).
		Type(binance.OrderTypeMarket).
		//	TimeInForce(binance.TimeInForceTypeGTC).
		Quantity(size).
		//	Price(price).
		Do(context.Background())
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(order.OrderID, 10), nil
}

func (c *Client) LimitOrder(symbol, side, price, size string) (string, error) {
	var s binance.SideType
	if side == base.BID {
		s = binance.SideTypeBuy
	} else if side == base.ASK {
		s = binance.SideTypeSell
	}

	order, err := c.client.NewCreateOrderService().
		Symbol(symbol).
		Side(s).
		Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).
		Quantity(size).
		Price(price).
		Do(context.Background())
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(order.OrderID, 10), nil
}

func (c *Client) TakerOrder(symbol, side, price, size string) (string, error) {
	var s binance.SideType
	if side == base.BID {
		s = binance.SideTypeBuy
	} else if side == base.ASK {
		s = binance.SideTypeSell
	}

	order, err := c.client.NewCreateOrderService().
		Symbol(symbol).
		Side(s).
		Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeIOC).
		Quantity(size).
		Price(price).
		Do(context.Background())
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(order.OrderID, 10), nil
}

func (c *Client) CancelOrder(symbol, id string) (bool, error) {
	parseInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return false, err
	}

	resp, err := c.client.NewCancelOrderService().
		Symbol(symbol).
		OrderID(parseInt).
		Do(context.Background())
	if err != nil {
		return false, err
	}

	if resp.Status != "CANCELED" {
		return false, err
	}
	return true, nil
}

func (c *Client) CancelOrders(symbol string) error {
	_, err := c.client.NewCancelOpenOrdersService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return err
	}

	return err
}

func (c *Client) GetOrder(symbol, id string) (models.OrderInfo, error) {

	oid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return models.OrderInfo{}, err
	}

	order, err := c.client.NewGetOrderService().
		Symbol(symbol).
		OrderID(oid).
		Do(context.Background())
	if err != nil {
		return models.OrderInfo{}, err
	}

	var side string

	if order.Side == "BUY" {
		side = base.BID
	} else if order.Side == "SELL" {
		side = base.ASK
	}
	var status string
	if order.Status == "NEW" || order.Status == "PARTIALLY_FILLED" {
		status = base.OPEN
	} else if order.Status == "CANCELED" {
		status = base.CANCELED
	} else if order.Status == "FILLED" {
		status = base.FILLED
	}

	o := models.OrderInfo{
		OrderID:  strconv.FormatInt(order.OrderID, 10),
		Symbol:   order.Symbol,
		Side:     side,
		Price:    order.Price,
		Quantity: order.OrigQuantity,
		Status:   status,
		Time:     order.Time,
	}

	return o, nil
}

func (c *Client) GetOpenOrders(symbol string) ([]models.OrderInfo, error) {
	orders, err := c.client.NewListOpenOrdersService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	if len(orders) == 0 {
		return nil, nil
	}

	var orderInfos []models.OrderInfo

	for _, order := range orders {

		var side string

		if order.Side == "BUY" {
			side = base.BID
		} else if order.Side == "SELL" {
			side = base.ASK
		}

		o := models.OrderInfo{
			OrderID:  strconv.FormatInt(order.OrderID, 10),
			Symbol:   order.Symbol,
			Side:     side,
			Price:    order.Price,
			Quantity: order.OrigQuantity,
			Time:     order.Time,
		}

		orderInfos = append(orderInfos, o)

	}

	return orderInfos, nil
}

func (c *Client) GetOpenOrdersWithSide(symbol, side string) ([]models.OrderInfo, error) {
	orders, err := c.GetOpenOrders(symbol)
	if err != nil {
		return nil, err
	}
	var info []models.OrderInfo
	for _, o := range orders {
		if o.Side == side {
			info = append(info, o)
		}
	}
	return info, err
}

func (c *Client) GetOpenSplitOrders(symbol string) ([]models.OrderInfo, []models.OrderInfo, error) {
	orders, err := c.client.NewListOpenOrdersService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return nil, nil, err
	}

	if len(orders) == 0 {
		return nil, nil, nil
	}

	var buys []models.OrderInfo
	var sells []models.OrderInfo

	for _, order := range orders {

		var side string

		if order.Side == "BUY" {
			side = base.BID

			o := models.OrderInfo{
				OrderID:  strconv.FormatInt(order.OrderID, 10),
				Symbol:   order.Symbol,
				Side:     side,
				Price:    order.Price,
				Quantity: order.OrigQuantity,
				Time:     order.Time,
			}

			buys = append(buys, o)

		} else if order.Side == "SELL" {
			side = base.ASK

			o := models.OrderInfo{
				OrderID:  strconv.FormatInt(order.OrderID, 10),
				Symbol:   order.Symbol,
				Side:     side,
				Price:    order.Price,
				Quantity: order.OrigQuantity,
				Time:     order.Time,
			}

			sells = append(sells, o)
		}

	}

	return buys, sells, nil
}

func (c *Client) GetMarketPrice(symbol string) (string, error) {
	prices, err := c.client.NewListPricesService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return "", err
	}

	return prices[0].Price, nil
}

func (c *Client) Depth(symbol, limit string) (models.WsData, error) {
	parseInt, err := strconv.Atoi(limit)
	if err != nil {
		return models.WsData{}, err
	}

	_, err = c.client.NewDepthService().
		Symbol(symbol).
		Limit(parseInt).
		Do(context.Background())

	if err != nil {
		return models.WsData{}, err
	}

	return models.WsData{}, nil
}

func (c *Client) LimitOrders(symbol string, ol []models.OrderList) ([]string, error) {
	l := len(ol)
	IdList := make([]string, 0, l)
	var err error
	var id string
	var Type string

	for i := 0; i < l; i++ {
		if ol[i].Side == base.BID {
			Type = "BUY"
		} else if ol[i].Side == base.ASK {
			Type = "SELL"
		}
		side := Type
		price := ol[i].Price
		size := ol[i].Size
		id, err = c.LimitOrder(side, symbol, price, size)
		if err != nil {
			return nil, err
		}
		IdList = append(IdList, id)
		time.Sleep(100 * time.Millisecond)
	}
	return IdList, err
}

func (c *Client) TakerOrders(symbol string, ol []models.OrderList) ([]string, error) {
	l := len(ol)
	IdList := make([]string, 0, l)
	var err error
	var id string
	var Type string

	for i := 0; i < l; i++ {
		if ol[i].Side == base.BID {
			Type = "BUY"
		} else if ol[i].Side == base.ASK {
			Type = "SELL"
		}
		side := Type
		price := ol[i].Price
		size := ol[i].Size
		id, err = c.TakerOrder(side, symbol, price, size)
		if err != nil {
			return nil, err
		}
		IdList = append(IdList, id)
		time.Sleep(100 * time.Millisecond)
	}
	return IdList, err
}

func (c *Client) MakerOrder(symbol, side, price, size string) (string, error) {
	var s binance.SideType
	if side == base.BID {
		s = binance.SideTypeBuy
	} else if side == base.ASK {
		s = binance.SideTypeSell
	}

	order, err := c.client.NewCreateOrderService().
		Symbol(symbol).
		Side(s).
		Type(binance.OrderTypeLimitMaker).
		TimeInForce(binance.TimeInForceTypeGTC).
		Quantity(size).
		Price(price).
		Do(context.Background())
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(order.OrderID, 10), nil
}

func (c *Client) MakerOrders(symbol string, ol []models.OrderList) ([]string, error) {
	l := len(ol)
	IdList := make([]string, 0, l)
	var err error
	var id string
	var Type string

	for i := 0; i < l; i++ {
		if ol[i].Side == base.BID {
			Type = "BUY"
		} else if ol[i].Side == base.ASK {
			Type = "SELL"
		}
		side := Type
		price := ol[i].Price
		size := ol[i].Size
		id, err = c.MakerOrder(side, symbol, price, size)
		if err != nil {
			return nil, err
		}
		IdList = append(IdList, id)
		time.Sleep(100 * time.Millisecond)
	}
	return IdList, err
}

func (c *Client) GetPairInfo(symbol string) (models.PairInfo, error) {
	pair, err := c.client.
		NewExchangeInfoService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return models.PairInfo{}, err
	}
	info := models.PairInfo{
		MinBaseAmount:   strconv.FormatFloat(float64(1/10^pair.Symbols[0].BaseAssetPrecision), 'f', 7, 64),
		MinQuoteAmount:  strconv.FormatFloat(float64(1/10^pair.Symbols[0].QuoteAssetPrecision), 'f', 7, 64),
		AmountPrecision: 0,
		Precision:       0,
	}
	return info, err

}

func (c *Client) GetTradingFee(symbol string) (models.TradingFee, error) {
	account, err := c.client.
		NewGetAccountService().
		Do(context.Background())
	if err != nil {
		return models.TradingFee{}, err
	}
	info := models.TradingFee{
		Symbol:                symbol,
		TakerFeeFromApi:       strconv.FormatInt(account.TakerCommission, 10),
		MakerFeeFromApi:       strconv.FormatInt(account.MakerCommission, 10),
		TakerFeeFromRealOrder: "",
		MakerFeeFromRealOrder: "",
		IfDiscount:            false,
	}
	return info, err
}
