package model

type DeploymentModel struct {
	Name        string
	Replicas    [3]int32
	Images      string
	NameSpace   string
	IsCompleted bool   // 是否就绪
	Message     string // 错误信息
	CreatedAt   string
	Pods        []*PodModel
}
