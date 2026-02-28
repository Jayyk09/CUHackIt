'use client'

import { motion, AnimatePresence } from 'framer-motion'
import type { GeneratedRecipe, SavedRecipe, Ingredient } from '@/lib/recipes-api'

type AnyRecipe = GeneratedRecipe | SavedRecipe

function MacroBar({
  label,
  value,
  unit,
  color,
}: {
  label: string
  value?: number
  unit: string
  color: string
}) {
  if (!value) return null
  return (
    <div>
      <p className="font-sans text-[9px] uppercase tracking-[0.2em] text-espresso/40 mb-1">
        {label}
      </p>
      <p className={`font-serif text-3xl leading-none ${color}`}>
        {Math.round(value)}
        <span className="font-sans text-[11px] ml-1 text-espresso/40">{unit}</span>
      </p>
    </div>
  )
}

function IngredientRow({ ingredient }: { ingredient: Ingredient }) {
  const amount = [ingredient.amount, ingredient.unit].filter(Boolean).join(' ')
  return (
    <div className="flex items-baseline gap-2 py-2 border-b border-espresso/6 last:border-0">
      <span
        className={`flex-shrink-0 text-[10px] font-sans mt-0.5 ${
          ingredient.from_pantry ? 'text-sage' : 'text-espresso/25'
        }`}
      >
        {ingredient.from_pantry ? '✓' : '○'}
      </span>
      <span className="font-sans text-[13px] text-espresso">{ingredient.name}</span>
      {amount && (
        <span className="ml-auto font-sans text-[11px] text-espresso/40 flex-shrink-0">
          {amount}
        </span>
      )}
    </div>
  )
}

