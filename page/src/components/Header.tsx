
export function Header() {
  return (
    <header className="fixed top-0 left-0 right-0 z-50 bg-background/80 backdrop-blur-md border-b border-border">
      <div className="container max-w-6xl mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <a href="/" className="flex items-center gap-2">
            <div className="w-24 h-8 rounded-lg bg-gradient-hero flex items-center justify-center">
              <span className="text-sm font-bold text-primary-foreground">hek3ster</span>
            </div>
          </a>

          {/* Navigation */}
          <nav className="hidden md:flex items-right gap-8">
            <a href="#architecture" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
              Architecture
            </a>
            <a href="#features" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
              Features
            </a>
            <a href="#comparison" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
              Compare
            </a>
            <a href="#pricing" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
              Savings
            </a>
            <a href="#docs" className="text-sm text-muted-foreground hover:text-foreground transition-colors">
              Deploy
            </a>
          </nav>

          {/* Actions */}
        </div>
      </div>
    </header>
  );
}
