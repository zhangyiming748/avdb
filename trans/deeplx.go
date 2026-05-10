package trans

var host = "https://api.deeplx.org/H4x25H69HomF0LS4-gVNX-SdEmFZksp_o7rjaXCxg6A/translate"
func DeepLX(src string) (dst string) {
	
	data := map[string]string{
		"text": src,
		"source_lang": "ja",
		"target_lang": "zh",
	}
	b, err := HttpGet(nil, data, host)
	if err != nil {
		panic(err.Error())
	}
	return string(b)
}
