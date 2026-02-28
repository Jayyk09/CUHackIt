'use client'

import { motion } from 'framer-motion'
import Image from 'next/image'
import { getFoodByName } from '@/lib/food-data'

const heroImages = [
  { item: getFoodByName('canned-tuna')!, className: 'top-[8%] left-[4%] w-28 md:w-40 lg:w-48', rotate: -6 },
  { item: getFoodByName('steak')!, className: 'top-[2%] right-[12%] w-32 md:w-44 lg:w-52', rotate: 4 },
  { item: getFoodByName('eggs')!, className: 'top-[38%] left-[8%] md:left-[12%] w-24 md:w-32 lg:w-40', rotate: 8 },
  { item: getFoodByName('blueberry')!, className: 'bottom-[22%] right-[6%] w-20 md:w-28 lg:w-32', rotate: -3 },
  { item: getFoodByName('broccoli')!, className: 'bottom-[8%] left-[2%] w-28 md:w-36 lg:w-44', rotate: 5 },
  { item: getFoodByName('salmon')!, className: 'top-[12%] left-[38%] md:left-[28%] w-24 md:w-32 lg:w-36', rotate: -8 },
  { item: getFoodByName('avocado')!, className: 'bottom-[14%] right-[28%] md:right-[22%] w-24 md:w-32 lg:w-40', rotate: 3 },
]

export function Hero() {
  return (
    <section className="relative min-h-screen flex items-center justify-center overflow-hidden border-b border-foreground">
      {heroImages.map((img, i) => (
        <motion.div
          key={img.item.id}
          className={`absolute ${img.className} z-0`}
          initial={{ opacity: 0, y: 40, rotate: 0 }}
          animate={{
            opacity: 1,
            y: [0, -10, 0],
            rotate: img.rotate,
          }}
          transition={{
            opacity: { duration: 0.9, delay: 0.15 + i * 0.1, ease: [0.22, 1, 0.36, 1] },
            rotate: { duration: 0.9, delay: 0.15 + i * 0.1, ease: [0.22, 1, 0.36, 1] },
            y: {
              duration: 4 + i * 0.5,
              repeat: Infinity,
              ease: 'easeInOut',
              delay: 0.9 + i * 0.3,
            },
          }}
        >
          <Image
            src={img.item.src}
            alt={img.item.label}
            width={400}
            height={400}
            className="w-full h-auto select-none pointer-events-none"
            priority={i < 3}
          />
        </motion.div>
      ))}

      <div className="relative z-10 flex flex-col items-center text-center px-6">
        <motion.h1
          className="font-serif text-[clamp(4rem,15vw,14rem)] leading-[0.85] tracking-[-0.03em] text-foreground"
          initial={{ opacity: 0, y: 30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8, delay: 0.05, ease: [0.22, 1, 0.36, 1] }}
        >
          Sift
        </motion.h1>

        <motion.p
          className="mt-4 md:mt-6 font-sans text-sm md:text-base tracking-[0.2em] uppercase text-muted-foreground max-w-md"
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7, delay: 0.3, ease: [0.22, 1, 0.36, 1] }}
        >
          Sift through the noise.
        </motion.p>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.7, delay: 0.5, ease: [0.22, 1, 0.36, 1] }}
          className="mt-8 md:mt-10"
        >
          <a
            href="#facts"
            className="border border-accent bg-accent text-accent-foreground px-8 py-3 font-sans text-xs tracking-[0.25em] uppercase hover:bg-background hover:text-accent transition-colors duration-300 cursor-pointer inline-block"
          >
            Learn Why It Matters
          </a>
        </motion.div>
      </div>
    </section>
  )
}
