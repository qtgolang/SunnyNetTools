package main

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type CreateRequest struct {
	URL      *url.URL
	Method   string
	Header   http.Header
	Body     []byte
	FuncName string
	Cookie   string
}

func CreateRequestCode(TheologyArray []int, Lang string, module string) string {
	index := 0
	code := ""
	for _, k := range TheologyArray {
		h := HashMap.GetRequest(k)
		if h != nil {
			u, e := url.Parse(h.URL)
			if e != nil {
				continue
			}
			mm := &CreateRequest{URL: u, Header: h.Header, Body: h.Body, Method: h.Method}
			index++
			mm.FuncName = "SunnyNetCreateRequest" + strconv.Itoa(index)
			mm.Cookie = mm.Header.Get("Cookie")
			mm.Header.Del("Cookie")
			if Lang == "E" {
				code += mm.ELang(module)
			} else if Lang == "C#" {
				code += mm.CSharp(module)
			} else if Lang == "Go" {
				code += mm.Go(module)
			} else if Lang == "Python" {
				code += mm.Python(module)
			}
		}
	}
	if Lang == "E" {
		code = ".版本 2\n.支持库 spec\n\n" + code
	}
	code = strings.ReplaceAll(code, "\r", "")
	code = strings.ReplaceAll(code, "\n", "\r\n")
	return code
}
func (e *CreateRequest) ELang(module string) string {
	code := ""
	if module == "WinInet" || module == "WinHttpW" || module == "WinHttpR" {
		code += ".子程序 " + e.FuncName + ", , 公开, 本子程序由Sunny中间件生成,请配合 [WinHttp模块] 使用。 " + e.URL.Path + "\n.局部变量 局_HTTP, " + module + "\n"
		BytesType := e.IsBytesType()
		if BytesType {
			code += ".局部变量 局_提交字节集, 字节集, , , \n"
		} else {
			code += ".局部变量 局_提交数据, 文本型, , , \n"
		}
		if e.Cookie != "" {
			code += ".局部变量 局_提交Cookie, 文本型, , , \n"
		}
		code += ".局部变量 局_响应字节集, 字节集, , , \n.局部变量 局_响应文本, 文本型, , , \n\n"

		if BytesType {
			code += "局_提交字节集 ＝ 编码_BASE64解码 (“" + base64.StdEncoding.EncodeToString(e.Body) + "”)\n\n"
		} else {
			ok, d := e.IsFormData()
			if ok {
				code += d + "\n"
			} else {
				code += "局_提交数据 ＝ " + convertELangFormat(string(e.Body)) + "\n"
			}
		}
		if e.Cookie != "" {
			code += "局_提交Cookie ＝ 子文本替换 (" + convertELangFormat(e.Cookie) + ", “'”, #引号, , , 真)\n\n"
		}
		code += "局_HTTP.Open (“" + e.Method + "”, “" + e.URL.String() + "”)\n"
		for k, v := range e.Header {
			if strings.ToUpper(k) == "CONTENT-LENGTH" {
				continue
			}
			if k == "Accept-Encoding" {
				if len(v) < 1 {
					code += "' 局_HTTP.SetHeader (“" + k + "”, “”)\n"
				} else {
					code += "' 局_HTTP.SetHeader (“" + k + "”, " + convertELangFormat(v[0]) + ")\n"
				}
				continue
			}
			if len(v) < 1 {
				code += "局_HTTP.SetHeader (“" + k + "”, “”)\n"
			} else {
				code += "局_HTTP.SetHeader (“" + k + "”, " + convertELangFormat(v[0]) + ")\n"
			}
		}
		if e.Cookie != "" {
			code += "局_HTTP.SetHeader (“Cookie”, 局_提交Cookie)\n\n"
		}
		if BytesType {
			code += "局_HTTP.SendBin (局_提交字节集)\n"
		} else {
			code += "局_HTTP.Send (局_提交数据)\n"
		}
		code += "局_响应字节集 ＝ 局_HTTP.GetBody ()\n"
		code += "局_响应文本 ＝ 编码_Utf8到Ansi (局_响应字节集)\n"
		code += "调试输出 (局_响应文本)\n\n"
	}
	if module == "e2ee" {
		code += ".子程序 " + e.FuncName + ", , 公开, 本子程序由Sunny中间件生成,请配合 [E2EE] 使用。 " + e.URL.Path + "\n.局部变量 局_HTTP, 网站客户端\n"
		BytesType := e.IsBytesType()
		code += ".局部变量 局_网址, 文本型, , , \n"
		if BytesType {
			code += ".局部变量 局_提交数据, 字节集, , , \n"
		} else {
			code += ".局部变量 局_提交数据, 文本型, , , \n"
		}
		if e.Cookie != "" {
			code += ".局部变量 局_提交Cookie, 文本型, , , \n"
		}
		code += ".局部变量 局_响应字节集, 字节集, , , \n.局部变量 局_响应文本, 文本型, , , \n\n"

		code += "局_网址 ＝ " + convertELangFormat(e.URL.String()) + "\n\n"
		if BytesType {
			code += "局_提交数据 ＝ 编码_BASE64解码 (“" + base64.StdEncoding.EncodeToString(e.Body) + "”)\n\n"
		} else {
			ok, d := e.IsFormData()
			if ok {
				code += d + "\n"
			} else {
				code += "局_提交数据 ＝ " + convertELangFormat(string(e.Body)) + "\n"
			}
		}
		if e.Cookie != "" {
			code += "局_提交Cookie ＝ 子文本替换 (" + convertELangFormat(e.Cookie) + ", “'”, #引号, , , 真)\n\n"
		}

		for k, v := range e.Header {
			if strings.ToUpper(k) == "CONTENT-LENGTH" {
				continue
			}
			if k == "Accept-Encoding" {
				if len(v) < 1 {
					code += "' 局_HTTP.置请求头 (“" + k + "”, “”)\n"
				} else {
					code += "' 局_HTTP.置请求头 (“" + k + "”, " + convertELangFormat(v[0]) + ")\n"
				}
				continue
			}
			if len(v) < 1 {
				code += "局_HTTP.置请求头 (“" + k + "”, “”)\n"
			} else {
				code += "局_HTTP.置请求头 (“" + k + "”, " + convertELangFormat(v[0]) + ")\n"
			}
		}
		if e.Cookie != "" {
			code += "局_HTTP.置请求头 (“Cookie”, 局_提交Cookie)\n\n"
		}
		if e.Method == "GET" {
			code += "网站客户端.执行GET (局_网址 , 局_结果, 真, )\n"
		} else {
			code += "网站客户端.执行POST (局_网址,局_提交数据, 局_响应字节集, 真, )\n"
		}
		code += "局_响应文本 ＝ 到文本 (局_响应字节集)\n"
		code += "调试输出 (局_响应文本)\n\n"
	}
	if module == "网页_访问" || module == "网页_访问_对象" {
		code += ".子程序 " + e.FuncName + ", , 公开, 本子程序由Sunny中间件生成,请配合 [精易模块] 使用。 " + e.URL.Path + "\n"
		BytesType := e.IsBytesType()
		code += ".局部变量 局_网址, 文本型, , , \n"
		if BytesType {
			code += ".局部变量 局_提交数据, 字节集, , , \n"
		} else {
			code += ".局部变量 局_提交数据, 文本型, , , \n"
		}
		code += ".局部变量 局_协议头, 类_POST数据类, , , \n"

		if e.Cookie != "" {
			code += ".局部变量 局_提交Cookie, 文本型, , , \n"
		}
		code += ".局部变量 局_响应字节集, 字节集, , , \n.局部变量 局_响应文本, 文本型, , , \n\n"

		code += "局_网址 ＝ " + convertELangFormat(e.URL.String()) + "\n\n"

		if BytesType {
			code += "局_提交数据 ＝ 编码_BASE64解码 (“" + base64.StdEncoding.EncodeToString(e.Body) + "”)\n\n"
		} else {
			ok, d := e.IsFormData()
			if ok {
				code += d + "\n"
			} else {
				code += "局_提交数据 ＝ " + convertELangFormat(string(e.Body)) + "\n"
			}
		}
		if e.Cookie != "" {
			code += "局_提交Cookie ＝ 子文本替换 (" + convertELangFormat(e.Cookie) + ", “'”, #引号, , , 真)\n\n"
		}

		for k, v := range e.Header {
			if strings.ToUpper(k) == "CONTENT-LENGTH" {
				continue
			}
			if k == "Accept-Encoding" {
				if len(v) < 1 {
					code += "' 局_协议头.添加 (“" + k + "”, “”)\n"
				} else {
					code += "' 局_协议头.添加 (“" + k + "”, " + convertELangFormat(v[0]) + ")\n"
				}
				continue
			}
			if len(v) < 1 {
				code += "局_协议头.添加 (“" + k + "”, “”)\n"
			} else {
				code += "局_协议头.添加 (“" + k + "”, " + convertELangFormat(v[0]) + ")\n"
			}
		}
		if e.Cookie != "" {
			code += "局_协议头.添加 (“Cookie”, 局_提交Cookie)\n\n"
		}
		if e.Method == "GET" {
			code += "局_响应字节集 ＝ " + module + " (局_网址, 0, , , , 局_协议头.获取协议头数据 ())\n"
		} else {
			mod := "1"
			if e.Method == "POST" {
				mod = "1"
			} else if e.Method == "HEAD" {
				mod = "2"
			} else if e.Method == "PUT" {
				mod = "3"
			} else if e.Method == "OPTIONS" {
				mod = "4"
			} else if e.Method == "DELETE" {
				mod = "5"
			} else if e.Method == "TRACE" {
				mod = "6"
			} else if e.Method == "CONNECT" {
				mod = "7"
			}
			if BytesType {
				if module == "网页_访问_对象" {
					code += "局_响应字节集 ＝ " + module + " (局_网址, " + mod + ", , , , 局_协议头.获取协议头数据 (),,,,局_提交数据)\n"
				} else {
					code += "局_响应字节集 ＝ " + module + " (局_网址, " + mod + ", , , , 局_协议头.获取协议头数据 (),,,局_提交数据)\n"
				}

			} else {
				code += "局_响应字节集 ＝ " + module + " (局_网址, " + mod + ", 局_提交数据, , , 局_协议头.获取协议头数据 ())\n"
			}
		}
		code += "局_响应文本 ＝ 到文本 (局_响应字节集)\n"
		code += "调试输出 (局_响应文本)\n\n"
	}
	return code
}

