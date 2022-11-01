package v1

import consts "github.com/easysoft/zagent/internal/pkg/const"

type ServiceCheckReq struct {
	Services string `json:"services"`
}

type ServiceCheckResp struct {
	Code       string                   `json:"code"`
	Kvm        consts.HostServiceStatus `json:"kvm"`        // Enums consts.HostServiceStatus
	Novnc      consts.HostServiceStatus `json:"novnc"`      // Enums consts.HostServiceStatus
	Websockify consts.HostServiceStatus `json:"websockify"` // Enums consts.HostServiceStatus
}
