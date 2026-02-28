"use client"

import { motion } from "framer-motion"
import Image from "next/image"
import { cn } from "@/lib/utils"
import type { FoodItem } from "@/lib/food-api"

interface SearchListProps {
  items: FoodItem[]
  onItemClick: (item: FoodItem) => void
  className?: string
}

interface SearchRowProps {
  item: FoodItem
  onClick: () => void
}

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.06,
      delayChildren: 0.1,
    },
  },
}

const rowVariants = {
  hidden: { 
    opacity: 0, 
    y: 8,
  },
  visible: { 
    opacity: 1, 
    y: 0,
    transition: {
      duration: 0.5,
      ease: [0.22, 1, 0.36, 1],
    },
  },
}

function SearchRow({ item, onClick }: SearchRowProps) {
  return (
    <motion.button
      variants={rowVariants}
      onClick={onClick}
      className={cn(
        "group w-full",
        "flex flex-row items-center justify-between",
        "py-5 px-0",
        "border-b border-[#2A2724]/10",
        "bg-transparent",
        "cursor-pointer",
        "text-left",
        "transition-colors duration-300 ease-out",
        "hover:bg-[#2A2724]/[0.02]",
        "focus:outline-none focus-visible:ring-1 focus-visible:ring-[#4A5D4E]/40 focus-visible:ring-offset-2 focus-visible:ring-offset-[#F9F8F6]"
      )}
    >
      {/* Left Column: Image + Text */}
      <div className="flex items-center gap-4">
        {/* Thumbnail */}
        <div 
          className={cn(
            "relative w-14 h-14 flex-shrink-0",
            "border border-[#2A2724]/15",
            "overflow-hidden",
            "bg-[#2A2724]/[0.03]"
          )}
        >
          {item.image_url ? (
            <Image
              src={item.image_url}
              alt={item.product_name}
              fill
              sizes="56px"
              className="object-cover"
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center">
              <span className="text-[10px] tracking-[0.15em] uppercase text-[#2A2724]/30">
                No img
              </span>
            </div>
          )}
        </div>

        {/* Text Content */}
        <div className="flex flex-col gap-1">
          <span 
            className={cn(
              "text-lg font-semibold text-[#2A2724] leading-tight",
              "transition-all duration-300 ease-out",
              "group-hover:text-[#4A5D4E]",
              "group-hover:translate-x-2"
            )}
            style={{ fontFamily: "var(--font-sans)" }}
          >
            {item.product_name}
          </span>
          <span 
            className={cn(
              "text-[10px] tracking-[0.2em] uppercase",
              "text-[#2A2724]/50",
              "transition-all duration-300 ease-out",
              "group-hover:translate-x-2"
            )}
          >
            {item.category}
          </span>
        </div>
      </div>

      {/* Right Column: Action */}
      <div 
        className={cn(
          "flex items-center gap-2",
          "transition-all duration-300 ease-out",
          "opacity-40 group-hover:opacity-100",
          "group-hover:translate-x-1"
        )}
      >
        <span 
          className={cn(
            "text-[9px] tracking-[0.25em] uppercase font-medium",
            "text-[#2A2724]/60",
            "transition-colors duration-300",
            "group-hover:text-[#4A5D4E]",
            "hidden sm:block"
          )}
        >
          View
        </span>
        <span 
          className={cn(
            "text-sm text-[#2A2724]/50",
            "transition-colors duration-300",
            "group-hover:text-[#4A5D4E]"
          )}
        >
          →
        </span>
      </div>
    </motion.button>
  )
}

export function SearchList({ items, onItemClick, className }: SearchListProps) {
  if (items.length === 0) {
    return (
      <div 
        className={cn(
          "py-16 text-center",
          className
        )}
      >
        <p className="text-[11px] tracking-[0.2em] uppercase text-[#2A2724]/40">
          No results found
        </p>
      </div>
    )
  }

  return (
    <motion.div
      className={cn(
        "w-full",
        "bg-[#F9F8F6]",
        className
      )}
      variants={containerVariants}
      initial="hidden"
      animate="visible"
    >
      {/* Optional Index Header */}
      <div className="flex items-center justify-between pb-3 mb-2 border-b border-[#2A2724]/20">
        <span className="text-[9px] tracking-[0.25em] uppercase text-[#2A2724]/40 font-medium">
          Product Index
        </span>
        <span className="text-[9px] tracking-[0.15em] uppercase text-[#2A2724]/30">
          {items.length} {items.length === 1 ? "item" : "items"}
        </span>
      </div>

      {/* List Container */}
      <div className="flex flex-col">
        {items.map((item) => (
          <SearchRow
            key={item.id}
            item={item}
            onClick={() => onItemClick(item)}
          />
        ))}
      </div>
    </motion.div>
  )
}

export default SearchList
