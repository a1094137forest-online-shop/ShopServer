package constant

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openlyinc/pointy"
	"github.com/spf13/viper"

	"ShopServer/config"
)

func ReadConfig(configPath string) {
	viper.SetConfigFile(configPath)
	viper.AddConfigPath(".")

	viper.SetDefault("PORT", ":1235")
	viper.SetDefault("RUN_MODE", "debug")
	viper.SetDefault("READ_TIME_OUT", 1000)
	viper.SetDefault("WRITE_TIMEOUT", 1000)
	viper.SetDefault("SHUTDOWN_TIMEOUT", 1000)

	envs := []string{
		"PORT",
		"RUN_MODE",
		"READ_TIMEOUT",
		"WRITE_TIMEOUT",
	}

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			log.Println(err)
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	}

	config.Initialize()
}

type RequestParam struct {
	Method  string
	URL     string
	Body    io.Reader
	Header  http.Header
	TimeOut time.Duration
}

type RequestOpts struct {
	UseNewContext bool
	TimeOut       time.Duration
}

func Request(ctx context.Context, param RequestParam) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, param.Method, param.URL, param.Body)
	if err != nil {
		return nil, err
	}

	for name, value := range param.Header {
		req.Header.Set(name, strings.Join(value, ","))
	}

	if param.TimeOut == 0 {
		param.TimeOut = time.Duration(viper.GetInt("REQUEST_TIMEOUT")) * time.Second
	}

	client := &http.Client{
		Timeout: param.TimeOut,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[%s]%s returns non-200 status : %s", param.Method, param.URL, resp.Status)
	}
	return io.ReadAll(resp.Body)
}

type Service string

const (
	ServiceOnlineShop Service = "onlineShop"
	ServiceUser       Service = "user"
	ServiceShop       Service = "shop"
)

func RequestToService(c *gin.Context, s Service, method, api string, body io.Reader, query url.Values, opts ...RequestOpts) ([]byte, error) {
	var serviceURL string

	switch s {
	case ServiceOnlineShop:
		serviceURL = config.OnlineShopServerUrl
	case ServiceUser:
		serviceURL = config.UserServerUrl
	case ServiceShop:
		serviceURL = config.ShopServerUrl
	}

	uService, err := url.Parse(serviceURL)
	if err != nil {
		return nil, err
	}

	u, err := uService.Parse(api)
	if err != nil {
		return nil, err
	}
	u.RawQuery = query.Encode()

	header := http.Header{}
	if c != nil {
		header.Set("TraceID", c.GetString("TraceID"))
		header.Set("SpanID", getNextSpanID(c))
	}

	reqParam := RequestParam{
		Method: method,
		URL:    u.String(),
		Body:   body,
		Header: header,
	}
	log.Println("[info] request to ", u.String())

	ctx := getRequestContext(c)
	if len(opts) > 0 {
		reqParam.TimeOut = opts[0].TimeOut
		if opts[0].UseNewContext {
			ctx = context.Background()
		}
	}

	result, err := Request(ctx, reqParam)
	return result, err
}

func getNextSpanID(c *gin.Context) string {
	countPtr, _ := c.Get("SpanCount")
	spanCount, ok := countPtr.(*int64)
	if !ok {
		spanCount = pointy.Int64(0)
	}
	newCount := atomic.AddInt64(spanCount, 1)
	return c.GetString("SpanID") + "." + strconv.FormatInt(newCount-1, 10)
}

func getRequestContext(c *gin.Context) context.Context {
	if c == nil || c.Request == nil {
		return context.Background()
	}
	return c.Request.Context()
}
