package ast

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

func requirePluginGVAFixture(t *testing.T, extra ...string) string {
	t.Helper()
	parts := []string{global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server, "plugin", "gva"}
	parts = append(parts, extra...)
	target := filepath.Join(parts...)
	if _, err := os.Stat(target); err != nil {
		t.Skipf("skip plugin gva fixture dependent test: %v", err)
	}
	return target
}
