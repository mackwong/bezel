package model

// Gloabl Config Section

type GlobalConfig struct {
	Diamond  *DiamondConfig
	Machines []*MachineConfig
}

// DiamondConfig Items Section

type DiamondConfig struct {
	Name           string `yaml:"name"`
	Arranger       string `yaml:"arranger"`
	UpstreamDNS    string `yaml:"upstreamDNS"`
	DockerRegistry string `yaml:"dockerRegistry"`
	MachineNum     int64  `yaml:"machine-num"`
	MasterNum      int64  `yaml:"masternum"`
	K8sMasterIP    string `yaml:"k8sMasterIP"`
}

func NewDefaultDiamondConfig() *DiamondConfig {
	return &DiamondConfig{
		Name:           "diamond-edge-ha",
		Arranger:       "edgesite",
		UpstreamDNS:    "114.114.114.114",
		DockerRegistry: "10.5.49.73",
		MachineNum:     4,
		MasterNum:      3,
		K8sMasterIP:    "10.4.72.231",
	}
}

// MachineConfig Items Section

type MachineConfig struct {
	Name      string `yaml:"name"`
	HostName  string `yaml:"hostname"`
	Role      string `yaml:"role"`
	IP        string `yaml:"ip"`
	GatewayIP string `yaml:"gatewayIP"`
	Netmask   string `yaml:"netmask"`
}

func NewDefaultMachineConfig() *MachineConfig {
	return &MachineConfig{
		Name:      "test-00",
		HostName:  "test-00",
		Role:      "master",
		IP:        "10.4.72.140",
		GatewayIP: "10.4.72.1",
		Netmask:   "255.255.255.0",
	}
}

// Sub Config Section
type SubConfig struct {
	Arranger       string `yaml:"arranger"`
	Role           string `yaml:"role"`
	UpstreamDNS    string `yaml:"upstreamDNS"`
	Hostname       string `yaml:"hostname"`
	IP             string `yaml:"ip"`
	GatewayIP      string `yaml:"gatewayIP"`
	Netmask        string `yaml:"netmask"`
	K8sMasterIP    string `yaml:"k8sMasterIP"`
	DockerRegistry string `yaml:"dockerRegistry"`
	HaPeer         []Peer `yaml:"ha-peer"`
}

// Ha Peer Section
type Peer struct {
	Hostname string `yaml:"hostname"`
	IP       string `yaml:"ip"`
}

type BezelConfig struct {
	Name           string     `yaml:"name"`
	MachineNum     int64      `yaml:"machine-num"`
	MasterNum      int64      `yaml:"master-num"`
	Arranger       string     `yaml:"arranger"`
	UpstreamDNS    string     `yaml:"upstream-dns"`
	DockerRegistry string     `yaml:"docker-registry"`
	K8sMasterIP    string     `yaml:"k8sMaster-ip"`
	IPRange        []IPConfig `yaml:"ip-range"`
	MasterIP       []string   `yaml:"master-ip"`
	NameFormat     string     `yaml:"name-format"`
	HostNameFormat string     `yaml:"hostname-format"`
}

type IPConfig struct {
	IPRange   string `yaml:"ipRange"`
	GatewayIP string `yaml:"gatewayIP"`
	Netmask   string `yaml:"netmask"`
}

func NewSampleBezelConfig() *BezelConfig {
	return &BezelConfig{
		Name:           "diamond-edge-ha",
		Arranger:       "edgesite",
		UpstreamDNS:    "114.114.114.114",
		DockerRegistry: "10.5.49.73",
		K8sMasterIP:    "10.4.72.231",
		MasterNum:      3,
		MachineNum:     4,
		NameFormat:     "node-{{.Role}}-{{.Index}}",
		HostNameFormat: "ubuntu-{{.Role}}-{{.Index}}",
		MasterIP: []string{
			"10.4.72.1",
			"10.4.72.2",
			"10.4.73.1",
		},
		IPRange: []IPConfig{
			{
				IPRange:   "10.4.72.1/24",
				GatewayIP: "10.4.72.254",
				Netmask:   "255.255.255.0",
			},
			{
				IPRange:   "10.4.73.1/32",
				GatewayIP: "10.4.73.254",
				Netmask:   "255.255.255.255",
			},
		},
	}
}

type Formatter struct {
	Role  string `json:"role"`
	IP    string `json:"ip"`
	Index int    `json:"index"`
}

func NewHaPeer(haPeer map[string]string) []Peer {
	peers := []Peer{}
	for hostName, ip := range haPeer {
		p := Peer{
			Hostname: hostName,
			IP:       ip,
		}
		peers = append(peers, p)
	}
	return peers
}