func (e *CreateRequest) Python(module string) string {
	code := ""
	if module == "requests" {
		_header := ""
		for k, v := range e.Header {
			if strings.ToUpper(k) == "CONTENT-LENGTH" {
				continue
			}
			if len(v) < 1 {
				_header += `        '` + k + `': "",` + "\n"
			} else {
				_header += `        '` + k + `': "` + strReplaceAll([]byte(v[0])) + `",` + "\n"
			}
		}
		if e.Cookie != "" {
			_header += `        'Cookie': "` + strReplaceAll([]byte(e.Cookie)) + `",` + "\n"
		}
		_t := `def ` + e.FuncName + `():
    """ 
    [ ` + e.URL.Path + ` ]
    本函数由SunnyNet网络中间件生成   
    """
    url = "` + e.URL.String() + `"
    payload = "` + strReplaceAll(e.Body) + `"
    headers = {
` + _header + `
    }
    response = requests.request("` + e.Method + `", url, data=payload, headers=headers)

    print(response.text)
`
		return _t
	}
	if module == "Flurl" {
		return code
	}
	return code
}
func strReplaceAll(abody []byte) string {
	ss := strings.ReplaceAll(string(abody), "\\", "\\\\")
	ss = strings.ReplaceAll(ss, "\"", "\\\"")
	return ss
}
func (e *CreateRequest) CSharp(module string) string {
	code := ""

	if module == "RestSharp" {
		templateData := ""
		templateData1 := ""
		if len(e.Body) > 0 {
			if e.IsBytesType() {
				templateData = `string base64String = "` + base64.StdEncoding.EncodeToString(e.Body) + `";  
            byte[] body = Convert.FromBase64String(base64String); `
			} else {
				templateData = `string String = "` + strReplaceAll(e.Body) + `";  
            byte[] body = Encoding.Default.GetBytes(String); `
			}
		}

		mod := "Method.Post"
		s := strings.ToUpper(e.Method)
		{
			if s == "POST" {
				mod = "Method.Post"
			}
			if s == "GET" {
				mod = "Method.Get"
			}
			if s == "PUT" {
				mod = "Method.Put"
			}
			if s == "DELETE" {
				mod = "Method.Delete"
			}
			if s == "HEAD" {
				mod = "Method.Head"
			}
			if s == "OPTIONS" {
				mod = "Method.Options"
			}
			if s == "PATCH" {
				mod = "Method.Patch"
			}
			if s == "MERGE" {
				mod = "Method.Merge"
			}
			if s == "COPY" {
				mod = "Method.Copy"
			}
			if s == "SEARCH" {
				mod = "Method.Search"
			}
		}
		header := ""
		CONTENTMENT := "application/x-www-form-urlencoded"
		for k, v := range e.Header {
			if strings.ToUpper(k) == "CONTENT-LENGTH" {
				continue
			}
			if strings.ToUpper(k) == "ACCEPT-ENCODING" {
				if len(v) > 0 {
					header += `            //request.AddHeader("` + k + `","` + strReplaceAll([]byte(v[0])) + `")` + ";\n"
				}
				continue
			}
			if strings.ToUpper(k) == "CONTENT-TYPE" {
				if len(v) > 0 {
					ss := strings.Split(v[0]+";", ";")
					if len(ss) > 0 {
						CONTENTMENT = ss[0]
					}
				}

			}
			if len(v) < 1 {
				//request.AddHeader("cookie", "xxxxx");
				header += `            request.AddHeader("` + k + `","")` + ";\n"
			} else {
				header += `            request.AddHeader("` + k + `","` + strReplaceAll([]byte(v[0])) + `")` + ";\n"
			}
		}
		if e.Cookie != "" {
			header += `            request.AddHeader("Cookie","` + strReplaceAll([]byte(e.Cookie)) + `")` + ";\n"
		}
		if len(e.Body) > 0 {
			templateData1 = `            request.AddParameter("` + CONTENTMENT + `", body, ParameterType.RequestBody); 
            `
		}
		_tmp := `/// <summary> 
        /// ` + e.FuncName + ` ` + e.URL.Path + `
        ///<para>本函数由SunnyNet网络中间件生成</para> 
        /// </summary> 
        public static void ` + e.FuncName + `()
        {
            string url = "` + e.URL.String() + `";
            ` + templateData + ` 
            var client = new RestClient(url);
            var request = new RestRequest("",` + mod + `);
` + header + templateData1 + `var response = client.Execute(request); 
            Trace.WriteLine("Response StateCode:" + ((int)response.StatusCode)); 
            Trace.WriteLine("Response Text:\n" + response.Content); 
        }
`
		return _tmp
	}
	if module == "HttpClient" {
		templateData := `byte[] data = Encoding.Default.GetBytes(""); `
		if len(e.Body) > 0 {
			if e.IsBytesType() {
				templateData = `string base64String = "` + base64.StdEncoding.EncodeToString(e.Body) + `";  
            byte[] data = Convert.FromBase64String(base64String); `
			} else {
				templateData = `string String = "` + strReplaceAll(e.Body) + `";  
            byte[] data = Encoding.Default.GetBytes(String); `
			}
		}
		mod := "Method.Post"
		CONTENTMENT := "application/x-www-form-urlencoded"
		s := strings.ToUpper(e.Method)
		{
			if s == "POST" {
				mod = "HttpMethod.Post"
			}
			if s == "GET" {
				mod = "HttpMethod.Get"
			}
			if s == "PUT" {
				mod = "HttpMethod.Put"
			}
			if s == "DELETE" {
				mod = "HttpMethod.Delete"
			}
			if s == "HEAD" {
				mod = "HttpMethod.Head"
			}
			if s == "OPTIONS" {
				mod = "HttpMethod.Options"
			}
			if s == "PATCH" {
				mod = "HttpMethod.Patch"
			}
			if s == "MERGE" {
				mod = "HttpMethod.Merge"
			}
			if s == "COPY" {
				mod = "HttpMethod.Copy"
			}
			if s == "SEARCH" {
				mod = "HttpMethod.Search"
			}
		}
		header := ""
		for k, v := range e.Header {
			if strings.ToUpper(k) == "CONTENT-LENGTH" {
				continue
			}
			if strings.ToUpper(k) == "CONTENT-TYPE" {
				if len(v) > 0 {
					ss := strings.Split(v[0]+";", ";")
					if len(ss) > 0 {
						CONTENTMENT = ss[0]
					}
				}
				continue
			}
			if strings.ToUpper(k) == "ACCEPT-ENCODING" {
				if len(v) > 0 {
					header += `            //client.DefaultRequestHeaders.Add("` + k + `","` + strReplaceAll([]byte(v[0])) + `")` + ";\n"
				}
				continue
			}
			if len(v) < 1 {
				header += `            client.DefaultRequestHeaders.Add("` + k + `","")` + ";\n"
			} else {
				header += `            client.DefaultRequestHeaders.Add("` + k + `","` + strReplaceAll([]byte(v[0])) + `")` + ";\n"
			}
		}
		if e.Cookie != "" {
			header += `            client.DefaultRequestHeaders.Add("Cookie","` + strReplaceAll([]byte(e.Cookie)) + `")` + ";\n"
		}
		_tmp := `/// <summary> 
        /// ` + e.FuncName + ` ` + e.URL.Path + `
        ///<para>本函数由SunnyNet网络中间件生成</para> 
        /// </summary> 
        public static async void ` + e.FuncName + `()
        {
            string url = "` + e.URL.String() + `";
            ` + templateData + `
            using (HttpClient client = new HttpClient())
            {
                ` + header + `
                HttpContent content = new ByteArrayContent(data);
                content.Headers.ContentType = new System.Net.Http.Headers.MediaTypeHeaderValue("` + CONTENTMENT + `");
                HttpRequestMessage request = new HttpRequestMessage(` + mod + `, url);
                request.Content = content;
                HttpResponseMessage response = client.SendAsync(request).Result;
 
                Trace.WriteLine("Response Status code: " + response.StatusCode);
                byte[] responseBytes = response.Content.ReadAsByteArrayAsync().Result;
                Trace.WriteLine("Response bytes: " + BitConverter.ToString(responseBytes));
                Trace.WriteLine("Response Content: " + response.Content.ReadAsStringAsync().Result);
                 
            }
        }
`
		return _tmp
	}

	return code
}

