package businesslogic

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
)

const (
	UsersTotal           = "users_total"
	UsersOnline          = "users_online"
	UsersDevClient       = "users_dev_client"
	UsersActivityHistory = "users_activity_history"
	UsersOnlinePeak      = "users_online_peaks"
	UsersOnlinePeriod    = "users_online_period"
	UsersOnlineTime      = "user_online_time"
	UsersOnlineTimeMin   = "user_online_time_min"
	UsersOnlineTimeMax   = "user_online_time_max"

	ProjectUsers           = "project_user"
	ProjectViewers         = "project_viewers"
	ProjectWatchlisters    = "project_watchlisters"
	ProjectObservers       = "project_observers"
	ProjectActiveUsers     = "project_active_users"
	ProjectActivityHistory = "project_activity_history"
	ProjectRatingHistory   = "project_rating_history"

	ProductViewers           = "product_viewers"
	ProductUsers             = "product_users"
	ProductWatchlisters      = "product_watchlister"
	ProductViewAndPurchase   = "product_success" // Returns the stats showing how many products have actually been purchased after viewing or watching
	ProductProjectGeneration = "product_project_generation"

	Timestamp = "timestamp"
	Value     = "value"
)

const (
	Min        = "min"
	MinPercent = "min_percent"
	MinTrend   = "min_trend"
	Max        = "max"
	MaxPercent = "max_percent"
	MaxTrend   = "max_trend"
	Avg        = "avg"
	AvgPercent = "avg_percent"
	AvgTrend   = "avg_trend"
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

func (c *Context) ProvideUserStats() ([]byte, error) {
	data := make(map[string]interface{})
	data[UsersTotal] = make([][]int, 0)
	data[UsersDevClient] = [4]int{10, 10, 20, 60}
	data[UsersOnline] = make([][]int, 0)
	data[UsersOnlinePeak] = make(map[string]interface{})
	data[UsersOnlinePeriod] = make([][]int, 0)
	data[UsersActivityHistory] = make([][]int, 0)
	data[UsersOnlineTime] = make([][]int, 0)
	data[UsersOnlineTimeMin] = make([][]int, 0)
	data[UsersOnlineTimeMax] = make([][]int, 0)
	timestamp := 1550197757000
	for i := 0; i < 300; i++ {
		dataPoint := []int{timestamp, genRandNum()}
		data[UsersTotal] = append(data[UsersTotal].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 10}
		data[UsersOnline] = append(data[UsersOnline].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 20}
		data[UsersActivityHistory] = append(data[UsersActivityHistory].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 30}
		data[UsersOnlineTime] = append(data[UsersOnlineTime].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 40}
		data[UsersOnlineTimeMin] = append(data[UsersOnlineTimeMin].([][]int), dataPoint)
		dataPoint = []int{timestamp, genRandNum() + 50}
		data[UsersOnlineTimeMax] = append(data[UsersOnlineTimeMax].([][]int), dataPoint)
		timestamp += 5000000
	}
	offset := 0
	for month := 0; month < 12; month++ {
		dataPoints := make([]int, 0)
		for hour := 0; hour < 24; hour++ {
			dataPoints = append(dataPoints, genRandNum()+offset)
		}
		offset += 40
		data[UsersOnlinePeriod] = append(data[UsersOnlinePeriod].([][]int), dataPoints)
	}
	usersCount := 600
	data[UsersOnlinePeak].(map[string]interface{})[Min] = 200
	data[UsersOnlinePeak].(map[string]interface{})[MinTrend] = "up"
	data[UsersOnlinePeak].(map[string]interface{})[MinPercent] = (100 * 200) / usersCount
	data[UsersOnlinePeak].(map[string]interface{})[Max] = 450
	data[UsersOnlinePeak].(map[string]interface{})[MaxTrend] = "down"
	data[UsersOnlinePeak].(map[string]interface{})[MaxPercent] = (100 * 450) / usersCount
	data[UsersOnlinePeak].(map[string]interface{})[Avg] = 320
	data[UsersOnlinePeak].(map[string]interface{})[AvgTrend] = "up"
	data[UsersOnlinePeak].(map[string]interface{})[AvgPercent] = (100 * 320) / usersCount
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (c *Context) ProvideProjectStats() ([]byte, error) {
	data := make(map[string]interface{})
	data[UsersTotal] = make(map[int64]int)
	data[UsersOnline] = make(map[int64]int)
	data[UsersActivityHistory] = make(map[int64]int)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
