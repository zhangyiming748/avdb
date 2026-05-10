package trans

import (
	"encoding/json"
	"log"
	"strings"
)

// DeepLXResponse DeepLX API 响应结构体
type DeepLXResponse struct {
	Code         int      `json:"code"`
	ID           int64    `json:"id"`
	Data         string   `json:"data"`
	Alternatives []string `json:"alternatives"`
}

/*
curl --location --request POST 'https://api.deeplx.org/{{ Authorization }}/translate' \
--header 'Authorization: {{ Authorization }}' \

	--header 'Content-Type: application/json' \
		--data-raw '{
		    "text": "hello",
		    "source_lang": "auto",
		    "target_lang": "zh"
		}'
*/
func DeepLX(src, Authorization string) (dst string) {

	query := strings.Join([]string{"https://api.deeplx.org", Authorization, "translate"}, "/")
	log.Printf("[DeepLX] 请求URL: %s", query)
	log.Printf("[DeepLX] Authorization: %s", Authorization)

	headers := map[string]string{
		"Authorization": Authorization,
	}
	data := map[string]string{
		"text":        src,
		"source_lang": "ja",
		"target_lang": "zh",
	}
	log.Printf("[DeepLX] 请求数据: %+v", data)

	b, err := HttpPostJson(headers, data, query)
	if err != nil {
		log.Printf("[DeepLX] HTTP请求失败: %v", err)
		panic(err.Error())
	}
	log.Printf("[DeepLX] 原始响应: %s", string(b))

	// 解析 JSON 响应
	var response DeepLXResponse
	if err := json.Unmarshal(b, &response); err != nil {
		log.Printf("[DeepLX] JSON解析失败: %v", err)
		panic(err.Error())
	}

	log.Printf("[DeepLX] 解析结果 - Code: %d, ID: %d, Data: %s, Alternatives: %v",
		response.Code, response.ID, response.Data, response.Alternatives)

	return response.Data
}
