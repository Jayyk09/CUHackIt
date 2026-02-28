'use client';

import { motion } from 'framer-motion';
import { cn } from '@/lib/utils';

interface BubbleSelectProps {
  options: string[];
  selected: string[];
  onChange: (selected: string[]) => void;
}

export function BubbleSelect({ options, selected, onChange }: BubbleSelectProps) {
  const toggleOption = (option: string) => {
    if (selected.includes(option)) {
      onChange(selected.filter((item) => item !== option));
    } else {
      onChange([...selected, option]);
    }
  };

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.04,
        delayChildren: 0.1,
      },
    },
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 10, scale: 0.95 },
    visible: {
      opacity: 1,
      y: 0,
      scale: 1,
      transition: {
        duration: 0.3,
        ease: [0.22, 1, 0.36, 1],
      },
    },
  };

  return (
    <motion.div
      variants={containerVariants}
      initial="hidden"
      animate="visible"
      className="flex flex-wrap gap-3"
    >
      {options.map((option) => {
        const isSelected = selected.includes(option);
        return (
          <motion.button
            key={option}
            type="button"
            variants={itemVariants}
            onClick={() => toggleOption(option)}
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
            className={cn(
              'px-5 py-2.5',
              'font-sans text-sm tracking-[0.05em]',
              'transition-all duration-200 ease-out',
              'cursor-pointer select-none',
              isSelected
                ? 'bg-[#4A5D4E] text-[#F9F8F6] border border-transparent'
                : 'bg-transparent text-[#2A2724]/70 border border-[#2A2724]/25 hover:border-[#2A2724]/50 hover:text-[#2A2724]'
            )}
          >
            {option}
          </motion.button>
        );
      })}
    </motion.div>
  );
}
