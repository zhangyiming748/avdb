package param


import(
	"github.com/zhangyiming748/goini"
)
var (
	conf    *goini.Config
)
func InitConfig(configPath string) {
	conf = goini.SetConfig(configPath)
}

func GetVal(section, name string) string {
	val, err := conf.GetValue(section, name)
	if err != nil {
		return "not found"
	}
	return val
}