'use client';

import { motion } from 'framer-motion';

interface NameInputProps {
  value: string;
  onChange: (value: string) => void;
}

export function NameInput({ value, onChange }: NameInputProps) {
  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ delay: 0.2, duration: 0.5 }}
      className="w-full"
    >
      <input
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder="Your name"
        autoFocus
        autoComplete="off"
        className="
          w-full
          bg-transparent
          border-0 border-b border-[#2A2724]/30
          font-serif text-4xl md:text-5xl lg:text-6xl
          text-[#2A2724]
          placeholder:text-[#2A2724]/25
          tracking-[-0.02em]
          py-4
          focus:outline-none focus:border-[#2A2724]
          transition-colors duration-300
          caret-[#4A5D4E]
        "
      />
      <p className="mt-6 font-sans text-xs tracking-[0.2em] uppercase text-[#2A2724]/40">
        Press Enter or click Next to continue
      </p>
    </motion.div>
  );
}
