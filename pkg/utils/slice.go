package utils

// Contains 检查切片中是否包含某个元素
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
