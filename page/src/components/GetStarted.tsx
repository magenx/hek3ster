import { Copy } from "lucide-react";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";

const steps = [
  {
    number: 1,
    title: "Install hek3ster",
    description: "Download the binary for your platform",
    code: `# macOS arm64 / Linux (arm64/amd64)
https://github.com/magenx/hek3ster/releases`,
  },
  {
    number: 2,
    title: "Create Configuration",
    description: "Define your cluster in a simple YAML file",
    code: `# cluster.yaml
---
hetzner_token: xxxx
cluster_name: &cluster_name demo
kubeconfig_path: "~/.kube/config"
k3s_version: v1.35.0+k3s1

domain: &domain example.com
location: &location nbg1
image: &image debian-13
autoscaling_image: *image

protect_against_deletion: true
create_load_balancer_for_the_kubernetes_api: true
k3s_upgrade_concurrency: 1
schedule_workloads_on_masters: false

datastore:
  mode: "etcd"
  embedded_etcd:
    snapshot_retention: 24
    snapshot_schedule_cron: "0 * * * *"
    s3_enabled: true
    s3_endpoint: "nbg1.your-objectstorage.com"
    s3_region: *location
    s3_bucket: *cluster_name
    s3_folder: "etcd-snapshot"
    s3_access_key: "xxxx"
    s3_secret_key: "xxxx"

networking:
  ssh:
    port: 22
    use_agent: false
    public_key_path: "~/.ssh/id_ed25519.pub"
    private_key_path: "~/.ssh/id_ed25519"

  allowed_networks:
    ssh:
      - 123.10.45.34/32
      - 72.40.0.0/16
    api:
      - 123.10.45.34/32
      - 72.40.0.0/16
      - 33.79.45.65/32
      
# Public network disabled
  public_network:
    use_local_firewall: false
    ipv4:
      enabled: false
    ipv6:
      enabled: false

# Private network only with NAT gateway
  private_network:
    enabled: true
    subnet: 10.0.0.0/16
    existing_network_name: ""
    nat_gateway:
      enabled: true
      instance_type: "cx23"
      location: *location

  cni:
    enabled: true
    mode: flannel
    encryption: true
    flannel:
      encryption: true
      disable_kube_proxy: false
      
# Global Load Balancer Configuration
# This load balancer will serve public traffic for your applications
load_balancer:
  # Optional: Custom name for the load balancer
  # If not specified, defaults to "{cluster_name}-global-lb"
  # name: *cluster_name
  enabled: true
  target_pools: ["varnish"]
  use_private_ip: true
  attach_to_network: true
  type: "lb11"
  location: *location
  algorithm:
    type: "round_robin"
  services:
    - protocol: "http"
      listen_port: 80
      destination_port: 80
      proxyprotocol: false
      health_check:
        protocol: "http"
        port: 80
        interval: 15
        timeout: 10
        retries: 3
        http:
          domain: *domain
          path: "/health_check.php"
          status_codes: ["2??", "3??"]
          tls: false

masters_pool:
  instance_type: cpx22
  instance_count: 2
  locations:
    - *location

worker_node_pools:
- name: varnish
  instance_type: cpx22
  instance_count: 1
  location: *location
 
- name: nginx
  instance_type: cpx22
  instance_count: 1
  location: *location

- name: php
  instance_type: cpx22
  location: *location
  autoscaling:
    enabled: true
    min_instances: 1
    max_instances: 3

- name: valkey
  instance_type: cpx22
  instance_count: 1
  location: *location

- name: rabbitmq
  instance_type: cpx22
  instance_count: 1
  location: *location
  
- name: opensearch
  instance_type: cpx22
  instance_count: 1
  location: *location

- name: mariadb
  instance_type: cpx22
  instance_count: 1
  location: *location

addons:
  metrics_server:
    enabled: true
  csi_driver:
    enabled: true
    manifest_url: "https://raw.githubusercontent.com/hetznercloud/csi-driver/
    v2.18.3/deploy/kubernetes/hcloud-csi.yml"
  cluster_autoscaler:
    enabled: true
    manifest_url: "https://raw.githubusercontent.com/kubernetes/autoscaler/master/
    cluster-autoscaler/cloudprovider/hetzner/examples/cluster-autoscaler-run-on-master.yaml"
    container_image_tag: "v1.34.2"
    scan_interval: "10s"                        
    scale_down_delay_after_add: "10m"
    scale_down_delay_after_delete: "10s"
    scale_down_delay_after_failure: "3m"
    max_node_provision_time: "5m"
  cloud_controller_manager:
    enabled: true
    manifest_url: "https://github.com/hetznercloud/
    hcloud-cloud-controller-manager/releases/download/v1.28.0/ccm-networks.yaml"
  system_upgrade_controller:
    enabled: true
    deployment_manifest_url: "https://github.com/rancher/
    system-upgrade-controller/releases/download/v0.18.0/system-upgrade-controller.yaml"
    crd_manifest_url: "https://github.com/rancher/
    system-upgrade-controller/releases/download/v0.18.0/crd.yaml"
  embedded_registry_mirror:
    enabled: true 

additional_packages:
  - ufw

additional_pre_k3s_commands:
  - apt autoremove -y hc-utils
  - apt purge -y hc-utils
  - echo "auto enp7s0" > /etc/network/interfaces
  - echo "iface enp7s0 inet dhcp" >> /etc/network/interfaces
  - echo "    post-up ip route add default via 10.0.0.1"  >> /etc/network/interfaces
  - echo "[Resolve]" > /etc/systemd/resolved.conf
  - echo "DNS=185.12.64.2 185.12.64.1" >> /etc/systemd/resolved.conf
  - ifdown enp7s0 2>/dev/null || true
  - ifup enp7s0 2>/dev/null || true
  - sleep 2
  - apt update
  - apt install -y resolvconf syslog-ng
  - systemctl enable --now resolvconf
  - echo "nameserver 185.12.64.2" >> /etc/resolvconf/resolv.conf.d/head
  - echo "nameserver 185.12.64.1" >> /etc/resolvconf/resolv.conf.d/head
  - resolvconf --enable-updates
  - resolvconf -u

additional_post_k3s_commands:
  - apt autoremove -y`,
  },
  {
    number: 3,
    title: "Create Your Cluster",
    description: "One command to deploy everything",
    code: `# Generate SSH key pair
ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519 -C "your@example.com"

# Create Your Cluster
hek3ster create --config cluster.yaml

# That's it! Your cluster is ready.
kubectl get nodes`,
  },
];

