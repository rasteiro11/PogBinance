package novadax

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"github.com/gorilla/websocket"
// 	"github.com/rasteiro11/PogCore/pkg/logger"
// 	"net/url"
// 	"pogbinance/models"
// 	"pogbinance/pkg/marketdata"
// 	"sync"
// 	"time"
// )
//
// type ws struct {
// 	url            *url.URL
// 	once           sync.Once
// 	retry          chan struct{}
// 	done           chan struct{}
//
// 	writer      chan []byte
// 	ready          chan struct{}
// 	subscriberChan chan []string
// 	priceStream    chan models.MarketData
// 	conn           *websocket.Conn
// }
//
// func NewMarketDataProvider() marketdata.MarketDataProvider {
// 	return &ws{
// 		done:        make(chan struct{}),
// 		priceStream: make(chan models.MarketData, 1),
// 		url: &url.URL{
// 			Scheme: "wss",
// 			Host:   "stream.binance.com:443",
// 			Path:   "/ws",
// 		},
// 	}
// }
//
// func (w *ws) subscriber(ctx context.Context) {
// 	select {
// 	case <-w.done:
// 		return
// 	case message := <-w.subscriberChan:
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case _, ok := <-w.done:
// 			if !ok {
// 				return
// 			}
// 		default:
// 		}
//
// 		<-w.ready
//
// 		logger.Of(ctx).Debugf("Sending message: %s", message)
//
// 		// if err := w.client.Emit("SUBSCRIBE", message); err != nil {
// 		// 	select {
// 		// 	case <-ctx.Done():
// 		// 		return
// 		// 	case _, ok := <-w.done:
// 		// 		if !ok {
// 		// 			return
// 		// 		}
// 		// 	default:
// 		// 	}
// 		//
// 		// 	logger.Of(ctx).Errorf("SUBSCRIBE ERROR: %+v", err)
// 		// 	w.retry <- struct{}{}
// 		// 	return
// 		// }
// 		//
// 		// for _, topic := range message {
// 		// 	log.Printf("TOPIC: %+v\n", topic)
// 		// 	err := w.client.On(topic, func(h *gosocketio.Channel, args models.MarketData) {
// 		// 		w.priceStream <- args
// 		// 	})
// 		// 	if err != nil {
// 		// 		select {
// 		// 		case <-ctx.Done():
// 		// 			return
// 		// 		case _, ok := <-w.done:
// 		// 			if !ok {
// 		// 				return
// 		// 			}
// 		// 		default:
// 		// 		}
// 		// 		logger.Of(ctx).Errorf("REGISTER HANDLER ERROR: %+v", err)
// 		// 		w.retry <- struct{}{}
// 		// 		return
// 		// 	}
// 		// }
// 	}
// }
//
// type subscribePayload struct {
// 	Method string   `json:"method"`
// 	Params []string `json:"params"`
// }
//
// func prepareSubscriberPayload(currency ...string) ([]byte, error) {
// 	payload := []string{}
//
// 	for _, c := range currency {
// 		payload = append(payload, fmt.Sprintf("%s@aggTrade", c))
// 	}
//
// 	return json.Marshal(subscribePayload{
// 		Method: "SUBSCRIBE",
// 		Params: payload,
// 	})
// }
//
// type AggTradeEvent struct {
// 	E            string `json:"e"`
// 	EventTime    int64  `json:"E"`
// 	Symbol       string `json:"s"`
// 	TradeID      int    `json:"a"`
// 	Price        string `json:"p"`
// 	Quantity     string `json:"q"`
// 	FirstID      int    `json:"f"`
// 	LastID       int    `json:"l"`
// 	TradeTime    int64  `json:"T"`
// 	IsBuyerMaker bool   `json:"m"`
// 	Ignore       bool   `json:"M"`
// }
//
// type binanceProvider struct {
// 	websocket   *Websocket
// 	messageChan <-chan models.MarketData
// }
//
// func NewInvestingDataProvider(ctx context.Context) (marketdata.MarketDataProvider, error) {
// 	i := &binanceProvider{
// 		messageChan: make(chan models.MarketData),
// 	}
//
// 	provider, err := i.websocket.Start(ctx, "GAMER Currency")
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	i.messageChan = provider
//
// 	return i, nil
// }
//
// type Payload struct {
// 	Message string `json:"message"`
// }
//
// type Websocket struct {
// 	url         *url.URL
// 	once        sync.Once
// 	conn        *websocket.Conn
// 	retry       chan struct{}
// 	done        chan struct{}
// 	ready       chan struct{}
// 	writer      chan []byte
// 	priceStream chan Message
// }
//
// func NewWebsocket() *Websocket {
// 	return &Websocket{
// 		done:        make(chan struct{}),
// 		priceStream: make(chan Message, 1),
// 		url: &url.URL{
// 			Scheme: "wss",
// 			Host:   "streaming.forexpros.com",
// 			Path:   "/echo/180/gf2gdirz/websocket",
// 		},
// 	}
// }
//
// func (w *Websocket) pingHandler(ctx context.Context) {
// 	ticker := time.NewTicker(time.Second)
// 	defer ticker.Stop()
//
// 	<-w.ready
//
// 	for {
// 		select {
// 		case <-w.done:
// 			return
// 		case <-w.retry:
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case _, ok := <-w.done:
// 				if !ok {
// 					return
// 				}
// 			default:
// 			}
//
// 			// logger.Of(ctx).Debugf("Sending ping message: %s", heartbeatPayload)
//
// 			// if err := w.conn.WriteMessage(websocket.TextMessage, []byte(heartbeatPayload)); err != nil {
// 			// 	select {
// 			// 	case <-ctx.Done():
// 			// 		return
// 			// 	case _, ok := <-w.done:
// 			// 		if !ok {
// 			// 			return
// 			// 		}
// 			// 	default:
// 			// 	}
// 			//
// 			// 	logger.Of(ctx).Errorf("pingMessage error: %+v", err)
// 			//
// 			// 	w.retry <- struct{}{}
// 			// 	return
// 			// }
// 		}
// 	}
// }
//
// func (w *Websocket) sendMessage(ctx context.Context) {
// 	for {
// 		select {
// 		case <-w.done:
// 			return
// 		case message := <-w.writer:
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case _, ok := <-w.done:
// 				if !ok {
// 					return
// 				}
// 			default:
// 			}
//
// 			logger.Of(ctx).Debugf("Sending message: %s", message)
//
// 			if err := w.conn.WriteMessage(message); err != nil {
// 				select {
// 				case <-ctx.Done():
// 					return
// 				case _, ok := <-w.done:
// 					if !ok {
// 						return
// 					}
// 				default:
// 				}
//
// 				logger.Of(ctx).Errorf("writeMessage error: %+v", err)
//
// 				w.retry <- struct{}{}
//
// 				return
// 			}
// 		}
// 	}
// }
//
// func (w *Websocket) Start(ctx context.Context, currency ...string) (<-chan Message, error) {
// 	w.priceStream = make(chan Message, 1)
//
// 	payload, err := prepareSubscriberPayload(currency...)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	go func() {
// 		defer func() {
// 			if w.conn != nil {
// 				logger.Of(ctx).Debug("closing database connection")
//
// 				if err := w.conn.Close(); err != nil {
// 					logger.Of(ctx).Warn("failed to closed websocket connection")
// 				}
// 			}
// 		}()
//
// 	Retry:
// 		w.retry = make(chan struct{}, 1)
// 		w.ready = make(chan struct{}, 1)
// 		w.writer = make(chan []byte)
//
// 		logger.Of(ctx).Debugf("Starting connection to: %s", w.url.String())
//
// 		c, _, err := websocket.DefaultDialer.Dial(w.url.String(), nil)
// 		if err != nil {
// 			goto Retry
// 		}
//
// 		logger.Of(ctx).Debug("Successfully connected")
//
// 		go w.sendMessage(ctx)
// 		go w.pingHandler(ctx)
//
// 		w.conn = c
// 		close(w.ready)
//
// 		go func() { w.writer <- []byte(payload) }()
//
// 		for {
// 			select {
// 			case <-w.done:
// 				return
// 			case <-w.retry:
// 				goto Retry
// 			default:
// 				select {
// 				case <-ctx.Done():
// 					return
// 				case _, ok := <-w.done:
// 					if !ok {
// 						return
// 					}
// 				default:
// 				}
//
// 				_, message, err := c.ReadMessage()
// 				if err != nil {
// 					logger.Of(ctx).Errorf("readMessage: %+v\n", err)
// 					goto Retry
// 				}
//
// 				ticker := &AggTradeEvent{}
// 				if err := json.Unmarshal(message, ticker); err != nil {
// 					continue
// 				}
//
// 				w.priceStream <- Message{
// 					PID:         "",
// 					LastDir:     "",
// 					LastNumeric: 0,
// 					Last:        "",
// 					Bid:         "",
// 					Ask:         "",
// 					High:        "",
// 					Low:         "",
// 					LastClose:   "",
// 					PC:          "",
// 					PCP:         "",
// 					PCCol:       "",
// 					Time:        "",
// 					Timestamp:   0,
// 				}
//
// 			}
// 		}
// 	}()
//
// 	return w.priceStream, nil
// }
//
// func (w *Websocket) Close(ctx context.Context) error {
// 	select {
// 	case <-ctx.Done():
// 		return ctx.Err()
// 	default:
// 	}
//
// 	logger.Of(ctx).Debug("Closing connection")
//
// 	w.once.Do(func() {
// 		close(w.done)
// 		close(w.retry)
// 		close(w.priceStream)
// 	})
//
// 	return nil
// }
//
// func (w *ws) sendMessage(ctx context.Context) {
// 	for {
// 		select {
// 		case <-w.done:
// 			return
// 		case message := <-w.writer:
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case _, ok := <-w.done:
// 				if !ok {
// 					return
// 				}
// 			default:
// 			}
//
// 			logger.Of(ctx).Debugf("Sending message: %s", message)
//
// 			if err := w.conn.WriteMessage(websocket.TextMessage, message); err != nil {
// 				select {
// 				case <-ctx.Done():
// 					return
// 				case _, ok := <-w.done:
// 					if !ok {
// 						return
// 					}
// 				default:
// 				}
//
// 				logger.Of(ctx).Errorf("writeMessage error: %+v", err)
//
// 				w.retry <- struct{}{}
//
// 				return
// 			}
// 		}
// 	}
// }
//
// func (w *ws) Start(ctx context.Context, currency ...string) (<-chan models.MarketData, error) {
// 	w.priceStream = make(chan models.MarketData, 1)
//
// 	subscriberPayload, err := prepareSubscriberPayload(currency...)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	go func() {
// 		defer func() {
// 			if w.conn != nil {
// 				logger.Of(ctx).Debug("closing database connection")
// 				if err := w.conn.Close(); err != nil {
// 					logger.Of(ctx).Warn("failed to closed websocket connection")
// 				}
// 			}
// 		}()
//
// 	Retry:
// 		w.retry = make(chan struct{}, 1)
// 		w.ready = make(chan struct{}, 1)
// 		w.subscriberChan = make(chan []string)
//
// 		c, _, err := websocket.DefaultDialer.Dial(w.url.String(), nil)
// 		if err != nil {
// 			goto Retry
// 		}
//
// 		logger.Of(ctx).Debug("Successfully connected")
//
// 		// go w.sendMessage(ctx)
// 		// go w.pingHandler(ctx)
//
// 		w.conn = c
// 		close(w.ready)
//
// 		go func() { w.subscriberChan <- subscriberPayload }()
//
// 		<-w.retry
// 		goto Retry
//
// 	}()
//
// 	return w.priceStream, nil
// }
//
// func (w *ws) Close(ctx context.Context) error {
// 	select {
// 	case <-ctx.Done():
// 		return ctx.Err()
// 	default:
// 	}
//
// 	logger.Of(ctx).Debug("Closing connection")
//
// 	w.once.Do(func() {
// 		close(w.done)
// 		close(w.retry)
// 		close(w.priceStream)
// 	})
//
// 	return nil
// }

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rasteiro11/PogBinance/models"
	"github.com/rasteiro11/PogBinance/pkg/marketdata"
	"github.com/rasteiro11/PogCore/pkg/logger"
	"github.com/shopspring/decimal"
)

