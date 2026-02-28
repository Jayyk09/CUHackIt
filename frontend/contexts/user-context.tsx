'use client'

import {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
  type ReactNode,
} from 'react'
import { useSearchParams, useRouter, usePathname } from 'next/navigation'
import { getUserByAuth0ID, type User } from '@/lib/user-api'

const STORAGE_KEY = 'sift_user'

interface UserContextValue {
  /** The fully-resolved internal user (null while loading or when logged out). */
  user: User | null
  /** True during the initial resolution of the auth0 sub → internal user. */
  isLoading: boolean
  /** Clears stored user state (call on logout). */
  clearUser: () => void
}

const UserContext = createContext<UserContextValue>({
  user: null,
  isLoading: true,
  clearUser: () => {},
})

export function useUser() {
  return useContext(UserContext)
}

export function UserProvider({ children }: { children: ReactNode }) {
  const searchParams = useSearchParams()
  const router = useRouter()
  const pathname = usePathname()

  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  const clearUser = useCallback(() => {
    setUser(null)
    try {
      localStorage.removeItem(STORAGE_KEY)
    } catch {
      // SSR or storage unavailable
    }
  }, [])

  useEffect(() => {
    async function resolve() {
      // 1. Check if we already have a cached user in localStorage
      try {
        const cached = localStorage.getItem(STORAGE_KEY)
        if (cached) {
          const parsed: User = JSON.parse(cached)
          if (parsed?.id) {
            setUser(parsed)
            setIsLoading(false)
            // Still check URL for a fresh uid (user may have re-logged in)
          }
        }
      } catch {
        // ignore parse errors
      }

      // 2. Check URL for ?uid= (passed by backend after OAuth callback)
      const uid = searchParams.get('uid')
      if (uid) {
        try {
          const resolved = await getUserByAuth0ID(uid)
          setUser(resolved)
          localStorage.setItem(STORAGE_KEY, JSON.stringify(resolved))

          // Clean the uid param from the URL so it doesn't linger
          const params = new URLSearchParams(searchParams.toString())
          params.delete('uid')
          const cleanURL = params.toString()
            ? `${pathname}?${params.toString()}`
            : pathname
          router.replace(cleanURL)
        } catch (err) {
          console.error('Failed to resolve user from auth0 uid:', err)
        }
      }

      setIsLoading(false)
    }

    resolve()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchParams])

  return (
    <UserContext.Provider value={{ user, isLoading, clearUser }}>
      {children}
    </UserContext.Provider>
  )
}
