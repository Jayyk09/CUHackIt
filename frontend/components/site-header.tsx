'use client'

import { motion } from 'framer-motion'

export function SiteHeader() {
  return (
    <motion.header
      className="fixed top-0 left-0 right-0 z-50 border-b border-foreground/15 bg-background/80 backdrop-blur-md"
      initial={{ opacity: 0, y: -10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5, ease: [0.22, 1, 0.36, 1] }}
    >
      <div className="flex items-center justify-between px-4 md:px-8 lg:px-16 h-12 md:h-14">
        <span className="font-serif text-lg md:text-xl tracking-[-0.02em] text-foreground">
          Sift
        </span>
        <nav className="hidden md:flex items-center gap-8">
          <a
            href="#facts"
            className="font-sans text-[10px] tracking-[0.25em] uppercase text-muted-foreground hover:text-foreground transition-colors duration-300"
          >
            The Facts
          </a>
          <a
            href="#pyramid"
            className="font-sans text-[10px] tracking-[0.25em] uppercase text-muted-foreground hover:text-foreground transition-colors duration-300"
          >
            Food Pyramid
          </a>
        </nav>
        <span className="font-sans text-[10px] tracking-[0.2em] uppercase text-muted-foreground">
          Est. 2026
        </span>
      </div>
    </motion.header>
  )
}
