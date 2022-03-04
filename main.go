package main

import (
	"edgeTD/config"
	tdM "edgeTD/model/tdengine"
	"encoding/json"
	"fmt"
	mqttG "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/stevenyao001/edgeCommon/http"
	"github.com/stevenyao001/edgeCommon/logger"
	"github.com/stevenyao001/edgeCommon/mqtt"
	"github.com/stevenyao001/edgeCommon/tdengine"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	conf     *config.Config
	dataChan chan *mqtt.Msg
)

func init() {
	dataChan = make(chan *mqtt.Msg, 10)
	conf = new(config.Config)
	config.InitConf("./conf/local.yaml", conf)
}

func main() {
	logger.InitLog(conf.Log.MainPath)

	mqttConnect()
	tdEngine()
}

func httpServer() {
	var wg sync.WaitGroup
	wg.Add(1)
	closeChan := make(chan struct{})
	httpConf := http.Conf{
		Addr:            ":18081",
		ShutdownTimeout: time.Second * 50,
		Router:          GetCommend,
		Wg:              &wg,
		Close:           closeChan,
	}
	http.InitHttp(httpConf)

	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	_ = <-quit

	closeChan <- struct{}{}
	wg.Wait()
}

func GetCommend(router *gin.Engine) {
	router.GET("/", funcTest)
}

func funcTest(gctx *gin.Context) {
	http.RespError(gctx, http.RespCodeSuccess, nil)
}

func tdEngine() {
	tdConfs := make([]tdengine.Conf, 0)
	for _, value := range conf.TDengine {
		tdConf := tdengine.Conf{
			InsName:      value.InsName,
			Driver:       "taosRestful",
			Network:      "http",
			Addr:         "192.168.56.5",
			Port:         0,
			Username:     value.Username,
			Password:     value.Password,
			Db:           value.DbName,
			MaxIdleConns: 10,
			MaxIdleTime:  0,
			MaxLifeTime:  0,
			MaxOpenConns: 10,
		}
		tdConfs = append(tdConfs, tdConf)
	}
	tdengine.InitTdEngine(tdConfs)
	userTd := tdM.NewUserTd()

	go func() {
		for {
			err := userTd.Find()
			if err != nil {
				return
			}
			time.Sleep(10 * time.Second)
		}
	}()

	for {
		select {
		case data := <-dataChan:
			buf, _ := json.Marshal(data.Content)
			input := &tdM.Input{}
			err := json.Unmarshal(buf, input)
			if err != nil {
				fmt.Println(err)
			}
			userTd.Insert(input)
		}
	}
}

func mqttConnect() {
	opt := mqtt.SubscribeOpts{
		//Topic: "topic/test",
		Topic:    "$ROOTEDGE/thing/realtimedata/123456",
		Qos:      0,
		Callback: messagePubHandler,
	}
	opts := make([]mqtt.SubscribeOpts, 0)
	opts = append(opts, opt)
	SubOpts := make(map[string][]mqtt.SubscribeOpts)
	mqttConfs := make([]mqtt.Conf, 0)
	for _, value := range conf.Mqtt {
		SubOpts[value.InsName] = opts
		mqttConf := mqtt.Conf{
			InsName:  value.InsName,
			ClientId: value.ClientId,
			Username: value.Username,
			Password: value.Password,
			Addr:     value.Addr,
			Port:     value.Port,
		}
		mqttConfs = append(mqttConfs, mqttConf)
	}
	mqtt.InitMqtt(mqttConfs, SubOpts)
}

func messagePubHandler(client mqttG.Client, msg mqttG.Message) {
	//fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	tmp := &mqtt.Msg{}
	json.Unmarshal(msg.Payload(), tmp)
	dataChan <- tmp
}
