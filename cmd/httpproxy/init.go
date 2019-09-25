package httpproxy

import (
	"github.com/chattarajoy/edge/internal/httputils/routers"
	"database/sql"
	"github.com/go-kit/kit/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/heptio/workgroup"
	"github.com/spf13/cobra"
	"os"
)

type Config struct {
	ProxyConfigFile	string
	HTTPPort 		int
	HTTPSPort		int
	ServerDrainTine	int
	Debug			float64
}

type HTTPProxy struct {
	Config 	*Config
	Logger	log.Logger
	Name	string
	Router	routers.Router
	Db      *sql.DB
}

var (

	httpProxy = &HTTPProxy{
		Config: &Config{
			// ProxyConfigFile:"./config/proxy_config.yml",
			HTTPPort:3000,
			HTTPSPort:3001,
			ServerDrainTine:5,
		},
		Name:"httpproxy",
		Router: routers.CreateRouter("httprouter"),
	}

	httpProxyCmd = &cobra.Command{
		Use: "httpproxy",
		Short: "launch httpproxy",
		Long: "launch httproxy",

		Run: func(cmd *cobra.Command, args []string) {
			_ = httpProxy.runServer()
		},
	}
)

func Init(cmd *cobra.Command) {
	cmd.AddCommand(httpProxyCmd)

	httpProxy.initLogger()
	httpProxy.initDb()
}

func (proxy *HTTPProxy) initLogger() {
	proxy.Logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	proxy.Logger = log.With(proxy.Logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller, "app", proxy.Name)
}

func (proxy *HTTPProxy) initDb(){
	db, err := sql.Open("mysql", "ruser:ruser@tcp(localhost:3306)/edge")
	proxy.Db = db

	if err != nil {
		panic(err.Error())
	}
}

func (proxy *HTTPProxy) runServer() error {
	var g workgroup.Group

	err := createAndRunServer(&g, &serverInput{
		Port:            proxy.Config.HTTPPort,
		Logger:          proxy.Logger,
		Router:          proxy.Router,
		ServerDrainTime: proxy.Config.ServerDrainTine,
		Db:              proxy.Db,
	})

	if err != nil {
		return err
	}
	return nil
}