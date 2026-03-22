import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { cn } from "@/lib/utils";
import Footer from "@/components/footer";
import Navbar from "@/components/navbar";
import Providers from "@/components/providers/providers";
import { Toaster } from "@/components/ui/sonner";
import { NuqsAdapter } from "nuqs/adapters/next/app";

const fontSans = Geist({
  subsets: ["latin"],
  variable: "--font-sans",
});

const fontMono = Geist_Mono({
  subsets: ["latin"],
  variable: "--font-mono",
});

export const metadata: Metadata = {
  title: "terraforge",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="en"
      suppressHydrationWarning
      className={cn(
        "antialiased",
        fontMono.variable,
        "font-sans",
        fontSans.variable,
      )}
    >
      <body>
        <NuqsAdapter>
          <Providers>
            <Navbar />
            <main className="container mx-auto min-h-screen max-w-6xl px-4">
              <div className="py-12">{children}</div>
            </main>
            <Toaster richColors />
            <Footer />
          </Providers>
        </NuqsAdapter>
      </body>
    </html>
  );
}
