package trans

import (
	"avdb/param"
	"testing"
)

func TestDeepLX(t *testing.T) {
	param.InitConfig("C:\\Users\\zhang\\Github\\AVDB\\avdb.ini")
	Authorization := param.GetVal("trans", "Authorization")
	t.Logf("从配置读取的 Authorization: %s", Authorization)

	if Authorization == "" || Authorization == "not found" {
		t.Fatal("Authorization 未配置或为空")
	}

	ret := DeepLX("hello", Authorization)
	t.Logf("翻译结果: %s", ret)

	if ret == "" {
		t.Error("翻译结果为空")
	}
}
