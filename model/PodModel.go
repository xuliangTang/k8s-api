package model

type PodModel struct {
	Name      string
	NameSpace string
	NodeName  string
	Images    string
	Ip        [2]string // 0 podIp 1 hostIp
	Phase     string    // 阶段
	IsReady   bool      // pod 是否就绪
	Message   string
	CreatedAt string
}
