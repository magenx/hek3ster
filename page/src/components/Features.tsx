import { Zap, FileCode, Activity, GitPullRequestArrow, Fullscreen, HeartHandshake, Server, Cuboid, BrainCircuit, Waypoints, Lock, Combine } from "lucide-react";

const features = [
  {
    icon: Zap,
    title: "Unmatched Velocity",
    description: "K3s HA cluster management tool written in Go, providing creation, management, and operations.",
    highlights: ["Hetzner Cloud Go SDK v2", "Binary starts in 10ms", "Full rebuild in 15 seconds", "Zero runtime dependencies", "Unit and integration tests"],
  },
  {
    icon: FileCode,
    title: "Uncomplicated Excellence",
    description: "Complete tool, one YAML config file. No programming or Kubernetes skills are required. All infra pre-installed.",
    highlights: ["AI Cluster Manager", "NAT Gateway for access", "Global and API load balancers", "S3 backups for etcd", "Automated upgrade plan generation"],
  },
  {
    icon: Lock,
    title: "Uncompromised Security",
    description: "Core architectural principle - security first. It is essential foundation built into every layer of this tool.",
    highlights: ["Private network by default", "Automated firewall configuration", "Tokens and ssh keys on your PC", "SSL/TLS encrypting data sent between", "Interconnected and hermetically sealed"],
  },
];

const capabilities = [
  { icon: Server, title: "kubectl", desc: "Command line tool for communicating with a Kubernetes cluster's control plane." },
  { icon: Combine, title: "helm", desc: "Helm helps you define, install, and upgrade even the most complex Kubernetes applications." },
  { icon: BrainCircuit, title: "kubectl-ai", desc: "Translating user intent into precise Kubernetes operations, making management more accessible and efficient." },
  { icon: Fullscreen, title: "k3s", desc: "Lightweight, certified Kubernetes distribution by Rancher. Lower resource footprint, single binary, production-ready." },
  { icon: HeartHandshake, title: "Hetzner Cloud Controller Manager", desc: "Automatic load balancer provisioning and node lifecycle management integrated with Hetzner Cloud." },
  { icon: Cuboid, title: "Hetzner CSI Driver", desc: "Dynamic volume provisioning for Hetzner Cloud volumes. Create PVCs and get automatically provisioned storage." },
  { icon: Activity, title: "Cluster Autoscaler", desc: "Automatically adjust the number of nodes based on pending pods and resource utilization." },
  { icon: GitPullRequestArrow, title: "System Upgrade Controller", desc: "Declarative upgrades for k3s. Define the target version and let the controller handle rolling updates." },
  { icon: Waypoints, title: "Flannel", desc: "Choose your CNI: Flannel for simplicity or Cilium for advanced features like eBPF and network policies." },
];

export function Features() {
  return (
    <section className="py-20 px-4">
      <div className="container max-w-6xl mx-auto">
        {/* Section Header */}
        <div className="text-center mb-16">
        </div>

        {/* Main Features */}
        <div className="grid md:grid-cols-3 gap-8 mb-20">
          {features.map(({ icon: Icon, title, description, highlights }) => (
            <div
              key={title}
              className="group relative p-6 rounded-2xl bg-card border border-border hover:border-primary/30 hover:shadow-lg transition-all duration-300"
            >
              <div className="w-12 h-12 rounded-xl bg-accent flex items-center justify-center mb-4 group-hover:bg-primary/10 transition-colors">
                <Icon className="w-6 h-6 text-primary" />
              </div>
              <h3 className="text-xl font-semibold mb-2">{title}</h3>
              <p className="text-muted-foreground mb-4">{description}</p>
              <ul className="space-y-2">
                {highlights.map((item) => (
                  <li key={item} className="flex items-center gap-2 text-sm text-muted-foreground">
                    <span className="w-1.5 h-1.5 rounded-full bg-primary" />
                    {item}
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        {/* Capabilities Grid */}
        <div className="text-center mb-12">
          <h2 className="text-3xl sm:text-4xl font-bold mb-4">
            What tools are installed and configured
          </h2>
          <p className="text-lg text-muted-foreground">
            Everything you need for production.
          </p>
        </div>

        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {capabilities.map(({ icon: Icon, title, desc }) => (
            <div
              key={title}
              className="flex gap-4 p-5 rounded-xl bg-muted/50 border border-border hover:bg-muted transition-colors"
            >
              <div className="w-10 h-10 rounded-lg bg-accent flex items-center justify-center flex-shrink-0">
                <Icon className="w-5 h-5 text-primary" />
              </div>
              <div>
                <h4 className="font-semibold mb-1">{title}</h4>
                <p className="text-sm text-muted-foreground">{desc}</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
