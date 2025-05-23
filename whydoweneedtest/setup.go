package whydoweneedtest

import "github.com/Applessr/hello-sekai-shop-tutorial/config"

func NewTestConfig() *config.Config {
	cfg := config.LoadConfig("../env/test/.env")
	return &cfg
}