var ErrAssetPidNotFound = errors.New("error asset pid not found")

type InvestingOption func(*investinProvider)

type investinProvider struct {
	websocket   *Websocket
	messageChan <-chan Message
	message     Message
}

type Message struct {
	PID         string  `json:"pid"`
	LastDir     string  `json:"last_dir"`
	LastNumeric float64 `json:"last_numeric"`
	Last        string  `json:"last"`
	Bid         string  `json:"bid"`
	Ask         string  `json:"ask"`
	High        string  `json:"high"`
	Low         string  `json:"low"`
	LastClose   string  `json:"last_close"`
	PC          string  `json:"pc"`
	PCP         string  `json:"pcp"`
	PCCol       string  `json:"pc_col"`
	Time        string  `json:"time"`
	Timestamp   int64   `json:"timestamp"`
}

type Payload struct {
	Message string `json:"message"`
}

type Websocket struct {
	url         *url.URL
	once        sync.Once
	conn        *websocket.Conn
	retry       chan struct{}
	done        chan struct{}
	ready       chan struct{}
	writer      chan []byte
	priceStream chan models.MarketData
}

func NewMarketDataProvider() marketdata.MarketDataProvider {
	return &Websocket{
		done:        make(chan struct{}),
		priceStream: make(chan models.MarketData, 1),
		url: &url.URL{
			Scheme: "wss",
			Host:   "stream.binance.com:443",
			Path:   "/ws",
		},
	}
}

