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
    use_local_firewall: true
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
    mode: cilium
    cilium:
      enabled: true                               
      version: v1.18.6 
      encryption_type: wireguard

  # DNS Zone Management (Required for SSL certificate)
  dns_zone:
    enabled: true
    name: *domain
    ttl: 3600

  # SSL Certificate (Requires DNS zone)
  ssl_certificate:
    enabled: true
    name: *domain
    domain: *domain
...
...
`,
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
