package model

import (
	"github.com/coroot/coroot-focus/timeseries"
)

type World struct {
	Ctx timeseries.Context

	Nodes        []*Node
	Applications []*Application
	Services     []*Service
}

func (w *World) GetApplication(id ApplicationId) *Application {
	for _, a := range w.Applications {
		if a.Id == id {
			return a
		}
	}
	return nil
}

func (w *World) GetOrCreateApplication(id ApplicationId) *Application {
	app := w.GetApplication(id)
	if app == nil {
		app = NewApplication(id)
		w.Applications = append(w.Applications, app)
	}
	return app
}

func (w *World) GetServiceForConnection(c *Connection) *Service {
	for _, s := range w.Services {
		if s.ClusterIP == c.ServiceRemoteIP {
			return s
		}
		for _, sc := range s.Connections {
			if sc.ActualRemoteIP == c.ActualRemoteIP {
				return s
			}
		}
	}
	return nil
}

func (w *World) FindInstanceByListen(ip, port string) *Instance {
	l := Listen{IP: ip, Port: port}
	for _, app := range w.Applications {
		for _, i := range app.Instances {
			if i.TcpListens[l] {
				return i
			}
		}
	}
	return nil
}

func (w *World) FindInstanceByPod(ns, pod string) *Instance {
	for _, app := range w.Applications {
		if app.Id.Namespace != ns {
			continue
		}
		for _, i := range app.Instances {
			if i.Pod != nil && i.Name == pod {
				return i
			}
		}
	}
	return nil
}
