'use client'

import { useState } from 'react'
import Image from 'next/image'
import { motion, AnimatePresence } from 'framer-motion'

interface FoodItem {
  product_name: string
  environmental_score: number
  nutriscore_score: number
  labels_en: string[]
  allergens_en: string[]
  traces_en: string[]
  image_url: string
  shelf_life: number
  category: string
  quantity: number
  units: string
  is_spoiled: boolean
}

const mockFoodItems: FoodItem[] = [
  {
    product_name: 'Organic Atlantic Salmon',
    environmental_score: 72,
    nutriscore_score: 85,
    labels_en: ['Organic', 'Sustainable', 'Wild Caught'],
    allergens_en: ['Fish'],
    traces_en: ['Shellfish'],
    image_url: '/realfood/029_salmon.webp.png',
    shelf_life: 14,
    category: 'Protein',
    quantity: 500,
    units: 'g',
    is_spoiled: false,
  },
  {
    product_name: 'Heirloom Tomatoes',
    environmental_score: 91,
    nutriscore_score: 95,
    labels_en: ['Organic', 'Local', 'Seasonal'],
    allergens_en: [],
    traces_en: [],
    image_url: '/realfood/010_tomatoes.webp.png',
    shelf_life: 7,
    category: 'Vegetable',
    quantity: 1,
    units: 'kg',
    is_spoiled: false,
  },
  {
    product_name: 'Free-Range Eggs',
    environmental_score: 68,
    nutriscore_score: 78,
    labels_en: ['Free-Range', 'Farm Fresh', 'Omega-3'],
    allergens_en: ['Eggs'],
    traces_en: ['Milk'],
    image_url: '/realfood/002_eggs.webp.png',
    shelf_life: 21,
    category: 'Protein',
    quantity: 12,
    units: 'count',
    is_spoiled: false,
  },
  {
    product_name: 'Extra Virgin Olive Oil',
    environmental_score: 84,
    nutriscore_score: 92,
    labels_en: ['Cold Pressed', 'First Cold Press', 'DOP'],
    allergens_en: [],
    traces_en: [],
    image_url: '/realfood/036_olive-oil.webp.png',
    shelf_life: 365,
    category: 'Fats',
    quantity: 500,
    units: 'ml',
    is_spoiled: false,
  },
  {
    product_name: 'Raw Honeycomb',
    environmental_score: 88,
    nutriscore_score: 90,
    labels_en: ['Raw', 'Unfiltered', 'Local'],
    allergens_en: [],
    traces_en: [],
    image_url: '/realfood/023_bowl-oats.webp.png',
    shelf_life: 730,
    category: 'Pantry',
    quantity: 350,
    units: 'g',
    is_spoiled: false,
  },
  {
    product_name: 'Aged Cheddar Cheese',
    environmental_score: 58,
    nutriscore_score: 65,
    labels_en: ['Aged', 'Artisanal', 'Grass-Fed'],
    allergens_en: ['Milk'],
    traces_en: ['Lactose'],
    image_url: '/realfood/006_cheese.webp.png',
    shelf_life: 180,
    category: 'Dairy',
    quantity: 250,
    units: 'g',
    is_spoiled: false,
  },
]

function SearchBar({
  value,
  onChange,
}: {
  value: string
  onChange: (value: string) => void
}) {
  return (
    <div className="relative w-full border-b border-espresso pb-6">
      <div className="flex items-baseline gap-4">
        <input
          type="text"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder="search recipe or ingredients"
          className="w-full bg-transparent font-serif text-2xl md:text-3xl text-espresso placeholder:text-espresso/40 placeholder:font-serif focus:outline-none"
        />
        <svg
          width="28"
          height="28"
          viewBox="0 0 24 24"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
          className="flex-shrink-0 text-sage"
        >
          <path
            d="M17 8C17 10.7614 14.7614 13 12 13C9.23858 13 7 10.7614 7 8C7 5.23858 9.23858 3 12 3C14.7614 3 17 5.23858 17 8Z"
            stroke="currentColor"
            strokeWidth="1.5"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
          <path
            d="M20.5 20.5L16 16"
            stroke="currentColor"
            strokeWidth="1.5"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
          <path
            d="M12 7C12.8284 7 13.5 6.32843 13.5 5.5C13.5 4.67157 12.8284 4 12 4C11.1716 4 10.5 4.67157 10.5 5.5C10.5 6.32843 11.1716 7 12 7Z"
            fill="currentColor"
          />
          <path
            d="M8.5 11.5C9.05228 11.5 9.5 11.0523 9.5 10.5C9.5 9.94772 9.05228 9.5 8.5 9.5C7.94772 9.5 7.5 9.94772 7.5 10.5C7.5 11.0523 7.94772 11.5 8.5 11.5Z"
            fill="currentColor"
          />
        </svg>
      </div>
    </div>
  )
}

