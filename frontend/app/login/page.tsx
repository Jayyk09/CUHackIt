'use client'

import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import Link from 'next/link'
import { useRouter } from 'next/navigation'

interface LoginTokenResponse {
  access_token: string
  id_token: string
  token_type: string
  expires_in: number
}

interface LoginErrorResponse {
  error: string
  error_description: string
}

async function loginWithCredentials(
  email: string,
  password: string
): Promise<LoginTokenResponse> {
  const res = await fetch('/api/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password }),
  })

  const data = await res.json()

  if (!res.ok) {
    const err = data as LoginErrorResponse
    throw new Error(err.error_description || err.error || 'Authentication failed')
  }

  return data as LoginTokenResponse
}

export default function LoginPage() {
  const router = useRouter()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)

    try {
      const tokens = await loginWithCredentials(email, password)

      // Store tokens for authenticated API calls
      localStorage.setItem('access_token', tokens.access_token)
      if (tokens.id_token) {
        localStorage.setItem('id_token', tokens.id_token)
      }

      // Redirect to dashboard on success
      router.push('/dashboard')
    } catch (err) {
      const message =
        err instanceof Error ? err.message : 'Something went wrong'

      // Map common Auth0 error descriptions to user-friendly messages
      if (message.includes('Wrong email or password')) {
        setError('Invalid email or password. Please try again.')
      } else if (message.includes('too many')) {
        setError('Too many failed attempts. Please wait and try again.')
      } else if (message.includes('blocked')) {
        setError('This account has been blocked. Please contact support.')
      } else {
        setError(message)
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-background flex">
      {/* Left panel - branding */}
      <div className="hidden lg:flex lg:w-1/2 bg-espresso relative overflow-hidden flex-col justify-between p-12 xl:p-16">
        {/* Decorative elements */}
        <div className="absolute inset-0 opacity-[0.03]">
          <div className="absolute top-[10%] left-[10%] w-64 h-64 border border-cream" />
          <div className="absolute top-[30%] right-[15%] w-48 h-48 border border-cream rotate-12" />
          <div className="absolute bottom-[15%] left-[20%] w-56 h-56 border border-cream -rotate-6" />
          <div className="absolute bottom-[35%] right-[10%] w-40 h-40 border border-cream rotate-45" />
        </div>

        <motion.div
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, ease: [0.22, 1, 0.36, 1] }}
        >
          <Link
            href="/"
            className="font-serif text-xl tracking-[-0.02em] text-cream/90 hover:text-cream transition-colors duration-300"
          >
            Sift
          </Link>
        </motion.div>

        <div className="relative z-10 flex-1 flex flex-col justify-center">
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.8, delay: 0.1, ease: [0.22, 1, 0.36, 1] }}
          >
            <span className="font-sans text-[10px] tracking-[0.3em] uppercase text-cream/30 block mb-8">
              Know Your Food
            </span>
            <h1 className="font-serif text-5xl xl:text-6xl 2xl:text-7xl tracking-[-0.03em] text-cream leading-[0.95]">
              Sift through
              <br />
              the noise.
            </h1>
          </motion.div>

          <motion.div
            className="mt-12 space-y-6"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.7, delay: 0.3, ease: [0.22, 1, 0.36, 1] }}
          >
            <div className="w-12 h-px bg-sage" />
            <p className="font-sans text-sm leading-relaxed text-cream/40 max-w-sm">
              Real food data, presented without compromise. Track what you eat,
              understand what it means, and take control of your health.
            </p>
          </motion.div>
        </div>

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ duration: 0.6, delay: 0.5 }}
        >
          <span className="font-sans text-[10px] tracking-[0.2em] uppercase text-cream/20">
            Est. 2026
          </span>
        </motion.div>
      </div>

      {/* Right panel - login form */}
      <div className="w-full lg:w-1/2 flex items-center justify-center px-6 sm:px-12 lg:px-16 xl:px-24">
        <motion.div
          className="w-full max-w-md"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7, delay: 0.15, ease: [0.22, 1, 0.36, 1] }}
        >
          {/* Mobile logo */}
          <div className="lg:hidden mb-12">
            <Link
              href="/"
              className="font-serif text-2xl tracking-[-0.02em] text-foreground"
            >
              Sift
            </Link>
          </div>

          <div className="mb-10">
            <span className="font-sans text-[10px] tracking-[0.3em] uppercase text-muted-foreground block mb-4">
              Account
            </span>
            <h2 className="font-serif text-4xl md:text-5xl tracking-[-0.03em] text-foreground leading-[0.95]">
              Welcome
              <br />
              back.
            </h2>
          </div>

          {/* Error message */}
          <AnimatePresence>
            {error && (
              <motion.div
                initial={{ opacity: 0, y: -8, height: 0 }}
                animate={{ opacity: 1, y: 0, height: 'auto' }}
                exit={{ opacity: 0, y: -8, height: 0 }}
                transition={{ duration: 0.3, ease: [0.22, 1, 0.36, 1] }}
                className="mb-6 overflow-hidden"
              >
                <div className="border border-destructive/30 bg-destructive/5 px-4 py-3 flex items-start gap-3">
                  <svg
                    className="w-4 h-4 text-destructive mt-0.5 shrink-0"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M12 9v3.75m9-.75a9 9 0 1 1-18 0 9 9 0 0 1 18 0Zm-9 3.75h.008v.008H12v-.008Z"
                    />
                  </svg>
                  <p className="font-sans text-xs text-destructive leading-relaxed">
                    {error}
                  </p>
                </div>
              </motion.div>
            )}
          </AnimatePresence>

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label
                htmlFor="email"
                className="block font-sans text-[10px] tracking-[0.25em] uppercase text-muted-foreground mb-3"
              >
                Email Address
              </label>
              <input
                id="email"
                type="email"
                value={email}
                onChange={(e) => {
                  setEmail(e.target.value)
                  if (error) setError(null)
                }}
                required
                disabled={loading}
                autoComplete="email"
                placeholder="you@example.com"
                className="w-full bg-transparent border border-foreground/15 px-4 py-3 font-sans text-sm text-foreground placeholder:text-muted-foreground/40 focus:outline-none focus:border-sage transition-colors duration-300 disabled:opacity-50 disabled:cursor-not-allowed"
              />
            </div>

            <div>
              <label
                htmlFor="password"
                className="block font-sans text-[10px] tracking-[0.25em] uppercase text-muted-foreground mb-3"
              >
                Password
              </label>
              <div className="relative">
                <input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => {
                    setPassword(e.target.value)
                    if (error) setError(null)
                  }}
                  required
                  disabled={loading}
                  autoComplete="current-password"
                  placeholder="Enter your password"
                  className="w-full bg-transparent border border-foreground/15 px-4 py-3 pr-16 font-sans text-sm text-foreground placeholder:text-muted-foreground/40 focus:outline-none focus:border-sage transition-colors duration-300 disabled:opacity-50 disabled:cursor-not-allowed"
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute right-4 top-1/2 -translate-y-1/2 font-sans text-[10px] tracking-[0.15em] uppercase text-muted-foreground hover:text-foreground transition-colors duration-300"
                >
                  {showPassword ? 'Hide' : 'Show'}
                </button>
              </div>
            </div>

            <div className="flex items-center justify-between pt-1">
              <label className="flex items-center gap-2 cursor-pointer group">
                <span className="relative w-4 h-4 border border-foreground/20 flex items-center justify-center group-hover:border-foreground/40 transition-colors duration-300">
                  <input type="checkbox" className="sr-only peer" />
                  <span className="hidden peer-checked:block w-2 h-2 bg-sage" />
                </span>
                <span className="font-sans text-[10px] tracking-[0.15em] uppercase text-muted-foreground">
                  Remember me
                </span>
              </label>
              <button
                type="button"
                className="font-sans text-[10px] tracking-[0.15em] uppercase text-muted-foreground hover:text-foreground transition-colors duration-300"
              >
                Forgot password?
              </button>
            </div>

            <div className="pt-4">
              <button
                type="submit"
                disabled={loading}
                className="w-full border border-accent bg-accent text-accent-foreground px-8 py-3.5 font-sans text-xs tracking-[0.25em] uppercase hover:bg-sage-dark transition-colors duration-300 cursor-pointer disabled:opacity-60 disabled:cursor-not-allowed flex items-center justify-center gap-2"
              >
                {loading ? (
                  <>
                    <svg
                      className="animate-spin h-3.5 w-3.5"
                      viewBox="0 0 24 24"
                      fill="none"
                    >
                      <circle
                        className="opacity-25"
                        cx="12"
                        cy="12"
                        r="10"
                        stroke="currentColor"
                        strokeWidth="2"
                      />
                      <path
                        className="opacity-75"
                        fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
                      />
                    </svg>
                    Signing In
                  </>
                ) : (
                  'Sign In'
                )}
              </button>
            </div>
          </form>

          <div className="mt-10 flex items-center gap-4">
            <div className="flex-1 h-px bg-foreground/10" />
            <span className="font-sans text-[10px] tracking-[0.2em] uppercase text-muted-foreground/50">
              Or
            </span>
            <div className="flex-1 h-px bg-foreground/10" />
          </div>

          <div className="mt-6 space-y-3">
            <button
              type="button"
              className="w-full border border-foreground/15 bg-transparent text-foreground px-8 py-3.5 font-sans text-xs tracking-[0.2em] uppercase hover:bg-secondary transition-colors duration-300 cursor-pointer flex items-center justify-center gap-3"
            >
              <svg className="w-4 h-4" viewBox="0 0 24 24" fill="none">
                <path
                  d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z"
                  fill="#4285F4"
                />
                <path
                  d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                  fill="#34A853"
                />
                <path
                  d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                  fill="#FBBC05"
                />
                <path
                  d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                  fill="#EA4335"
                />
              </svg>
              Continue with Google
            </button>
          </div>

          <p className="mt-10 text-center font-sans text-xs tracking-[0.1em] text-muted-foreground/60">
            Don&apos;t have an account?{' '}
            <button
              type="button"
              className="text-foreground hover:text-sage transition-colors duration-300 tracking-[0.15em] uppercase"
            >
              Create one
            </button>
          </p>

          {/* Back to home link for mobile */}
          <div className="mt-8 lg:hidden text-center">
            <Link
              href="/"
              className="font-sans text-[10px] tracking-[0.2em] uppercase text-muted-foreground hover:text-foreground transition-colors duration-300"
            >
              Back to Home
            </Link>
          </div>
        </motion.div>
      </div>
    </div>
  )
}
