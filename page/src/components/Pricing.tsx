import { Check, X } from "lucide-react";
import { Badge } from "@/components/ui/badge";
const pricingPlans = [{
  name: "hek3ster + Hetzner Cloud",
  badge: "Best Value",
  price: "€142",
  subtitle: "Infrastructure only - 3 masters + 10 workers",
  features: [{
    text: "hek3ster is 100% free",
    included: true
  }, {
    text: "Pay only for Hetzner servers",
    included: true
  }, {
    text: "No management service fees",
    included: true
  }, {
    text: "Firewall protection included",
    included: true
  }, {
    text: "Nat Gateway and LoadBalancer included",
    included: true
  }],
  highlighted: true
}, {
  name: "AWS Kubernetes",
  badge: "10x More",
  price: "€1,495+",
  subtitle: "Same cluster size, different cloud",
  features: [{
    text: "€73/month control plane fee",
    included: false
  }, {
    text: "~€1,440/month for 10 EC2 instances",
    included: false
  }, {
    text: "~€75/month for Load Balancer",
    included: false
  }, {
    text: "Data transfer and EBS",
    included: false
  }, {
    text: "Other services",
    included: false
  }],
  highlighted: false
}, {
  name: "Managed Service + Hetzner",
  badge: "2x More",
  price: "€375+",
  subtitle: "Same infrastructure + platform fees",
  features: [{
    text: "Per CPU management fees",
    included: false
  }, {
    text: "Control plane fee",
    included: false
  }, {
    text: "Account sign-up required",
    included: false
  }, {
    text: "Shared API token with platform",
    included: false
  }],
  highlighted: false
}];
const clusterPricing = [{
  type: "Language",
  config: "GO 1.25",
  cost: "Modern, efficient, compiled language"
}, {
  type: "Startup Time",
  config: "~10ms",
  cost: "Instant binary startup"
}, {
  type: "Binary Size",
  config: "~12MB",
  cost: "Compact single executable"
}, {
  type: "Build Time",
  config: "~15sec",
  cost: "Fast development iteration"
}, {
  type: "Dependencies",
  config: "Static binary",
  cost: "Zero runtime dependencies"
}, {
  type: "Test Coverage",
  config: "Comprehensive",
  cost: "Unit and integration tests included"
}, {
  type: "Configuration",
  config: "YAML",
  cost: "Full syntax support"
}];
export function Pricing() {
  return <section className="py-20 px-4 bg-muted/30">
      <div className="container max-w-6xl mx-auto">
        {/* Section Header */}
        <div className="text-center mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold mb-4">
            Experience major savings versus top cloud providers
          </h2>
          <p className="text-lg text-muted-foreground">
            hek3ster is free open-source software that deploys to Hetzner Cloud
          </p>
        </div>

        {/* Pricing Cards */}
        <div className="grid lg:grid-cols-3 gap-6 mb-16">
          {pricingPlans.map(plan => <div key={plan.name} className={`relative p-6 rounded-2xl border ${plan.highlighted ? "bg-card border-primary shadow-lg shadow-primary/10" : "bg-card/50 border-border"}`}>
              {plan.highlighted && <div className="absolute -top-3 left-6">
                  <Badge className="bg-gradient-hero border-0 text-primary-foreground shadow-glow">
                    {plan.badge}
                  </Badge>
                </div>}
              {!plan.highlighted && <Badge variant="secondary" className="mb-4">
                  {plan.badge}
                </Badge>}
              
              <h3 className="text-xl font-semibold mt-2 mb-1">{plan.name}</h3>
              <div className="flex items-baseline gap-1 mb-2">
                <span className={`text-3xl font-bold ${plan.highlighted ? "text-gradient" : ""}`}>
                  {plan.price}
                </span>
                <span className="text-muted-foreground">/month</span>
              </div>
              <p className="text-sm text-muted-foreground mb-6">{plan.subtitle}</p>

              <ul className="space-y-3">
                {plan.features.map((feature, idx) => <li key={idx} className="flex items-start gap-3 text-sm">
                    {feature.included ? <Check className="w-4 h-4 text-primary mt-0.5 flex-shrink-0" /> : <X className="w-4 h-4 text-muted-foreground mt-0.5 flex-shrink-0" />}
                    <span className={feature.included ? "text-foreground" : "text-muted-foreground"}>
                      {feature.text}
                    </span>
                  </li>)}
              </ul>
            </div>)}
        </div>

        {/* Pricing Table */}
        <div className="bg-card rounded-2xl border border-border overflow-hidden">
          <div className="p-6 border-b border-border">
            <h3 className="text-xl font-semibold">Project Technical Data</h3>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-border bg-muted/50">
                  
                  
                  
                </tr>
              </thead>
              <tbody>
                {clusterPricing.map((row, idx) => <tr key={row.type} className={idx !== clusterPricing.length - 1 ? "border-b border-border" : ""}>
                    <td className="p-4 font-medium">{row.type}</td>
                    <td className="p-4 text-muted-foreground font-mono text-sm">{row.config}</td>
                    <td className="p-4 text-right font-semibold">{row.cost}</td>
                  </tr>)}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </section>;
}
