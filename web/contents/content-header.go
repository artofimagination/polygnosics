package contents

const (
	ParentPageKey      = "parent_page"
	ContentNameKey     = "content_name"
	QuickStats1Key     = "quick_stats_1"
	QuickStats1DataKey = "quick_stats_1_data"
	QuickStats2Key     = "quick_stats_2"
	QuickStats2DataKey = "quick_stats_2_data"
)

type QuickStat struct {
	statType string
	data     []int
}

func (c *ContentController) prepareContentHeader(content map[string]interface{}, parentPageName string, contentName string) map[string]interface{} {
	content[ParentPageKey] = parentPageName
	content[ContentNameKey] = contentName
	content[QuickStats1Key] = "users"
	content[QuickStats1DataKey] = []int{1, 2, 3, 5, 7, 11, 12, 14, 13, 12}
	content[QuickStats2Key] = "revenue"
	content[QuickStats2DataKey] = []int{1, 20, 3, 5, 7, 11, 12, 14, 15, 16}
	return content
}
