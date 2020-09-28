package drago

import (
	"net/url"
	"path"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/clientv3"
	embed "go.etcd.io/etcd/embed"
)

func (s *Server) setupEtcdServer() error {

	cfg := embed.NewConfig()

	// Advertise peer URLs
	apURLs, err := parseUrls(s.config.Etcd.InitialAdvertisePeerURLs)
	if err != nil {
		return err
	}

	// Listen peer URLs
	lpURLs, err := parseUrls(s.config.Etcd.ListenPeerURLs)
	if err != nil {
		return err
	}

	// Advertise client URLs
	acURLs, err := parseUrls(s.config.Etcd.InitialAdvertiseClientURLs)
	if err != nil {
		return err
	}

	// Listen client URLs
	lcURLs, err := parseUrls(s.config.Etcd.ListenClientURLs)
	if err != nil {
		return err
	}

	cfg.Name = s.config.Etcd.Name
	cfg.Dir = path.Join(s.config.DataDir, "/etcd")
	cfg.WalDir = path.Join(s.config.DataDir, "/etcd", "/wal")
	cfg.Logger = "zap"

	cfg.APUrls = apURLs
	cfg.LPUrls = lpURLs
	cfg.ACUrls = acURLs
	cfg.LCUrls = lcURLs

	cfg.LogOutputs = []string{"stderr", path.Join(s.config.DataDir, "/etcd.log")}
	cfg.LogLevel = strings.ToLower(s.config.LogLevel)

	s.logger.Infof("starting etcd server")

	etcdServer, err := embed.StartEtcd(cfg)
	if err != nil {
		return err
	}

	s.etcdServer = etcdServer

	return nil
}

func (s *Server) setupEtcdClient() error {

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:        s.config.Etcd.InitialAdvertiseClientURLs,
		AutoSyncInterval: time.Second * 5,
		DialTimeout:      5 * time.Second,
	})

	if err != nil {
		return err
	}

	s.etcdClient = etcdClient

	return nil
}

func parseUrls(urls []string) ([]url.URL, error) {
	res := []url.URL{}
	for _, v := range urls {
		url, err := url.Parse(v)
		if err != nil {
			return nil, err
		}
		res = append(res, *url)
	}
	return res, nil
}
