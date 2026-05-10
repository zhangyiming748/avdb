package trans

import (
	"encoding/json"
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
	headers := map[string]string{
		"Authorization": Authorization,
	}
	data := map[string]string{
		"text":        src,
		"source_lang": "ja",
		"target_lang": "zh",
	}
	b, err := HttpPostJson(headers, data, query)
	if err != nil {
		panic(err.Error())
	}

	// 解析 JSON 响应
	var response DeepLXResponse
	if err := json.Unmarshal(b, &response); err != nil {
		panic(err.Error())
	}

	return response.Data
}
