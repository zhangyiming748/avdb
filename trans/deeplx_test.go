package trans

import "testing"

func TestDeepLX(t *testing.T) {
	Authorization := "{{ Authorization }}"
	ret := DeepLX("hello", Authorization)
	t.Log(ret)
}
