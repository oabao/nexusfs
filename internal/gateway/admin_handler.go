package gateway

import (
	"encoding/json"
	"net/http"
)

// AdminDashboardData represents the data structure for the admin dashboard.
type AdminDashboardData struct {
	DashboardData DashboardData `json:"dashboardData"`
	SystemStatus  []SystemStatus `json:"systemStatus"`
	HealingStatus HealingStatus `json:"healingStatus"`
	Buckets       []Bucket       `json:"buckets"`
	Configuration Configuration `json:"configuration"`
}

type DashboardData struct {
	TotalCapacity float64 `json:"totalCapacity"`
	UsedCapacity  float64 `json:"usedCapacity"`
}

type SystemStatus struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Count   int    `json:"count"`
}

type HealingStatus struct {
	LastScrub   string `json:"lastScrub"`
	ErrorsFound int    `json:"errorsFound"`
	Status      string `json:"status"`
}

type Bucket struct {
	Name          string   `json:"name"`
	CreatedAt     string   `json:"createdAt"`
	StoragePolicy string   `json:"storagePolicy"`
	Objects       []Object `json:"objects"`
}

type Object struct {
	Name         string `json:"name"`
	Size         string `json:"size"`
	LastModified string `json:"lastModified"`
}

type Configuration struct {
	Gateway       GatewayConfig       `json:"gateway"`
	Storage       StorageConfig       `json:"storage"`
	Log           LogConfig           `json:"log"`
}

type GatewayConfig struct {
	Address      string `json:"address"`
	AuthProvider string `json:"authProvider"`
}

type StorageConfig struct {
	DefaultErasurePolicy string `json:"defaultErasurePolicy"`
}

type LogConfig struct {
	Level string `json:"level"`
}

// AdminDashboardHandler returns a JSON response with mock data for the admin dashboard.
func AdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	data := AdminDashboardData{
		DashboardData: DashboardData{
			TotalCapacity: 100,
			UsedCapacity:  45.5,
		},
		SystemStatus: []SystemStatus{
			{Name: "Gateway Nodes", Status: "Healthy", Count: 3},
			{Name: "Storage Nodes", Status: "Healthy", Count: 12},
			{Name: "Metadata DB", Status: "Healthy", Count: 1},
			{Name: "Cluster Operator", Status: "Healthy", Count: 1},
		},
		HealingStatus: HealingStatus{
			LastScrub:   "2 hours ago",
			ErrorsFound: 0,
			Status:      "Idle",
		},
		Buckets: []Bucket{
			{
				Name:          "customer-invoices",
				CreatedAt:     "2023-01-15",
				StoragePolicy: "EC 4+2",
				Objects: []Object{
					{Name: "invoice-2023-01.pdf", Size: "1.2 MB", LastModified: "2023-01-15 10:30"},
					{Name: "invoice-2023-02.pdf", Size: "1.4 MB", LastModified: "2023-02-15 11:00"},
				},
			},
			{
				Name:          "archived-logs",
				CreatedAt:     "2023-03-20",
				StoragePolicy: "EC 6+3",
				Objects: []Object{
					{Name: "gateway-2023-03-20.log.gz", Size: "55.8 MB", LastModified: "2023-03-20 18:00"},
					{Name: "storage-node-1-2023-03-20.log.gz", Size: "78.2 MB", LastModified: "2023-03-20 18:00"},
				},
			},
		},
		Configuration: Configuration{
			Gateway: GatewayConfig{
				Address:      "0.0.0.0:9000",
				AuthProvider: "https://keycloak.example.com/auth/realms/nexusfs",
			},
			Storage: StorageConfig{
				DefaultErasurePolicy: "EC 4+2",
			},
			Log: LogConfig{
				Level: "info",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(data)
}
