'use client'

import Image from 'next/image'
import { motion } from 'framer-motion'
import { FoodItem } from '@/lib/food-api'

interface SearchResultItemProps {
  item: FoodItem
  onSelect: (item: FoodItem) => void
  index: number
}

function SearchResultItem({ item, onSelect, index }: SearchResultItemProps) {
  return (
    <motion.button
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.2, delay: index * 0.05 }}
      onClick={() => onSelect(item)}
      className="group w-full flex items-center gap-6 py-4 px-4 border-b border-espresso/10 hover:bg-espresso/5 transition-colors duration-200 text-left"
    >
      {/* Product Image */}
      <div className="relative w-16 h-16 flex-shrink-0 overflow-hidden bg-espresso/5 rounded">
        <Image
          src={item.image_url}
          alt={item.product_name}
          fill
          className="object-contain"
          unoptimized
        />
      </div>

      {/* Product Info */}
      <div className="flex-1 min-w-0">
        <p className="font-sans font-bold text-sm text-espresso truncate">
          {item.product_name}
        </p>
        <p className="font-sans text-xs text-espresso/50 uppercase tracking-wider mt-1">
          {item.category}
        </p>
      </div>

      {/* Scores */}
      <div className="flex items-center gap-6 flex-shrink-0">
        <div className="text-center">
          <p className="font-sans text-[10px] uppercase tracking-wider text-espresso/40">
            Env
          </p>
          <p className="font-serif text-lg text-sage">
            {item.environmental_score}
          </p>
        </div>
        <div className="text-center">
          <p className="font-sans text-[10px] uppercase tracking-wider text-espresso/40">
            Nutri
          </p>
          <p className="font-serif text-lg text-sage">
            {item.nutriscore_score}
          </p>
        </div>
      </div>

      {/* Arrow indicator */}
      <div className="flex-shrink-0 text-espresso/30 group-hover:text-espresso transition-colors">
        <svg
          width="16"
          height="16"
          viewBox="0 0 16 16"
          fill="none"
          className="transform group-hover:translate-x-1 transition-transform"
        >
          <path
            d="M6 3L11 8L6 13"
            stroke="currentColor"
            strokeWidth="1.5"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        </svg>
      </div>
    </motion.button>
  )
}

interface SearchResultsProps {
  items: FoodItem[]
  onSelect: (item: FoodItem) => void
  isLoading?: boolean
}

export function SearchResults({ items, onSelect, isLoading }: SearchResultsProps) {
  if (isLoading) {
    return (
      <div className="space-y-4">
        {[...Array(5)].map((_, i) => (
          <div
            key={i}
            className="flex items-center gap-6 py-4 px-4 border-b border-espresso/10 animate-pulse"
          >
            <div className="w-16 h-16 bg-espresso/10 rounded" />
            <div className="flex-1 space-y-2">
              <div className="h-4 bg-espresso/10 rounded w-3/4" />
              <div className="h-3 bg-espresso/10 rounded w-1/4" />
            </div>
            <div className="flex gap-6">
              <div className="w-10 h-10 bg-espresso/10 rounded" />
              <div className="w-10 h-10 bg-espresso/10 rounded" />
            </div>
          </div>
        ))}
      </div>
    )
  }

  if (items.length === 0) {
    return null
  }

  return (
    <div className="border-t border-espresso/10">
      {items.map((item, index) => (
        <SearchResultItem
          key={item.id}
          item={item}
          onSelect={onSelect}
          index={index}
        />
      ))}
    </div>
  )
}
