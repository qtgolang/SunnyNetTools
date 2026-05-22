package Service

import (
	"changeme/Service/Config"
	"changeme/Service/Session"
	"fmt"
	"github.com/qtgolang/SunnyNet/SunnyNet"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// 操作锁
var lock sync.Mutex

// 插入列表_http/udp/tcp
var InsertList = make([]Session.Insert, 0)

// 要更新的列表
var updateDoneList = make([]updateHTTPDone, 0)
var updateSendList = make([]httpUpdateSend, 0)
var updateErrorList = make([]updateHTTPError, 0)
var updateSocket_List = make([]updateSocket, 0)
var updateSocketStreamList = make([]Session.UpdateSocketStream, 0)

type httpUpdateSend struct {
	Theology  int
	Ico       string
	Filter    bool
	URL       string
	Method    string
	Note      string
	BreakMode uint32 //断点模式
}
type insertHTTP struct {
	Method           string
	URL              string
	ProcessName      string
	ClientIP         string
	Theology         int
	Time             string
	Ico              string
	Filter           bool
	UserName         string //身份验证的账号
	BreakMode        uint32 //断点模式
	GuaranteeDisplay bool   //保证显示此消息
}
type updateHTTPDone struct {
	Theology  int
	Code      string
	Length    int
	Type      string
	Time      string
	Ico       string
	BreakMode uint32
	Filter    bool
	IP        string
	Note      string
}
type updateHTTPError struct {
	Theology int
	Code     string
	Length   int
	Time     string
	Ico      string
	Filter   bool
	Note     string
}
type updateSocket struct {
	Method    string
	Theology  int
	Code      string //状态
	RecLength int    //接收数据长度
	SenLength int    //发送数据长度
	Ico       string
	Filter    bool
	Note      string
	Stream    *Session.Stream `json:"-"`
}

func getTime() string {
	now := time.Now()
	// 获取当前时间的毫秒部分
	milliseconds := now.UnixNano() / int64(time.Millisecond)
	// 格式化为时-分-秒:毫秒
	formattedTime := now.Format("15-04-05") + ":" + fmt.Sprintf("%03d", milliseconds%1000)
	return formattedTime
}
func UpdateIco(Conn SunnyNet.ConnHTTP, _ContentType string) string {
	ContentType := strings.ToLower(_ContentType)
	Method := strings.ToUpper(Conn.Method())
	if strings.Contains(ContentType, "image/") {
		return "img"
	}
	if strings.Contains(ContentType, "/javascript") {
		return "js"
	}
	if strings.Contains(ContentType, "/x-javascript") {
		return "js"
	}
	if strings.Contains(ContentType, "/css") {
		return "css"
	}
	if strings.Contains(ContentType, "/xml") {
		return "XML"
	}
	if strings.Contains(ContentType, "/json") {
		return "JSON"
	}
	if strings.Contains(ContentType, "/html") {
		return "HTML"
	}
	if strings.Contains(ContentType, "audio/") {
		return "audio"
	}
	if strings.Contains(ContentType, "video/") {
		return "video"
	}
	{
		URL := strings.ToLower(Conn.URL())
		array := strings.Split(URL, "?")
		if len(array) > 1 {
			URL = array[0]
		}
		{
			if strings.HasSuffix(URL, ".woff2") ||
				strings.HasSuffix(URL, ".woff") ||
				strings.HasSuffix(URL, ".eot") ||
				strings.HasSuffix(URL, ".otf") ||
				strings.HasSuffix(URL, ".fon") ||
				strings.HasSuffix(URL, ".font") ||
				strings.HasSuffix(URL, ".ttc") ||
				strings.HasSuffix(URL, ".eotz") ||
				strings.HasSuffix(URL, ".dfont") ||
				strings.HasSuffix(URL, ".suit") ||
				strings.HasSuffix(URL, ".pfb") ||
				strings.HasSuffix(URL, ".ttf") {
				return "font"
			}
		}
		{
			if strings.HasSuffix(URL, ".swf") || strings.Contains(ContentType, "flash") {
				return "Flash"
			}
		}
	}
	if Method == "POST" || Method == "PUT" {
		return "POST"
	}
	StateCode := Conn.GetResponseCode()
	if StateCode == 302 {
		return "302"
	}
	if StateCode == 401 {
		return "401"
	}
	if StateCode == 403 || StateCode == 404 || StateCode == 405 {
		return "stop"
	}
	return "generic"
}

var insertLock int32

func (g *AppMain) AppInsertDone() {
	atomic.StoreInt32(&insertLock, 0)
}
func (g *AppMain) AppStartInsert() {
	for !atomic.CompareAndSwapInt32(&insertLock, 0, 1) {
		runtime.Gosched()
	}
}
func parseTime(t string) (time.Time, error) {
	// 将 "05-45-18:746" 转为 "05:45:18.746"
	t = strings.Replace(t, ":", ".", 1)
	t = strings.Replace(t, "-", ":", 2)
	return time.Parse("15:04:05.000", t)
}

func init() {
	go func() {
		for {
			if len(InsertList) == 0 && len(updateSocketStreamList) == 0 {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			for !atomic.CompareAndSwapInt32(&insertLock, 0, 1) {
				runtime.Gosched()
			}
			lock.Lock()
			for len(InsertList) > 0 {
				sort.Slice(InsertList, func(i, j int) bool {
					ti, _ := parseTime(InsertList[i].Time)
					tj, _ := parseTime(InsertList[j].Time)
					return ti.Before(tj)
				})
				batch := InsertList
				if len(InsertList) > batchSize {
					batch = InsertList[:batchSize]
					InsertList = InsertList[batchSize:]
					Config.AppList["Main"].EmitEvent("insert", batch, false)
				} else {
					InsertList = make([]Session.Insert, 0)
					Config.AppList["Main"].EmitEvent("insert", batch, true)
				}
			}
			for len(updateSocketStreamList) > 0 {
				batch := updateSocketStreamList
				if len(updateSocketStreamList) > batchSize {
					batch = updateSocketStreamList[:batchSize]
					updateSocketStreamList = updateSocketStreamList[batchSize:]
					Config.AppList["Main"].EmitEvent("updateSocketStreamList", batch, false)
				} else {
					updateSocketStreamList = make([]Session.UpdateSocketStream, 0)
					Config.AppList["Main"].EmitEvent("updateSocketStreamList", batch, true)
				}
			}
			lock.Unlock()
		}
	}()
	go func() {
		for {
			if len(updateDoneList) == 0 && len(updateSendList) == 0 && len(updateErrorList) == 0 && len(updateSocket_List) == 0 {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			lock.Lock()
			for len(updateSocket_List) > 0 && len(updateDoneList) == 0 {
				batch := updateSocket_List
				if len(updateSocket_List) > batchSize {
					batch = updateSocket_List[:batchSize]
					updateSocket_List = updateSocket_List[batchSize:]
				} else {
					updateSocket_List = make([]updateSocket, 0)
				}
				Config.AppList["Main"].EmitEvent("updateWebsocket_tcp_udp_List", batch)
			}
			for len(updateDoneList) > 0 {
				batch := updateDoneList
				if len(updateDoneList) > batchSize {
					batch = updateDoneList[:batchSize]
					updateDoneList = updateDoneList[batchSize:]
				} else {
					updateDoneList = make([]updateHTTPDone, 0)
				}
				Config.AppList["Main"].EmitEvent("updateDoneHTTP", batch)
			}

			for len(updateSendList) > 0 {
				batch := updateSendList
				if len(updateSendList) > batchSize {
					batch = updateSendList[:batchSize]
					updateSendList = updateSendList[batchSize:]
				} else {
					updateSendList = make([]httpUpdateSend, 0)
				}
				Config.AppList["Main"].EmitEvent("updateSendHTTP", batch)
			}

			for len(updateErrorList) > 0 {
				batch := updateErrorList
				if len(updateErrorList) > batchSize {
					batch = updateErrorList[:batchSize]
					updateErrorList = updateErrorList[batchSize:]
				} else {
					updateErrorList = make([]updateHTTPError, 0)
				}
				Config.AppList["Main"].EmitEvent("updateErrorHTTP", batch)
			}
			lock.Unlock()

		}
	}()
}

const batchSize = 100
