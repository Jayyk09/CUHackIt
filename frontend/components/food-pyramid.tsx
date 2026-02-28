'use client'

import { useRef } from 'react'
import { motion, useInView } from 'framer-motion'
import Image from 'next/image'
import { getFoodByName, FoodItem } from '@/lib/food-data'

// Images are now LARGE and OVERLAPPING — like a pile, not floating dots
const tiers = [
  {
    number: '01',
    title: 'Protein & Fats',
    description:
      'Meat, seafood, eggs, full-fat dairy, nuts, seeds, avocados, olive oil. The foundation of cellular repair and sustained energy.',
    images: [
      // Steak — center-left, large, dominant
      {
        item: getFoodByName('steak') || { id: '1', src: '/placeholder.png', label: 'steak' },
        className: 'w-52 md:w-64 lg:w-72 top-[25%] left-[5%]',
        rotate: -6,
      },
      // Chicken — top right, bleeds out of frame
      {
        item: getFoodByName('chicken') || { id: '43', src: '/placeholder.png', label: 'chicken' },
        className: 'w-56 md:w-72 lg:w-80 top-[-8%] left-[40%]',
        rotate: 8,
      },
      // Salmon — overlapping steak bottom area
      {
        item: getFoodByName('salmon') || { id: '3', src: '/placeholder.png', label: 'salmon' },
        className: 'w-44 md:w-56 lg:w-64 top-[55%] left-[28%]',
        rotate: -4,
      },
      // Avocado — far right, partially cut off
      {
        item: getFoodByName('avocado') || { id: '4', src: '/placeholder.png', label: 'avocado' },
        className: 'w-36 md:w-44 lg:w-52 top-[30%] left-[72%]',
        rotate: 12,
      },
      // Eggs — bottom left cluster
      {
        item: getFoodByName('eggs') || { id: '2', src: '/placeholder.png', label: 'eggs' },
        className: 'w-40 md:w-52 lg:w-60 top-[70%] left-[55%]',
        rotate: -3,
      },
    ],
  },
  {
    number: '02',
    title: 'Vegetables & Fruits',
    description:
      'Whole, colorful, minimally processed. Fiber, vitamins, and phytonutrients that regulate everything.',
    images: [
      {
        item: getFoodByName('broccoli') || { id: '5', src: '/placeholder.png', label: 'broccoli' },
        className: 'w-48 md:w-60 lg:w-72 top-[10%] left-[5%]',
        rotate: -8,
      },
      {
        item: getFoodByName('tomatoes') || { id: '6', src: '/placeholder.png', label: 'tomatoes' },
        className: 'w-44 md:w-56 lg:w-64 top-[45%] left-[30%]',
        rotate: 5,
      },
      {
        item: getFoodByName('blueberry') || { id: '7', src: '/placeholder.png', label: 'blueberry' },
        className: 'w-36 md:w-48 lg:w-56 top-[65%] left-[60%]',
        rotate: -5,
      },
      {
        item: getFoodByName('carrots') || { id: '8', src: '/placeholder.png', label: 'carrots' },
        className: 'w-40 md:w-52 lg:w-60 top-[5%] left-[55%]',
        rotate: 10,
      },
    ],
  },
  {
    number: '03',
    title: 'Whole Grains',
    description:
      'Oats, rice, true sourdough. Complex carbohydrates for sustained fuel. Refined carbs discouraged.',
    images: [
      {
        item: getFoodByName('bowl-oats') || { id: '9', src: '/placeholder.png', label: 'oats' },
        className: 'w-52 md:w-64 lg:w-72 top-[15%] left-[8%]',
        rotate: -5,
      },
      {
        item: getFoodByName('bread') || { id: '10', src: '/placeholder.png', label: 'bread' },
        className: 'w-48 md:w-60 lg:w-68 top-[50%] left-[40%]',
        rotate: 7,
      },
      {
        item: getFoodByName('rice-beans') || { id: '11', src: '/placeholder.png', label: 'rice' },
        className: 'w-40 md:w-52 lg:w-60 top-[10%] left-[65%]',
        rotate: -10,
      },
    ],
  },
]

const avoid = [
  'Packaged salty/sweet snacks with added sugars and salt',
  'Foods with artificial flavors, petroleum-based dyes, or non-nutritive sweeteners',
  'Sugar-sweetened beverages: sodas, fruit drinks, energy drinks',
]

interface TierImage {
  item: FoodItem | any
  className: string
  rotate: number
}

interface TierData {
  number: string
  title: string
  description: string
  images: TierImage[]
}

