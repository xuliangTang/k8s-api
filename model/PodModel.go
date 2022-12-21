package model

type PodModel struct {
	Name      string
	NodeName  string
	Images    string
	Phase     string // 阶段
	IsReady   bool   // pod 是否就绪
	Message   string
	CreatedAt string
}
