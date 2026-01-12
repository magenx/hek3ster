package config

// NodePool represents common node pool configuration
type NodePool struct {
	Name                           *string      `yaml:"name,omitempty"`
	LegacyInstanceType             string       `yaml:"legacy_instance_type,omitempty"`
	InstanceType                   string       `yaml:"instance_type"`
	Image                          interface{}  `yaml:"image,omitempty"` // Can be string or int64
	InstanceCount                  int          `yaml:"instance_count,omitempty"`
	Labels                         []Label      `yaml:"labels,omitempty"`
	Taints                         []Taint      `yaml:"taints,omitempty"`
	Autoscaling                    *Autoscaling `yaml:"autoscaling,omitempty"`
	AdditionalPreK3sCommands       []string     `yaml:"additional_pre_k3s_commands,omitempty"`
	AdditionalPostK3sCommands      []string     `yaml:"additional_post_k3s_commands,omitempty"`
	AdditionalPackages             []string     `yaml:"additional_packages,omitempty"`
	IncludeClusterNameAsPrefix     bool         `yaml:"include_cluster_name_as_prefix,omitempty"`
	GrowRootPartitionAutomatically *bool        `yaml:"grow_root_partition_automatically,omitempty"`
}

// AutoscalingEnabled returns true if autoscaling is enabled for this pool
func (n *NodePool) AutoscalingEnabled() bool {
	return n.Autoscaling != nil && n.Autoscaling.Enabled
}

// EffectiveGrowRootPartitionAutomatically returns the effective value for grow root partition
func (n *NodePool) EffectiveGrowRootPartitionAutomatically(globalValue bool) bool {
	if n.GrowRootPartitionAutomatically != nil {
		return *n.GrowRootPartitionAutomatically
	}
	return globalValue
}

// MasterNodePool represents master node pool configuration
type MasterNodePool struct {
	NodePool  `yaml:",inline"`
	Locations []string `yaml:"locations,omitempty"`
}

// SetDefaults sets default values for master pool
func (m *MasterNodePool) SetDefaults() {
	if len(m.Locations) == 0 {
		m.Locations = []string{"fsn1"}
	}
	if m.InstanceCount == 0 {
		m.InstanceCount = 1
	}
	if !m.IncludeClusterNameAsPrefix {
		m.IncludeClusterNameAsPrefix = true
	}
}

// WorkerNodePool represents worker node pool configuration
type WorkerNodePool struct {
	NodePool `yaml:",inline"`
	Location string `yaml:"location,omitempty"`
}

// SetDefaults sets default values for worker pool
func (w *WorkerNodePool) SetDefaults() {
	if w.Location == "" {
		w.Location = "fsn1"
	}
	// Only set instance_count default for non-autoscaling pools
	// Autoscaling pools should have instance_count = 0 by default
	if w.InstanceCount == 0 && !w.AutoscalingEnabled() {
		w.InstanceCount = 1
	}
	if !w.IncludeClusterNameAsPrefix {
		w.IncludeClusterNameAsPrefix = true
	}
}

// BuildNodePoolName builds the node pool name for this worker pool
// This is used by the cluster autoscaler and for finding autoscaled instances
func (w *WorkerNodePool) BuildNodePoolName(clusterName string) string {
	poolName := "default"
	if w.Name != nil {
		poolName = *w.Name
	}

	if w.IncludeClusterNameAsPrefix {
		return clusterName + "-" + poolName
	}
	return poolName
}

// Label represents a Kubernetes label
type Label struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// Taint represents a Kubernetes taint
type Taint struct {
	Key    string `yaml:"key"`
	Value  string `yaml:"value,omitempty"`
	Effect string `yaml:"effect"`
}

// Autoscaling represents autoscaling configuration
type Autoscaling struct {
	Enabled      bool `yaml:"enabled,omitempty"`
	MinInstances int  `yaml:"min_instances,omitempty"`
	MaxInstances int  `yaml:"max_instances,omitempty"`
}
