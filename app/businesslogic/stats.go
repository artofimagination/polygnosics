package businesslogic

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
)

const (
	userDataChannel      = "user-data"
	usersTotal           = "users_total"
	usersDeleted         = "users_deleted"
	usersOnline          = "users_online"
	usersDevClient       = "users_dev_client"
	usersActivityHistory = "users_activity_history"
	usersOnlinePeak      = "users_online_peaks"
	usersOnlinePeriod    = "users_online_period"
	usersOnlineTime      = "user_online_time"
)
const (
	projectDataChannel = "project-data"
	projectViewers     = "project_viewers"
	projectConnections = "project_connections"
	projectPins        = "project_pins"
	projectPopularity  = "project_popularity" // Returns how many people are viewing/pining, and connecting compared to overall
	projectSuccess     = "project_success"    // Returns the stats showing how many projects have actually been connected to after viewing or pinning
	projectUsageTime   = "project_usage_time"
	projectRating      = "project_rating"
)

const (
	productDataChannel       = "product-data"
	productViewers           = "product_viewers"
	productPurchase          = "product_purchase"
	productPins              = "product_pins"
	productPopularity        = "product_popularity" // Returns how many people are viewing/pining, and purchasing compared to overall
	productSuccess           = "product_success"    // Returns the stats showing how many products have actually been purchased after viewing or watching
	productProjectGeneration = "product_project_generation"
	productRating            = "product_rating"
)

const (
	itemDataChannel            = "item-data"
	itemsUsersProjectActivity  = "users_project_activity"
	itemsProjectLength         = "project_length"
	itemsProjectCount          = "project_count"
	itemsProductCount          = "product_count"
	itemsCurrentProductCount   = "product_current_count"
	itemsCurrentProjectCount   = "project_current_count"
	itemsProductPerUser        = "products_per_user"
	itemsProjectPerUser        = "prodjets_per_user"
	itemsCurrentProductPerUser = "current_products_per_user"
	itemsCurrentProjectPerUser = "current_prodjets_per_user"
)

const (
	Min            = "min"
	MinPercent     = "min_percent"
	MinTrend       = "min_trend"
	Max            = "max"
	MaxPercent     = "max_percent"
	MaxTrend       = "max_trend"
	Avg            = "avg"
	AvgPercent     = "avg_percent"
	AvgTrend       = "avg_trend"
	Current        = "current"
	CurrentPercent = "current_percent"
	CurrentTrend   = "current_trend"
)

func genRandNum() int {
	// calculate the max we will be using
	bg := big.NewInt(100)

	// get big.Int between 0 and bg
	// in this case 0 to 20
	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		panic(err)
	}

	// add n to min to support the passed in range
	return int(n.Int64())
}

func (c *Context) GetDataChannelProvider(channelType string) (func() ([]byte, error), error) {
	switch channelType {
	case userDataChannel:
		return c.provideUserStats, nil
	case itemDataChannel:
		return c.provideItemStats, nil
	case productDataChannel:
		return c.provideProductStats, nil
	case projectDataChannel:
		return c.provideProjectStats, nil
	default:
		return nil, fmt.Errorf("Unknown data channel %s", channelType)
	}
}

