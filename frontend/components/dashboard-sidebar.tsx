'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'

const navItems = [
  { href: '/dashboard', label: 'Dashboard' },
  { href: '/dashboard/pantry', label: 'Pantry' },
]

export function DashboardSidebar() {
  const pathname = usePathname()

  return (
    <aside className="fixed left-0 top-0 h-screen w-72 border-r border-sage/30 bg-sage-dark">
      <div className="flex h-full flex-col px-8 pt-16">
        <Link 
          href="/dashboard" 
          className="font-serif text-4xl text-cream tracking-tight mb-16"
        >
          Sift
        </Link>

        <nav className="space-y-6">
          {navItems.map((item) => {
            const isActive = pathname === item.href
            return (
              <Link
                key={item.href}
                href={item.href}
                className={`block text-xs uppercase tracking-[0.2em] font-sans transition-colors duration-200 ${
                  isActive
                    ? 'text-cream'
                    : 'text-cream/50 hover:text-cream/70'
                }`}
              >
                {isActive ? '— ' : ''}{item.label}
              </Link>
            )
          })}
        </nav>

        <div className="mt-auto pb-8">
          <Link
            href="/"
            className="block text-xs uppercase tracking-[0.2em] font-sans text-cream/50 hover:text-cream/70 transition-colors duration-200"
          >
            Back to Home
          </Link>
        </div>
      </div>
    </aside>
  )
}
