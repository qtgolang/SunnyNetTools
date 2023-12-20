package main

import (
	"encoding/json"
	"github.com/qtgolang/SunnyNet/src/GoWinHttp"
	"io"
	"math/rand"
	"net/url"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const LanZouDomain = "wwa.lanzoub.com\nwwa.lanzoue.com\nwwa.lanzouf.com\nwwa.lanzouh.com\nwwa.lanzouj.com\nwwa.lanzoul.com\nwwa.lanzoum.com\nwwa.lanzouo.com\nwwa.lanzoup.com\nwwa.lanzouq.com\nwwa.lanzout.com\nwwa.lanzouu.com\nwwa.lanzouv.com\nwwa.lanzouw.com\nwwa.lanzoux.com\nwwa.lanzouy.com\nwwb.lanzoub.com\nwwb.lanzoue.com\nwwb.lanzouf.com\nwwb.lanzouh.com\nwwb.lanzoui.com\nwwb.lanzouj.com\nwwb.lanzoul.com\nwwb.lanzoum.com\nwwb.lanzouo.com\nwwb.lanzoup.com\nwwb.lanzouq.com\nwwb.lanzout.com\nwwb.lanzouu.com\nwwb.lanzouv.com\nwwb.lanzouw.com\nwwc.lanzoub.com\nwwc.lanzoue.com\nwwc.lanzouf.com\nwwc.lanzouh.com\nwwc.lanzoui.com\nwwc.lanzouj.com\nwwc.lanzoul.com\nwwc.lanzoum.com\nwwc.lanzouo.com\nwwc.lanzoup.com\nwwc.lanzouq.com\nwwc.lanzout.com\nwwc.lanzouu.com\nwwc.lanzouv.com\nwwc.lanzouw.com\nwwc.lanzoux.com\nwwc.lanzouy.com\nwwd.lanzoub.com\nwwd.lanzoue.com\nwwd.lanzouh.com\nwwd.lanzoui.com\nwwd.lanzouj.com\nwwd.lanzoul.com\nwwd.lanzoum.com\nwwd.lanzouo.com\nwwd.lanzoup.com\nwwd.lanzouq.com\nwwd.lanzout.com\nwwd.lanzouu.com\nwwd.lanzouv.com\nwwd.lanzouw.com\nwwd.lanzoux.com\nwwd.lanzouy.com\nwwe.lanzoub.com\nwwe.lanzoue.com\nwwe.lanzouf.com\nwwe.lanzouh.com\nwwe.lanzoui.com\nwwe.lanzouj.com\nwwe.lanzoul.com\nwwe.lanzoum.com\nwwe.lanzouo.com\nwwe.lanzoup.com\nwwe.lanzouq.com\nwwe.lanzout.com\nwwe.lanzouu.com\nwwe.lanzouv.com\nwwe.lanzouw.com\nwwe.lanzoux.com\nwwe.lanzouy.com\nwwf.lanzoub.com\nwwf.lanzoue.com\nwwf.lanzouf.com\nwwf.lanzouh.com\nwwf.lanzoui.com\nwwf.lanzouj.com\nwwf.lanzoul.com\nwwf.lanzoum.com\nwwf.lanzouo.com\nwwf.lanzoup.com\nwwf.lanzouq.com\nwwf.lanzout.com\nwwf.lanzouu.com\nwwf.lanzouv.com\nwwf.lanzouw.com\nwwf.lanzoux.com\nwwf.lanzouy.com\nwwg.lanzoub.com\nwwg.lanzoue.com\nwwg.lanzouf.com\nwwg.lanzouh.com\nwwg.lanzoui.com\nwwg.lanzouj.com\nwwg.lanzoul.com\nwwg.lanzoum.com\nwwg.lanzouo.com\nwwg.lanzoup.com\nwwg.lanzouq.com\nwwg.lanzout.com\nwwg.lanzouu.com\nwwg.lanzouv.com\nwwg.lanzouw.com\nwwg.lanzoux.com\nwwg.lanzouy.com\nwwh.lanzoub.com\nwwh.lanzoue.com\nwwh.lanzouf.com\nwwh.lanzouh.com\nwwh.lanzoui.com\nwwh.lanzouj.com\nwwh.lanzoul.com\nwwh.lanzoum.com\nwwh.lanzouo.com\nwwh.lanzoup.com\nwwh.lanzouq.com\nwwh.lanzout.com\nwwh.lanzouu.com\nwwh.lanzouv.com\nwwh.lanzouw.com\nwwh.lanzoux.com\nwwh.lanzouy.com\nwwi.lanzoub.com\nwwi.lanzoue.com\nwwi.lanzouf.com\nwwi.lanzouh.com\nwwi.lanzoui.com\nwwi.lanzouj.com\nwwi.lanzoul.com\nwwi.lanzoum.com\nwwi.lanzouo.com\nwwi.lanzoup.com\nwwi.lanzouq.com\nwwi.lanzout.com\nwwi.lanzouu.com\nwwi.lanzouv.com\nwwi.lanzouw.com\nwwi.lanzoux.com\nwwi.lanzouy.com\nwwj.lanzoub.com\nwwj.lanzoue.com\nwwj.lanzouf.com\nwwj.lanzouh.com\nwwj.lanzoui.com\nwwj.lanzouj.com\nwwj.lanzoul.com\nwwj.lanzoum.com\nwwj.lanzouo.com\nwwj.lanzoup.com\nwwj.lanzouq.com\nwwj.lanzout.com\nwwj.lanzouu.com\nwwj.lanzouv.com\nwwj.lanzouw.com\nwwj.lanzoux.com\nwwj.lanzouy.com\nwwk.lanzoub.com\nwwk.lanzoue.com\nwwk.lanzouf.com\nwwk.lanzouh.com\nwwk.lanzoui.com\nwwk.lanzouj.com\nwwk.lanzoul.com\nwwk.lanzoum.com\nwwk.lanzouo.com\nwwk.lanzoup.com\nwwk.lanzouq.com\nwwk.lanzout.com\nwwk.lanzouu.com\nwwk.lanzouv.com\nwwk.lanzouw.com\nwwk.lanzoux.com\nwwk.lanzouy.com\nwwl.lanzoub.com\nwwl.lanzoue.com\nwwl.lanzouf.com\nwwl.lanzouh.com\nwwl.lanzoui.com\nwwl.lanzouj.com\nwwl.lanzoul.com\nwwl.lanzoum.com\nwwl.lanzouo.com\nwwl.lanzoup.com\nwwl.lanzouq.com\nwwl.lanzout.com\nwwl.lanzouu.com\nwwl.lanzouv.com\nwwl.lanzouw.com\nwwl.lanzoux.com\nwwl.lanzouy.com\nwwm.lanzoub.com\nwwm.lanzoue.com\nwwm.lanzouf.com\nwwm.lanzouh.com\nwwm.lanzoui.com\nwwm.lanzouj.com\nwwm.lanzoul.com\nwwm.lanzoum.com\nwwm.lanzouo.com\nwwm.lanzoup.com\nwwm.lanzouq.com\nwwm.lanzout.com\nwwm.lanzouu.com\nwwm.lanzouv.com\nwwm.lanzouw.com\nwwm.lanzoux.com\nwwm.lanzouy.com\nwwn.lanzoub.com\nwwn.lanzoue.com\nwwn.lanzouf.com\nwwn.lanzouh.com\nwwn.lanzoui.com\nwwn.lanzouj.com\nwwn.lanzoul.com\nwwn.lanzoum.com\nwwn.lanzouo.com\nwwn.lanzoup.com\nwwn.lanzouq.com\nwwn.lanzout.com\nwwn.lanzouu.com\nwwn.lanzouv.com\nwwn.lanzouw.com\nwwn.lanzoux.com\nwwn.lanzouy.com\nwwo.lanzoub.com\nwwo.lanzoue.com\nwwo.lanzouf.com\nwwo.lanzouh.com\nwwo.lanzoui.com\nwwo.lanzouj.com\nwwo.lanzoul.com\nwwo.lanzoum.com\nwwo.lanzouo.com\nwwo.lanzoup.com\nwwo.lanzouq.com\nwwo.lanzout.com\nwwo.lanzouu.com\nwwo.lanzouv.com\nwwo.lanzouw.com\nwwo.lanzoux.com\nwwo.lanzouy.com\nwwp.lanzoub.com\nwwp.lanzoue.com\nwwp.lanzouf.com\nwwp.lanzouh.com\nwwp.lanzoui.com\nwwp.lanzoul.com\nwwp.lanzoum.com\nwwp.lanzouo.com\nwwp.lanzoup.com\nwwp.lanzouq.com\nwwp.lanzout.com\nwwp.lanzouu.com\nwwp.lanzouw.com\nwwp.lanzoux.com\nwwp.lanzouy.com\nwwq.lanzoub.com\nwwq.lanzoue.com\nwwq.lanzouf.com\nwwq.lanzouh.com\nwwq.lanzoui.com\nwwq.lanzouj.com\nwwq.lanzoul.com\nwwq.lanzoum.com\nwwq.lanzouo.com\nwwq.lanzoup.com\nwwq.lanzouq.com\nwwq.lanzout.com\nwwq.lanzouu.com\nwwq.lanzouv.com\nwwq.lanzouw.com\nwwq.lanzoux.com\nwwq.lanzouy.com\nwwr.lanzoub.com\nwwr.lanzoue.com\nwwr.lanzouf.com\nwwr.lanzouh.com\nwwr.lanzoui.com\nwwr.lanzouj.com\nwwr.lanzoul.com\nwwr.lanzoum.com\nwwr.lanzouo.com\nwwr.lanzoup.com\nwwr.lanzouq.com\nwwr.lanzout.com\nwwr.lanzouu.com\nwwr.lanzouv.com\nwwr.lanzouw.com\nwwr.lanzoux.com\nwwr.lanzouy.com\nwws.lanzoub.com\nwws.lanzoue.com\nwws.lanzouf.com\nwws.lanzouh.com\nwws.lanzoui.com\nwws.lanzouj.com\nwws.lanzoul.com\nwws.lanzoum.com\nwws.lanzouo.com\nwws.lanzoup.com\nwws.lanzouq.com\nwws.lanzout.com\nwws.lanzouu.com\nwws.lanzouv.com\nwws.lanzouw.com\nwws.lanzoux.com\nwws.lanzouy.com\nwwt.lanzoub.com\nwwt.lanzoue.com\nwwt.lanzouf.com\nwwt.lanzouh.com\nwwt.lanzoui.com\nwwt.lanzouj.com\nwwt.lanzoul.com\nwwt.lanzoum.com\nwwt.lanzouo.com\nwwt.lanzoup.com\nwwt.lanzouq.com\nwwt.lanzout.com\nwwt.lanzouu.com\nwwt.lanzouv.com\nwwt.lanzouw.com\nwwt.lanzoux.com\nwwt.lanzouy.com\nwwu.lanzoub.com\nwwu.lanzoue.com\nwwu.lanzouf.com\nwwu.lanzouh.com\nwwu.lanzoui.com\nwwu.lanzouj.com\nwwu.lanzoul.com\nwwu.lanzoum.com\nwwu.lanzouo.com\nwwu.lanzoup.com\nwwu.lanzouq.com\nwwu.lanzout.com\nwwu.lanzouu.com\nwwu.lanzouv.com\nwwu.lanzouw.com\nwwu.lanzoux.com\nwwu.lanzouy.com\nwwv.lanzoub.com\nwwv.lanzoue.com\nwwv.lanzouf.com\nwwv.lanzouh.com\nwwv.lanzoui.com\nwwv.lanzouj.com\nwwv.lanzoul.com\nwwv.lanzoum.com\nwwv.lanzouo.com\nwwv.lanzoup.com\nwwv.lanzouq.com\nwwv.lanzout.com\nwwv.lanzouu.com\nwwv.lanzouv.com\nwwv.lanzouw.com\nwwv.lanzoux.com\nwwv.lanzouy.com\nwww.lanzoub.com\nwww.lanzoue.com\nwww.lanzouf.com\nwww.lanzouh.com\nwww.lanzoui.com\nwww.lanzoul.com\nwww.lanzoum.com\nwww.lanzoup.com\nwww.lanzouq.com\nwww.lanzout.com\nwww.lanzouu.com\nwww.lanzouv.com\nwww.lanzouw.com\nwww.lanzoux.com\nwww.lanzouy.com\nwwx.lanzoub.com\nwwx.lanzoue.com\nwwx.lanzouf.com\nwwx.lanzouh.com\nwwx.lanzoui.com\nwwx.lanzouj.com\nwwx.lanzoul.com\nwwx.lanzoum.com\nwwx.lanzouo.com\nwwx.lanzoup.com\nwwx.lanzouq.com\nwwx.lanzout.com\nwwx.lanzouu.com\nwwx.lanzouv.com\nwwx.lanzouw.com\nwwx.lanzoux.com\nwwx.lanzouy.com\nwwy.lanzoub.com\nwwy.lanzoue.com\nwwy.lanzouf.com\nwwy.lanzouh.com\nwwy.lanzoui.com\nwwy.lanzouj.com\nwwy.lanzoul.com\nwwy.lanzoum.com\nwwy.lanzouo.com\nwwy.lanzoup.com\nwwy.lanzouq.com\nwwy.lanzout.com\nwwy.lanzouu.com\nwwy.lanzouw.com\nwwy.lanzoux.com\nwwy.lanzouy.com\nwwz.lanzoub.com\nwwz.lanzoue.com\nwwz.lanzouf.com\nwwz.lanzouh.com\nwwz.lanzoui.com\nwwz.lanzouj.com\nwwz.lanzoul.com\nwwz.lanzoum.com\nwwz.lanzouo.com\nwwz.lanzoup.com\nwwz.lanzouq.com\nwwz.lanzout.com\nwwz.lanzouu.com\nwwz.lanzouw.com\nwwz.lanzoux.com\nwwz.lanzouy.com"

const Version = 1006

func lanZouUpdate() {
	if runtime.GOOS == "windows" {
		_lanZouUpdate("https://wwxa.lanzouj.com/b0cior9kb", "2brf")
	} else {
		_lanZouUpdate("https://wwxa.lanzouj.com/b0ciopv1c", "2oxf")
	}
}
func _lanZouUpdate(u, pass string) {
	rand.Seed(time.Now().UnixNano()) // 设置随机数种子
	DoMainList := strings.Split(LanZouDomain, "\n")
	_Url, _ := url.Parse(u)
	for _, _ = range DoMainList {
		l := rand.Intn(len(DoMainList))
		if l >= len(DoMainList) {
			continue
		}
		_Url.Host = DoMainList[l]
		data := _lanZouParse(_Url)
		if data == "" {
			continue
		}
		ok, name, id := _lanZouParseId(data, pass, _Url)
		if !ok {
			continue
		}
		if name == "无数据" {
			return
		}
		NewVersion := _lanZouParseConnect(id, _Url)
		if NewVersion > Version {
			s := `当前有新版本：`
			s += `<span class="el-text el-text--success mx-1 el-tooltip__trigger el-tooltip__trigger" style="cursor: pointer; text-decoration: underline;" onclick="WindowOpenURL('https://` + _Url.Host + `/` + id + `')">点击去下载</span>`
			CallJs("有新版本", s)
		}
		if NewVersion != -1 {
			return
		}
	}
}
func _lanZouParseConnect(id string, u *url.URL) int {
	a := GoWinHttp.WinHttp{}
	a.Open("GET", "https://"+u.Host+"/"+id)
	a.SetOutTime(1000, 1000, 1000)
	a.SetHeader("Host", u.Host)
	a.SetHeader("Sec-Fetch-Dest", "document")
	a.SetHeader("Cache-Control", "max-age=0")
	a.SetHeader("Upgrade-Insecure-Requests", "1")
	a.SetHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	a.SetHeader("Sec-Fetch-Mode", "navigate")
	a.SetHeader("Connection", "keep-alive")
	a.SetHeader("sec-ch-ua-mobile", "?0")
	a.SetHeader("Sec-Fetch-Site", "none")
	a.SetHeader("Sec-Fetch-User", "?1")
	a.SetHeader("Referer", u.String())
	a.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	a.SetHeader("Accept-Language", "zh-CN,zh;q=0.9")
	Response, _ := a.Send("")
	if Response == nil {
		return -1
	}
	if Response.Body == nil {
		return -1
	}
	body, _ := io.ReadAll(Response.Body)
	_ = Response.Body.Close()
	res := string(body)
	ms := SubString(res, "文件描述：</span><br>", "<")
	ms = strings.ReplaceAll(ms, "【", "[")
	ms = strings.ReplaceAll(ms, "】", "]")
	ms = SubString(ms, "[", "]")
	iii, e := strconv.Atoi(ms)
	if e != nil {
		return -1
	}
	return iii
}
func _lanZouParseId(data string, pass string, u *url.URL) (bool, string, string) {
	a := GoWinHttp.WinHttp{}
	a.Open("POST", "https://"+u.Host+"/filemoreajax.php")
	a.SetOutTime(1000, 1000, 1000)
	a.SetHeader("Host", u.Host)
	a.SetHeader("Sec-Fetch-Dest", "document")
	a.SetHeader("Cache-Control", "max-age=0")
	a.SetHeader("Upgrade-Insecure-Requests", "1")
	a.SetHeader("Accept", "application/json, text/javascript, */*")
	a.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	a.SetHeader("Sec-Fetch-Mode", "navigate")
	a.SetHeader("Connection", "keep-alive")
	a.SetHeader("sec-ch-ua-mobile", "?0")
	a.SetHeader("Sec-Fetch-Site", "none")
	a.SetHeader("Sec-Fetch-User", "?1")
	a.SetHeader("Referer", u.String())
	a.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	a.SetHeader("Accept-Language", "zh-CN,zh;q=0.9")
	Response, _ := a.Send(data + "&pwd=" + pass)
	if Response == nil {
		return false, "", ""
	}
	if Response.Body == nil {
		return false, "", ""
	}
	body, _ := io.ReadAll(Response.Body)
	_ = Response.Body.Close()
	JsonBody := &JSONData{}
	json.Unmarshal(body, JsonBody)
	if JsonBody == nil {
		return false, "", ""
	}
	if len(JsonBody.Text) < 1 {
		return true, "无数据", ""
	}
	name := JsonBody.Text[0].NameAll
	id := JsonBody.Text[0].ID
	if name != "" && id != "" {
		return true, name, id
	}
	return false, "", ""
}

type JSONData struct {
	Zt   int    `json:"zt"`
	Info string `json:"info"`
	Text []Text `json:"text"`
}
type Text struct {
	Icon    string `json:"icon"`
	T       int    `json:"t"`
	ID      string `json:"id"`
	NameAll string `json:"name_all"`
	Size    string `json:"size"`
	Time    string `json:"time"`
	Duan    string `json:"duan"`
	PIco    int    `json:"p_ico"`
}

func _lanZouParse(u *url.URL) string {
	a := GoWinHttp.WinHttp{}
	a.Open("GET", u.String())
	a.SetOutTime(1000, 1000, 1000)
	a.SetHeader("Host", u.Host)
	a.SetHeader("Sec-Fetch-Dest", "document")
	a.SetHeader("Cache-Control", "max-age=0")
	a.SetHeader("Upgrade-Insecure-Requests", "1")
	a.SetHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	a.SetHeader("Sec-Fetch-Mode", "navigate")
	a.SetHeader("Connection", "keep-alive")
	a.SetHeader("sec-ch-ua-mobile", "?0")
	a.SetHeader("Sec-Fetch-Site", "none")
	a.SetHeader("Sec-Fetch-User", "?1")
	a.SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")
	a.SetHeader("Accept-Language", "zh-CN,zh;q=0.9")
	Response, _ := a.Send("")
	if Response != nil {
		if Response.Body != nil {
			body, _ := io.ReadAll(Response.Body)
			_ = Response.Body.Close()
			_res := string(body)
			res := SubString(_res, "document.getElementById", "document.getElementById")
			arr := strings.Split(res, ";")
			t := ""
			sign := ""
			for _, v := range arr {
				str1 := SubString(v, "'", "'")
				if len(str1) == 10 {
					t = str1
				}
				if len(str1) == 32 {
					sign = str1
				}
			}
			lx := SubString(_res, "'lx':", ",")
			fid := SubString(_res, "'fid':", ",")
			uid := SubString(_res, "'uid':", ",")
			uid = strings.ReplaceAll(uid, "'", "")
			if uid == "" || t == "" || sign == "" {
				return ""
			}
			up := SubString(_res, "'up':", ",")
			ls := SubString(_res, "'ls':", ",")
			pg := "1"

			rep := SubString(_res, "'rep':", ",")
			rep = strings.ReplaceAll(rep, "'", "")
			str := "lx=" + lx + "&fid=" + fid + "&uid=" + uid + "&pg=" + pg + "&rep=" + rep + "&t=" + t + "&k=" + sign + "&up=" + up + "&ls=" + ls
			return str
		}

	}

	return ""
}
func SubString(str, left, Right string) string {
	s := strings.Index(str, left)
	if s < 0 {
		return ""
	}
	s += len(left)
	e := strings.Index(str[s:], Right)
	if e+s <= s {
		return ""
	}
	return str[s : s+e]

}
