package main

import (
	"errors"
	"fmt"
	"goGrpcConn/cms/handler"
	"goGrpcConn/svcUtils/logging"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/yookoala/realpath"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log := logging.NewLogger().WithFields(logrus.Fields{
		"service": "goGrpcConn",
		"version": "1.0",
	})
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)

	config.SetConfigFile("env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	switch config.GetString("runtime.loglevel") {
	case "trace":
		log.Logger.SetLevel(logrus.TraceLevel)
	case "debug":
		log.Logger.SetLevel(logrus.DebugLevel)
	default:
		log.Logger.SetLevel(logrus.InfoLevel)
	}
	log.WithField("log level", log.Logger.Level).Info("starting cms service")
	// dialing API microservice
	api := config.GetString("api.url")
	optsApi := getGRPCOpts(config, false)
	log.Info("dialing api service url :", api)
	apiCon, err := grpc.Dial(api, optsApi...)
	if err != nil {
		log.Infof("dialing api..")
		logging.WithError(err, log).Fatal("unable to connect api")
	}
	defer apiCon.Close()
	s, err := newServer(log, config, apiCon)
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", ":"+config.GetString("server.port"))
	if err != nil {
		return err
	}
	fmt.Printf("\n**********************************************************\n")
	fmt.Println("CMS - Client Running on Port : ", config.GetString("server.port"))
	fmt.Printf("**********************************************************\n\n")
	if err := http.Serve(l, s); err != nil {
		return err
	}
	return nil
}

func getGRPCOpts(_ *viper.Viper, withTLS bool) []grpc.DialOption {
	var opts []grpc.DialOption
	if withTLS {
		creds := credentials.NewClientTLSFromCert(nil, "")
		opts = []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	} else {
		opts = []grpc.DialOption{grpc.WithInsecure()}
	}
	opts = append(opts, grpc.WithBlock())
	return opts
}

func newServer(log *logrus.Entry, config *viper.Viper, apiCon *grpc.ClientConn) (*mux.Router, error) {
	env := config.GetString("runtime.environment")
	log.WithField("environment", env).Info("configuring service")
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	assetPath, err := realpath.Realpath(filepath.Join(wd, "assets"))
	if err != nil {
		return nil, err
	}
	asst := afero.NewIOFS(afero.NewBasePathFs(afero.NewOsFs(), assetPath))
	cookieSecret := config.GetString("auth.cookieSecret")
	if cookieSecret == "" {
		return nil, errors.New("missing cookie secret")
	}

	cookies := sessions.NewCookieStore([]byte(cookieSecret))
	cookies.Options.HttpOnly = true
	cookies.Options.MaxAge = config.GetInt("auth.cookieMaxAge")
	cookies.Options.Secure = config.GetBool("auth.cookieSecure")

	srv, err := handler.NewServer(env, config, log, cookies, decoder, asst, apiCon)
	if err != nil {
		log.Fatalf("error in connecting Server : %v", err)
	}
	return srv, err
}
