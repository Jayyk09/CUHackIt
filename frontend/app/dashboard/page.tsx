'use client'

import { useState, useMemo } from 'react'
import Image from 'next/image'
import { motion, AnimatePresence } from 'framer-motion'
import { toast } from 'sonner'
import { FoodItem, mockPantryItems } from '@/lib/mock-data'
import { getCategoryImage } from '@/lib/category-images'
import { useDebouncedSearch } from '@/hooks/use-debounced-search'

function SearchBar({
  value,
  onChange,
}: {
  value: string
  onChange: (value: string) => void
}) {
  return (
    <div className="relative w-full border-b border-espresso pb-6">
      <div className="flex items-center gap-4">
        <input
          type="text"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder="search recipe or ingredients"
          className="w-full bg-transparent font-serif text-2xl md:text-3xl text-espresso placeholder:text-espresso/40 placeholder:font-serif focus:outline-none"
        />
        <Image
          src="/leaf.png"
          alt="Search"
          width={28}
          height={28}
          unoptimized
          className="flex-shrink-0"
        />
      </div>
    </div>
  )
}

function FoodGridItem({
  item,
  categoryImage,
  onSelect,
}: {
  item: FoodItem
  categoryImage: string
  onSelect: (item: FoodItem) => void
}) {
  return (
    <button
      onClick={() => onSelect(item)}
      className="group text-left w-full focus:outline-none"
    >
      <div className="relative aspect-square mb-4 overflow-hidden">
        <Image
          src={categoryImage}
          alt={item.product_name}
          fill
          className="object-contain transition-transform duration-500 group-hover:scale-105"
        />
      </div>
      <div className="border-t border-espresso pt-3">
        <p className="font-sans font-bold text-sm text-espresso truncate">
          {item.product_name}
        </p>
      </div>
    </button>
  )
}

function FoodGrid({
  items,
  onSelect,
}: {
  items: FoodItem[]
  onSelect: (item: FoodItem) => void
}) {
  // Generate category images once on mount, stable for each item
  const categoryImages = useMemo(() => {
    const images: Record<string, string> = {}
    items.forEach((item) => {
      if (!images[item.id]) {
        images[item.id] = getCategoryImage(item.category)
      }
    })
    return images
  }, [items])

  return (
    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-x-6 gap-y-12">
      {items.map((item) => (
        <FoodGridItem
          key={item.id}
          item={item}
          categoryImage={categoryImages[item.id] || getCategoryImage(item.category)}
          onSelect={onSelect}
        />
      ))}
    </div>
  )
}

