package service

import (
	"github.com/harveytvt/movie-reservation-system/internal/config"
	"github.com/harveytvt/movie-reservation-system/internal/data/r2"
)

var (
	r2Client r2.Client
)

func init() {
	cfConf := config.Get().Cloudflare
	r2Client = r2.NewClient(cfConf.AccountID, cfConf.AccessKeyID, cfConf.AccessKeySecret)
}
