import { Button } from "@/components/ui/button";
import { Github, ArrowRight, BrainCircuit, ShieldCheck, Euro } from "lucide-react";

const stats = [
  { icon: BrainCircuit, value: "AI Manager", label: "Intelligent Interface" },
  { icon: ShieldCheck, value: "Security First", label: "Air-gapped Environment" },
  { icon: Euro, value: "85% Savings", label: "Compared to Others" },
];

export function Hero() {
  return (
    <section className="relative pt-24 pb-16 px-4 overflow-hidden">
      {/* Background decoration */}
      <div className="absolute inset-0 -z-10">
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[600px] bg-primary/5 rounded-full blur-3xl" />
        <div className="absolute top-40 right-0 w-[400px] h-[400px] bg-primary/3 rounded-full blur-3xl" />
      </div>

      <div className="container max-w-6xl mx-auto">
        {/* Header */}
        <div className="text-center mb-12">
          <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-accent text-accent-foreground text-sm font-medium mb-6">
            <span className="w-2 h-2 rounded-full bg-primary animate-pulse" />
            Open Source CLI Tool
          </div>
          
          <h1 className="text-4xl sm:text-5xl lg:text-6xl font-extrabold tracking-tight mb-6">
            <span className="text-gradient">Deploy secure k3s cluster on</span>
            <br />
            <span className="text-foreground">Hetzner Cloud</span>
          </h1>
          
          <p className="text-lg sm:text-xl text-muted-foreground max-w-2xl mx-auto mb-8">
            Production-ready private Kubernetes clusters created with a single command. 
            No programming required, no complexity.
          </p>
          <br/>
          <br/>
          {/* Stats */}
          <div className="flex flex-wrap justify-center gap-8 mb-10">
            {stats.map(({ icon: Icon, value, label }) => (
              <div key={label} className="flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-accent flex items-center justify-center">
                  <Icon className="w-5 h-5 text-primary" />
                </div>
                <div className="text-left">
                  <div className="text-xl font-bold text-foreground">{value}</div>
                  <div className="text-sm text-muted-foreground">{label}</div>
                </div>
              </div>
            ))}
          </div>

          {/* CTAs */}
          <div className="flex flex-wrap justify-center gap-4">
            <Button size="lg" className="gap-2 shadow-glow" asChild>
              <a href="https://github.com/magenx/hek3ster/wiki">
                Read the Docs
                <ExternalLink className="w-4 h-4" />
              </a>
            </Button>
            <Button variant="outline" size="lg" className="gap-2" asChild>
              <a href="https://github.com/magenx/hek3ster" target="_blank" rel="noopener noreferrer">
                <Github className="w-4 h-4" />
                View on GitHub
              </a>
            </Button>
          </div>
        </div>
      </div>
    </section>
  );
}
