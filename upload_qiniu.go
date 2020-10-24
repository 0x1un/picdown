package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

func PostFileToStorage(filename string) (string, error) {
	localFile := filename
	key := ComputeFileSHA(filename)
	putPolicy := storage.PutPolicy{
		Scope: cfg.Default["QiNiuBucket"],
	}
	mac := qbox.NewMac(cfg.Default["QiNiuAK"], cfg.Default["QiNiuSK"])
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuanan
	cfg.UseHTTPS = false
	cfg.UseCdnDomains = false

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "github logo",
		},
	}
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return ret.Key, nil
}

func ComputeFileSHA(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file %s, reason: %s", filename, err.Error())
	}
	sha := md5.New()
	_, err = io.Copy(sha, file)
	if err != nil {
		log.Fatalf("failed compute %s, reason: %s", filename, err.Error())
	}
	return hex.EncodeToString(sha.Sum(nil))
}
