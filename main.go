/*
 * Copyright 2018 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/SmartEnergyPlatform/platform-connector-lib"
)

func main() {

	configLocation := flag.String("config", "config.json", "configuration file")
	flag.Parse()

	err := LoadConfig(*configLocation)
	if err != nil {
		log.Fatal(err)
	}

	switch Config.MqttLogLevel {
	case "critical":
		paho.CRITICAL = log.New(os.Stderr, "[paho] ", log.LstdFlags)
	case "error":
		paho.CRITICAL = log.New(os.Stderr, "[paho] ", log.LstdFlags)
		paho.ERROR = log.New(os.Stderr, "[paho] ", log.LstdFlags)
	case "warn":
		paho.CRITICAL = log.New(os.Stderr, "[paho] ", log.LstdFlags)
		paho.ERROR = log.New(os.Stderr, "[paho] ", log.LstdFlags)
		paho.WARN = log.New(os.Stderr, "[paho] ", log.LstdFlags)
	case "debug":
		paho.CRITICAL = log.New(os.Stderr, "[paho] ", log.LstdFlags)
		paho.ERROR = log.New(os.Stderr, "[paho] ", log.LstdFlags)
		paho.WARN = log.New(os.Stderr, "[paho] ", log.LstdFlags)
		paho.DEBUG = log.New(os.Stdout, "[paho] ", log.LstdFlags)
	}

	libConf, err := platform_connector_lib.LoadConfig(*configLocation)
	if err != nil {
		log.Fatal(err)
	}

	connector := platform_connector_lib.Init(libConf, func(endpoint string, protocolParts map[string]string) (responseParts map[string]string, err error) {
		responseParts = map[string]string{}
		err = MqttPublish(endpoint, protocolParts["payload"])
		return
	})

	defer connector.Stop()

	go AuthWebhooks(connector)

	err = MqttStart()
	if err != nil {
		panic(err)
	}
	defer MqttClose()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	sig := <-shutdown
	log.Println("received shutdown signal", sig)
}
