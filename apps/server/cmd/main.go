package main

import (
	"net/http"

	"github.com/DSAwithGautam/CodeConquerers/internal/common/logger"
	config "github.com/DSAwithGautam/CodeConquerers/internal/config"
)

func main() {

	cfg := config.LoadConfig()
	logger := logger.NewLogger(cfg)

	
	http.ListenAndServe(":"+cfg.Port, nil)
}
