'use client'

import { useRef } from 'react'
import { motion, useInView } from 'framer-motion'

const SAGE = '#4A5D4E'

interface FactRow {
  stat: string
  statValue: number
  description: string
  graphType: 'ring' | 'bars' | 'arc' | 'blocks'
}

const healthFacts: FactRow[] = [
  {
    stat: '50%',
    statValue: 50,
    description: 'of Americans have prediabetes or diabetes. Half the country is metabolically broken before they feel a single symptom.',
    graphType: 'ring',
  },
  {
    stat: '75%',
    statValue: 75,
    description: 'of adults have at least one chronic condition. Three in four people are fighting something that food could help prevent.',
    graphType: 'bars',
  },
  {
    stat: '90%',
    statValue: 90,
    description: 'of U.S. healthcare spending goes toward treating chronic disease, much of which is diet-linked. We treat symptoms, not causes.',
    graphType: 'arc',
  },
  {
    stat: '70%',
    statValue: 70,
    description: "of an American child's diet is ultra-processed. In many other countries, that figure is below 20%.",
    graphType: 'blocks',
  },
]

function RingGraph({ value, inView }: { value: number; inView: boolean }) {
  const size = 100
  const strokeWidth = 3
  const radius = (size - strokeWidth) / 2
  const circumference = 2 * Math.PI * radius
  const offset = circumference - (value / 100) * circumference

  return (
    <svg width={size} height={size} className="transform -rotate-90">
      <circle
        cx={size / 2}
        cy={size / 2}
        r={radius}
        fill="none"
        stroke="currentColor"
        strokeWidth={strokeWidth}
        className="text-foreground/10"
      />
      <motion.circle
        cx={size / 2}
        cy={size / 2}
        r={radius}
        fill="none"
        stroke={SAGE}
        strokeWidth={strokeWidth}
        strokeLinecap="round"
        strokeDasharray={circumference}
        initial={{ strokeDashoffset: circumference }}
        animate={inView ? { strokeDashoffset: offset } : { strokeDashoffset: circumference }}
        transition={{ duration: 1.4, delay: 0.5, ease: [0.22, 1, 0.36, 1] }}
      />
    </svg>
  )
}

function BarsGraph({ value, inView }: { value: number; inView: boolean }) {
  const bars = [
    { h: value * 0.5, delay: 0.4 },
    { h: value * 0.7, delay: 0.5 },
    { h: value, delay: 0.6 },
    { h: value * 0.85, delay: 0.7 },
    { h: value * 0.6, delay: 0.8 },
  ]

  return (
    <svg width={80} height={60} className="overflow-visible">
      {bars.map((bar, i) => (
        <motion.rect
          key={i}
          x={i * 16}
          y={60}
          width={12}
          height={0}
          fill={SAGE}
          initial={{ height: 0, y: 60 }}
          animate={inView ? { height: (bar.h / 100) * 50, y: 60 - (bar.h / 100) * 50 } : { height: 0, y: 60 }}
          transition={{ duration: 0.7, delay: bar.delay, ease: [0.22, 1, 0.36, 1] }}
        />
      ))}
    </svg>
  )
}

function ArcGraph({ value, inView }: { value: number; inView: boolean }) {
  const width = 100
  const height = 50
  const strokeWidth = 3
  const radius = 45
  const circumference = Math.PI * radius
  const offset = circumference - (value / 100) * circumference

  return (
    <svg width={width} height={height + 5} className="overflow-visible">
      <path
        d={`M 5 ${height} A ${radius} ${radius} 0 0 1 ${width - 5} ${height}`}
        fill="none"
        stroke="currentColor"
        strokeWidth={strokeWidth}
        className="text-foreground/10"
      />
      <motion.path
        d={`M 5 ${height} A ${radius} ${radius} 0 0 1 ${width - 5} ${height}`}
        fill="none"
        stroke={SAGE}
        strokeWidth={strokeWidth}
        strokeLinecap="round"
        strokeDasharray={circumference}
        initial={{ strokeDashoffset: circumference }}
        animate={inView ? { strokeDashoffset: offset } : { strokeDashoffset: circumference }}
        transition={{ duration: 1.2, delay: 0.5, ease: [0.22, 1, 0.36, 1] }}
      />
    </svg>
  )
}

