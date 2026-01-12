import { Server, Database, Shield, Globe, Network, Cpu, Box, Router } from "lucide-react";

export function Architecture() {
  return (
    <section className="py-16 bg-muted/20">
      <div className="container max-w-6xl mx-auto px-4">
        <br/>
        <h2 className="text-2xl md:text-3xl font-bold text-center mb-4">
          Secure Kubernetes Architecture
        </h2>
        <p className="text-muted-foreground text-center mb-12 max-w-2xl mx-auto">
          High-availability K3s cluster with private networking and NAT gateway for secure outbound traffic
        </p>
        <div className="relative overflow-x-auto">
          <div className="min-w-[900px] flex items-stretch justify-center gap-4 p-4">
            
            {/* External Traffic / Load Balancer */}
            <div className="flex flex-col items-center justify-center">
              <div className="bg-primary/10 border-2 border-primary rounded-xl p-4 text-center">
                <Globe className="w-8 h-8 mx-auto mb-2 text-primary" />
                <span className="text-sm font-semibold">External<br/>Traffic</span>
              </div>
              <div className="h-full flex items-center">
                <div className="w-8 border-t-2 border-dashed border-primary" />
              </div>
            </div>
            
            {/* Load Balancers Column */}
            <div className="flex flex-col items-center justify-center gap-3">
              {/* Global Load Balancer */}
              <div className="flex items-center">
                <div className="bg-accent/30 border-2 border-accent-foreground/30 rounded-xl p-4 text-center min-w-[100px]">
                  <Network className="w-8 h-8 mx-auto mb-2 text-accent-foreground" />
                  <span className="text-sm font-semibold">Global LB</span>
                  <div className="text-[10px] text-muted-foreground mt-1">
                    (HTTPS)
                  </div>
                </div>
                <div className="w-6 border-t-2 border-dashed border-muted-foreground" />
              </div>

              {/* API Load Balancer - Separate below */}
              <div className="flex items-center">
                <div className="bg-purple-100 dark:bg-purple-900/30 border-2 border-purple-500 rounded-xl p-4 text-center min-w-[100px]">
                  <Network className="w-8 h-8 mx-auto mb-2 text-purple-600 dark:text-purple-400" />
                  <span className="text-sm font-semibold text-purple-800 dark:text-purple-200">API LB</span>
                  <div className="text-[10px] text-purple-600 dark:text-purple-400 mt-1">
                    (6443/TCP)
                  </div>
                </div>
                <div className="w-6 border-t-2 border-dashed border-purple-500" />
              </div>
            </div>

            {/* Private Network Box */}
            <div className="flex-1 border-2 border-dashed border-primary/50 rounded-2xl bg-card/50 p-4 relative">
              <div className="absolute -top-3 left-4 bg-background px-2">
                <span className="text-xs font-semibold text-primary flex items-center gap-1">
                  <Shield className="w-3 h-3" /> Private Network (No Public IPs)
                </span>
              </div>
              
              <div className="flex gap-4 h-full">
                {/* Server Nodes (Control Plane) */}
                <div className="flex-1">
                  <div className="text-xs font-semibold text-muted-foreground mb-2 text-center">Control Plane (HA)</div>
                  <div className="grid grid-cols-1 gap-2">
                    {[1, 2, 3].map((num) => (
                      <div 
                        key={num}
                        className="bg-amber-100 dark:bg-amber-900/30 border border-amber-300 dark:border-amber-700 rounded-lg p-3"
                      >
                        <div className="flex items-center gap-2 mb-2">
                          <Server className="w-4 h-4 text-amber-600 dark:text-amber-400" />
                          <span className="text-xs font-semibold text-amber-800 dark:text-amber-200">
                            Server Node {num}
                          </span>
                        </div>
                        <div className="grid grid-cols-2 gap-1">
                          <div className="bg-amber-50 dark:bg-amber-950/50 rounded px-2 py-1 text-[10px] text-center">
                            <Cpu className="w-3 h-3 mx-auto mb-0.5 text-amber-600" />
                            K3s Server
                          </div>
                          <div className="bg-amber-50 dark:bg-amber-950/50 rounded px-2 py-1 text-[10px] text-center">
                            <Database className="w-3 h-3 mx-auto mb-0.5 text-amber-600" />
                            etcd
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>

                {/* Connection Lines */}
                <div className="flex flex-col items-center justify-center gap-2">
                  <div className="h-full border-l-2 border-dotted border-muted-foreground/50" />
                  <span className="text-[10px] text-muted-foreground rotate-90 whitespace-nowrap">kubectl</span>
                  <div className="h-full border-l-2 border-dotted border-muted-foreground/50" />
                </div>

                {/* Agent Nodes */}
                <div className="flex-1">
                  <div className="text-xs font-semibold text-muted-foreground mb-2 text-center">Worker Nodes</div>
                  <div className="grid grid-cols-1 gap-2">
                    {[1, 2, 3].map((num) => (
                      <div 
                        key={num}
                        className="bg-blue-100 dark:bg-blue-900/30 border border-blue-300 dark:border-blue-700 rounded-lg p-3"
                      >
                        <div className="flex items-center gap-2 mb-2">
                          <Box className="w-4 h-4 text-blue-600 dark:text-blue-400" />
                          <span className="text-xs font-semibold text-blue-800 dark:text-blue-200">
                            Agent Node {num}
                          </span>
                        </div>
                        <div className="bg-blue-50 dark:bg-blue-950/50 rounded px-2 py-1 text-[10px] text-center">
                          <Cpu className="w-3 h-3 mx-auto mb-0.5 text-blue-600" />
                          K3s Agent + Workloads
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>

            {/* NAT Gateway */}
            <div className="flex flex-col items-center justify-center">
              <div className="h-full flex items-center">
                <div className="w-8 border-t-2 border-dashed border-muted-foreground" />
              </div>
              <div className="bg-green-100 dark:bg-green-900/30 border-2 border-green-500 rounded-xl p-4 text-center min-w-[120px]">
                <Router className="w-8 h-8 mx-auto mb-2 text-green-600 dark:text-green-400" />
                <span className="text-sm font-semibold text-green-800 dark:text-green-200">NAT<br/>Gateway</span>
                <div className="text-[10px] text-green-600 dark:text-green-400 mt-1">
                  Outbound Only
                </div>
              </div>
              <div className="h-full flex items-center">
                <div className="w-8 border-t-2 border-dashed border-muted-foreground" />
              </div>
            </div>

            {/* Internet */}
            <div className="flex flex-col items-center justify-center">
              <div className="bg-muted border-2 border-border rounded-xl p-4 text-center">
                <Globe className="w-8 h-8 mx-auto mb-2 text-muted-foreground" />
                <span className="text-sm font-semibold text-muted-foreground">Internet</span>
                <div className="text-[10px] text-muted-foreground mt-1">
                  (APIs, Registries)
                </div>
              </div>
            </div>
          </div>

          {/* Legend */}
          <div className="flex flex-wrap justify-center gap-4 mt-6 text-xs text-muted-foreground">
            <div className="flex items-center gap-1">
              <div className="w-3 h-3 rounded bg-amber-200 dark:bg-amber-800 border border-amber-400" />
              <span>Control Plane</span>
            </div>
            <div className="flex items-center gap-1">
              <div className="w-3 h-3 rounded bg-blue-200 dark:bg-blue-800 border border-blue-400" />
              <span>Worker Nodes</span>
            </div>
            <div className="flex items-center gap-1">
              <div className="w-3 h-3 rounded bg-purple-200 dark:bg-purple-800 border border-purple-500" />
              <span>API Load Balancer</span>
            </div>
            <div className="flex items-center gap-1">
              <div className="w-3 h-3 rounded bg-green-200 dark:bg-green-800 border border-green-500" />
              <span>NAT Gateway</span>
            </div>
            <div className="flex items-center gap-1">
              <div className="w-4 border-t-2 border-dashed border-primary" />
              <span>Inbound Traffic</span>
            </div>
            <div className="flex items-center gap-1">
              <div className="w-4 border-t-2 border-dotted border-muted-foreground" />
              <span>Internal / Outbound</span>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