function DetailModal({
  item,
  onClose,
  onAddToPantry,
  isInPantry,
}: {
  item: FoodItem
  onClose: () => void
  onAddToPantry: (item: FoodItem) => void
  isInPantry: boolean
}) {
  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      transition={{ duration: 0.2 }}
      className="fixed inset-0 z-50 flex items-center justify-center bg-espresso/30"
      onClick={onClose}
    >
      <motion.div
        initial={{ scale: 0.95, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        exit={{ scale: 0.95, opacity: 0 }}
        transition={{ type: 'spring', damping: 30, stiffness: 300 }}
        className="relative w-full max-w-2xl mx-4 bg-cream border border-espresso/10 max-h-[90vh] overflow-y-auto"
        onClick={(e) => e.stopPropagation()}
      >
        <button
          onClick={onClose}
          className="absolute top-6 right-6 text-xs uppercase tracking-[0.2em] font-sans text-espresso/50 hover:text-espresso transition-colors z-10"
        >
          — Close
        </button>

        <div className="p-8 md:p-10">
          {/* Open Food Facts image */}
          <div className="relative aspect-video w-full mb-8 overflow-hidden bg-espresso/5">
            <Image
              src={item.image_url}
              alt={item.product_name}
              fill
              className="object-contain"
              unoptimized
            />
          </div>

          <h2 className="font-serif text-3xl md:text-4xl text-espresso mb-8 leading-tight">
            {item.product_name}
          </h2>

          <div className="flex items-stretch gap-8 mb-10">
            <div className="flex-1">
              <p className="font-sans text-xs uppercase tracking-[0.2em] text-espresso/50 mb-2">
                Environmental
              </p>
              <p className="font-serif text-5xl md:text-6xl text-sage leading-none">
                {item.environmental_score}
              </p>
            </div>
            <div className="w-px bg-espresso/10" />
            <div className="flex-1">
              <p className="font-sans text-xs uppercase tracking-[0.2em] text-espresso/50 mb-2">
                Nutri-Score
              </p>
              <p className="font-serif text-5xl md:text-6xl text-sage leading-none">
                {item.nutriscore_score}
              </p>
            </div>
          </div>

          <div className="space-y-6">
            {item.labels_en.length > 0 && (
              <div>
                <p className="font-sans text-xs uppercase tracking-[0.2em] text-espresso/50 mb-3">
                  Labels
                </p>
                <div className="flex flex-wrap gap-2">
                  {item.labels_en.map((label) => (
                    <span
                      key={label}
                      className="px-3 py-1 border border-espresso/20 font-sans text-xs uppercase tracking-wider text-espresso"
                    >
                      {label}
                    </span>
                  ))}
                </div>
              </div>
            )}

            {item.allergens_en.length > 0 && (
              <div>
                <p className="font-sans text-xs uppercase tracking-[0.2em] text-espresso/50 mb-3">
                  Allergens
                </p>
                <div className="flex flex-wrap gap-2">
                  {item.allergens_en.map((allergen) => (
                    <span
                      key={allergen}
                      className="px-3 py-1 border border-sage/40 font-sans text-xs uppercase tracking-wider text-sage"
                    >
                      {allergen}
                    </span>
                  ))}
                </div>
              </div>
            )}

            {item.traces_en.length > 0 && (
              <div>
                <p className="font-sans text-xs uppercase tracking-[0.2em] text-espresso/50 mb-3">
                  May Contain
                </p>
                <div className="flex flex-wrap gap-2">
                  {item.traces_en.map((trace) => (
                    <span
                      key={trace}
                      className="px-3 py-1 border border-espresso/10 font-sans text-xs uppercase tracking-wider text-espresso/60"
                    >
                      {trace}
                    </span>
                  ))}
                </div>
              </div>
            )}

            <div className="pt-6 border-t border-espresso/10">
              <div className="grid grid-cols-2 gap-6">
                <div>
                  <p className="font-sans text-xs uppercase tracking-[0.2em] text-espresso/50 mb-1">
                    Category
                  </p>
                  <p className="font-sans text-sm text-espresso capitalize">{item.category}</p>
                </div>
                <div>
                  <p className="font-sans text-xs uppercase tracking-[0.2em] text-espresso/50 mb-1">
                    Shelf Life
                  </p>
                  <p className="font-sans text-sm text-espresso">
                    {item.shelf_life} days
                  </p>
                </div>
                <div>
                </div>
                <div>
                  <p className="font-sans text-xs uppercase tracking-[0.2em] text-espresso/50 mb-1">
                    Status
                  </p>
                  <p
                    className={`font-sans text-sm ${
                      item.is_spoiled ? 'text-red-700' : 'text-sage'
                    }`}
                  >
                    {item.is_spoiled ? 'Spoiled' : 'Fresh'}
                  </p>
                </div>
              </div>
            </div>

            {/* Add to Pantry Button */}
            {!isInPantry && (
              <div className="pt-6 border-t border-espresso/10">
                <button
                  onClick={() => onAddToPantry(item)}
                  className="w-full py-3 px-6 border border-sage text-sage font-sans text-xs uppercase tracking-[0.2em] hover:bg-sage hover:text-cream transition-colors duration-200"
                >
                  Add to Pantry
                </button>
              </div>
            )}
          </div>
        </div>
      </motion.div>
    </motion.div>
  )
}

export default function DashboardPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedItem, setSelectedItem] = useState<FoodItem | null>(null)
  const [pantryItems, setPantryItems] = useState<FoodItem[]>(mockPantryItems)
  
  const { results: searchResults, isLoading, isSearching } = useDebouncedSearch(searchQuery)

  // Determine which items to display
  const displayItems = isSearching ? searchResults : pantryItems
  const isShowingSearchResults = isSearching

  const handleAddToPantry = (item: FoodItem) => {
    // Check if item already in pantry
    const exists = pantryItems.some((p) => p.id === item.id)
    if (exists) {
      toast.info(`${item.product_name} is already in your pantry`)
      return
    }

    setPantryItems((prev) => [...prev, item])
    toast.success(`${item.product_name} added to pantry`)
    setSelectedItem(null)
  }

  const isItemInPantry = (item: FoodItem) => {
    return pantryItems.some((p) => p.id === item.id)
  }

  return (
    <>
      <div className="border-b border-espresso/10">
        <div className="px-8 py-6">
          <h1 className="font-serif text-2xl text-espresso">Dashboard</h1>
        </div>
      </div>
      <div className="p-8">
        <SearchBar value={searchQuery} onChange={setSearchQuery} />

        <div className="mt-8">
          <p className="font-sans text-xs uppercase tracking-[0.2em] text-espresso/40 mb-8">
            {isShowingSearchResults ? (
              isLoading ? (
                'Searching...'
              ) : (
                `${displayItems.length} ${displayItems.length === 1 ? 'result' : 'results'} for "${searchQuery}"`
              )
            ) : (
              `${displayItems.length} ${displayItems.length === 1 ? 'item' : 'items'} in pantry`
            )}
          </p>

          {displayItems.length === 0 && !isLoading ? (
            <div className="py-16 text-center">
              <p className="font-serif text-xl text-espresso/50">
                {isShowingSearchResults
                  ? `No results found for "${searchQuery}"`
                  : 'Your pantry is empty'}
              </p>
            </div>
          ) : (
            <FoodGrid items={displayItems} onSelect={setSelectedItem} />
          )}
        </div>
      </div>

      <AnimatePresence>
        {selectedItem && (
          <DetailModal
            item={selectedItem}
            onClose={() => setSelectedItem(null)}
            onAddToPantry={handleAddToPantry}
            isInPantry={isItemInPantry(selectedItem)}
          />
        )}
      </AnimatePresence>
    </>
  )
}
