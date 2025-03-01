package model

import "github.com/coroot/coroot/timeseries"

type Pod struct {
	Phase     string
	Reason    string
	Scheduled bool

	Running  timeseries.TimeSeries
	Ready    timeseries.TimeSeries
	LifeSpan timeseries.TimeSeries

	ReplicaSet string

	InitContainers map[string]*Container
}

func (pod *Pod) IsRunning() bool {
	return pod.Phase == "Running"
}

func (pod *Pod) IsPending() bool {
	return pod.Phase == "Pending"
}

func (pod *Pod) IsObsolete() bool {
	return pod.Phase == ""
}

func (pod *Pod) IsReady() bool {
	return pod.Ready != nil && pod.Ready.Last() > 0
}
