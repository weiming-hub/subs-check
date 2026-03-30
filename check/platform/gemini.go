package platform

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/biter777/countries"
)

// alpha3ToAlpha2 使用 countries 库将三字码转换为二字码
func alpha3ToAlpha2(alpha3 string) string {
	code := strings.ToUpper(alpha3)
	country := countries.ByName(code)
	if country == countries.Unknown {
		return ""
	}
	return country.Alpha2()
}

// https://github.com/clash-verge-rev/clash-verge-rev/blob/c894a15d13d5bcce518f8412cc393b56272a9afa/src-tauri/src/cmd/media_unlock_checker.rs#L241
func CheckGemini(httpClient *http.Client) (string, error) {
	req, err := http.NewRequest("GET", "https://gemini.google.com/", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	bodyStr := string(body)

	// 检查是否可用
	if !strings.Contains(bodyStr, "45631641,null,true") {
		return "", nil
	}

	// 使用正则表达式提取国家代码
	re := regexp.MustCompile(`,2,1,200,"([A-Z]{3})"`)
	matches := re.FindStringSubmatch(bodyStr)

	if len(matches) > 1 {
		alpha3Code := matches[1] // 三字码，如 "USA"
		alpha2Code := alpha3ToAlpha2(alpha3Code)
		if alpha2Code == "" {
			return "N/A", nil
		}
		return alpha2Code, nil
	}

	return "N/A", nil
}
