'use client'

import { useRef } from 'react'
import { motion, useInView } from 'framer-motion'

const SAGE = '#4A5D4E'
const ESPRESSO = '#2A2724'

interface HealthStat {
  percentage: string
  value: number
  description: string
  graphType: 'ring' | 'bars' | 'arc' | 'blocks'
}

const healthStats: HealthStat[] = [
  {
    percentage: '50%',
    value: 50,
    description: 'of Americans have prediabetes or diabetes. Half the country is metabolically broken before they feel a single symptom.',
    graphType: 'ring',
  },
  {
    percentage: '75%',
    value: 75,
    description: 'of adults have at least one chronic condition. Three in four people are fighting something that food could help prevent.',
    graphType: 'bars',
  },
  {
    percentage: '90%',
    value: 90,
    description: 'of U.S. healthcare spending goes toward treating chronic disease, much of which is diet-linked. We treat symptoms, not causes.',
    graphType: 'arc',
  },
  {
    percentage: '70%',
    value: 70,
    description: "of an American child's diet is ultra-processed. In many other countries, that figure is below 20%.",
    graphType: 'blocks',
  },
]

function RingChart({ value, inView }: { value: number; inView: boolean }) {
  const size = 240
  const strokeWidth = 1.5
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
        stroke={`${ESPRESSO}1A`}
        strokeWidth={strokeWidth}
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
        transition={{ duration: 1.8, delay: 0.3, ease: [0.22, 1, 0.36, 1] }}
      />
    </svg>
  )
}

function BarsChart({ value, inView }: { value: number; inView: boolean }) {
  const bars = [
    { h: value * 0.6, delay: 0.3 },
    { h: value * 0.8, delay: 0.4 },
    { h: value, delay: 0.5 },
    { h: value * 0.9, delay: 0.6 },
    { h: value * 0.7, delay: 0.7 },
  ]

  const maxHeight = 180

  return (
    <svg width={240} height={200} viewBox="0 0 240 200">
      {bars.map((bar, i) => {
        const barHeight = (bar.h / 100) * maxHeight
        return (
          <motion.rect
            key={i}
            x={i * 48}
            y={maxHeight - barHeight + 10}
            width={36}
            height={barHeight}
            fill={SAGE}
            rx={1}
            initial={{ height: 0, y: maxHeight + 10 }}
            animate={inView ? { height: barHeight, y: maxHeight - barHeight + 10 } : { height: 0, y: maxHeight + 10 }}
            transition={{ duration: 0.8, delay: bar.delay, ease: [0.22, 1, 0.36, 1] }}
          />
        )
      })}
    </svg>
  )
}

function ArcChart({ value, inView }: { value: number; inView: boolean }) {
  const width = 240
  const height = 140
  const strokeWidth = 1.5
  const radius = 110
  const circumference = Math.PI * radius
  const offset = circumference - (value / 100) * circumference

  return (
    <svg width={width} height={height} className="overflow-visible">
      <path
        d={`M 10 ${height - 10} A ${radius} ${radius} 0 0 1 ${width - 10} ${height - 10}`}
        fill="none"
        stroke={`${ESPRESSO}1A`}
        strokeWidth={strokeWidth}
      />
      <motion.path
        d={`M 10 ${height - 10} A ${radius} ${radius} 0 0 1 ${width - 10} ${height - 10}`}
        fill="none"
        stroke={SAGE}
        strokeWidth={strokeWidth}
        strokeLinecap="round"
        strokeDasharray={circumference}
        initial={{ strokeDashoffset: circumference }}
        animate={inView ? { strokeDashoffset: offset } : { strokeDashoffset: circumference }}
        transition={{ duration: 1.6, delay: 0.3, ease: [0.22, 1, 0.36, 1] }}
      />
    </svg>
  )
}

function BlocksChart({ value, inView }: { value: number; inView: boolean }) {
  const totalBlocks = 10
  const filledBlocks = Math.round((value / 100) * totalBlocks)

  return (
    <div className="flex gap-3" style={{ width: 240, justifyContent: 'center' }}>
      {Array.from({ length: totalBlocks }).map((_, i) => (
        <motion.div
          key={i}
          className="h-32"
          style={{ 
            width: 16,
            backgroundColor: i < filledBlocks ? SAGE : `${ESPRESSO}1A`,
            borderRadius: 1
          }}
          initial={{ opacity: 0, scaleY: 0 }}
          animate={inView ? { opacity: 1, scaleY: 1 } : { opacity: 0, scaleY: 0 }}
          transition={{ duration: 0.5, delay: 0.3 + i * 0.06, ease: [0.22, 1, 0.36, 1] }}
        />
      ))}
    </div>
  )
}

function ChartDisplay({ type, value, inView }: { type: HealthStat['graphType']; value: number; inView: boolean }) {
  switch (type) {
    case 'ring':
      return <RingChart value={value} inView={inView} />
    case 'bars':
      return <BarsChart value={value} inView={inView} />
    case 'arc':
      return <ArcChart value={value} inView={inView} />
    case 'blocks':
      return <BlocksChart value={value} inView={inView} />
  }
}

function StatRow({ stat, index }: { stat: HealthStat; index: number }) {
  const ref = useRef<HTMLDivElement>(null)
  const isInView = useInView(ref, { once: true, margin: '-100px' })
  const isEven = index % 2 === 0

  return (
    <div
      ref={ref}
      className="border-b border-foreground/10 last:border-b-0 py-24"
    >
      <div className={`flex flex-col md:flex-row ${!isEven ? 'md:flex-row-reverse' : ''} gap-8 md:gap-16`}>
        <motion.div
          className="flex-1 flex flex-col justify-center"
          initial={{ opacity: 0, y: 40 }}
          animate={isInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 40 }}
          transition={{ duration: 0.8, delay: 0.1, ease: [0.22, 1, 0.36, 1] }}
        >
          <span className="font-sans text-xs tracking-[0.25em] uppercase opacity-30 block mb-6">
            {String(index + 1).padStart(2, '0')}
          </span>
          
          <h3 className="font-serif text-7xl md:text-8xl tracking-[-0.03em] text-foreground leading-none mb-6">
            {stat.percentage}
          </h3>

          <p className="font-sans text-base leading-relaxed opacity-70 max-w-lg" style={{ color: ESPRESSO }}>
            {stat.description}
          </p>
        </motion.div>

        <motion.div
          className="flex-1 flex items-center justify-center"
          initial={{ opacity: 0, scale: 0.8 }}
          animate={isInView ? { opacity: 1, scale: 1 } : { opacity: 0, scale: 0.8 }}
          transition={{ duration: 0.8, delay: 0.2, ease: [0.22, 1, 0.36, 1] }}
        >
          <ChartDisplay type={stat.graphType} value={stat.value} inView={isInView} />
        </motion.div>
      </div>
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
        className="border-b border-foreground pb-6 mb-0"
        initial={{ opacity: 0, y: 30 }}
        animate={headerInView ? { opacity: 1, y: 0 } : { opacity: 0, y: 30 }}
        transition={{ duration: 0.8, ease: [0.22, 1, 0.36, 1] }}
      >
        <h2 className="font-serif text-4xl md:text-6xl lg:text-7xl tracking-[-0.02em] text-foreground">
          The State of American Health
        </h2>
      </motion.div>

      <div>
        {healthStats.map((stat, index) => (
          <StatRow key={index} stat={stat} index={index} />
        ))}
      </div>
    </section>
  )
}
