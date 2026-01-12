import { Check, X } from "lucide-react";

const comparisonData = [
  {
    factor: "Setup time",
    hek3ster: { value: "~5 minutes", status: "good" },
    hetznerK3s: { value: "~5 minutes", status: "good" },
    managed: { value: "30+ minutes", status: "bad" },
    terraform: { value: "30+ minutes", status: "bad" },
  },
  {
    factor: "Dependencies",
    hek3ster: { value: "CLI tool only", status: "good" },
    hetznerK3s: { value: "helm, kubectl, homebrew", status: "bad" },
    managed: { value: "Third-party account*", status: "bad" },
    terraform: { value: "Terraform, Packer, helm, kubectl ", status: "bad" },
  },
  {
    factor: "Data privacy",
    hek3ster: { value: "Full control", status: "good" },
    hetznerK3s: { value: "Full control", status: "good" },
    managed: { value: "Third-party access", status: "bad" },
    terraform: { value: "Platform dependent", status: "bad" },
  },
  {
    factor: "Monthly cost",
    hek3ster: { value: "Infrastructure", status: "good" },
    hetznerK3s: { value: "Infrastructure", status: "good" },
    managed: { value: "Infra + platform fees", status: "bad" },
    terraform: { value: "Infra + platform fees", status: "bad" },
  },
  {
    factor: "Credential exposure",
    hek3ster: { value: "None", status: "good" },
    hetznerK3s: { value: "None", status: "good" },
    managed: { value: "API tokens", status: "bad" },
    terraform: { value: "Setup dependent", status: "bad" },
  },
  {
    factor: "Learning curve",
    hek3ster: { value: "Low", status: "good" },
    hetznerK3s: { value: "Low", status: "good" },
    managed: { value: "Medium", status: "bad" },
    terraform: { value: "Medium-High", status: "bad" },
  },
  {
    factor: "Secure by default",
    hek3ster: { value: "Yes", status: "good" },
    hetznerK3s: { value: "Manual setup", status: "bad" },
    managed: { value: "Varies", status: "bad" },
    terraform: { value: "Manual setup", status: "bad" },
  },
  {
    factor: "Private network",
    hek3ster: { value: "Yes", status: "good" },
    hetznerK3s: { value: "Manual config", status: "bad" },
    managed: { value: "Manual config", status: "bad" },
    terraform: { value: "Manual config", status: "bad" },
  },
  {
    factor: "Public IP exposure",
    hek3ster: { value: "No public IPs", status: "good" },
    hetznerK3s: { value: "Public by default", status: "bad" },
    managed: { value: "Public by default", status: "bad" },
    terraform: { value: "Public by default", status: "bad" },
  },
  {
    factor: "NAT gateway",
    hek3ster: { value: "Included", status: "good" },
    hetznerK3s: { value: "Manual setup", status: "bad" },
    managed: { value: "Manual setup", status: "bad" },
    terraform: { value: "Manual setup", status: "bad" },
  },
  {
    factor: "Global load balancer",
    hek3ster: { value: "Included", status: "good" },
    hetznerK3s: { value: "Manual setup", status: "bad" },
    managed: { value: "Manual setup", status: "bad" },
    terraform: { value: "Manual setup", status: "bad" },
  },
  {
    factor: "Configuration",
    hek3ster: { value: "AI", status: "good" },
    hetznerK3s: { value: "Manual setup + YAML", status: "bad" },
    managed: { value: "Web UI / API", status: "bad" },
    terraform: { value: "HCL files", status: "bad" },
  },
];

function StatusIcon({ status }: { status: string }) {
  if (status === "good") {
    return <Check className="w-4 h-4 text-primary" />;
  }
  if (status === "bad") {
    return <X className="w-4 h-4 text-destructive" />;
  }
}

export function Comparison() {
  return (
    <section className="py-20 px-4">
      <div className="container max-w-6xl mx-auto">
        {/* Section Header */}
        <div className="text-center mb-12">
          <h2 className="text-3xl sm:text-4xl font-bold mb-4">
            How <span className="text-gradient">hek3ster</span> Compares
          </h2>
          <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
            See how hek3ster stacks up against managed services and Terraform-based solutions
          </p>
        </div>

        {/* Comparison Table */}
        <div className="bg-card rounded-2xl border border-border overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-border bg-muted/50">
                  <th className="text-left p-4 text-sm font-semibold">Factor</th>
                  <th className="text-left p-4 text-sm font-semibold text-primary">hek3ster</th>
                  <th className="text-left p-4 text-sm font-semibold">hetzner-k3s</th>
                  <th className="text-left p-4 text-sm font-semibold">Managed Services</th>
                  <th className="text-left p-4 text-sm font-semibold">Terraform-based</th>
                </tr>
              </thead>
              <tbody>
                {comparisonData.map((row, idx) => (
                  <tr
                    key={row.factor}
                    className={idx !== comparisonData.length - 1 ? "border-b border-border" : ""}
                  >
                    <td className="p-4 font-medium">{row.factor}</td>
                    <td className="p-4">
                      <div className="flex items-center gap-2">
                        <StatusIcon status={row.hek3ster.status} />
                        <span className="text-sm">{row.hek3ster.value}</span>
                      </div>
                    </td>
                    <td className="p-4">
                      <div className="flex items-center gap-2">
                        <StatusIcon status={row.hetznerK3s.status} />
                        <span className="text-sm text-muted-foreground">{row.hetznerK3s.value}</span>
                      </div>
                    </td>
                    <td className="p-4">
                      <div className="flex items-center gap-2">
                        <StatusIcon status={row.managed.status} />
                        <span className="text-sm text-muted-foreground">{row.managed.value}</span>
                      </div>
                    </td>
                    <td className="p-4">
                      <div className="flex items-center gap-2">
                        <StatusIcon status={row.terraform.status} />
                        <span className="text-sm text-muted-foreground">{row.terraform.value}</span>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <div className="p-4 border-t border-border bg-muted/30">
            <p className="text-xs text-muted-foreground">
              *Managed services require signing up for their platform in addition to Hetzner Cloud.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
