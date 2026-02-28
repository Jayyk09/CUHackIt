'use client'

import { useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { Checkbox } from '@/components/ui/checkbox'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'

const LABELS = ['Organic', 'Vegetarian', 'Vegan', 'Pescatarian'] as const
const ALLERGENS = [
  'Milk',
  'Gluten',
  'Soybeans',
  'Eggs',
  'Nuts',
  'Peanuts',
  'Fish',
  'Shellfish',
] as const

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080'

export default function OnboardingPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const userId = searchParams.get('uid')

  const [selectedLabels, setSelectedLabels] = useState<string[]>([])
  const [selectedAllergens, setSelectedAllergens] = useState<string[]>([])
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)

  function toggle(
    list: string[],
    setList: (v: string[]) => void,
    value: string,
  ) {
    setList(
      list.includes(value)
        ? list.filter((v) => v !== value)
        : [...list, value],
    )
  }

  async function handleSubmit() {
    if (!userId) {
      setError('Missing user ID — please log in again.')
      return
    }

    setSaving(true)
    setError(null)

    try {
      // Save preferences using the uid from the query string.
      const res = await fetch(`${API_URL}/users/${encodeURIComponent(userId)}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          labels: selectedLabels,
          allergens: selectedAllergens,
        }),
      })

      if (!res.ok) {
        const text = await res.text()
        throw new Error(text || 'Failed to save preferences.')
      }

      router.push('/')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Something went wrong.')
    } finally {
      setSaving(false)
    }
  }

  return (
    <main className="min-h-screen flex items-center justify-center px-4">
      <div className="w-full max-w-lg space-y-10">
        {/* Header */}
        <div className="text-center space-y-2">
          <h1 className="font-serif text-3xl md:text-4xl tracking-[-0.02em]">
            Welcome to Sift
          </h1>
          <p className="text-muted-foreground text-sm md:text-base">
            Tell us about your dietary preferences so we can personalize your
            experience.
          </p>
        </div>

        {/* Dietary Labels */}
        <section className="space-y-4">
          <h2 className="font-sans text-xs tracking-[0.2em] uppercase text-muted-foreground">
            Dietary Preferences
          </h2>
          <div className="grid grid-cols-2 gap-3">
            {LABELS.map((label) => (
              <label
                key={label}
                className="flex items-center gap-3 rounded-md border px-4 py-3 cursor-pointer hover:bg-accent/50 transition-colors"
              >
                <Checkbox
                  checked={selectedLabels.includes(label)}
                  onCheckedChange={() =>
                    toggle(selectedLabels, setSelectedLabels, label)
                  }
                />
                <Label className="cursor-pointer">{label}</Label>
              </label>
            ))}
          </div>
        </section>

        {/* Allergens */}
        <section className="space-y-4">
          <h2 className="font-sans text-xs tracking-[0.2em] uppercase text-muted-foreground">
            Allergens
          </h2>
          <div className="grid grid-cols-2 gap-3">
            {ALLERGENS.map((allergen) => (
              <label
                key={allergen}
                className="flex items-center gap-3 rounded-md border px-4 py-3 cursor-pointer hover:bg-accent/50 transition-colors"
              >
                <Checkbox
                  checked={selectedAllergens.includes(allergen)}
                  onCheckedChange={() =>
                    toggle(selectedAllergens, setSelectedAllergens, allergen)
                  }
                />
                <Label className="cursor-pointer">{allergen}</Label>
              </label>
            ))}
          </div>
        </section>

        {/* Error */}
        {error && (
          <p className="text-sm text-destructive text-center">{error}</p>
        )}

        {/* Submit */}
        <div className="flex flex-col gap-3">
          <Button size="lg" onClick={handleSubmit} disabled={saving}>
            {saving ? 'Saving…' : 'Save Preferences'}
          </Button>
          <Button
            variant="ghost"
            size="lg"
            onClick={() => router.push('/')}
            disabled={saving}
          >
            Skip for now
          </Button>
        </div>
      </div>
    </main>
  )
}
