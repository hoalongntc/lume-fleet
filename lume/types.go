package lume

// VM represents a VM as returned by the Lume API.
type VM struct {
	Name                  string    `json:"name"`
	Status                string    `json:"status"`
	CPUCount              int       `json:"cpuCount"`
	MemorySize            int64     `json:"memorySize"` // bytes
	DiskSize              *DiskSize `json:"diskSize"`
	IPAddress             *string   `json:"ipAddress"`
	OS                    string    `json:"os"`
	Display               string    `json:"display"`
	VNCUrl                *string   `json:"vncUrl"`
	SSHAvailable          *bool     `json:"sshAvailable"`
	LocationName          string    `json:"locationName"`
	SharedDirectories     []string  `json:"sharedDirectories"`
	ProvisioningOperation *string   `json:"provisioningOperation"`
}

type DiskSize struct {
	Total     int64 `json:"total"`
	Allocated int64 `json:"allocated"`
}

// CreateRequest is the POST /lume/vms body.
type CreateRequest struct {
	Name       string `json:"name"`
	OS         string `json:"os"`
	CPU        int    `json:"cpu"`
	Memory     string `json:"memory"`
	DiskSize   string `json:"diskSize"`
	Display    string `json:"display"`
	IPSW       string `json:"ipsw,omitempty"`
	Unattended string `json:"unattended,omitempty"`
	VNCPort    int    `json:"vncPort,omitempty"`
	Storage    string `json:"storage,omitempty"`
	Network    string `json:"network,omitempty"`
}

// RunRequest is the POST /lume/vms/{name}/run body.
type RunRequest struct {
	NoDisplay bool   `json:"noDisplay"`
	SharedDir string `json:"sharedDir,omitempty"`
}