func (c *Context) provideItemStats() ([]byte, error) {
	data := make(map[string]interface{})
	data[itemsUsersProjectActivity] = make([][]int, 0)
	data[itemsProjectLength] = make([][]int, 0)
	data[itemsProjectCount] = make([][]int, 0)
	data[itemsProductCount] = make([][]int, 0)
	data[itemsCurrentProductCount] = make(map[string]interface{})
	data[itemsCurrentProjectCount] = make(map[string]interface{})
	data[itemsCurrentProductPerUser] = make(map[string]interface{})
	data[itemsCurrentProjectPerUser] = make(map[string]interface{})
	data[itemsProductPerUser] = make([][]int, 0)
	data[itemsProjectPerUser] = make([][]int, 0)
	timestamp := 1550197757000
	for i := 0; i < 300; i++ {
		dataPoint := []int{timestamp, genRandNum()}
		data[itemsUsersProjectActivity] = append(data[itemsUsersProjectActivity].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum()}
		data[itemsProjectLength] = append(data[itemsProjectLength].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum()}
		data[itemsProjectCount] = append(data[itemsProjectCount].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum()}
		data[itemsProductCount] = append(data[itemsProductCount].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum()}
		data[itemsProductPerUser] = append(data[itemsProductPerUser].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum()}
		data[itemsProjectPerUser] = append(data[itemsProjectPerUser].([][]int), dataPoint)
		timestamp += 5000000
	}
	total := 600
	data[itemsCurrentProductCount].(map[string]interface{})[Current] = 200
	data[itemsCurrentProductCount].(map[string]interface{})[CurrentTrend] = "up"
	data[itemsCurrentProductCount].(map[string]interface{})[CurrentPercent] = (100 * 200) / total
	data[itemsCurrentProjectCount].(map[string]interface{})[Current] = 450
	data[itemsCurrentProjectCount].(map[string]interface{})[CurrentTrend] = "down"
	data[itemsCurrentProjectCount].(map[string]interface{})[CurrentPercent] = (100 * 450) / total
	data[itemsCurrentProductPerUser].(map[string]interface{})[Current] = 320
	data[itemsCurrentProductPerUser].(map[string]interface{})[CurrentTrend] = "up"
	data[itemsCurrentProductPerUser].(map[string]interface{})[CurrentPercent] = (100 * 320) / total
	data[itemsCurrentProjectPerUser].(map[string]interface{})[Current] = 320
	data[itemsCurrentProjectPerUser].(map[string]interface{})[CurrentTrend] = "up"
	data[itemsCurrentProjectPerUser].(map[string]interface{})[CurrentPercent] = (100 * 100) / total
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (c *Context) provideUserStats() ([]byte, error) {
	data := make(map[string]interface{})
	data[usersTotal] = make([][]int, 0)
	data[usersDeleted] = make([][]int, 0)
	data[usersDevClient] = [4]int{10, 10, 20, 60}
	data[usersOnline] = make([][]int, 0)
	data[usersOnlinePeak] = make(map[string]interface{})
	data[usersOnlinePeriod] = make([][]int, 0)
	data[usersActivityHistory] = make([][]int, 0)
	data[usersOnlineTime] = make([][]int, 0)
	timestamp := 1550197757000
	for i := 0; i < 300; i++ {
		dataPoint := []int{timestamp, genRandNum()}
		data[usersTotal] = append(data[usersTotal].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 10}
		data[usersOnline] = append(data[usersOnline].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 60}
		data[usersDeleted] = append(data[usersDeleted].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 20}
		data[usersActivityHistory] = append(data[usersActivityHistory].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 30}
		data[usersOnlineTime] = append(data[usersOnlineTime].([][]int), dataPoint)
		timestamp += 5000000
	}
	offset := 0
	for month := 0; month < 12; month++ {
		dataPoints := make([]int, 0)
		for hour := 0; hour < 24; hour++ {
			dataPoints = append(dataPoints, genRandNum()+offset)
		}
		offset += 40
		data[usersOnlinePeriod] = append(data[usersOnlinePeriod].([][]int), dataPoints)
	}
	usersCount := 600
	data[usersOnlinePeak].(map[string]interface{})[Min] = 200
	data[usersOnlinePeak].(map[string]interface{})[MinTrend] = "up"
	data[usersOnlinePeak].(map[string]interface{})[MinPercent] = (100 * 200) / usersCount
	data[usersOnlinePeak].(map[string]interface{})[Max] = 450
	data[usersOnlinePeak].(map[string]interface{})[MaxTrend] = "down"
	data[usersOnlinePeak].(map[string]interface{})[MaxPercent] = (100 * 450) / usersCount
	data[usersOnlinePeak].(map[string]interface{})[Avg] = 320
	data[usersOnlinePeak].(map[string]interface{})[AvgTrend] = "up"
	data[usersOnlinePeak].(map[string]interface{})[AvgPercent] = (100 * 320) / usersCount
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (c *Context) provideProductStats() ([]byte, error) {
	data := make(map[string]interface{})
	data[productViewers] = make([][]int, 0)
	data[productPins] = make([][]int, 0)
	data[productPurchase] = make([][]int, 0)
	data[productPopularity] = make([][]int, 0)
	data[productSuccess] = make([][]int, 0)
	data[productProjectGeneration] = make([][]int, 0)
	data[productRating] = make([][]int, 0)
	timestamp := 1550197757000
	for i := 0; i < 300; i++ {
		dataPoint := []int{timestamp, genRandNum()}
		data[productViewers] = append(data[productViewers].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 10}
		data[productPurchase] = append(data[productPurchase].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 60}
		data[productPins] = append(data[productPins].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 80}
		data[productPopularity] = append(data[productPopularity].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 100}
		data[productSuccess] = append(data[productSuccess].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 120}
		data[productProjectGeneration] = append(data[productProjectGeneration].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 120}
		data[productRating] = append(data[productRating].([][]int), dataPoint)
		timestamp += 5000000
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (c *Context) provideProjectStats() ([]byte, error) {
	data := make(map[string]interface{})
	data[projectViewers] = make([][]int, 0)
	data[projectPins] = make([][]int, 0)
	data[projectConnections] = make([][]int, 0)
	data[projectPopularity] = make([][]int, 0)
	data[projectSuccess] = make([][]int, 0)
	data[projectUsageTime] = make([][]int, 0)
	data[projectRating] = make([][]int, 0)
	timestamp := 1550197757000
	for i := 0; i < 300; i++ {
		dataPoint := []int{timestamp, genRandNum()}
		data[projectViewers] = append(data[projectViewers].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 10}
		data[projectConnections] = append(data[projectConnections].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 60}
		data[projectPins] = append(data[projectPins].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 80}
		data[projectPopularity] = append(data[projectPopularity].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 100}
		data[projectSuccess] = append(data[projectSuccess].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 120}
		data[projectUsageTime] = append(data[projectUsageTime].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 120}
		data[projectRating] = append(data[projectRating].([][]int), dataPoint)
		timestamp += 5000000
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
