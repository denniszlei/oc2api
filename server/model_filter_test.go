package main

import (
	"testing"
	"time"
)

// 目录解析: 只认 cost 输入/输出都是数字 0; 缺 cost、只有一边为 0、cost 为 null/字符串 都不算免费。
func TestExtractCatalogFreeModels(t *testing.T) {
	raw := `{"opencode":{"models":{
		"a-zero":{"cost":{"input":0,"output":0}},
		"b-half":{"cost":{"input":0,"output":1}},
		"c-nocost":{"id":"c-nocost"},
		"d-null-cost":{"cost":null},
		"e-null-input":{"cost":{"input":null,"output":0}},
		"f-string":{"cost":{"input":"0","output":"0"}}
	}}}`
	free := extractCatalogFreeModels(raw)
	if len(free) != 1 || !free["a-zero"] {
		t.Fatalf("只有 a-zero 应被判为免费, 实际 %v", free)
	}

	// 解析失败返回 nil(调用方据此退避并沿用旧集合)
	for _, bad := range []string{`not json`, `{}`, `{"opencode":{}}`} {
		if got := extractCatalogFreeModels(bad); got != nil {
			t.Fatalf("输入 %q 应返回 nil, 实际 %v", bad, got)
		}
	}

	// 解析成功但一个免费模型都没有: 必须是空 map 而不是 nil, 否则会被当成抓取失败
	empty := extractCatalogFreeModels(`{"opencode":{"models":{}}}`)
	if empty == nil || len(empty) != 0 {
		t.Fatalf("空目录应返回非 nil 空集合, 实际 %v", empty)
	}
}

// 判定优先级: BLOCK 最高, 其次 EXTRA, 再静态规则, 最后目录规则。
func TestIsAllowedModelId(t *testing.T) {
	extra := map[string]bool{"union-alpha": true}
	blocked := map[string]bool{"mimo-v2.5-free": true}
	catalogFree := map[string]bool{"grok-code": true}

	cases := []struct {
		id   string
		want bool
	}{
		{"big-pickle", true},      // 静态规则
		{"hy3-free", true},        // 静态规则
		{"grok-code", true},       // 目录规则: 免费但名字不带 -free
		{"union-alpha", true},     // EXTRA_MODELS 兜底限免隐身模型
		{"mimo-v2.5-free", false}, // BLOCK_MODELS 优先于静态规则
		{"gpt-5.5", false},        // 目录里 cost 非 0, 付费
		{"", false},               // 空 id
	}

	for _, c := range cases {
		if got := isAllowedModelId(c.id, catalogFree, extra, blocked); got != c.want {
			t.Errorf("isAllowedModelId(%q) = %v, 期望 %v", c.id, got, c.want)
		}
	}
}

// 缓存 TTL 解析: 空值/空格/非十进制整数/负值都回退默认 10 分钟。
func TestParseModelCacheTTL(t *testing.T) {
	cases := []struct {
		raw  string
		want time.Duration
	}{
		{"", defaultModelCacheTTL},
		{"   ", defaultModelCacheTTL},
		{"0", 0},
		{"1000", time.Second},
		{"600000", 10 * time.Minute},
		{"1e3", defaultModelCacheTTL},
		{"0x10", defaultModelCacheTTL},
		{"-1", defaultModelCacheTTL},
		{"abc", defaultModelCacheTTL},
	}

	for _, c := range cases {
		if got := parseModelCacheTTL(c.raw); got != c.want {
			t.Errorf("parseModelCacheTTL(%q) = %v, 期望 %v", c.raw, got, c.want)
		}
	}
}
