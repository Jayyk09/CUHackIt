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
    icon: [
      {
        url: '/icon-light-32x32.png',
        media: '(prefers-color-scheme: light)',
      },
      {
        url: '/icon-dark-32x32.png',
        media: '(prefers-color-scheme: dark)',
      },
      {
        url: '/icon.svg',
        type: 'image/svg+xml',
      },
    ],
    apple: '/apple-icon.png',
  },
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" className={`${playfair.variable} ${spaceGrotesk.variable}`}>
      <body className="font-sans antialiased">
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