func (e *CreateRequest) Go(module string) string {
	templateData := ""
	if e.IsBytesType() {
		templateData = `	_Base64 := "` + base64.StdEncoding.EncodeToString(e.Body) + `"
	_Data, _ := base64.StdEncoding.DecodeString(_Base64)
	Body := io.NopCloser(bytes.NewBuffer(_Data))`
	} else {
		templateData = `	Body := io.NopCloser(bytes.NewBuffer([]byte("` + strReplaceAll(e.Body) + `")))`
	}
	header := ""

	for k, v := range e.Header {
		if strings.ToUpper(k) == "CONTENT-LENGTH" {
			continue
		}
		if len(v) < 1 {
			header += `	e.Header.Set("` + k + `","")` + "\n"
		} else {
			header += `	e.Header.Set("` + k + `","` + strReplaceAll([]byte(v[0])) + `")` + "\n"
		}
	}
	if e.Cookie != "" {
		header += `	e.Header.Set("Cookie","` + strReplaceAll([]byte(e.Cookie)) + `")` + "\n"
	}
	template := `// ` + e.FuncName + ` 本函数由SunnyNet网络中间件生成  //` + e.URL.Path + `
func ` + e.FuncName + `() {
` + templateData + `
	defer func() { _ = Body.Close() }()
	req, err := http.NewRequest("` + e.Method + `", "` + e.URL.String() + `", Body)
	if err != nil {
		panic(err)
	}
	` + header + `
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() { _ = res.Body.Close() }()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("响应状态码", res.StatusCode)
	for k, v := range res.Header {
		fmt.Println(k, v)
	}
	fmt.Println(string(body))
}
`
	return template
}
func (e *CreateRequest) IsBytesType() bool {
	for _, v := range e.Body {
		if v < 9 {
			return true
		}
	}
	return false
}
func (e *CreateRequest) IsFormData() (bool, string) {
	p := string(e.Body)
	if !strings.Contains(p, "&") {
		return false, ""
	}
	if !strings.Contains(p, "=") {
		return false, ""
	}
	if strings.HasPrefix(p, "{") || strings.HasPrefix(p, "[") {
		return false, ""
	}
	Array := strings.Split(p, "&")
	Code := ""
	for index, v := range Array {
		Array1 := strings.Split(v, "=")
		if len(Array1) == 1 {
			Code += "局_提交数据 ＝ “" + Array1[0] + "=”\n"
		} else if len(Array1) == 2 {
			value, ex := url.QueryUnescape(Array1[1])
			if ex != nil {
				if index == 0 {
					Code += "局_提交数据 ＝ “" + Array1[0] + "=" + Array1[1] + "”\n"
				} else {
					Code += "局_提交数据 ＝ 局_提交数据 ＋ “" + Array1[0] + "=" + Array1[1] + "”\n"
				}
				continue
			}
			if index == 0 {
				if isChinese(value) || value != Array1[1] {
					Code += "局_提交数据 ＝ “" + Array1[0] + "=” ＋ 编码_URL编码 (" + convertELangFormat(value) + ",真,真)\n"
				} else {
					Code += "局_提交数据 ＝ “" + Array1[0] + "=” ＋ " + convertELangFormat(value) + "\n"
				}
			} else {
				if isChinese(value) || value != Array1[1] {
					Code += "局_提交数据 ＝ 局_提交数据 ＋ “" + Array1[0] + "=” ＋ 编码_URL编码 (" + convertELangFormat(value) + ",真,真)\n"
				} else {
					Code += "局_提交数据 ＝ 局_提交数据 ＋ “" + Array1[0] + "=” ＋ " + convertELangFormat(value) + "\n"
				}
			}
		}
	}

	return true, Code
}
func isChinese(str string) bool {
	for _, v := range str {
		if v > 255 {
			return true
		}
	}
	return false
}

