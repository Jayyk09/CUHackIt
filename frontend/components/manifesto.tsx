'use client'

import { useRef } from 'react'
import { motion, useInView } from 'framer-motion'

const statements = [
  { number: '01', text: 'Every ingredient has a story. We make it legible.' },
  { number: '02', text: 'No proprietary scores. No sponsored placements. Just data.' },
  { number: '03', text: 'Built on USDA FoodData Central. Open. Verifiable. Yours.' },
]

export function Manifesto() {
  const ref = useRef<HTMLDivElement>(null)
  const isInView = useInView(ref, { once: true, margin: '-100px' })

  return (
    <section
      ref={ref}
      className="border-t border-foreground bg-accent text-accent-foreground"
    >
      <div className="px-4 md:px-8 lg:px-16 py-20 md:py-32">
        <motion.p
          className="font-serif text-2xl md:text-4xl lg:text-5xl leading-snug tracking-[-0.01em] max-w-4xl text-balance"
          initial={{ opacity: 0, y: 40 }}
          animate={isInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 40 }}
          transition={{ duration: 0.8, ease: [0.22, 1, 0.36, 1] }}
        >
          Sift through the noise. Real food data, presented without compromise.
        </motion.p>

        <div className="mt-16 md:mt-24 grid grid-cols-1 md:grid-cols-3 gap-0">
          {statements.map((s, i) => (
            <motion.div
              key={s.number}
              className="border-t border-accent-foreground/20 pt-6 pb-8 md:pr-8"
              initial={{ opacity: 0, y: 30 }}
              animate={isInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 30 }}
              transition={{
                duration: 0.6,
                delay: 0.3 + i * 0.12,
                ease: [0.22, 1, 0.36, 1],
              }}
            >
              <span className="font-sans text-[10px] tracking-[0.3em] uppercase text-accent-foreground/40 block mb-3">
                {s.number}
              </span>
              <p className="font-sans text-sm md:text-base leading-relaxed text-accent-foreground/80">
                {s.text}
              </p>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  )
}
