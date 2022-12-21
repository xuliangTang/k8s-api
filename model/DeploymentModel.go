package model

type DeploymentModel struct {
	Name      string
	Replicas  [3]int32
	Images    string
	NameSpace string
	CreatedAt string
	Pods      []*PodModel
}
