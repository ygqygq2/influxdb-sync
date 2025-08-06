package utils

// ContainsString 判断字符串是否在列表中
func ContainsString(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