func (w *Websocket) pingHandler(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	<-w.ready

	for {
		select {
		case <-w.done:
			return
		case <-w.retry:
			select {
			case <-ctx.Done():
				return
			case _, ok := <-w.done:
				if !ok {
					return
				}
			default:
			}
		case <-ticker.C:
			logger.Of(ctx).Debugf("Sending ping message")

			if err := w.conn.WriteMessage(websocket.PongMessage, nil); err != nil {
				select {
				case <-ctx.Done():
					return
				case _, ok := <-w.done:
					if !ok {
						return
					}
				default:
				}

				logger.Of(ctx).Errorf("pongMessage error: %+v", err)

				w.retry <- struct{}{}
				return
			}
		}
	}
}

func (w *Websocket) sendMessage(ctx context.Context) {
	for {
		select {
		case <-w.done:
			return
		case message := <-w.writer:
			select {
			case <-ctx.Done():
				return
			case _, ok := <-w.done:
				if !ok {
					return
				}
			default:
			}

			logger.Of(ctx).Debugf("Sending message: %s", message)

			if err := w.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				select {
				case <-ctx.Done():
					return
				case _, ok := <-w.done:
					if !ok {
						return
					}
				default:
				}

				logger.Of(ctx).Errorf("writeMessage error: %+v", err)

				w.retry <- struct{}{}

				return
			}
		}
	}
}

