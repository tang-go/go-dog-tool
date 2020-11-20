package gateway

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/swaggo/swag"
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-gw/ws"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/pkg/config"
	"github.com/tang-go/go-dog/pkg/context"
	"github.com/tang-go/go-dog/pkg/service"
	"github.com/tang-go/go-dog/plugins"
)

//Gateway 服务发现
type Gateway struct {
	service   plugins.Service
	ws        *ws.Ws
	discovery *GoDogDiscovery
}

//NewGateway  新建发现服务
func NewGateway() *Gateway {
	gateway := new(Gateway)
	//初始化配置
	cfg := config.NewConfig()
	//初始化服务发现
	gateway.discovery = NewGoDogDiscovery(cfg.GetDiscovery())
	//初始化rpc服务
	gateway.service = service.CreateService(define.SvcGateWay, cfg, gateway.discovery)
	//设置服务端最大访问量
	gateway.service.GetLimit().SetLimit(define.MaxServiceRequestCount)
	//设置客户端最大访问量
	gateway.service.GetClient().GetLimit().SetLimit(define.MaxClientRequestCount)
	//初始化websocket客户端
	gateway.ws = ws.NewWs(gateway.service)
	//推送消息
	gateway.service.RPC("Push", 3, false, "推送消息", gateway.ws.Push)
	//初始化文档
	swag.Register(swag.Name, gateway)
	return gateway
}

//Run 启动
func (g *Gateway) Run() {
	go func() {
		gin.SetMode(gin.ReleaseMode)
		router := gin.New()
		router.Use(g.cors())
		router.Use(g.logger())
		//静态文件夹
		//router.StaticFS("/", http.Dir("./static"))
		//websocket路由
		router.GET("/ws", func(c *gin.Context) {
			g.ws.Connect(c.Writer, c.Request, c)
		})
		//swagger 文档
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		//添加路由
		router.POST("/api/*router", g.routerPostResolution)
		//GET请求
		router.GET("/api/*router", g.routerGetResolution)
		log.Tracef("网管启动 0.0.0.0:8080")
		err := router.Run(":8080")
		if err != nil {
			panic(err.Error())
		}
	}()
	err := g.service.Run()
	if err != nil {
		log.Warnln(err.Error())
	}
}

//routerGetResolution get路由解析
func (g *Gateway) routerGetResolution(c *gin.Context) {
	url := "/api" + c.Param("router")
	apiservice, ok := g.discovery.GetAPIByURL(url)
	if !ok {
		c.JSON(http.StatusNotFound, customerror.EnCodeError(http.StatusNotFound, "路由错误"))
		return
	}
	if c.Request.Method != apiservice.Method.Kind {
		c.JSON(http.StatusNotFound, customerror.EnCodeError(http.StatusNotFound, "路由错误"))
		return
	}

	timeoutstr := c.Request.Header.Get("TimeOut")
	if timeoutstr == "" {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "timeout不能为空"))
		return
	}
	timeout, err := strconv.Atoi(timeoutstr)
	if err != nil {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
		return
	}
	if timeout <= 0 {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "timeout必须大于0"))
		return
	}
	istest := c.Request.Header.Get("IsTest")
	if istest == "" {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "istest不能为空"))
		return
	}
	traceID := c.Request.Header.Get("TraceID")
	if traceID == "" {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "traceID不能为空"))
		return
	}
	isTest, err := strconv.ParseBool(istest)
	if err != nil {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
		return
	}
	token := ""
	if apiservice.Method.IsAuth {
		token = c.Request.Header.Get("Token")
		if token == "" {
			c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "token不能为空"))
			return
		}
	}
	p := make(map[string]interface{})
	for key, value := range apiservice.Method.Request {
		data := c.Query(key)
		if data == "" {
			c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "参数不正确"))
			return
		}
		vali, ok := value.(map[string]interface{})
		if !ok {
			c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "参数不正确"))
			return
		}
		tp, ok2 := vali["type"].(string)
		if !ok2 {
			c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "参数不正确"))
			return
		}
		v, err := _Transformation(tp, data)
		if err != nil {
			c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
			return
		}
		p[key] = v
	}
	body, _ := g.service.GetClient().GetCodec().EnCode("json", p)
	ctx := context.Background()
	ctx.SetAddress(c.ClientIP())
	ctx.SetIsTest(isTest)
	ctx.SetTraceID(traceID)
	ctx.SetToken(token)
	ctx.SetData("URL", url)
	ctx = context.WithTimeout(ctx, int64(time.Second*time.Duration(timeout)))
	back, err := g.service.GetClient().SendRequest(ctx, plugins.RandomMode, apiservice.Name, apiservice.Method.Name, "json", body)
	if err != nil {
		e := customerror.DeCodeError(err)
		c.JSON(http.StatusOK, e)
		return
	}
	resp := make(map[string]interface{})
	g.service.GetClient().GetCodec().DeCode("json", back, &resp)
	c.JSON(http.StatusOK, gin.H{
		"code": define.SuccessCode,
		"body": resp,
		"time": time.Now().Unix(),
	})
	return
}

