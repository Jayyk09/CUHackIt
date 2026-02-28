'use client'

import { useState, useEffect, useMemo } from 'react'
import Image from 'next/image'
import { motion, AnimatePresence } from 'framer-motion'
import { toast } from 'sonner'
import { useUser } from '@/contexts/user-context'
import { listPantryItems, getCategorySummary, type PantryItem } from '@/lib/pantry-api'
import { getCategoryImage, getAvailableCategories } from '@/lib/category-images'
import { cn } from '@/lib/utils'

const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080'

// ─── Delete helper ────────────────────────────────────────────────────────────

async function deletePantryItem(userId: string, itemId: number): Promise<void> {
  const res = await fetch(`${API_BASE}/users/${userId}/pantry/${itemId}`, {
    method: 'DELETE',
  })
  if (!res.ok && res.status !== 204) {
    const err = await res.json().catch(() => ({}))
    throw new Error(err.error || 'Failed to delete item')
  }
}

// ─── Animations ───────────────────────────────────────────────────────────────

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: { staggerChildren: 0.08, delayChildren: 0.1 },
  },
}

const itemVariants = {
  hidden: { opacity: 0, y: 10 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.5, ease: [0.22, 1, 0.36, 1] },
  },
  exit: {
    opacity: 0,
    x: -20,
    transition: { duration: 0.3, ease: [0.22, 1, 0.36, 1] },
  },
}

const categoryVariants = {
  hidden: { opacity: 0, y: 16 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.6, ease: [0.22, 1, 0.36, 1] },
  },
}

// ─── Category filter pill ─────────────────────────────────────────────────────

function CategoryPill({
  label,
  count,
  isActive,
  onClick,
}: {
  label: string
  count: number
  isActive: boolean
  onClick: () => void
}) {
  return (
    <button
      onClick={onClick}
      className={cn(
        'px-4 py-2 font-sans text-[10px] uppercase tracking-[0.2em] transition-colors duration-200 border',
        isActive
          ? 'border-sage bg-sage text-cream'
          : 'border-espresso/15 text-espresso/50 hover:border-espresso/30 hover:text-espresso/70'
      )}
    >
      {label}
      <span className="ml-2 opacity-60">{count}</span>
    </button>
  )
}

// ─── Pantry item row ──────────────────────────────────────────────────────────

function PantryItemRow({
  item,
  onDelete,
  isDeleting,
}: {
  item: PantryItem
  onDelete: (id: number) => void
  isDeleting: boolean
}) {
  return (
    <motion.div
      variants={itemVariants}
      exit="exit"
      layout
      className="group flex items-center gap-4 py-5 border-b border-espresso/10 hover:bg-espresso/[0.02] transition-colors duration-300"
    >
      {/* Thumbnail */}
      <div className="relative w-14 h-14 flex-shrink-0 border border-espresso/15 overflow-hidden bg-espresso/[0.03]">
        {item.image_url || item.image_small_url ? (
          <Image
            src={item.image_url ?? item.image_small_url ?? ''}
            alt={item.product_name}
            fill
            sizes="56px"
            className="object-cover"
            unoptimized
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center">
            <span className="text-[9px] tracking-[0.15em] uppercase text-espresso/30">
              No img
            </span>
          </div>
        )}
      </div>

      {/* Name + category */}
      <div className="flex-1 min-w-0">
        <p className="font-sans font-semibold text-base text-espresso leading-tight truncate">
          {item.product_name}
        </p>
        <p className="text-[10px] tracking-[0.2em] uppercase text-espresso/50 mt-1">
          {item.category ?? 'Uncategorized'}
        </p>
      </div>

      {/* Badges */}
      <div className="flex items-center gap-3 flex-shrink-0">
        {/* Quantity */}
        <span className="px-3 py-1 border border-espresso/15 font-sans text-[10px] uppercase tracking-[0.15em] text-espresso/60">
          qty {item.quantity}
        </span>

        {/* Frozen badge */}
        {item.is_frozen && (
          <span className="px-3 py-1 border border-sky-400/30 bg-sky-50 font-sans text-[10px] uppercase tracking-[0.15em] text-sky-600">
            Frozen
          </span>
        )}

        {/* Shelf life */}
        {item.shelf_life != null && (
          <span
            className={cn(
              'px-3 py-1 border font-sans text-[10px] uppercase tracking-[0.15em]',
              item.shelf_life <= 3
                ? 'border-red-300/40 text-red-600/70'
                : item.shelf_life <= 7
                  ? 'border-amber-300/40 text-amber-600/70'
                  : 'border-sage/30 text-sage/70'
            )}
          >
            {item.shelf_life}d shelf
          </span>
        )}
      </div>

      {/* Scores */}
      <div className="hidden md:flex items-center gap-4 flex-shrink-0">
        {item.nutriscore_score != null && (
          <div className="text-right">
            <p className="text-[9px] tracking-[0.15em] uppercase text-espresso/35 mb-0.5">
              Nutri
            </p>
            <p className="font-serif text-xl text-espresso/70 leading-none">
              {Math.round(item.nutriscore_score)}
            </p>
          </div>
        )}
        {item.norm_environmental_score != null && (
          <div className="text-right">
            <p className="text-[9px] tracking-[0.15em] uppercase text-espresso/35 mb-0.5">
              Eco
            </p>
            <p className="font-serif text-xl text-sage leading-none">
              {Math.round(item.norm_environmental_score)}
            </p>
          </div>
        )}
      </div>

      {/* Delete */}
      <button
        onClick={() => onDelete(item.id)}
        disabled={isDeleting}
        className={cn(
          'flex-shrink-0 px-3 py-1.5 font-sans text-[9px] uppercase tracking-[0.2em] transition-all duration-200',
          'opacity-0 group-hover:opacity-100',
          isDeleting
            ? 'text-espresso/30 cursor-not-allowed'
            : 'text-espresso/40 hover:text-red-600'
        )}
      >
        {isDeleting ? '...' : 'Remove'}
      </button>
    </motion.div>
  )
}

