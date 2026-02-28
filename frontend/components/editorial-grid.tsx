'use client'

import { useRef } from 'react'
import { motion, useInView } from 'framer-motion'
import Image from 'next/image'
import type { FoodItem } from '@/lib/food-data'

function GridItem({
  item,
  index,
  size = 'default',
}: {
  item: FoodItem
  index: number
  size?: 'large' | 'default' | 'small'
}) {
  const ref = useRef<HTMLDivElement>(null)
  const isInView = useInView(ref, { once: true, margin: '-80px' })

  const sizeClasses = {
    large: 'col-span-2 row-span-2',
    default: 'col-span-1 row-span-1',
    small: 'col-span-1 row-span-1',
  }

  return (
    <motion.div
      ref={ref}
      className={`${sizeClasses[size]} group relative`}
      initial={{ opacity: 0, y: 50 }}
      animate={isInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 50 }}
      transition={{
        duration: 0.7,
        delay: index * 0.08,
        ease: [0.22, 1, 0.36, 1],
      }}
    >
      <div className="border border-foreground/10 hover:border-foreground/40 transition-colors duration-500 h-full">
        <div className="relative aspect-square flex items-center justify-center p-8 md:p-12 bg-card">
          <Image
            src={item.src}
            alt={item.label}
            width={600}
            height={600}
            className="w-full h-auto max-w-[75%] object-contain select-none transition-transform duration-700 group-hover:scale-105"
          />
        </div>
        <div className="border-t border-foreground/10 px-4 py-3 md:px-6 md:py-4 flex items-baseline justify-between">
          <span className="font-sans text-xs md:text-sm tracking-[0.15em] uppercase text-foreground">
            {item.label}
          </span>
          <span className="font-sans text-[10px] md:text-xs tracking-[0.1em] uppercase text-muted-foreground">
            {item.category}
          </span>
        </div>
      </div>
    </motion.div>
  )
}

export function EditorialGrid({ items }: { items: FoodItem[] }) {
  return (
    <section className="px-4 md:px-8 lg:px-16 py-16 md:py-24">
      {/* Section header */}
      <div className="border-b border-foreground pb-4 mb-12 md:mb-16 flex items-end justify-between">
        <h2 className="font-serif text-3xl md:text-5xl lg:text-6xl tracking-[-0.02em] text-foreground">
          The Index
        </h2>
        <span className="font-sans text-[10px] md:text-xs tracking-[0.2em] uppercase text-muted-foreground">
          {items.length} Items
        </span>
      </div>

      {/* Asymmetric grid */}
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-px bg-foreground/5">
        {items.map((item, index) => (
          <GridItem
            key={item.id}
            item={item}
            index={index}
            size={index === 0 || index === 5 ? 'large' : 'default'}
          />
        ))}
      </div>
    </section>
  )
}
