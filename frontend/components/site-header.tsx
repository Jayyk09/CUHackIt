'use client'

import { motion } from 'framer-motion'
import { useUser } from '@/contexts/user-context'

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080'

export function SiteHeader() {
  const { user, isLoading, clearUser } = useUser()

  function handleSignOut() {
    clearUser()
    window.location.href = `${API_URL}/logout`
  }

  return (
    <motion.header
      className="fixed top-0 left-0 right-0 z-50 border-b border-foreground/10 bg-background/80 backdrop-blur-md"
      initial={{ opacity: 0, y: -10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5, ease: [0.22, 1, 0.36, 1] }}
    >
      <div className="flex items-center justify-between px-6 md:px-10 lg:px-16 h-14">

        {/* Brand */}
        <span className="font-serif text-xl tracking-[-0.02em] text-foreground">
          Sift
        </span>

        {/* Center nav */}
        <nav className="hidden md:flex items-center gap-10">
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

        {/* Right: Est. label + auth action */}
        <div className="flex items-center gap-6">
          <span className="hidden lg:block font-sans text-[10px] tracking-[0.2em] uppercase text-muted-foreground/40 select-none">
            Est. 2026
          </span>

          {!isLoading && (
            user ? (
              <div className="flex items-center gap-6">
                <button
                  onClick={handleSignOut}
                  className="font-sans text-[10px] tracking-[0.25em] uppercase text-muted-foreground hover:text-foreground transition-colors duration-300 cursor-pointer"
                >
                  Sign Out
                </button>
                <a
                  href="/dashboard"
                  className="font-sans text-[10px] tracking-[0.25em] uppercase px-4 py-2 border border-foreground/20 text-foreground/70 hover:border-foreground/60 hover:text-foreground transition-all duration-300"
                >
                  Home
                </a>
              </div>
            ) : (
              <a
                href={`${API_URL}/login`}
                className="font-sans text-[10px] tracking-[0.25em] uppercase px-4 py-2 border border-foreground/20 text-foreground/70 hover:border-foreground/60 hover:text-foreground transition-all duration-300"
              >
                Sign In / Up
              </a>
            )
          )}
        </div>

      </div>
    </motion.header>
  )
}
