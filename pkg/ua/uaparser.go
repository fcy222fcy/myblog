package ua

import (
	"regexp"
	"strings"
)

type Info struct {
	OS             string
	OSVersion      string
	Browser        string
	BrowserVersion string
}

var emptyInfo = Info{OS: "未知", OSVersion: "", Browser: "未知", BrowserVersion: ""}

// Parse 从 User-Agent 字符串解析操作系统（含版本）和浏览器（含版本）
func Parse(userAgent string) Info {
	if userAgent == "" {
		return emptyInfo
	}
	ua := userAgent
	os, osVer := parseOS(ua)
	browser, browserVer := parseBrowser(ua)
	return Info{
		OS:             os,
		OSVersion:      toTwoPartVersion(osVer),
		Browser:        browser,
		BrowserVersion: toTwoPartVersion(browserVer),
	}
}

// toTwoPartVersion 把版本号统一成「主版本.次版本」两位格式
//
//	11          → 11.0
//	126.0.6099  → 126.0
//	17.5        → 17.5
//	XP / (iOS) / 14 (Sonoma)  → 非纯数字或带名字的：保留原名+版本号按前两位处理
func toTwoPartVersion(ver string) string {
	ver = strings.TrimSpace(ver)
	if ver == "" {
		return ""
	}
	// 首字符非数字（XP、(iOS)）原样返回
	if ver[0] < '0' || ver[0] > '9' {
		return ver
	}
	// 带空格/括号拼接的版本名（如 "14 (Sonoma)" / "iOS 150.1"）：拆分处理
	if idx := strings.IndexAny(ver, " ("); idx > 0 {
		numPart := strings.TrimSpace(ver[:idx])
		suffix := ver[idx:]
		two := toTwoPartVersion(numPart)
		if two == "" {
			return ver
		}
		return two + suffix
	}
	dot := strings.Index(ver, ".")
	if dot < 0 {
		// 纯整数 → 补 .0
		for _, ch := range ver {
			if ch < '0' || ch > '9' {
				return ver
			}
		}
		return ver + ".0"
	}
	major := ver[:dot]
	rest := ver[dot+1:]
	for _, ch := range major {
		if ch < '0' || ch > '9' {
			return ver
		}
	}
	var minor strings.Builder
	for _, ch := range rest {
		if ch >= '0' && ch <= '9' {
			minor.WriteRune(ch)
		} else {
			break
		}
	}
	minorStr := minor.String()
	if minorStr == "" {
		minorStr = "0"
	}
	return major + "." + minorStr
}

// FormatOS 返回格式化的操作系统 + 版本（如果有版本号）
func (i Info) FormatOS() string {
	name := strings.TrimSpace(i.OS)
	ver := strings.TrimSpace(i.OSVersion)
	if name == "" || name == "未知" {
		return ""
	}
	if ver == "" {
		return name
	}
	return name + " " + ver
}

// FormatBrowser 返回格式化的浏览器 + 版本（如果有版本号）
func (i Info) FormatBrowser() string {
	name := strings.TrimSpace(i.Browser)
	ver := strings.TrimSpace(i.BrowserVersion)
	if name == "" || name == "未知" {
		return ""
	}
	if ver == "" {
		return name
	}
	return name + " " + ver
}

// 正则工具：提取某个 token 后面的版本号（如 "Chrome/120.0.6099" -> "120"）
func extractVersion(ua, prefix string, takeFirstNumOnly bool) string {
	idx := strings.Index(ua, prefix)
	if idx < 0 {
		return ""
	}
	rest := ua[idx+len(prefix):]
	// 取到空格、分号、右括号前的内容
	end := len(rest)
	for j, ch := range rest {
		if ch == ' ' || ch == ';' || ch == ')' {
			end = j
			break
		}
	}
	version := strings.TrimSpace(rest[:end])
	if version == "" {
		return ""
	}
	if takeFirstNumOnly {
		// 只取小数点前的主版本号
		if dot := strings.Index(version, "."); dot > 0 {
			return version[:dot]
		}
	}
	return version
}

