package zapmarshaller

const securedTagName = "secured"

func secured(tags map[string]string) bool {
	for k, _ := range tags {
		if k == securedTagName {
			return true
		}
	}
	return false
}
