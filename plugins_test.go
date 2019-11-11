package gomplate

import (
	"context"
	"testing"
	"text/template"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"

	"github.com/hairyhenderson/gomplate/v3/internal/config"
)

func TestBindPlugins(t *testing.T) {
	ctx := context.TODO()
	fm := template.FuncMap{}
	cfg := &config.Config{
		Plugins: map[string]string{},
	}
	err := bindPlugins(ctx, cfg, fm)
	assert.NilError(t, err)
	assert.DeepEqual(t, template.FuncMap{}, fm)

	cfg.Plugins = map[string]string{"foo": "bar"}
	err = bindPlugins(ctx, cfg, fm)
	assert.NilError(t, err)
	assert.Check(t, cmp.Contains(fm, "foo"))

	err = bindPlugins(ctx, cfg, fm)
	assert.ErrorContains(t, err, "already bound")
}

func TestBuildCommand(t *testing.T) {
	ctx := context.TODO()
	data := []struct {
		name, path string
		args       []string
		expected   []string
	}{
		{"foo", "foo", nil, []string{"foo"}},
		{"foo", "foo", []string{"bar"}, []string{"foo", "bar"}},
		{"foo", "foo.bat", nil, []string{"cmd.exe", "/c", "foo.bat"}},
		{"foo", "foo.cmd", []string{"bar"}, []string{"cmd.exe", "/c", "foo.cmd", "bar"}},
		{"foo", "foo.ps1", []string{"bar", "baz"}, []string{"pwsh", "-File", "foo.ps1", "bar", "baz"}},
	}
	for _, d := range data {
		p := &plugin{
			ctx:  ctx,
			name: d.name,
			path: d.path,
		}
		name, args := p.buildCommand(d.args)
		actual := append([]string{name}, args...)
		assert.DeepEqual(t, d.expected, actual)
	}
}
