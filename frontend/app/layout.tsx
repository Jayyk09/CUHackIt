import type { Metadata } from 'next'
import { Suspense } from 'react'
import { Playfair_Display, Space_Grotesk } from 'next/font/google'
import { Analytics } from '@vercel/analytics/next'
import { UserProvider } from '@/contexts/user-context'
import './globals.css'

const playfair = Playfair_Display({
  subsets: ['latin'],
  variable: '--font-playfair',
  display: 'swap',
})

const spaceGrotesk = Space_Grotesk({
  subsets: ['latin'],
  variable: '--font-space-grotesk',
  display: 'swap',
})

export const metadata: Metadata = {
  title: 'Sift — Know Your Food',
  description: 'Sift through the noise. Real food data, presented without compromise.',
  icons: {
    icon: '/leaf.png',
    apple: '/leaf.png',
  },
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" className={`${playfair.variable} ${spaceGrotesk.variable}`}>
      <body className="font-sans antialiased" suppressHydrationWarning>
        <Suspense>
          <UserProvider>
            {children}
          </UserProvider>
        </Suspense>
        <Analytics />
      </body>
    </html>
  )
}