function FoodGrid({
  items,
  onSelect,
}: {
  items: FoodItem[]
  onSelect: (item: FoodItem) => void
}) {
  return (
    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-x-6 gap-y-12">
      {items.map((item) => (
        <button
          key={item.product_name}
          onClick={() => onSelect(item)}
          className="group text-left w-full focus:outline-none"
        >
          <div className="relative aspect-square mb-4 overflow-hidden">
            <Image
              src={item.image_url}
              alt={item.product_name}
              fill
              className="object-contain transition-transform duration-500 group-hover:scale-105"
            />
          </div>
          <div className="border-t border-espresso pt-3">
            <p className="font-sans font-bold text-sm text-espresso truncate">
              {item.product_name}
            </p>
            <p className="font-sans text-xs text-espresso/50 mt-1">
              {item.quantity} {item.units}
            </p>
          </div>
        </button>
      ))}
    </div>
  )
}

function DetailSlideOver({
  item,
  onClose,
}: {
  item: FoodItem
  onClose: () => void
}) {
  return (
    <motion.div
      initial={{ x: '100%' }}
      animate={{ x: 0 }}
      exit={{ x: '100%' }}
      transition={{ type: 'spring', damping: 30, stiffness: 300 }}
      className="fixed right-0 top-0 h-screen w-full md:w-[480px] bg-cream border-l border-espresso/10 z-50 overflow-y-auto"
    >
      <div className="p-8 md:p-12">
        <button
          onClick={onClose}
          className="absolute top-8 right-8 text-xs uppercase tracking-[0.2em] font-sans text-espresso/50 hover:text-espresso transition-colors"
        >
          — Close
        </button>

        <div className="relative aspect-square w-full mb-10 overflow-hidden">
          <Image
            src={item.image_url}
            alt={item.product_name}
            fill
            className="object-contain"
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
                <p className="font-sans text-sm text-espresso">{item.category}</p>
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
                <p className="font-sans text-xs uppercase tracking-[0.2em] text-espresso/50 mb-1">
                  Quantity
                </p>
                <p className="font-sans text-sm text-espresso">
                  {item.quantity} {item.units}
                </p>
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
        </div>
      </div>
    </motion.div>
  )
}

export default function DashboardPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedItem, setSelectedItem] = useState<FoodItem | null>(null)

  const filteredItems = mockFoodItems.filter((item) =>
    item.product_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    item.category.toLowerCase().includes(searchQuery.toLowerCase()) ||
    item.labels_en.some((label) =>
      label.toLowerCase().includes(searchQuery.toLowerCase())
    )
  )

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
            {filteredItems.length} {filteredItems.length === 1 ? 'item' : 'items'}
          </p>

          {filteredItems.length === 0 ? (
            <div className="py-16 text-center">
              <p className="font-serif text-xl text-espresso/50">
                No items found for "{searchQuery}"
              </p>
            </div>
          ) : (
            <>
              <AnimatePresence mode="wait">
                {selectedItem ? (
                  <motion.div
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    exit={{ opacity: 0 }}
                    className="fixed inset-0 bg-espresso/5 z-40"
                    onClick={() => setSelectedItem(null)}
                  />
                ) : null}
              </AnimatePresence>

              <FoodGrid items={filteredItems} onSelect={setSelectedItem} />
            </>
          )}
        </div>
      </div>

      <AnimatePresence>
        {selectedItem && (
          <DetailSlideOver
            item={selectedItem}
            onClose={() => setSelectedItem(null)}
          />
        )}
      </AnimatePresence>
    </>
  )
}
