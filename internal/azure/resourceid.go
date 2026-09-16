package azure

func nthSlash(s string, n int) int {
	count := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '/' {
			count++
			if count == n {
				return i
			}
		}
	}
	return -1
}

func SubscriptionFromResourceID(resourceID string) string {
	if len(resourceID) >= 52 && resourceID[51] == '/' {
		return resourceID[15:51]
	}
	s2 := nthSlash(resourceID, 2)
	s3 := nthSlash(resourceID, 3)
	if s2 < 0 || s3 < 0 {
		return ""
	}
	return resourceID[s2+1 : s3]
}

func ResourceGroupFromResourceID(resourceID string) string {
	s4 := nthSlash(resourceID, 4)
	s5 := nthSlash(resourceID, 5)
	if s4 < 0 || s5 < 0 {
		return ""
	}
	return resourceID[s4+1 : s5]
}

func ResourceGroupIDFromResourceID(resourceID string) string {
	s5 := nthSlash(resourceID, 5)
	if s5 < 0 {
		return ""
	}
	return resourceID[:s5]
}

func ResourceTypeFromResourceID(resourceID string) string {
	s6 := nthSlash(resourceID, 6)
	s7 := nthSlash(resourceID, 7)
	if s6 < 0 || s7 < 0 {
		return ""
	}
	s8 := nthSlash(resourceID, 8)
	end := len(resourceID)
	if s8 >= 0 {
		end = s8
	}
	return resourceID[s6+1 : end]
}

func ResourceNameFromResourceID(resourceID string) string {
	s8 := nthSlash(resourceID, 8)
	if s8 < 0 {
		return ""
	}
	s9 := nthSlash(resourceID, 9)
	if s9 < 0 {
		return resourceID[s8+1:]
	}
	return resourceID[s8+1 : s9]
}
