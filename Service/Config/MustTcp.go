package Config

type MustTcpType int

const (
	MustTcpTypeAll MustTcpType = 0 //全部强制TCP
	MustTcpTypeLei MustTcpType = 1 //规则内走TCP
	MustTcpTypeWai MustTcpType = 2 //规则外走TCP
)

func (f *config) initMustTcp() {
	//恢复上游代理网关列表
	{
		f.MustTcp.Type = MustTcpTypeLei
		if f.MustTcp.Roles == "" {
			f.MustTcp.Roles = `// 多个请用 “;” 分号分割  
// 或 换行(一行一个) 
// 使用 “//” 开头表示注释,不会生效

// 修改完成后,请记得右键 “保存/应用修改” 哦

// 以下为示例

text.com;*.bb.com;61.61.*.*;
//1.2.3.*;
1.2.4.5
6.6.6.6`
		}
	}
}
