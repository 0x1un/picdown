package main

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	action "picdown/actions"
	"strings"
	"time"

	"github.com/0x1un/omtools/zbxgraph"
	"github.com/chromedp/chromedp"
	gim "github.com/ozankasikci/go-image-merge"
	"github.com/sirupsen/logrus"
)

var (
	cfg = ReadConfigFile("conf/conf.ini")
)

func init() {
	cfg.Init()
	file, err := os.OpenFile(cfg.Default["LogFile"], os.O_APPEND|os.O_CREATE, os.ModeAppend)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetOutput(io.MultiWriter(os.Stdout, file))
	logrus.SetReportCaller(true)
}

type ChromeDPCtx struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewChromeDPContext() *ChromeDPCtx {
	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.Flag("mute-audio", false),
		chromedp.Flag("headless", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.UserAgent(UserAgent),
		chromedp.WindowSize(cfg.Width, cfg.Height),
	)
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel = context.WithTimeout(ctx, time.Duration(cfg.TimeOut)*(time.Second))
	ctx, cancel = chromedp.NewContext(ctx)
	return &ChromeDPCtx{
		ctx:    ctx,
		cancel: cancel,
	}
}

func downloadSangForPics() []string {
	// capture screenshot of an element
	chromeDP := NewChromeDPContext()
	defer chromeDP.cancel()
	picName := make([]string, 0)
	for _, profile := range cfg.Another {
		var buf []byte // 保证每次循环完成后都为空
		sfAcn := action.SangForLogin("https://"+profile["SangForURI"], profile["SangForUser"], profile["SangForPass"], 5, 15)
		if err := chromedp.Run(chromeDP.ctx, fullScreenshot(100, sfAcn, &buf)); err != nil {
			logrus.Errorln(err)
		}
		target := cfg.Output + profile["SangForURI"] + ".png"
		// picName: picture path
		if len(buf) != 0 {
			picName = append(picName, target)
			if err := ioutil.WriteFile(target, buf, 0644); err != nil {
				logrus.Errorln(err)
			}
		}
	}
	return picName
}

func MergePIC(pics []string, xNum, yNum int, target string) error {
	grids := make([]*gim.Grid, 0)
	for _, pic := range pics {
		grids = append(grids, &gim.Grid{
			ImageFilePath: pic,
		})
	}
	rgba, err := gim.New(grids, xNum, yNum).Merge()
	if err != nil {
		return err
	}
	fileHandle, err := os.Create(cfg.Default["Output"] + target)
	if err != nil {
		return err
	}
	if err := png.Encode(fileHandle, rgba); err != nil {
		return err
	}
	return nil
}

func downloadZBXPicAndMerge(cfgPath string, output string) error {
	omap, err := zbxgraph.Run(cfgPath, cfg.Default["Output"], true)
	if err != nil {
		return err
	}
	grids := make([]*gim.Grid, 0)
	for _, picPath := range omap {
		grids = append(grids, &gim.Grid{
			ImageFilePath: func(path []string) string {
				if len(path) <= 1 {
					return path[0]
				}
				return ""
			}(picPath),
		})
	}
	rgba, err := mergePicture(grids)
	if err != nil {
		return err
	}
	fHandle, err := os.Create(cfg.Default["Output"] + output)
	if err != nil {
		return err
	}
	return png.Encode(fHandle, rgba)
}

func mergePicture(grids []*gim.Grid) (*image.RGBA, error) {
	n := len(grids)
	rgba, err := gim.New(grids, n/2, n/2).Merge()
	if err != nil {
		return rgba, err
	}
	return rgba, nil
}

// call 具体的实现
func call() {
	CreateDirectory(cfg.Output)
	CreateDirectory(cfg.Default["Output"])
	pics := downloadSangForPics()
	err := downloadZBXPicAndMerge("conf/zbx.ini", "zbx_merged.png")
	if err != nil {
		logrus.Fatalf("download pictures failure: %s\n", err.Error())
	}

	target := "sf_merged.png"
	if err := MergePIC(pics, 2, 1, target); err != nil {
		logrus.Fatalf("merge the pictures failure: %s\n", err.Error())
	}

	strBuf := strings.Builder{}
	// upload file to QiNiu
	grids := make([]*gim.Grid, 0)
	for _, tg := range []string{"zbx_merged.png", "sf_merged.png"} {
		grids = append(grids, &gim.Grid{
			ImageFilePath: cfg.Default["Output"] + tg,
		})
	}
	mergeRGBA, err := gim.New(grids, 1, 2).Merge()
	if err != nil {
		logrus.Fatalf("Error encountered while merging pictures: %s\n", err.Error())
	}
	resTarget := cfg.Default["Output"] + "result.png"
	fHandler, err := os.Create(resTarget)
	if err != nil {
		logrus.Fatalf("failed to create result.png: %s\n", err.Error())
	}
	if err := png.Encode(fHandler, mergeRGBA); err != nil {
		logrus.Fatalf("failed to save merged pictures: %s", err.Error())
	}
	resName, err := PostFileToStorage(resTarget)
	if err != nil {
		logrus.Fatalf("uploading pic to QiNiu failure: %s\n", err.Error())
	}
	if resName != "" {
		url := cfg.Default["QiNiuURL"] + resName
		strBuf.WriteString(fmt.Sprintf("![](%s)\n", url))
	}
	Send(strings.Split(cfg.Default["DingTokens"], ","), nil, false, strBuf.String(), "screenshot")
}
