package larkclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkboard "github.com/larksuite/oapi-sdk-go/v3/service/board/v1"
	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
)

type larkClient struct {
	appId  string
	secret string
	client *lark.Client

	token      string
	expireTime int64 // 过期时间

	ticker *time.Ticker
}

func (l *larkClient) loop() {
	go func() {
		l.ticker = time.NewTicker(time.Second * 5)
		for range l.ticker.C {
			if err := l.refreshAccessToken(); err != nil {
				fmt.Printf("refreshAccessToken error %v\n", err)
			}
		}
	}()
}

func (l *larkClient) refreshAccessToken() error {
	curTime := time.Now().Unix()
	// 无需刷新
	if l.expireTime > curTime {
		return nil
	}

	rsp, err := l.client.GetTenantAccessTokenBySelfBuiltApp(context.Background(), &larkcore.SelfBuiltTenantAccessTokenReq{
		AppID:     l.appId,
		AppSecret: l.secret,
	})

	if err != nil {
		return err
	}

	if rsp.Code != 0 {
		return errors.New(rsp.Msg)
	}

	l.token = rsp.TenantAccessToken
	// 过期时间，提前5分钟刷新
	l.expireTime = curTime + int64(rsp.Expire) - 600

	fmt.Printf("tenantAccessToken: %s, expire: %d\n", l.token, l.expireTime)

	return nil
}

var client *larkClient

// Init 初始化
func Init(appId, secret string) {
	client = &larkClient{
		appId:  appId,
		secret: secret,
		client: lark.NewClient(appId, secret), // 默认配置为自建应用
	}

	if err := client.refreshAccessToken(); err != nil {
		panic(err)
	}
	client.loop()
}

// Stop 停止
func Stop() {
	if client != nil {
		client.ticker.Stop()
	}
}

// WhiteboardImg 下载白板图片
func WhiteboardImg(documentId, recordId string) (data []byte, err error) {
	// 创建请求对象
	req := larkdocx.NewGetDocumentBlockReqBuilder().DocumentId(documentId).BlockId(recordId)
	rsp, err := client.client.Docx.V1.DocumentBlock.Get(
		context.Background(),
		req.Build(),
		larkcore.WithTenantAccessToken(client.token),
	)
	if err != nil {
		return
	}
	if rsp.CodeError.Code != 0 {
		err = fmt.Errorf("code: %d, msg: %s", rsp.CodeError.Code, rsp.CodeError.Msg)
		return
	}

	req2 := larkboard.NewDownloadAsImageWhiteboardReqBuilder().WhiteboardId(*rsp.Data.Block.Board.Token)
	rsp2, err := client.client.Board.V1.Whiteboard.DownloadAsImage(
		context.Background(),
		req2.Build(),
		larkcore.WithTenantAccessToken(client.token),
	)
	if err != nil {
		return
	}
	if rsp2.CodeError.Code != 0 {
		err = fmt.Errorf("code: %d, msg: %s", rsp.CodeError.Code, rsp.CodeError.Msg)
		return
	}
	if !rsp2.Success() {
		err = fmt.Errorf("logId: %s, error response: \n%s", rsp.RequestId(), larkcore.Prettify(rsp.CodeError))
		return
	}

	data, err = io.ReadAll(rsp2.File)
	os.WriteFile("a.png", data, 0777)

	return
}