func (w *Websocket) Start(ctx context.Context, currency ...string) (<-chan models.MarketData, error) {
	w.priceStream = make(chan models.MarketData, 1)

	payload, err := prepareSubscriberPayload(currency...)
	if err != nil {
		return nil, err
	}

	logger.Of(ctx).Debugf("PAYLOAD: %+v\n", string(payload))

	go func() {
		defer func() {
			if w.conn != nil {
				logger.Of(ctx).Debug("closing websocket connection")

				if err := w.conn.Close(); err != nil {
					logger.Of(ctx).Warn("failed to close websocket connection")
				}
			}
		}()

	Retry:
		w.retry = make(chan struct{}, 1)
		w.ready = make(chan struct{}, 1)
		w.writer = make(chan []byte)

		logger.Of(ctx).Debugf("Starting connection to: %s", w.url.String())

		c, _, err := websocket.DefaultDialer.Dial(w.url.String(), nil)
		if err != nil {
			goto Retry
		}

		logger.Of(ctx).Debug("Successfully connected")

		go w.sendMessage(ctx)
		go w.pingHandler(ctx)

		w.conn = c
		close(w.ready)

		go func() { w.writer <- []byte(payload) }()

		for {
			select {
			case <-w.done:
				return
			case <-w.retry:
				goto Retry
			default:
				select {
				case <-ctx.Done():
					return
				case _, ok := <-w.done:
					if !ok {
						return
					}
				default:
				}

				_, message, err := c.ReadMessage()
				if err != nil {
					logger.Of(ctx).Errorf("readMessage: %+v\n", err)
					goto Retry
				}

				ticker := &TickerEvent{}
				if err := json.Unmarshal(message, ticker); err != nil {
					continue
				}

				w.priceStream <- models.MarketData{
					Ask:            ticker.BestAskPrice,
					BaseVolume24h:  ticker.TotalTradedBaseVolume,
					Bid:            ticker.BestBidPrice,
					High24h:        ticker.HighPrice,
					LastPrice:      ticker.LastPrice,
					Low24h:         ticker.LowPrice,
					Open24h:        ticker.OpenPrice,
					QuoteVolume24h: ticker.TotalTradedQuoteVolume,
					Symbol:         ticker.Symbol,
					Timestamp:      ticker.EventTime,
				}

				// output := message[3 : len(message)-2]
				// processedJson := backslashMatcher.ReplaceAll([]byte(output), []byte(""))
				// match := payloadMatcher.FindStringSubmatch(string(processedJson))

				// if len(match) >= 2 {
				// 	jsonPayload := match[1]

				// 	var message Message
				// 	if err := json.Unmarshal([]byte(jsonPayload), &message); err != nil {
				// 		continue
				// 	}

				// 	logger.Of(ctx).Debugf("PAYLOAD: %s\n", jsonPayload)

				// 	select {
				// 	case <-ctx.Done():
				// 		return
				// 	case _, ok := <-w.done:
				// 		if !ok {
				// 			return
				// 		}
				// 	default:
				// 	}

				// 	w.priceStream <- message
				// }
			}
		}
	}()

	return w.priceStream, nil
}

