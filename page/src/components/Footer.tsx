import { Button } from "@/components/ui/button";
import { Github, Heart, MessageCircle, ExternalLink } from "lucide-react";

export function Footer() {
  return (
    <footer className="py-16 px-4 border-t border-border bg-muted/20">
      <div className="container max-w-6xl mx-auto">
        {/* CTA Section */}
        <div className="text-center mb-12">
          <h2 className="text-2xl sm:text-3xl font-bold mb-4">
            Ready to Get Started?
          </h2>
          <p className="text-muted-foreground mb-6 max-w-xl mx-auto">
            Join thousands of developers running production Kubernetes on Hetzner Cloud
          </p>
          <div className="flex flex-wrap justify-center gap-4">
            <Button size="lg" className="gap-2 shadow-glow" asChild>
			  <a href="https://deepwiki.com/magenx/hek3ster">
                Read the Docs
                <ExternalLink className="w-4 h-4" />
              </a>
            </Button>
            <Button variant="outline" size="lg" className="gap-2" asChild>
			   <a href="https://github.com/magenx/hek3ster/discussions" target="_blank" rel="noopener noreferrer">
                <MessageCircle className="w-4 h-4" />
                Join Discussions
              </a>
            </Button>
          </div>
        </div>

        {/* Divider */}
        <div className="border-t border-border my-8" />

        {/* Bottom Section */}
        <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2">
            <div className="w-24 h-8 rounded-lg bg-gradient-hero flex items-center justify-center">
              <span className="text-sm font-bold text-primary-foreground">hek3ster</span>
            </div>
          </div>

          <p className="text-sm text-muted-foreground flex items-center gap-1">
            Made with <Heart className="w-4 h-4 text-red-500 fill-red-500" /> by the community
          </p>

          <div className="flex items-center gap-4">
            <a
              href="https://github.com"
              target="_blank"
              rel="noopener noreferrer"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              <Github className="w-5 h-5" />
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}
