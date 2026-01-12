import { Blocks, Users, Scale, Shield } from "lucide-react";

const trustItems = [
  { icon: Blocks, label: "An independent open source project, not affiliated with Hetzner Online GmbH." },
];

export function TrustBar() {
  return (
    <section className="py-6 border-y border-border bg-muted/30">
      <div className="container max-w-6xl mx-auto px-4">
        <div className="flex flex-wrap justify-center gap-2 md:gap-4">
          {trustItems.map(({ icon: Icon, label }) => (
            <div
              key={label}
              className="flex items-center gap-2 text-sm font-medium text-muted-foreground"
            >
              <Icon className="w-4 h-4" />
              <span>{label}</span>
            </div>
          ))}
        </div>
      </div>
</section>
  );
}
