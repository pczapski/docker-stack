package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/hashicorp/consul/api"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type envConfig struct {
	ConsulDataCenter string `default:"localhost"`
	ServiceName      string `default:"api-gateway"`
	ReloadInterval   string `default:"5s"`
}
type appConfig struct {
	Port    string
	Welcome string
}

func main() {
	var s envConfig
	var cfg appConfig
	var runtime_viper = viper.New()

	err := envconfig.Process("", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	// Create a Consul Config
	consulConfig := api.DefaultConfig()
	consulConfig.Datacenter = s.ConsulDataCenter

	runtime_viper.AddRemoteProvider("consul", consulConfig.Address, fmt.Sprintf("/app/%s", s.ServiceName))
	runtime_viper.SetConfigType("json")
	if err = runtime_viper.ReadRemoteConfig(); err != nil {
		log.Fatal(err.Error())
	}
	// unmarshal config
	runtime_viper.Unmarshal(&cfg)
	go func() {
		for {
			t, err := time.ParseDuration(s.ReloadInterval)
			if err != nil {
				log.Fatalf("interval time  env parse error: %v", err)
			}
			time.Sleep(t)

			if err := runtime_viper.WatchRemoteConfig(); err != nil {
				log.Printf("unable to read remote config: %v", err)
				continue
			}
			runtime_viper.Unmarshal(&cfg)
		}
	}()

	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("connected to consul")

	p, err := strconv.Atoi(cfg.Port)
	if err != nil {
		log.Fatal(err.Error())
	}

	reg := &api.AgentServiceRegistration{
		Name:    s.ServiceName,
		Port:    p,
		Address: getLocalIP(),

		Check: &api.AgentServiceCheck{
			Interval:                       "10s",
			HTTP:                           fmt.Sprintf("http://%s:%s/health", getLocalIP(), cfg.Port),
			DeregisterCriticalServiceAfter: "1m",
		},
	}
	err = client.Agent().ServiceRegister(reg)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("registed %s service in consul\n", s.ServiceName)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router(&cfg, client),
	}

	log.Printf("registed %s service in http\n", cfg.Port)
	if err = server.ListenAndServe(); err != nil {
		log.Fatalln(err.Error())
	}
}

func router(cfg *appConfig, cl *api.Client) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {

		fmt.Fprintf(w, "OK")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		a, err := addr(cl, "worker")
		if err != nil {
			log.Println(err.Error)
		} else {
			sendEvent(fmt.Sprintf("http://%s", a))
		}

		fmt.Fprintf(w, cfg.Welcome)
	})
	return mux
}
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func sendEvent(targetUrl string) {
	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Printf("failed to create client, %v", err)
	}
	event := cloudevents.NewEvent()
	event.SetSource("api-gateway")
	event.SetType("example.type")
	event.SetData(cloudevents.ApplicationJSON, map[string]string{"hello": "world"})

	// Set a target.
	ctx := cloudevents.ContextWithTarget(context.Background(), targetUrl)

	// Send that Event.
	if result := c.Send(ctx, event); !cloudevents.IsACK(result) {
		log.Printf("failed to send, %v", result)
	}
}
func addr(cl *api.Client, serviceName string) (string, error) {
	addrs, _, err := cl.Health().Service(serviceName, "", false, nil)
	if len(addrs) == 0 && err == nil {
		return "", fmt.Errorf("service ( %s ) was not found", serviceName)
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", addrs[0].Service.Address, addrs[0].Service.Port), nil
}
