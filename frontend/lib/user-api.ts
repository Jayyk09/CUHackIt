const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080'

export interface User {
  id: string
  auth0_id: string
  email: string
  name?: string
  allergens: string[]
  dietary_preferences: string[]
  nutritional_goals: string[]
  cooking_skill: string
  cuisine_preferences: string[]
  onboarding_completed: boolean
  created_at: string
  updated_at: string
}

/**
 * Resolve an Auth0 sub (oauth ID) to the full internal user object.
 */
export async function getUserByAuth0ID(auth0Id: string): Promise<User> {
  const res = await fetch(`${API_BASE}/auth0-users/${encodeURIComponent(auth0Id)}`)
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Failed to resolve user')
  }
  return res.json()
}