// —— 操作系统解析 ——
func parseOS(ua string) (string, string) {
	// Windows 版本
	if strings.Contains(ua, "Windows") {
		switch {
		case strings.Contains(ua, "Windows NT 10.0"):
			// Win10/11 都是 NT 10.0，靠 Build 号区分：22000+ 是 Win11
			if build := extractVersion(ua, "Windows NT 10.0; Win64; x64; rv:", false); build != "" {
			}
			if strings.Contains(ua, "Windows NT 10.0; Win64; x64") && (strings.Contains(ua, "rv:") || strings.Contains(ua, "Gecko/") || strings.Contains(ua, "Chrome/")) {
				// 粗略：尝试找 build 号，没找到就统一返回 10/11
			}
			build := extractBuildFromNT10(ua)
			if build != "" && compareBuild(build, 22000) >= 0 {
				return "Windows", "11"
			}
			return "Windows", "10"
		case strings.Contains(ua, "Windows NT 6.3"):
			return "Windows", "8.1"
		case strings.Contains(ua, "Windows NT 6.2"):
			return "Windows", "8"
		case strings.Contains(ua, "Windows NT 6.1"):
			return "Windows", "7"
		case strings.Contains(ua, "Windows NT 6.0"):
			return "Windows", "Vista"
		case strings.Contains(ua, "Windows NT 5.1"):
			return "Windows", "XP"
		default:
			return "Windows", ""
		}
	}

	// macOS
	if strings.Contains(ua, "Mac OS X") && !strings.Contains(ua, "iPhone") && !strings.Contains(ua, "iPad") {
		// Mac OS X 10_15_7 -> 10.15.7 (Catalina) ; macOS 11_0 -> Big Sur / 12 -> Monterey 等
		if ver := extractVersion(ua, "Mac OS X ", false); ver != "" {
			ver = strings.ReplaceAll(ver, "_", ".")
			// 截取前两个段：14.4.1 -> 14
			parts := strings.Split(ver, ".")
			major := ""
			if len(parts) >= 1 {
				major = parts[0]
			}
			switch major {
			case "10":
				// 10.15 Catalina, 10.14 Mojave, ...
				if len(parts) >= 2 {
					return "macOS", major + "." + parts[1]
				}
				return "macOS", ver
			case "11":
				return "macOS", "11 (Big Sur)"
			case "12":
				return "macOS", "12 (Monterey)"
			case "13":
				return "macOS", "13 (Ventura)"
			case "14":
				return "macOS", "14 (Sonoma)"
			case "15":
				return "macOS", "15 (Sequoia)"
			default:
				if major != "" {
					return "macOS", major
				}
				return "macOS", ver
			}
		}
		return "macOS", ""
	}

	// iOS (iPhone/iPod)
	if strings.Contains(ua, "iPhone") || strings.Contains(ua, "iPod") {
		if ver := extractVersion(ua, "iPhone OS ", false); ver != "" {
			ver = strings.ReplaceAll(ver, "_", ".")
			parts := strings.Split(ver, ".")
			if len(parts) >= 1 {
				return "iOS", parts[0]
			}
			return "iOS", ver
		}
		if ver := extractVersion(ua, "OS ", false); ver != "" {
			ver = strings.ReplaceAll(ver, "_", ".")
			parts := strings.Split(ver, ".")
			if len(parts) >= 1 {
				return "iOS", parts[0]
			}
			return "iOS", ver
		}
		return "iOS", ""
	}

	// iPadOS
	if strings.Contains(ua, "iPad") {
		if ver := extractVersion(ua, "OS ", false); ver != "" {
			ver = strings.ReplaceAll(ver, "_", ".")
			parts := strings.Split(ver, ".")
			if len(parts) >= 1 {
				return "iPadOS", parts[0]
			}
			return "iPadOS", ver
		}
		return "iPadOS", ""
	}

	// HarmonyOS / 鸿蒙
	if strings.Contains(ua, "HarmonyOS") || strings.Contains(ua, "Harmony") {
		if ver := extractVersion(ua, "HarmonyOS/", true); ver != "" {
			return "鸿蒙", ver
		}
		if ver := extractVersion(ua, "Harmony/", true); ver != "" {
			return "鸿蒙", ver
		}
		return "鸿蒙", ""
	}

	// Android
	if strings.Contains(ua, "Android") {
		if ver := extractVersion(ua, "Android ", true); ver != "" {
			return "Android", ver
		}
		// "Android/13"
		if ver := extractVersion(ua, "Android/", true); ver != "" {
			return "Android", ver
		}
		return "Android", ""
	}

	// Linux 发行版
	switch {
	case strings.Contains(ua, "Ubuntu"):
		if ver := extractVersion(ua, "Ubuntu/", true); ver != "" {
			return "Ubuntu", ver
		}
		return "Ubuntu", ""
	case strings.Contains(ua, "Fedora"):
		return "Fedora", ""
	case strings.Contains(ua, "Debian"):
		return "Debian", ""
	case strings.Contains(ua, "FreeBSD"):
		return "FreeBSD", ""
	case strings.Contains(ua, "OpenBSD"):
		return "OpenBSD", ""
	case strings.Contains(ua, "Linux"):
		return "Linux", ""
	}

	return "未知", ""
}