func convertELangFormat(v string) string {
	str := v
	str = strings.ReplaceAll(str, "“", "\"")
	str = strings.ReplaceAll(str, "”", "\"")
	str = strings.ReplaceAll(str, "\r\n", "\n")
	str = strings.ReplaceAll(str, "\n", "\r\n")
	if v == "\r\n" || v == "\n" || v == "\r" {
		return "#换行符"
	}
	if !strings.Contains(str, "\"") && !strings.Contains(str, "\r\n") {
		return "“" + str + "”"
	}
	if strings.Contains(str, "\r\n") {
		arr := strings.Split(str, "\r\n")
		str = ""
		a := false
		for _, va := range arr {
			if !strings.Contains(va, "gzip, ") && !strings.Contains(va, "Content-Length:") && !strings.Contains(va, "Accept-Encoding: gzip") {
				str += va + "\r\n"
				a = true
			}
		}
		if a && strings.HasSuffix(str, "\r\n") {
			str = str[:len(str)-2]
		}
	}
	fh := ""
	issc := false
	if strings.Contains(str, "\"") {
		if !strings.Contains(str, "'") {
			fh = "'"
		} else if !strings.Contains(str, "#") {
			fh = "#"
		} else if !strings.Contains(str, "~") {
			fh = "~"
		} else if !strings.Contains(str, "!") {
			fh = "!"
		} else if !strings.Contains(str, "|") {
			fh = "|"
		} else if !strings.Contains(str, "/") {
			fh = "/"
		} else if !strings.Contains(str, "\\") {
			fh = "\\"
		} else if !strings.Contains(str, "&") {
			fh = "&"
		} else if !strings.Contains(str, "*") {
			fh = "*"
		} else {
			issc = true
			str = "\"" + strings.ReplaceAll(str, "\"", "\"+#引号+\"") + "\""
			fh = ""
		}
		if fh != "" {
			str = strings.ReplaceAll(str, "\"", fh)
		}
	}

	if strings.Contains(str, "\r\n") {
		if issc {
			str = strings.ReplaceAll(str, "\r\n", "\"+#换行符+\"")
		} else {
			str = "\"" + strings.ReplaceAll(str, "\r\n", "\"+#换行符+\"") + "\""
		}
	}
	str = strings.ReplaceAll(str, "\"\"", "")
	if strings.HasPrefix(str, "+") {
		str = str[1:]
	}
	if strings.HasPrefix(str, "+") {
		str = str[1:]
	}
	if strings.Contains(str, "++") {
		str = strings.ReplaceAll(str, "++", "+")
	}
	if fh != "" {
		str = "子文本替换 (“" + str + "”, \"" + fh + "\",#引号, , , 真)"
	}
	return str
}
