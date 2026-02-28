'use client';

import { motion } from 'framer-motion';
import { cn } from '@/lib/utils';

interface StackedSelectProps {
  options: string[];
  selected: string;
  onChange: (selected: string) => void;
}

export function StackedSelect({ options, selected, onChange }: StackedSelectProps) {
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.08,
        delayChildren: 0.1,
      },
    },
  };

  const itemVariants = {
    hidden: { opacity: 0, x: -20 },
    visible: {
      opacity: 1,
      x: 0,
      transition: {
        duration: 0.4,
        ease: [0.22, 1, 0.36, 1],
      },
    },
  };

  return (
    <motion.div
      variants={containerVariants}
      initial="hidden"
      animate="visible"
      className="w-full border border-[#2A2724]/20"
    >
      {options.map((option, index) => {
        const isSelected = selected === option;
        const isLast = index === options.length - 1;

        return (
          <motion.button
            key={option}
            type="button"
            variants={itemVariants}
            onClick={() => onChange(option)}
            whileTap={{ scale: 0.995 }}
            className={cn(
              'w-full px-6 py-5',
              'flex items-center justify-between',
              'font-sans text-base tracking-[0.02em]',
              'transition-all duration-200 ease-out',
              'cursor-pointer select-none',
              !isLast && 'border-b border-[#2A2724]/10',
              isSelected
                ? 'bg-[#4A5D4E] text-[#F9F8F6]'
                : 'bg-transparent text-[#2A2724]/70 hover:bg-[#2A2724]/[0.03] hover:text-[#2A2724]'
            )}
          >
            <span>{option}</span>
            {isSelected && (
              <motion.span
                initial={{ opacity: 0, scale: 0.5 }}
                animate={{ opacity: 1, scale: 1 }}
                className="font-sans text-xs tracking-[0.15em] uppercase opacity-70"
              >
                Selected
              </motion.span>
            )}
          </motion.button>
        );
      })}
    </motion.div>
  );
}
