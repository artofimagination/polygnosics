package contents

const (
	NewsKey     = "news"
	GroupKey    = "group_index"
	PostKey     = "post_key"
	TextKey     = "news_text"
	YearKey     = "news_year"
	MonthKey    = "news_month"
	DayKey      = "news_day"
	CommentsKey = "comments"
)

// prepareNewsFeed creates the UI content from DB content
// Structure:
//		Content->News content map
//									-> Year key
//												-> Post key
//													-> Month
//													-> Day
//													-> Text
func (*ContentController) prepareNewsFeed(content map[string]interface{}) map[string]interface{} {
	groupIndex := 0
	content[NewsKey] = make(map[int]interface{})
	groups := content[NewsKey].(map[int]interface{})
	groups[groupIndex] = make(map[string]interface{})
	group := groups[groupIndex].(map[string]interface{})
	group[YearKey] = "2015"
	group[PostKey] = make(map[string]interface{})
	post := group[PostKey].(map[string]interface{})
	post["0"] = make(map[string]interface{})
	post["0"].(map[string]interface{})[TextKey] = "This is first news"
	post["0"].(map[string]interface{})[MonthKey] = "Feb"
	post["0"].(map[string]interface{})[DayKey] = "5"
	post["1"] = make(map[string]interface{})
	post["1"].(map[string]interface{})[TextKey] = "This is second news"
	post["1"].(map[string]interface{})[MonthKey] = "Mar"
	post["1"].(map[string]interface{})[DayKey] = "15"
	groupIndex++
	groups[groupIndex] = make(map[string]interface{})
	group = groups[groupIndex].(map[string]interface{})
	group[YearKey] = "2019"
	group[PostKey] = make(map[string]interface{})
	post = group[PostKey].(map[string]interface{})
	post["2"] = make(map[string]interface{})
	post["2"].(map[string]interface{})[TextKey] = "This is third news"
	post["2"].(map[string]interface{})[MonthKey] = "Jun"
	post["2"].(map[string]interface{})[DayKey] = "4"
	post["3"] = make(map[string]interface{})
	post["3"].(map[string]interface{})[TextKey] = "This is first news"
	post["3"].(map[string]interface{})[MonthKey] = "Sep"
	post["3"].(map[string]interface{})[DayKey] = "9"
	return content
}
