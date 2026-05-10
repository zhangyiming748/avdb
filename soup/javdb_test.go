package soup

import (
	"log"
	"os"
	"strings"
	"testing"

	"avdb/param"
	"avdb/storage"
	"github.com/anaskhan96/soup"
)

func TestJavdb(t *testing.T) {
	param.InitConfig("C:\\Users\\zhang\\Github\\AVDB\\avdb.ini") // 请根据实际情况修改配置文件路径
	storage.SetSqlite()
	result, err := Javdb("OEA-002")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("Result: %v", result)
}

func TestParse(t *testing.T) {
	// 读取HTML文件内容
	data, err := os.ReadFile("/Users/zen/github/avdb/soup/tmp.html")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
		return
	}

	// 解析HTML内容
	root := soup.HTMLParse(string(data))
	log.Printf("root is : %v\n", root)
	//<div class="item">
	items := root.FindAll("div", "class", "item")
	log.Printf("found %d items\n", len(items))
	for i, item := range items {
		log.Printf("item %d: %v\n", i, item)
		videoTitleDiv := item.Find("div", "class", "video-title")
		if videoTitleDiv.Pointer != nil {
			titleText := videoTitleDiv.Text()
			t.Logf("完整标题行: %v\n", titleText)
			// 查找 strong 标签获取番号
			strongTag := videoTitleDiv.Find("strong")
			if strongTag.Pointer != nil {
				idNumber := strongTag.Text()
				t.Logf("番号: %v\n", idNumber)
				// 提取标题部分（去掉番号的部分）
				titleWithoutId := strings.TrimSpace(strings.TrimPrefix(titleText, idNumber))
				t.Logf("标题: %v\n", titleWithoutId)
			} else {
				t.Logf("未找到番号\n")
			}
		} else {
			t.Logf("未找到 video-title 元素\n")
		}
	}
}
