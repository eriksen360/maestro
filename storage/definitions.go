package storage

type HostConfig struct {
	Hosts []Host
}

type Host struct {
	Fqdn   string   `yaml:"fqdn"`
	IpAddr string   `yaml:"ipaddr"`
	PubKey string   `yaml:"pubkey"`
	Groups []string `yaml:"groups"`
}

type TaskConfig struct {
	Tasks []Task
}

type Task struct {
	Name  string
	User  string   `yaml:"user"`
	Hosts []string `yaml:"hosts"`
	/* groups []string
	   schedule time.Time
	   persistSession bool
	   logs: Log */
	Instructions []Instruction `mapstructure:"instructions"`
}

type Instruction struct {
	Name    string `yaml:"name"`
	Command string `yaml:"command"`
	/*type string
	  fileSrc string
	  fileDst string
	  dependencies []Dependency
	  retries int */
}
