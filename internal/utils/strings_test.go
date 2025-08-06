package utils

import "testing"

func TestContainsString(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []string
		str      string
		expected bool
	}{
		{"找到字符串", []string{"a", "b", "c"}, "b", true},
		{"未找到字符串", []string{"a", "b", "c"}, "d", false},
		{"空切片", []string{}, "a", false},
		{"空字符串在非空切片", []string{"a", "b", "c"}, "", false},
		{"空字符串在包含空字符串的切片", []string{"a", "", "c"}, "", true},
		{"单元素切片匹配", []string{"test"}, "test", true},
		{"单元素切片不匹配", []string{"test"}, "other", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ContainsString(tc.slice, tc.str)
			if result != tc.expected {
				t.Errorf("ContainsString(%v, %s) = %v, 期望 %v",
					tc.slice, tc.str, result, tc.expected)
			}
		})
	}
}
