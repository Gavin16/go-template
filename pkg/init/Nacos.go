package init

import (
	"bytes"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"log"
	_ "minsky/go-template/pkg/conf"
)

// Service Register and Reading Dynamic Config
// merge nacos dynamic config into viper

func init() {
	register()
	downloadConfig()
}

// Register to Nacos
func register() {

	serviceIp := viper.GetString("register" + ".serviceIP")

	clientConfig := constant.ClientConfig{
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "error",
	}

	if namespaceId := viper.GetString("register" + ".namespaceId"); namespaceId != "" {
		clientConfig.NamespaceId = namespaceId
	}

	var serverConfigs []constant.ServerConfig
	serverConfig := constant.ServerConfig{
		ContextPath: "/nacos",
		Scheme:      "http",
	}
	serverConfig.IpAddr = viper.GetString("register" + ".ipAddr")
	serverConfig.Port = 8848
	if port := viper.GetUint64("register" + ".port"); port != 0 {
		serverConfig.Port = port
	}
	serverConfigs = append(serverConfigs, serverConfig)

	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs},
	)
	if err != nil {
		panic(fmt.Sprintf("create nameing client error:%v", err))
	}
	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          serviceIp,
		Port:        8288,
		ServiceName: "algorithm-integrate",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	if !success {
		panic(fmt.Errorf("register to nacos center failed: %w", err))
	} else {
		log.Println(">>> ... successfully register to nacos! ...")
	}
}

// Download Dynamic Config
// Merge Into Viper
func downloadConfig() {

	ipAddr := viper.GetString("config" + ".ipAddr")
	finalPort := uint64(8848)
	if port := viper.GetUint64("config" + ".port"); port != 0 {
		finalPort = port
	}
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(ipAddr, finalPort, constant.WithContextPath("/nacos")),
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)
	if namespaceId := viper.GetString("config" + ".namespaceId"); namespaceId != "" {
		cc.NamespaceId = namespaceId
	}

	// create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(fmt.Errorf("nacos config center register failed:%w", err))
	}
	log.Println(">>> ... connected to dynamic config successfully! ...")

	// Merge config into Viper
	mergeConfigToViper(client)
}

func mergeConfigToViper(client config_client.IConfigClient) {
	config, err := client.GetConfig(vo.ConfigParam{
		DataId: "go-template",
		Group:  "DEFAULT_GROUP",
	})

	if err != nil {
		panic(fmt.Errorf("get dynamic config error: %w", err))
	}

	log.Printf("remote dynamic config: {%v}\n", config)

	// merge config from nacos into viper
	err = viper.MergeConfig(bytes.NewBufferString(config))
	if err != nil {
		return
	}
}