func (w *Websocket) Close(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	logger.Of(ctx).Debug("Closing connection")

	w.once.Do(func() {
		close(w.done)
		close(w.retry)
		close(w.priceStream)
	})

	return nil
}

type subscribePayload struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
}

func prepareSubscriberPayload(currency ...string) ([]byte, error) {
	payload := []string{}

	for _, c := range currency {
		payload = append(payload, fmt.Sprintf("%s@ticker", c))
	}

	return json.Marshal(subscribePayload{
		Method: "SUBSCRIBE",
		Params: payload,
	})
}

type TickerEvent struct {
	EventTime              int64           `json:"E"`
	Event                  string          `json:"e"`
	Symbol                 string          `json:"s"`
	PriceChange            string          `json:"p"`
	PriceChangePercent     string          `json:"P"`
	WeightedAveragePrice   string          `json:"w"`
	FirstTradePrice        string          `json:"x"`
	LastPrice              decimal.Decimal `json:"c"`
	LastQuantity           string          `json:"Q"`
	BestBidPrice           decimal.Decimal `json:"b"`
	BestBidQuantity        string          `json:"B"`
	BestAskPrice           decimal.Decimal `json:"a"`
	BestAskQuantity        string          `json:"A"`
	OpenPrice              decimal.Decimal `json:"o"`
	HighPrice              decimal.Decimal `json:"h"`
	LowPrice               decimal.Decimal `json:"l"`
	TotalTradedBaseVolume  decimal.Decimal `json:"v"`
	TotalTradedQuoteVolume decimal.Decimal `json:"q"`
	StatisticsOpenTime     int64           `json:"O"`
	StatisticsCloseTime    int64           `json:"C"`
	FirstTradeID           int64           `json:"F"`
	LastTradeID            int64           `json:"L"`
	TotalNumberTrades      int64           `json:"n"`
}