function BlocksGraph({ value, inView }: { value: number; inView: boolean }) {
  const totalBlocks = 10
  const filledBlocks = Math.round((value / 100) * totalBlocks)

  return (
    <div className="flex gap-1.5">
      {Array.from({ length: totalBlocks }).map((_, i) => (
        <motion.div
          key={i}
          className="w-2 h-8"
          style={{ backgroundColor: i < filledBlocks ? SAGE : 'currentColor' }}
          initial={{ opacity: 0, scaleY: 0 }}
          animate={inView ? { opacity: i < filledBlocks ? 1 : 0.1, scaleY: 1 } : { opacity: 0, scaleY: 0 }}
          transition={{ duration: 0.4, delay: 0.4 + i * 0.05, ease: [0.22, 1, 0.36, 1] }}
        />
      ))}
    </div>
  )
}

function GraphDisplay({ type, value, inView }: { type: FactRow['graphType']; value: number; inView: boolean }) {
  switch (type) {
    case 'ring':
      return <RingGraph value={value} inView={inView} />
    case 'bars':
      return <BarsGraph value={value} inView={inView} />
    case 'arc':
      return <ArcGraph value={value} inView={inView} />
    case 'blocks':
      return <BlocksGraph value={value} inView={inView} />
  }
}

function FactRow({ fact, index }: { fact: FactRow; index: number }) {
  const ref = useRef<HTMLDivElement>(null)
  const isInView = useInView(ref, { once: true, margin: '-100px' })

  return (
    <div
      ref={ref}
      className="border-b border-foreground/10 last:border-b-0 py-16 md:py-20"
    >
      <motion.div
        className="flex flex-col md:flex-row md:items-center gap-8 md:gap-12"
        initial={{ opacity: 0, y: 40 }}
        animate={isInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 40 }}
        transition={{ duration: 0.8, delay: 0.2, ease: [0.22, 1, 0.36, 1] }}
      >
        <div className="shrink-0">
          <GraphDisplay type={fact.graphType} value={fact.statValue} inView={isInView} />
        </div>

        <div className="flex-1">
          <span className="font-sans text-xs tracking-[0.25em] uppercase text-muted-foreground/50 block mb-4">
            {String(index + 1).padStart(2, '0')}
          </span>
          
          <div className="flex items-baseline gap-4 mb-4">
            <span className="font-serif text-5xl md:text-6xl lg:text-7xl tracking-[-0.03em] text-foreground leading-none">
              {fact.stat}
            </span>
          </div>

          <p className="font-sans text-sm md:text-base leading-relaxed text-muted-foreground max-w-lg">
            {fact.description}
          </p>
        </div>
      </motion.div>
    </div>
  )
}

export function HealthFacts() {
  const headerRef = useRef<HTMLDivElement>(null)
  const headerInView = useInView(headerRef, { once: true, margin: '-60px' })

  return (
    <section className="px-6 md:px-12 lg:px-20 xl:px-32 py-20 md:py-28">
      <motion.div
        ref={headerRef}
        className="border-b border-foreground pb-6 mb-0 flex flex-col md:flex-row md:items-end justify-between gap-4"
        initial={{ opacity: 0, y: 30 }}
        animate={headerInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 30 }}
        transition={{ duration: 0.8, ease: [0.22, 1, 0.36, 1] }}
      >
        <h2 className="font-serif text-4xl md:text-6xl lg:text-7xl tracking-[-0.02em] text-foreground">
          The State of American Health
        </h2>
        <span className="font-sans text-xs tracking-[0.25em] uppercase text-muted-foreground/50 shrink-0 pb-2">
          Why it matters
        </span>
      </motion.div>

      <div>
        {healthFacts.map((fact, index) => (
          <FactRow key={index} fact={fact} index={index} />
        ))}
      </div>
    </section>
  )
}
