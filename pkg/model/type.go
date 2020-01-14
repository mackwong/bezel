package model

// Gloabl Config Section

type GlobalConfig struct {
	Diamond  Diamond
	Machines []Machine
}

// Diamond Items Section

type Diamond struct {
	Name           string `yaml:"name"`
	Arranger       string `yaml:"arranger"`
	UpstreamDNS    string `yaml:"upstreamDNS"`
	DockerRegistry string `yaml:"dockerRegistry"`
	MachineNum     string `yaml:"machine-num"`
	MasterNum      string `yaml:"masternum"`
	K8sMasterIP    string `yaml:"k8sMasterIP"`
}

func NewGlobalConfig(di map[string]string, mal map[int]map[string]string) *GlobalConfig {
	return &GlobalConfig{
		Diamond: Diamond{
			Name:           di["name"],
			Arranger:       di["arranger"],
			UpstreamDNS:    di["upstreamDNS"],
			DockerRegistry: di["dockerRegistry"],
			MachineNum:     di["machine-num"],
			MasterNum:      di["master-num"],
			K8sMasterIP:    di["k8sMasterIP"],
		},
		Machines: MachineInfoFromMap(mal),
	}
}

// Machine Items Section

type Machine struct {
	Name      string `yaml:"name"`
	Hostname  string `yaml:"hostname"`
	Role      string `yaml:"role"`
	IP        string `yaml:"ip"`
	GatewayIP string `yaml:"gatewayIP"`
	Netmask   string `yaml:"netmask"`
}

func MachineInfoFromMap(mal map[int]map[string]string) []Machine {
	machines := []Machine{}
	for _, ma := range mal {
		machine := &Machine{
			Name:      ma["name"],
			Hostname:  ma["hostname"],
			Role:      ma["role"],
			IP:        ma["ip"],
			GatewayIP: ma["gatewayIP"],
			Netmask:   ma["netmask"],
		}
		machines = append(machines, *machine)
	}
	return machines
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

func NewSubConfig(sc map[string]string, hp []Peer) *SubConfig {
	return &SubConfig{
		Arranger:       sc["arranger"],
		Role:           sc["role"],
		UpstreamDNS:    sc["upstreamDNS"],
		Hostname:       sc["hostname"],
		IP:             sc["ip"],
		GatewayIP:      sc["gatewayIP"],
		Netmask:        sc["netmask"],
		K8sMasterIP:    sc["k8sMasterIP"],
		DockerRegistry: sc["dockerRegistry"],
		HaPeer:         hp,
	}
}

// Ha Peer Section

type Peer struct {
	Hostname string `yaml:"hostname"`
	IP       string `yaml:"ip"`
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