function TierRow({ tier, index }: { tier: TierData; index: number }) {
  const ref = useRef<HTMLDivElement>(null)
  const isInView = useInView(ref, { once: true, margin: '-80px' })

  return (
    <div ref={ref} className="border-b border-[#282622]/10 last:border-b-0">
      {/* 
        Key layout change: 
        - Right column is overflow-visible so images can bleed out
        - min-h is much taller to accommodate large images
        - Images use absolute positioning relative to the right column
      */}
      <div className="grid grid-cols-1 lg:grid-cols-[40%_60%] gap-0 py-20 md:py-28">
        
        {/* LEFT: Typography */}
        <motion.div
          className="flex flex-col gap-5 lg:pr-16 z-10 self-center"
          initial={{ opacity: 0, y: 30 }}
          animate={isInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 30 }}
          transition={{ duration: 0.8, delay: 0.1, ease: [0.22, 1, 0.36, 1] }}
        >
          <span className="font-sans text-xs tracking-[0.25em] uppercase text-[#282622]/50">
            Tier {tier.number}
          </span>
          <h3 className="font-serif text-5xl md:text-6xl lg:text-7xl tracking-[-0.03em] text-[#282622] leading-[1.05]">
            {tier.title}
          </h3>
          <p className="font-sans text-sm md:text-base leading-relaxed text-[#282622]/70 max-w-sm mt-3">
            {tier.description}
          </p>
        </motion.div>

        {/* RIGHT: Large overlapping food pile */}
        {/* overflow-visible is critical — images bleed beyond this container */}
        <div className="relative w-full min-h-[420px] md:min-h-[500px] overflow-visible">
          {tier.images.map((img, i) => (
            <motion.div
              key={img.item.id || i}
              className={`absolute ${img.className}`}
              style={{ rotate: img.rotate }}
              initial={{ opacity: 0, scale: 0.75, y: 40 }}
              animate={
                isInView
                  ? { opacity: 1, scale: 1, y: 0 }
                  : { opacity: 0, scale: 0.75, y: 40 }
              }
              transition={{
                duration: 0.9,
                delay: 0.2 + i * 0.12,
                ease: [0.22, 1, 0.36, 1],
              }}
            >
              {/* Subtle float — gentle, desynchronized per-image */}
              <motion.div
                animate={isInView ? { y: [0, -6, 0] } : {}}
                transition={{
                  duration: 3.5 + i * 0.7,
                  repeat: Infinity,
                  ease: 'easeInOut',
                  delay: i * 0.5,
                }}
              >
                <Image
                  src={img.item.src}
                  alt={img.item.label}
                  width={400}
                  height={400}
                  className="w-full h-auto select-none drop-shadow-md"
                  style={{ objectFit: 'contain' }}
                />
              </motion.div>
            </motion.div>
          ))}
        </div>
      </div>
    </div>
  )
}

export function FoodPyramid() {
  const headerRef = useRef<HTMLDivElement>(null)
  const headerInView = useInView(headerRef, { once: true, margin: '-60px' })
  const avoidRef = useRef<HTMLDivElement>(null)
  const avoidInView = useInView(avoidRef, { once: true, margin: '-60px' })

  return (
    <section className="bg-[#F4F3EF] text-[#282622] overflow-x-hidden">
      <div className="px-6 md:px-12 lg:px-20 xl:px-32 py-20 md:py-32">

        {/* HEADER */}
        <motion.div
          ref={headerRef}
          className="border-b border-[#282622]/20 pb-12 mb-0"
          initial={{ opacity: 0, y: 30 }}
          animate={headerInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 30 }}
          transition={{ duration: 0.8, ease: [0.22, 1, 0.36, 1] }}
        >
          <span className="font-sans text-xs tracking-[0.25em] uppercase text-[#282622]/50 block mb-6">
            The Framework
          </span>
          <h2 className="font-serif text-5xl md:text-7xl lg:text-[5rem] tracking-[-0.03em] text-[#282622]">
            The New Food Pyramid
          </h2>
        </motion.div>

        {/* TIERS */}
        <div>
          {tiers.map((tier, index) => (
            <TierRow key={tier.number} tier={tier} index={index} />
          ))}
        </div>

        {/* AVOID */}
        <motion.div
          ref={avoidRef}
          className="mt-24 md:mt-32 pt-20 border-t border-[#282622]/10"
          initial={{ opacity: 0, y: 40 }}
          animate={avoidInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 40 }}
          transition={{ duration: 0.8, ease: [0.22, 1, 0.36, 1] }}
        >
          <div className="grid grid-cols-1 lg:grid-cols-[40%_60%] gap-16 lg:gap-0">
            <div>
              <span className="font-sans text-xs tracking-[0.25em] uppercase text-[#282622]/50 block mb-4">
                What to Avoid
              </span>
              <h3 className="font-serif text-4xl md:text-5xl lg:text-6xl tracking-[-0.02em] text-[#282622]">
                Highly Processed
              </h3>
            </div>
            <div className="flex flex-col">
              {avoid.map((item, i) => (
                <motion.div
                  key={i}
                  className="border-t border-[#282622]/10 py-8 first:border-t-0 first:pt-0"
                  initial={{ opacity: 0, x: 20 }}
                  animate={avoidInView ? { opacity: 1, x: 0 } : { opacity: 0, x: 20 }}
                  transition={{
                    duration: 0.5,
                    delay: 0.2 + i * 0.1,
                    ease: [0.22, 1, 0.36, 1],
                  }}
                >
                  <div className="flex gap-8 items-start">
                    <span className="font-sans text-sm tracking-[0.2em] text-[#282622]/40 pt-1">
                      {String(i + 1).padStart(2, '0')}
                    </span>
                    <p className="font-sans text-base md:text-lg leading-relaxed text-[#282622]/80 max-w-lg">
                      {item}
                    </p>
                  </div>
                </motion.div>
              ))}
            </div>
          </div>
        </motion.div>

      </div>
    </section>
  )
}
