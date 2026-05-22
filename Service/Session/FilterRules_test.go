package Session

import (
	"fmt"
	"testing"
)

func Test_FilterRules(t *testing.T) {
	jsonStr := `
{"filterType":"text","colId":"数据","type":"contains","filter":"11"}
`
	filter := parseAgGridFilterRules(jsonStr)
	testData := map[string]interface{}{
		"数据": "11x23ABC",
		"长度": 120.0,
		"时间": "123121",
	}
	//			([数据] 包含 "12" 且 [数据] 不包含 "322"  ) 或 ([长度] > 100  且 [时间] 包含 "123")
	result := matchFilters(testData, filter)
	fmt.Println("匹配结果:", result) // 预期输出: true
}
