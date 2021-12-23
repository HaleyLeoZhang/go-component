package helper

import "testing"

func TestGetLocalIpV4(t *testing.T) {
	ip := GetLocalIpV4()
	t.Logf("ip %v", ip)
}