function CodeBlock({ code }: { code: string }) {
  const copyToClipboard = () => {
    navigator.clipboard.writeText(code);
    toast.success("Copied to clipboard!");
  };

  return (
    <div className="relative group">
      <pre className="bg-terminal-bg rounded-lg p-4 overflow-x-auto font-mono text-sm text-terminal-text">
        <code>{code}</code>
      </pre>
      <Button
        variant="ghost"
        size="icon"
        className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity h-8 w-8 bg-terminal-header hover:bg-terminal-comment/30"
        onClick={copyToClipboard}
      >
        <Copy className="h-4 w-4 text-terminal-text" />
      </Button>
    </div>
  );
}

export function GetStarted() {
  return (
    <section className="py-20 px-4">
      <div className="container max-w-4xl mx-auto">
        {/* Section Header */}
        <div className="text-center mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold mb-4">
            Get Started in <span className="text-gradient">3 Steps</span>
          </h2>
          <p className="text-lg text-muted-foreground">
            From zero to production-ready cluster in minutes
          </p>
        </div>

        {/* Steps */}
        <div className="space-y-12">
          {steps.map((step) => (
            <div key={step.number} className="relative">
              {/* Step connector line */}
              {step.number < steps.length && (
                <div className="absolute left-6 top-16 w-0.5 h-[calc(100%-2rem)] bg-border hidden md:block" />
              )}

              <div className="flex gap-6">
                {/* Step number */}
                <div className="w-12 h-12 rounded-full bg-gradient-hero flex items-center justify-center flex-shrink-0 shadow-glow">
                  <span className="text-lg font-bold text-primary-foreground">{step.number}</span>
                </div>

                {/* Step content */}
                <div className="flex-1">
                  <h3 className="text-xl font-semibold mb-1">{step.title}</h3>
                  <p className="text-muted-foreground mb-4">{step.description}</p>
                  <CodeBlock code={step.code} />
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
