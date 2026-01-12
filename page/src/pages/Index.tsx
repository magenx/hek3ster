import { Header } from "@/components/Header";
import { Hero } from "@/components/Hero";
import { TrustBar } from "@/components/TrustBar";
import { Architecture } from "@/components/Architecture";
import { Features } from "@/components/Features";
import { Comparison } from "@/components/Comparison";
import { Pricing } from "@/components/Pricing";
import { GetStarted } from "@/components/GetStarted";
import { Footer } from "@/components/Footer";

const Index = () => {
  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main>
        <Hero />
        <TrustBar />
        <section id="architecture">
          <Architecture />
        </section>
        <section id="features">
          <Features />
        </section>
        <section id="comparison">
          <Comparison />
        </section>
        <section id="pricing">
          <Pricing />
        </section>
        <section id="docs">
          <GetStarted />
        </section>
      </main>
      <Footer />
    </div>
  );
};

export default Index;