// —— 浏览器解析（带版本） ——
func parseBrowser(ua string) (string, string) {
	// 微信
	if strings.Contains(ua, "MicroMessenger") {
		if ver := extractVersion(ua, "MicroMessenger/", true); ver != "" {
			return "微信", ver
		}
		return "微信", ""
	}
	// QQ
	if strings.Contains(ua, "QQBrowser") {
		if ver := extractVersion(ua, "QQBrowser/", true); ver != "" {
			return "QQ浏览器", ver
		}
		return "QQ浏览器", ""
	}
	if strings.Contains(ua, " QQ/") {
		if ver := extractVersion(ua, " QQ/", true); ver != "" {
			return "QQ", ver
		}
	}
	// UC
	if strings.Contains(ua, "UCBrowser") {
		if ver := extractVersion(ua, "UCBrowser/", true); ver != "" {
			return "UC浏览器", ver
		}
		return "UC浏览器", ""
	}
	// 钉钉
	if strings.Contains(ua, "DingTalk") {
		if ver := extractVersion(ua, "DingTalk/", true); ver != "" {
			return "钉钉", ver
		}
		return "钉钉", ""
	}
	// 支付宝
	if strings.Contains(ua, "AlipayClient") {
		if ver := extractVersion(ua, "AlipayClient/", true); ver != "" {
			return "支付宝", ver
		}
		return "支付宝", ""
	}
	// 搜狗
	if strings.Contains(ua, "Sogou") || strings.Contains(ua, "SogouMobileBrowser") {
		return "搜狗浏览器", ""
	}
	// 360
	if strings.Contains(ua, "360browser") || strings.Contains(ua, "QihooBrowser") || strings.Contains(ua, "QIHU") {
		return "360浏览器", ""
	}
	// Edge（新）
	if strings.Contains(ua, "Edg/") {
		if ver := extractVersion(ua, "Edg/", true); ver != "" {
			return "Edge", ver
		}
		return "Edge", ""
	}
	if strings.Contains(ua, "Edge/") {
		if ver := extractVersion(ua, "Edge/", true); ver != "" {
			return "Edge", ver
		}
		return "Edge", ""
	}
	// Vivaldi
	if strings.Contains(ua, "Vivaldi/") {
		if ver := extractVersion(ua, "Vivaldi/", true); ver != "" {
			return "Vivaldi", ver
		}
		return "Vivaldi", ""
	}
	// Opera
	if strings.Contains(ua, "OPR/") {
		if ver := extractVersion(ua, "OPR/", true); ver != "" {
			return "Opera", ver
		}
		return "Opera", ""
	}
	if strings.Contains(ua, "Opera/") {
		if ver := extractVersion(ua, "Opera/", true); ver != "" {
			return "Opera", ver
		}
		return "Opera", ""
	}
	// Firefox
	if strings.Contains(ua, "Firefox/") {
		if ver := extractVersion(ua, "Firefox/", true); ver != "" {
			return "Firefox", ver
		}
		return "Firefox", ""
	}
	if strings.Contains(ua, "FxiOS/") {
		if ver := extractVersion(ua, "FxiOS/", true); ver != "" {
			return "Firefox", "iOS " + ver
		}
		return "Firefox", "(iOS)"
	}
	// Chrome iOS
	if strings.Contains(ua, "CriOS/") {
		if ver := extractVersion(ua, "CriOS/", true); ver != "" {
			return "Chrome", "iOS " + ver
		}
		return "Chrome", "(iOS)"
	}
	// Brave
	if strings.Contains(ua, "Brave/") {
		return "Brave", ""
	}
	// IE
	if strings.Contains(ua, "MSIE ") {
		if ver := extractVersion(ua, "MSIE ", true); ver != "" {
			return "IE", ver
		}
		return "IE", ""
	}
	if strings.Contains(ua, "Trident/") {
		// IE11: Trident/7.0; rv:11.0
		if ver := extractVersion(ua, "rv:", true); ver != "" {
			return "IE", ver
		}
		return "IE", "11"
	}
	// Chrome（最后判断，因为很多浏览器 UA 里都带 Chrome）
	if strings.Contains(ua, "Chrome/") && !strings.Contains(ua, "Edg/") {
		if ver := extractVersion(ua, "Chrome/", true); ver != "" {
			return "Chrome", ver
		}
		return "Chrome", ""
	}
	// Safari（也要在 Chrome/Edg 之后判断）
	if strings.Contains(ua, "Safari/") && !strings.Contains(ua, "Chrome/") && !strings.Contains(ua, "Edg/") {
		if strings.Contains(ua, "Version/") {
			if ver := extractVersion(ua, "Version/", true); ver != "" {
				return "Safari", ver
			}
		}
		return "Safari", ""
	}

	return "未知", ""
}

// —— Windows Build 号辅助工具 ——
var nt10BuildRe = regexp.MustCompile(`Windows NT 10\.0;[^)]*Build\s+(\d+)`)

func extractBuildFromNT10(ua string) string {
	sub := nt10BuildRe.FindStringSubmatch(ua)
	if len(sub) >= 2 {
		return sub[1]
	}
	// Edge UA 常见格式：Windows NT 10.0; Win64; x64
	// Firefox: Windows NT 10.0; Win64; x64; rv:
	// 没 Build 号就用另一种方式：Chrome/120.0.0.0 Windows NT 10.0
	// 统一返回空，外面兜底
	return ""
}

func compareBuild(buildStr string, threshold int) int {
	var b int
	for _, ch := range buildStr {
		if ch >= '0' && ch <= '9' {
			b = b*10 + int(ch-'0')
		} else {
			break
		}
	}
	if b > threshold {
		return 1
	}
	if b == threshold {
		return 0
	}
	return -1
}
