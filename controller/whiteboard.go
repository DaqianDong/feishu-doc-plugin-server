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
		rsp, err = ocr.OCR("data:image/jpeg;base64,"+encoded, "请识别画板图片中的所有文本内容，并以纯 JSON 格式返回。必须严格按照以下 JSON 结构输出：\n\n{\n  \"标题\": \"提取的标题或主题\",\n  \"描述\": \"详细描述内容，包含完整的段落和说明\",\n  \"类型\": \"分类或类型\",\n  \"核心概念\": \"主要概念或主题\",\n  \"关键元素\": \"元素1, 元素2, 元素3\",\n  \"重要机制\": \"机制或流程说明\",\n  \"触发条件\": \"触发条件说明\",\n  \"标签\": \"标签1, 标签2, 标签3\"\n}\n\n重要要求：\n1. 只返回 JSON 数据本身，不要使用 markdown 代码块标记（如 ```json 或 ```）\n2. 不要在 JSON 前后添加任何说明文字或换行符\n3. 字段内容为空时使用空字符串 \"\"\n4. 保持原文的层级结构和逻辑关系\n5. 关键元素和标签用逗号分隔\n6. JSON 字段名必须与上述示例完全一致（使用中文键名）")
		log.Printf("ocr 耗时：%v", time.Since(start))
	}

	return
}
