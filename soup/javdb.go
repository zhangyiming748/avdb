package soup

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"avdb/param"
	"avdb/storage"
	"avdb/trans"
	"avdb/util"

	"github.com/anaskhan96/soup"
)

var (
	html  = "https://javdb.com/search?q"
	proxy = "http://127.0.0.1:8889"
)

// loadCookiesFromDisk 加载磁盘上的 cookies
func loadCookiesFromDisk() ([]*http.Cookie, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("获取当前用户失败: %w", err)
	}
	cookieFilePath := filepath.Join(home, "javdb.cookie")
	// 读取 cookie 文件
	data, err := os.ReadFile(cookieFilePath)
	if err != nil {
		return nil, fmt.Errorf("读取cookie文件失败: %w", err)
	}

	// 解析 cookies
	var cookies []*http.Cookie
	err = json.Unmarshal(data, &cookies)
	if err != nil {
		return nil, fmt.Errorf("解析cookie文件失败: %w", err)
	}

	return cookies, nil
}

// createProxyClient 创建带代理的 HTTP 客户端
func createProxyClient(proxyURL string) (*http.Client, error) {
	proxyParsed, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("解析代理地址失败: %w", err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyParsed),
	}

	// 创建 cookie jar 来持久化 cookie
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("创建cookie jar失败: %w", err)
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
		Jar:       cookieJar, // 添加 cookie jar 以持久化 cookie
	}

	// 尝试加载已存在的 cookies
	cookies, err := loadCookiesFromDisk()
	if err == nil && len(cookies) > 0 {
		// 获取 JavDB 域名
		targetURL, _ := url.Parse("https://javdb.com")
		// 将 cookies 添加到 cookie jar
		cookieJar.SetCookies(targetURL, cookies)
	}

	return client, nil
}

func Javdb(keyword string) (dst string, err error) {
	//soup.SetDebug(true)
	// 创建带代理和 cookie jar 的客户端
	client, err := createProxyClient(proxy)
	if err != nil {
		return "", fmt.Errorf("创建代理客户端失败: %w", err)
	}

	// 直接进行搜索
	query := strings.Join([]string{html, keyword}, "=")
	log.Printf("请求的网址为: %v\n", query)
	resp, err := soup.GetWithClient(query, client)
	if err != nil {
		return "", fmt.Errorf("搜索请求失败: %w", err)
	}

	//在这里把这次请求到的网址保存为tmp.html文件
	f, err := os.OpenFile("tmp.html", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("打开HTML文件失败: %v\n", err)
		// 文件操作失败不影响返回结果，继续执行
	} else {
		defer f.Close()

		_, err = f.Write([]byte(resp))
		if err != nil {
			log.Printf("写入HTML文件失败: %v\n", err)
		} else {
			log.Println("HTML已保存到 tmp.html")
		}
	}
	root := soup.HTMLParse(resp)
	//log.Printf("root is : %v\n", root)
	//<div class="item">
	items := root.FindAll("div", "class", "item")
	log.Printf("found %d items\n", len(items))
	for i, item := range items {
		log.Printf("item %d: %v\n", i, item)
		videoTitleDiv := item.Find("div", "class", "video-title")
		if videoTitleDiv.Pointer != nil {

			titleText := videoTitleDiv.Text()
			log.Printf("完整标题行: %v\n", titleText)

			// 查找 strong 标签获取番号
			strongTag := videoTitleDiv.Find("strong")
			if strongTag.Pointer != nil {
				idNumber := strongTag.Text()
				log.Printf("番号: %v\n", idNumber)

				// 提取标题部分（去掉番号的部分）
				titleWithoutId := strings.TrimSpace(strings.TrimPrefix(titleText, idNumber))
				log.Printf("标题: %v\n", titleWithoutId)
				// 翻译为中文标题'
				authorization := param.GetVal("trans", "Authorization")
				zhcnTitle := trans.DeepLX(titleWithoutId, authorization)
				log.Printf("中文标题: %v\n", zhcnTitle)
				pretty := strings.Join([]string{idNumber, zhcnTitle}, " ")
				if i == 0 {
					dst = pretty
				}
				avdb := storage.AVDB{
					NO:        idNumber,
					Title:     titleWithoutId,
					ZhCnTitle: zhcnTitle,
					Pretty:    pretty,
				}
				log.Printf("准备插入数据库的AVDB: %+v\n", avdb)
				err := avdb.Insert()
				if err != nil {
					log.Printf("插入数据库失败: %v\n", err)
				} else {
					log.Printf("插入数据库成功: %+v\n", avdb)
				}
			} else {
				log.Printf("未找到番号\n")
			}
		} else {
			log.Printf("未找到 video-title 元素\n")
		}
	}
	return dst, nil
}

func SearchByJavdb(src, dst string) (string, error) {
	// 首先判断dst是一个目录还是文件
	//如果是文件形式 即使这个文件不存在 直接用作输出文件，覆盖写入
	//如果是目录 就在这个目录下新建一个export.txt作为输出文件
	info, err := os.Stat(dst)
	var outputPath string
	if err != nil {
		// 如果路径不存在，将其视为文件路径
		if os.IsNotExist(err) {
			// 确保父目录存在
			parentDir := filepath.Dir(dst)
			err = os.MkdirAll(parentDir, 0755)
			if err != nil {
				return "", fmt.Errorf("创建目录失败: %w", err)
			}
			outputPath = dst
		} else {
			return "", fmt.Errorf("检查路径失败: %w", err)
		}
	} else {
		if info.IsDir() {
			// 如果dst是目录，则在该目录下创建export.txt文件
			outputPath = filepath.Join(dst, "export.txt")
		} else {
			// 如果dst是文件，则直接使用该文件
			outputPath = dst
		}
	}

	keys := util.ReadByLine(src)
	s := []string{}
	for _, keyword := range keys {
		if pretty, err := Javdb(keyword); err == nil {
			s = append(s, pretty)
		} else {
			log.Printf("搜索 %v 失败: %v\n", keyword, err)
		}
	}
	util.WriteByLine(outputPath, s)
	return outputPath, nil
}