// ─── Category section ─────────────────────────────────────────────────────────

function CategorySection({
  category,
  items,
  onDelete,
  deletingIds,
}: {
  category: string
  items: PantryItem[]
  onDelete: (id: number) => void
  deletingIds: Set<number>
}) {
  const categoryImage = getCategoryImage(category)

  return (
    <motion.section variants={categoryVariants} className="mb-12">
      {/* Category header */}
      <div className="flex items-center gap-5 mb-2">
        <div className="relative w-10 h-10 flex-shrink-0">
          <Image
            src={categoryImage}
            alt={category}
            fill
            className="object-contain"
          />
        </div>
        <div className="flex items-baseline gap-3 flex-1 border-b border-espresso/15 pb-3">
          <h2 className="font-serif text-2xl text-espresso capitalize">
            {category.toLowerCase()}
          </h2>
          <span className="text-[10px] tracking-[0.2em] uppercase text-espresso/35 font-sans">
            {items.length} {items.length === 1 ? 'item' : 'items'}
          </span>
        </div>
      </div>

      {/* Items */}
      <motion.div variants={containerVariants} initial="hidden" animate="visible">
        <AnimatePresence mode="popLayout">
          {items.map((item) => (
            <PantryItemRow
              key={item.id}
              item={item}
              onDelete={onDelete}
              isDeleting={deletingIds.has(item.id)}
            />
          ))}
        </AnimatePresence>
      </motion.div>
    </motion.section>
  )
}

// ─── Page ─────────────────────────────────────────────────────────────────────

