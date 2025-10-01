package httpclient_test

import (
	"testing"

	"github.com/ProImpact/service-check/pkg/httpclient"
)

func TestMakeRequest(t *testing.T) {
	resp := httpclient.MakeRequest("http://localhost:3001")
	if resp == nil {
		t.Fatal("there is no socket listening on that port")
	}
	t.Logf("%+v", resp)
}
