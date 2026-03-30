package check

import (
	"fmt"
	"log/slog"
	"regexp"

	"github.com/beck-8/subs-check/config"
)

// FilterResults 根据配置的正则表达式过滤节点
// 只有节点名称匹配任一正则表达式的节点才会被保留
func FilterResults(results []Result) []Result {
	// 如果没有配置过滤规则，直接返回所有结果
	if len(config.GlobalConfig.Filter) == 0 {
		return results
	}

	// 编译所有正则表达式
	var patterns []*regexp.Regexp
	for _, pattern := range config.GlobalConfig.Filter {
		re, err := regexp.Compile(pattern)
		if err != nil {
			slog.Warn(fmt.Sprintf("过滤正则表达式编译失败，已跳过: %s, 错误: %v", pattern, err))
			continue
		}
		patterns = append(patterns, re)
	}

	// 如果所有正则都编译失败，返回所有结果
	if len(patterns) == 0 {
		slog.Warn("所有过滤正则表达式编译失败，跳过过滤")
		return results
	}

	slog.Info(fmt.Sprintf("应用节点过滤规则，共 %d 个正则表达式", len(patterns)))

	// 过滤结果
	var filtered []Result
	for _, result := range results {
		if result.Proxy == nil {
			continue
		}

		name, ok := result.Proxy["name"].(string)
		if !ok {
			continue
		}

		// 检查节点名称是否匹配任一正则表达式
		for _, re := range patterns {
			if re.MatchString(name) {
				filtered = append(filtered, result)
				break
			}
		}
	}

	slog.Info(fmt.Sprintf("过滤后节点数量: %d (过滤前: %d)", len(filtered), len(results)))
	return filtered
}
