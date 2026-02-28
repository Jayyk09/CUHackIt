'use client'

import { motion } from 'framer-motion'
import type { GeneratedRecipe, SavedRecipe } from '@/lib/recipes-api'

type AnyRecipe = GeneratedRecipe | SavedRecipe

function isSaved(r: AnyRecipe): r is SavedRecipe {
  return 'id' in r
}

function DifficultyDots({ level }: { level: string }) {
  const map: Record<string, number> = { easy: 1, medium: 2, hard: 3 }
  const filled = map[level] ?? 1
  return (
    <span className="flex items-center gap-[3px]">
      {[1, 2, 3].map((i) => (
        <span
          key={i}
          className={`block w-[5px] h-[5px] rounded-full transition-colors ${
            i <= filled ? 'bg-espresso' : 'bg-espresso/20'
          }`}
        />
      ))}
    </span>
  )
}

function SourceBadge({ source }: { source: AnyRecipe['source'] }) {
  const isPantry = source === 'pantry_only'
  return (
    <span
      className={`inline-block text-[9px] font-sans uppercase tracking-[0.2em] border px-2 py-[2px] leading-none ${
        isPantry
          ? 'border-sage/50 text-sage'
          : 'border-espresso/20 text-espresso/50'
      }`}
    >
      {isPantry ? 'Pantry Only' : source === 'flexible' ? 'Flexible' : 'Custom'}
    </span>
  )
}

export function RecipeCard({
  recipe,
  index = 0,
  onClick,
}: {
  recipe: AnyRecipe
  index?: number
  onClick?: () => void
}) {
  const pantryCount = recipe.ingredients?.filter((i) => i.from_pantry).length ?? 0
  const missingCount = recipe.missing_ingredients?.length ?? 0
  const totalTime = recipe.total_time_minutes
    ?? ((recipe.prep_time_minutes ?? 0) + (recipe.cook_time_minutes ?? 0))
  const isFav = isSaved(recipe) && recipe.is_favorite

  return (
    <motion.button
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.35, delay: index * 0.06, ease: [0.22, 1, 0.36, 1] }}
      onClick={onClick}
      className="group w-full text-left border-t border-espresso/12 pt-6 pb-7 focus:outline-none"
    >
      {/* top row */}
      <div className="flex items-center justify-between mb-3">
        <div className="flex items-center gap-3">
          <SourceBadge source={recipe.source} />
          {isFav && (
            <span className="text-[9px] font-sans uppercase tracking-[0.2em] text-espresso/30">
              ♥ Saved
            </span>
          )}
        </div>
        {recipe.difficulty && (
          <div className="flex items-center gap-2">
            <span className="font-sans text-[9px] uppercase tracking-[0.2em] text-espresso/35">
              {recipe.difficulty}
            </span>
            <DifficultyDots level={recipe.difficulty} />
          </div>
        )}
      </div>

      {/* title */}
      <h3 className="font-serif text-[1.65rem] leading-[1.15] text-espresso mb-2 group-hover:opacity-60 transition-opacity duration-200">
        {recipe.title}
      </h3>

      {/* description */}
      {recipe.description && (
        <p className="font-sans text-[13px] text-espresso/50 leading-relaxed mb-4 line-clamp-2">
          {recipe.description}
        </p>
      )}

      {/* meta chips */}
      <div className="flex items-center gap-5 mb-4">
        {totalTime > 0 && (
          <div>
            <p className="font-sans text-[9px] uppercase tracking-[0.18em] text-espresso/35 mb-0.5">
              Time
            </p>
            <p className="font-sans text-[13px] font-medium text-espresso">
              {totalTime} min
            </p>
          </div>
        )}
        {recipe.servings != null && recipe.servings > 0 && (
          <div>
            <p className="font-sans text-[9px] uppercase tracking-[0.18em] text-espresso/35 mb-0.5">
              Serves
            </p>
            <p className="font-sans text-[13px] font-medium text-espresso">
              {recipe.servings}
            </p>
          </div>
        )}
        {pantryCount > 0 && (
          <div>
            <p className="font-sans text-[9px] uppercase tracking-[0.18em] text-espresso/35 mb-0.5">
              Have
            </p>
            <p className="font-sans text-[13px] font-medium text-sage">
              {pantryCount}
            </p>
          </div>
        )}
        {missingCount > 0 && (
          <div>
            <p className="font-sans text-[9px] uppercase tracking-[0.18em] text-espresso/35 mb-0.5">
              Need
            </p>
            <p className="font-sans text-[13px] font-medium text-espresso/50">
              {missingCount}
            </p>
          </div>
        )}
        {recipe.calories_per_serving != null && recipe.calories_per_serving > 0 && (
          <div>
            <p className="font-sans text-[9px] uppercase tracking-[0.18em] text-espresso/35 mb-0.5">
              Cal
            </p>
            <p className="font-sans text-[13px] font-medium text-espresso">
              {Math.round(recipe.calories_per_serving)}
            </p>
          </div>
        )}
      </div>

      {/* tags */}
      {recipe.tags && recipe.tags.length > 0 && (
        <div className="flex flex-wrap gap-1.5 mb-4">
          {recipe.tags.slice(0, 5).map((tag) => (
            <span
              key={tag}
              className="font-sans text-[9px] uppercase tracking-[0.15em] text-espresso/40 border border-espresso/10 px-2 py-[3px]"
            >
              {tag}
            </span>
          ))}
        </div>
      )}

      {/* hover cta */}
      <div className="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-all duration-200 translate-x-0 group-hover:translate-x-1">
        <span className="font-sans text-[9px] uppercase tracking-[0.2em] text-espresso/40">
          View recipe
        </span>
        <span className="text-espresso/40 text-[10px]">→</span>
      </div>
    </motion.button>
  )
}