export default function PantryPage() {
  const { user, isLoading: isUserLoading } = useUser()
  const [pantryItems, setPantryItems] = useState<PantryItem[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [activeCategory, setActiveCategory] = useState<string | null>(null)
  const [deletingIds, setDeletingIds] = useState<Set<number>>(new Set())

  // Fetch pantry items
  useEffect(() => {
    if (!user) return
    setIsLoading(true)
    listPantryItems(user.id)
      .then(setPantryItems)
      .catch((err) => {
        console.error('Failed to load pantry:', err)
        toast.error('Could not load pantry items')
      })
      .finally(() => setIsLoading(false))
  }, [user])

  // Group items by category
  const groupedItems = useMemo(() => {
    const groups: Record<string, PantryItem[]> = {}
    const filtered = activeCategory
      ? pantryItems.filter(
          (p) => (p.category ?? 'uncategorized').toLowerCase() === activeCategory.toLowerCase()
        )
      : pantryItems

    filtered.forEach((item) => {
      const cat = (item.category ?? 'Uncategorized').toUpperCase()
      if (!groups[cat]) groups[cat] = []
      groups[cat].push(item)
    })

    // Sort categories alphabetically
    const sorted: Record<string, PantryItem[]> = {}
    Object.keys(groups)
      .sort()
      .forEach((key) => {
        sorted[key] = groups[key]
      })

    return sorted
  }, [pantryItems, activeCategory])

  // Category counts for filter pills
  const categoryCounts = useMemo(() => {
    const counts: Record<string, number> = {}
    pantryItems.forEach((item) => {
      const cat = (item.category ?? 'Uncategorized').toLowerCase()
      counts[cat] = (counts[cat] ?? 0) + 1
    })
    return counts
  }, [pantryItems])

  const totalItems = pantryItems.length

  const handleDelete = async (itemId: number) => {
    if (!user) return

    setDeletingIds((prev) => new Set(prev).add(itemId))

    try {
      await deletePantryItem(user.id, itemId)
      setPantryItems((prev) => prev.filter((p) => p.id !== itemId))
      toast.success('Item removed from pantry')
    } catch (err) {
      console.error('Failed to delete item:', err)
      toast.error(err instanceof Error ? err.message : 'Failed to remove item')
    } finally {
      setDeletingIds((prev) => {
        const next = new Set(prev)
        next.delete(itemId)
        return next
      })
    }
  }

  // Loading state
  if (isUserLoading || isLoading) {
    return (
      <>
        <div className="border-b border-espresso/10">
          <div className="px-8 py-6">
            <h1 className="font-serif text-2xl text-espresso">Pantry</h1>
          </div>
        </div>
        <div className="p-8">
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="py-16 text-center"
          >
            <p className="font-serif text-xl text-espresso/30 mb-2">
              Loading your pantry
            </p>
            <p className="font-sans text-[11px] uppercase tracking-[0.2em] text-espresso/20">
              fetching items
            </p>
          </motion.div>
        </div>
      </>
    )
  }

  // Not logged in
  if (!user) {
    return (
      <>
        <div className="border-b border-espresso/10">
          <div className="px-8 py-6">
            <h1 className="font-serif text-2xl text-espresso">Pantry</h1>
          </div>
        </div>
        <div className="p-8">
          <div className="py-16 text-center">
            <p className="font-serif text-xl text-espresso/50">
              Sign in to view your pantry
            </p>
          </div>
        </div>
      </>
    )
  }

  return (
    <>
      {/* Header */}
      <div className="border-b border-espresso/10">
        <div className="px-8 py-6 flex items-baseline justify-between">
          <h1 className="font-serif text-2xl text-espresso">Pantry</h1>
          <span className="font-sans text-[10px] uppercase tracking-[0.2em] text-espresso/40">
            {totalItems} {totalItems === 1 ? 'item' : 'items'}
          </span>
        </div>
      </div>

      <div className="p-8">
        {/* Category filters */}
        {Object.keys(categoryCounts).length > 1 && (
          <div className="flex flex-wrap gap-2 mb-10">
            <CategoryPill
              label="All"
              count={totalItems}
              isActive={activeCategory === null}
              onClick={() => setActiveCategory(null)}
            />
            {Object.entries(categoryCounts)
              .sort(([a], [b]) => a.localeCompare(b))
              .map(([cat, count]) => (
                <CategoryPill
                  key={cat}
                  label={cat}
                  count={count}
                  isActive={activeCategory === cat}
                  onClick={() =>
                    setActiveCategory(activeCategory === cat ? null : cat)
                  }
                />
              ))}
          </div>
        )}

        {/* Empty state */}
        {totalItems === 0 ? (
          <motion.div
            initial={{ opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, ease: [0.22, 1, 0.36, 1] }}
            className="py-20 text-center"
          >
            <p className="font-serif text-2xl text-espresso/40 mb-3">
              Your pantry is empty
            </p>
            <p className="font-sans text-[11px] uppercase tracking-[0.2em] text-espresso/25">
              Search for food on the dashboard to add items
            </p>
          </motion.div>
        ) : Object.keys(groupedItems).length === 0 ? (
          <motion.div
            initial={{ opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            className="py-16 text-center"
          >
            <p className="font-serif text-xl text-espresso/40">
              No items in this category
            </p>
          </motion.div>
        ) : (
          /* Category sections */
          <motion.div
            variants={containerVariants}
            initial="hidden"
            animate="visible"
          >
            {Object.entries(groupedItems).map(([category, items]) => (
              <CategorySection
                key={category}
                category={category}
                items={items}
                onDelete={handleDelete}
                deletingIds={deletingIds}
              />
            ))}
          </motion.div>
        )}
      </div>
    </>
  )
}