// routerPostResolution post路由解析
func (g *Gateway) routerPostResolution(c *gin.Context) {
	//路由解析
	url := c.Request.URL.String()
	apiservice, ok := g.discovery.GetAPIByURL(url)
	if !ok {
		c.JSON(http.StatusNotFound, customerror.EnCodeError(http.StatusNotFound, "路由错误"))
		return
	}
	if c.Request.Method != apiservice.Method.Kind {
		c.JSON(http.StatusNotFound, customerror.EnCodeError(http.StatusNotFound, "路由错误"))
		return
	}
	timeoutstr := c.Request.Header.Get("TimeOut")
	if timeoutstr == "" {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "timeout不能为空"))
		return
	}
	timeout, err := strconv.Atoi(timeoutstr)
	if err != nil {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
		return
	}
	if timeout <= 0 {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "timeout必须大于0"))
		return
	}
	istest := c.Request.Header.Get("IsTest")
	if istest == "" {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "istest不能为空"))
		return
	}
	traceID := c.Request.Header.Get("TraceID")
	if traceID == "" {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "traceID不能为空"))
		return
	}
	isTest, err := strconv.ParseBool(istest)
	if err != nil {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
		return
	}
	token := ""
	if apiservice.Method.IsAuth {
		token = c.Request.Header.Get("Token")
		if token == "" {
			c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "token不能为空"))
			return
		}
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
		return
	}
	body, err = g.validation(string(body), apiservice.Method.Request)
	if err != nil {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
		return
	}
	ctx := context.Background()
	ctx.SetAddress(c.ClientIP())
	ctx.SetIsTest(isTest)
	ctx.SetTraceID(traceID)
	ctx.SetToken(token)
	ctx.SetData("URL", url)
	ctx = context.WithTimeout(ctx, int64(time.Second*time.Duration(timeout)))
	back, err := g.service.GetClient().SendRequest(ctx, plugins.RandomMode, apiservice.Name, apiservice.Method.Name, "json", body)
	if err != nil {
		e := customerror.DeCodeError(err)
		c.JSON(http.StatusOK, e)
		return
	}
	resp := make(map[string]interface{})
	g.service.GetClient().GetCodec().DeCode("json", back, &resp)
	c.JSON(http.StatusOK, gin.H{
		"code": define.SuccessCode,
		"body": resp,
		"time": time.Now().Unix(),
	})
	return
}

//validation 验证参数
func (g *Gateway) validation(param string, tem map[string]interface{}) ([]byte, error) {
	p := make(map[string]interface{})
	if err := g.service.GetClient().GetCodec().DeCode("json", []byte(param), &p); err != nil {
		return nil, err
	}
	for key := range p {
		if _, ok := tem[key]; !ok {
			log.Traceln("模版", tem, "传入参数", p)
			return nil, errors.New("参数内容不正确")
		}
	}
	return g.service.GetClient().GetCodec().EnCode("json", p)
}

//logger 自定义日志输出
func (g *Gateway) logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		//执行时间
		latency := end.Sub(start)
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		log.Tracef("| %3d | %13v | %15s | %s  %s \n",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}

//cors 处理跨域问题
func (g *Gateway) cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,TraceID, IsTest, Token,TimeOut")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