export function RecipeDetail({
  recipe,
  onClose,
}: {
  recipe: AnyRecipe | null
  onClose: () => void
}) {
  return (
    <AnimatePresence>
      {recipe && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.18 }}
          className="fixed inset-0 z-50 flex items-center justify-center bg-espresso/25 backdrop-blur-[2px]"
          onClick={onClose}
        >
          <motion.div
            initial={{ scale: 0.97, opacity: 0, y: 8 }}
            animate={{ scale: 1, opacity: 1, y: 0 }}
            exit={{ scale: 0.97, opacity: 0, y: 8 }}
            transition={{ type: 'spring', damping: 32, stiffness: 340 }}
            className="relative w-full max-w-2xl mx-4 bg-cream border border-espresso/12 max-h-[90vh] overflow-y-auto"
            onClick={(e) => e.stopPropagation()}
          >
            {/* header band */}
            <div className="bg-sage-dark px-8 pt-10 pb-8">
              <button
                onClick={onClose}
                className="absolute top-6 right-6 font-sans text-[10px] uppercase tracking-[0.2em] text-cream/50 hover:text-cream transition-colors"
              >
                — Close
              </button>

              {/* source + cuisine */}
              <div className="flex items-center gap-3 mb-4">
                <span className="font-sans text-[9px] uppercase tracking-[0.2em] text-cream/40 border border-cream/20 px-2 py-[2px]">
                  {recipe.source === 'pantry_only'
                    ? 'Pantry Only'
                    : recipe.source === 'flexible'
                    ? 'Flexible'
                    : 'Custom'}
                </span>
                {recipe.cuisine && (
                  <span className="font-sans text-[9px] uppercase tracking-[0.2em] text-cream/40">
                    {recipe.cuisine}
                  </span>
                )}
              </div>

              <h2 className="font-serif text-3xl md:text-4xl text-cream leading-tight mb-5">
                {recipe.title}
              </h2>

              {/* meta row */}
              <div className="flex items-center gap-6 flex-wrap">
                {(recipe.total_time_minutes ?? 0) > 0 && (
                  <div>
                    <p className="font-sans text-[9px] uppercase tracking-[0.2em] text-cream/35 mb-1">
                      Total time
                    </p>
                    <p className="font-sans text-sm text-cream/80">
                      {recipe.total_time_minutes} min
                    </p>
                  </div>
                )}
                {recipe.prep_time_minutes != null && recipe.prep_time_minutes > 0 && (
                  <div>
                    <p className="font-sans text-[9px] uppercase tracking-[0.2em] text-cream/35 mb-1">
                      Prep
                    </p>
                    <p className="font-sans text-sm text-cream/80">
                      {recipe.prep_time_minutes} min
                    </p>
                  </div>
                )}
                {recipe.cook_time_minutes != null && recipe.cook_time_minutes > 0 && (
                  <div>
                    <p className="font-sans text-[9px] uppercase tracking-[0.2em] text-cream/35 mb-1">
                      Cook
                    </p>
                    <p className="font-sans text-sm text-cream/80">
                      {recipe.cook_time_minutes} min
                    </p>
                  </div>
                )}
                {recipe.servings != null && recipe.servings > 0 && (
                  <div>
                    <p className="font-sans text-[9px] uppercase tracking-[0.2em] text-cream/35 mb-1">
                      Serves
                    </p>
                    <p className="font-sans text-sm text-cream/80">{recipe.servings}</p>
                  </div>
                )}
                {recipe.difficulty && (
                  <div>
                    <p className="font-sans text-[9px] uppercase tracking-[0.2em] text-cream/35 mb-1">
                      Difficulty
                    </p>
                    <p className="font-sans text-sm text-cream/80 capitalize">
                      {recipe.difficulty}
                    </p>
                  </div>
                )}
              </div>
            </div>

            {/* body */}
            <div className="px-8 py-8 space-y-8">
              {/* description */}
              {recipe.description && (
                <p className="font-sans text-[13px] text-espresso/60 leading-relaxed">
                  {recipe.description}
                </p>
              )}

              {/* ingredients */}
              {recipe.ingredients && recipe.ingredients.length > 0 && (
                <section>
                  <div className="flex items-center gap-4 mb-4">
                    <p className="font-sans text-[10px] uppercase tracking-[0.2em] text-espresso/40">
                      Ingredients
                    </p>
                    <div className="flex-1 h-px bg-espresso/8" />
                    <div className="flex items-center gap-4 text-[9px] font-sans uppercase tracking-[0.15em]">
                      <span className="text-sage">
                        ✓ {recipe.ingredients.filter((i) => i.from_pantry).length} from pantry
                      </span>
                      {(recipe.missing_ingredients?.length ?? 0) > 0 && (
                        <span className="text-espresso/35">
                          ○ {recipe.missing_ingredients!.length} needed
                        </span>
                      )}
                    </div>
                  </div>
                  <div>
                    {recipe.ingredients.map((ing, i) => (
                      <IngredientRow key={i} ingredient={ing} />
                    ))}
                  </div>
                </section>
              )}

              {/* instructions */}
              {recipe.instructions && recipe.instructions.length > 0 && (
                <section>
                  <div className="flex items-center gap-4 mb-5">
                    <p className="font-sans text-[10px] uppercase tracking-[0.2em] text-espresso/40">
                      Instructions
                    </p>
                    <div className="flex-1 h-px bg-espresso/8" />
                  </div>
                  <div className="space-y-5">
                    {recipe.instructions.map((step, i) => (
                      <div key={i} className="flex gap-5">
                        <span className="font-serif text-3xl text-espresso/12 leading-none flex-shrink-0 w-9 text-right">
                          {String(i + 1).padStart(2, '0')}
                        </span>
                        <p className="font-sans text-[13.5px] text-espresso/75 leading-relaxed pt-1">
                          {step}
                        </p>
                      </div>
                    ))}
                  </div>
                </section>
              )}

              {/* nutrition */}
              {(recipe.calories_per_serving ||
                recipe.protein_g ||
                recipe.carbs_g ||
                recipe.fat_g) && (
                <section>
                  <div className="flex items-center gap-4 mb-5">
                    <p className="font-sans text-[10px] uppercase tracking-[0.2em] text-espresso/40">
                      Nutrition per serving
                    </p>
                    <div className="flex-1 h-px bg-espresso/8" />
                  </div>
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
                    <MacroBar
                      label="Calories"
                      value={recipe.calories_per_serving}
                      unit="kcal"
                      color="text-espresso"
                    />
                    <MacroBar
                      label="Protein"
                      value={recipe.protein_g}
                      unit="g"
                      color="text-sage"
                    />
                    <MacroBar
                      label="Carbs"
                      value={recipe.carbs_g}
                      unit="g"
                      color="text-espresso"
                    />
                    <MacroBar
                      label="Fat"
                      value={recipe.fat_g}
                      unit="g"
                      color="text-espresso/60"
                    />
                  </div>
                </section>
              )}

              {/* tags */}
              {recipe.tags && recipe.tags.length > 0 && (
                <section className="border-t border-espresso/8 pt-6">
                  <div className="flex flex-wrap gap-2">
                    {recipe.tags.map((tag) => (
                      <span
                        key={tag}
                        className="font-sans text-[9px] uppercase tracking-[0.18em] text-espresso/40 border border-espresso/10 px-2 py-[3px]"
                      >
                        {tag}
                      </span>
                    ))}
                  </div>
                </section>
              )}
            </div>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  )
}
