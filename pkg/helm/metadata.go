package helm

type Metadata struct {
	AppVersion        string `json:"app_version"`
	Chart             string `json:"chart"`
	FirstDeployed     float64 `json:"first_deployed"`
	LastDeployed      float64 `json:"last_deployed"`
	Name              string `json:"name"`
	Namespace         string `json:"namespace"`
	Notes             string `json:"notes"`
	Revision          int    `json:"revision"`
	Values            string `json:"values"`
	Version           string `json:"version"`
}

