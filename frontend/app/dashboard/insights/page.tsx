'use client'

import { useEffect, useMemo, useState } from 'react'
import { toast } from 'sonner'
import { useUser } from '@/contexts/user-context'
import {
  getCategorySummary,
  listPantryItems,
  type PantryItem,
} from '@/lib/pantry-api'

const benefitHighlights = [
  { label: 'Items Rescued', value: '24', detail: 'past 30 days' },
  { label: 'Waste Avoided', value: '6.3 kg', detail: 'food saved' },
  { label: 'Money Saved', value: '$42', detail: 'estimated' },
  { label: 'CO2 Offset', value: '9.1 kg', detail: 'impact score' },
]

const pantrySignalLabels = {
  avgDays: 'Avg Days to Expiry',
  avgNutri: 'Avg Nutri-Score',
  avgEco: 'Avg Eco Score',
  freshness: 'Freshness Index',
}

const fallbackActions = [
  'Use in a skillet',
  'Blend into sauce',
  'Roast tonight',
  'Add to salad',
]

export default function InsightsPage() {
  const { user, isLoading: isUserLoading } = useUser()
  const [pantryItems, setPantryItems] = useState<PantryItem[]>([])
  const [categorySummary, setCategorySummary] = useState<Record<string, number>>({})
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    if (!user) return
    setIsLoading(true)
    Promise.all([listPantryItems(user.id), getCategorySummary(user.id)])
      .then(([items, summary]) => {
        setPantryItems(items)
        setCategorySummary(summary)
      })
      .catch((err) => {
        console.error('Failed to load insights:', err)
        toast.error(err instanceof Error ? err.message : 'Failed to load insights')
      })
      .finally(() => setIsLoading(false))
  }, [user])

  const pantrySignals = useMemo(() => {
    const now = Date.now()
    const expiryDays = pantryItems
      .map((item) => {
        if (!item.shelf_life) return null
        const addedAt = new Date(item.added_at).getTime()
        if (Number.isNaN(addedAt)) return null
        const expiry = addedAt + item.shelf_life * 24 * 60 * 60 * 1000
        const daysLeft = Math.ceil((expiry - now) / (24 * 60 * 60 * 1000))
        return daysLeft
      })
      .filter((value): value is number => value !== null)

    const avgDays = expiryDays.length
      ? (expiryDays.reduce((sum, d) => sum + d, 0) / expiryDays.length).toFixed(1)
      : '--'

    const ecoScores = pantryItems
      .map((item) => item.norm_environmental_score)
      .map((value) => (typeof value === 'string' ? Number.parseFloat(value) : value))
      .filter((value): value is number => typeof value === 'number' && !Number.isNaN(value))

    const avgEco = ecoScores.length
      ? Math.round(ecoScores.reduce((sum, v) => sum + v, 0) / ecoScores.length).toString()
      : '--'

    const nutriScores = pantryItems
      .map((item) => item.norm_nutriscore_score ?? item.nutriscore_score)
      .map((value) => (typeof value === 'string' ? Number.parseFloat(value) : value))
      .filter((value): value is number => typeof value === 'number' && !Number.isNaN(value))

    const avgNutri = nutriScores.length
      ? Math.round(nutriScores.reduce((sum, v) => sum + v, 0) / nutriScores.length)
      : null

    const nutriLetter = avgNutri === null
      ? '--'
      : avgNutri >= 80
        ? 'A'
        : avgNutri >= 65
          ? 'B'
          : avgNutri >= 50
            ? 'C'
            : avgNutri >= 35
              ? 'D'
              : 'E'

    const freshnessIndex = expiryDays.length
      ? Math.round((expiryDays.filter((d) => d > 3).length / expiryDays.length) * 100)
      : null

    return [
      { label: pantrySignalLabels.avgDays, value: avgDays, detail: 'across pantry' },
      { label: pantrySignalLabels.avgNutri, value: nutriLetter, detail: 'quality signal' },
      { label: pantrySignalLabels.avgEco, value: avgEco, detail: 'sustainability' },
      {
        label: pantrySignalLabels.freshness,
        value: freshnessIndex === null ? '--' : `${freshnessIndex}%`,
        detail: 'current mix',
      },
    ]
  }, [pantryItems])

  const expiringSoon = useMemo(() => {
    const now = Date.now()
    const items = pantryItems
      .map((item) => {
        if (!item.shelf_life) return null
        const addedAt = new Date(item.added_at).getTime()
        if (Number.isNaN(addedAt)) return null
        const expiry = addedAt + item.shelf_life * 24 * 60 * 60 * 1000
        const daysLeft = Math.ceil((expiry - now) / (24 * 60 * 60 * 1000))
        return {
          name: item.product_name || 'Item',
          days: daysLeft,
          action: fallbackActions[Math.abs(item.id) % fallbackActions.length],
        }
      })
      .filter((value): value is { name: string; days: number; action: string } => value !== null)
      .sort((a, b) => a.days - b.days)
      .slice(0, 3)

    if (items.length > 0) return items

    return [
      { name: 'No expiring items', days: 0, action: 'You are in good shape' },
    ]
  }, [pantryItems])

  const pantryMix = useMemo(() => {
    const total = Object.values(categorySummary).reduce((sum, value) => sum + value, 0)
    if (!total) {
      return [{ label: 'No items yet', value: '0%' }]
    }
    return Object.entries(categorySummary)
      .map(([label, count]) => ({
        label: label || 'Uncategorized',
        value: `${Math.round((count / total) * 100)}%`,
      }))
      .sort((a, b) => Number.parseInt(b.value, 10) - Number.parseInt(a.value, 10))
      .slice(0, 5)
  }, [categorySummary])

  if (isUserLoading || isLoading) {
    return (
      <div className="min-h-screen">
        <div className="border-b border-espresso/10">
          <div className="px-8 py-6">
            <h1 className="font-serif text-2xl text-espresso">Insights</h1>
          </div>
        </div>
        <div className="px-8 py-10">
          <p className="font-serif text-xl text-espresso/30">Loading your insights...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen">
      <div className="border-b border-espresso/10">
        <div className="px-8 py-6">
          <p className="font-sans text-xs uppercase tracking-[0.3em] text-espresso/40">
            Insights
          </p>
          <h1 className="font-serif text-3xl text-espresso">Your Impact and Pantry Health</h1>
        </div>
      </div>

      <div className="px-8 py-10 space-y-10">
        <section className="grid gap-6 lg:grid-cols-4">
          {benefitHighlights.map((item) => (
            <div
              key={item.label}
              className="rounded-3xl border border-espresso/10 bg-cream p-6"
            >
              <p className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/50">
                {item.label}
              </p>
              <p className="font-serif text-3xl text-espresso mt-4">{item.value}</p>
              <p className="text-espresso/60 font-sans text-sm mt-2">{item.detail}</p>
            </div>
          ))}
        </section>

        <section className="rounded-3xl border border-espresso/10 bg-gradient-to-br from-cream via-cream to-sage/10 p-8">
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div>
              <p className="text-xs uppercase tracking-[0.3em] font-sans text-espresso/40">
                Pantry Signals
              </p>
              <h2 className="font-serif text-2xl text-espresso">At-a-Glance Pantry Metrics</h2>
            </div>
            <button className="rounded-full border border-espresso/20 px-5 py-2 text-xs uppercase tracking-[0.25em] font-sans text-espresso hover:border-espresso/40 transition">
              Optimize Pantry
            </button>
          </div>
          <div className="mt-6 grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            {pantrySignals.map((item) => (
              <div
                key={item.label}
                className="rounded-2xl border border-espresso/10 bg-cream/80 p-5"
              >
                <p className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/50">
                  {item.label}
                </p>
                <p className="font-serif text-3xl text-espresso mt-3">{item.value}</p>
                <p className="text-espresso/60 font-sans text-sm mt-2">{item.detail}</p>
              </div>
            ))}
          </div>
        </section>

        <section className="grid gap-6 lg:grid-cols-[1.1fr_0.9fr]">
          <div className="rounded-3xl border border-espresso/10 bg-cream p-8">
            <div className="flex items-center justify-between">
              <h2 className="font-serif text-2xl text-espresso">Expiring Soon</h2>
              <span className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/40">
                next 72 hours
              </span>
            </div>
            <div className="mt-6 space-y-4">
              {expiringSoon.map((item) => (
                <div
                  key={item.name}
                  className="flex items-center justify-between rounded-2xl border border-espresso/10 px-4 py-3"
                >
                  <div>
                    <p className="font-sans text-espresso">{item.name}</p>
                    <p className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/40 mt-1">
                      {item.action}
                    </p>
                  </div>
                  <span className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/50">
                    {item.days} days
                  </span>
                </div>
              ))}
            </div>
            <button className="mt-6 w-full rounded-full border border-espresso/20 px-5 py-2 text-xs uppercase tracking-[0.25em] font-sans text-espresso hover:border-espresso/40 transition">
              Build a rescue recipe
            </button>
          </div>

          <div className="rounded-3xl border border-espresso/10 bg-cream p-8">
            <div className="flex items-center justify-between">
              <h2 className="font-serif text-2xl text-espresso">Pantry Mix</h2>
              <span className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/40">
                balance view
              </span>
            </div>
            <div className="mt-6 space-y-4">
              {pantryMix.map((item) => (
                <div key={item.label} className="space-y-2">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-sans text-espresso/70">{item.label}</span>
                    <span className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/50">
                      {item.value}
                    </span>
                  </div>
                  <div className="h-2 rounded-full bg-espresso/10">
                    <div
                      className="h-2 rounded-full bg-sage"
                      style={{ width: item.value }}
                    />
                  </div>
                </div>
              ))}
            </div>
          </div>
        </section>

      </div>
    </div>
  )
}
