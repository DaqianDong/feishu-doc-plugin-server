package controller

import (
	"encoding/base64"
	"fmt"
	"log"
	"oauth-test/infra/larkclient"
	"oauth-test/infra/ocr"
	"time"

	"github.com/gin-gonic/gin"
)

func WhiteboardController(c *gin.Context) (rsp any, err error) {
	documentId := c.Query("documentId")
	recordId := c.Query("recordId")
	if documentId == "" {
		err = fmt.Errorf("documentId invalid: %s", documentId)
		return
	}
	if recordId == "" {
		err = fmt.Errorf("recordId invalid: %s", recordId)
		return
	}

	data, err := larkclient.WhiteboardImg(documentId, recordId)
	// 下载到图片了？
	if err == nil && len(data) > 0 {
		// 使用 StdEncoding 进行编码
		encoded := base64.StdEncoding.EncodeToString(data)
		start := time.Now()
		rsp, err = ocr.OCR("data:image/jpeg;base64,"+encoded, "请按阅读顺序识别图片中的所有文本，以纯文本形式输出，不要包含任何位置信息、JSON 格式或 HTML 代码。保留标题、列表的层级结构。")
		log.Printf("ocr 耗时：%v", time.Since(start))
	}

	return
}
