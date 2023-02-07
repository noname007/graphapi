package graphapi

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	NewClient(Conf{
		TenantId:     "",
		ClientId:     "",
		ClientSecret: "",
	})
}

func TestClient_DeleteEmail(t *testing.T) {
	c := NewClient(Conf{
		TenantId:     "",
		ClientId:     "",
		ClientSecret: "",
	})
	c.FetchEmail(context.TODO(), "xx@xx.com", "域账号密码到期提醒")
}
