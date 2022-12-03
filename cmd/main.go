package main

import (
	"github.com/djedjethai/generation/internal/agent"
	"log"
	"os"
)

func main() {

	cfg, err := setupSrv()
	if err != nil {
		os.Exit(1)
	}

	switch cfg.Protocol {
	case "http":
		// runHTTP(&services, loggerFacade, cfg.Port)
	case "grpc":
		_, err := agent.New(cfg)
		if err != nil {
			log.Println("the err from setting the agent: ", err)
		}
		// srv, l := runGRPC(&services, loggerFacade, cfg.PortGRPC)
		// if err := srv.Serve(l); err != nil {
		// 	log.Fatal("Error run grpc server: ", err)
		// }
		// defer srv.Stop()
		// defer l.Close()

	default:
		log.Fatalln("Invalid protocol...")
	}

}

// func runGRPC(services *config.Services, loggerFacade *lgr.LoggerFacade, port string) (*gglGrpc.Server, net.Listener) {
//
// 	l, err := net.Listen("tcp", fmt.Sprintf("%s%s", "127.0.0.1", port))
// 	if err != nil {
// 		log.Fatalf("Failed to listen: %v", err)
// 	}
//
// 	// set tls
// 	serverTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
// 		CertFile:      config.ServerCertFile,
// 		KeyFile:       config.ServerKeyFile,
// 		CAFile:        config.CAFile,
// 		ServerAddress: l.Addr().String(),
// 	})
//
// 	serverCreds := credentials.NewTLS(serverTLSConfig)
//
// 	server, err := grpc.NewGRPCServer(services, loggerFacade, gglGrpc.Creds(serverCreds))
// 	if err != nil {
// 		log.Fatal("Error create GRPC server: ", err)
// 	}
//
// 	return server, l
// }
//
// func runHTTP(services *config.Services, loggerFacade *lgr.LoggerFacade, port string) {
// 	// handler(application layer)
// 	hdl := rest.NewHandler(services, loggerFacade)
// 	router := hdl.Multiplex()
//
// 	fmt.Printf("***** Service listening on port %s *****", port)
// 	log.Fatal(http.ListenAndServe(port, router))
// }
