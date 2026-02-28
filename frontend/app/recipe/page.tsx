'use client'

import { useState, useEffect, useCallback, useRef } from 'react'
import Image from 'next/image'
import { useSearchParams } from 'next/navigation'
import { motion, AnimatePresence } from 'framer-motion'
import { toast } from 'sonner'
import {
  generateRecipes,
  saveRecipe,
  listRecipes,
  toggleFavorite,
} from '@/lib/recipes-api'
import type { GeneratedRecipe, SavedRecipe, GenerateResult } from '@/lib/recipes-api'
import { RecipeCard } from '@/components/recipe-card'
import { RecipeDetail } from '@/components/recipe-detail'
import { useUser } from '@/contexts/user-context'

// ─── Search / Generate bar ────────────────────────────────────────────────────

function RecipeSearchBar({
  onGenerate,
  isGenerating,
  initialMode,
}: {
  onGenerate: (mode: 'pantry_only' | 'flexible' | 'spoiling', prompt: string) => void
  isGenerating: boolean
  initialMode?: 'pantry_only' | 'flexible' | 'spoiling'
}) {
  const [mode, setMode] = useState<'pantry_only' | 'flexible' | 'spoiling'>(initialMode ?? 'flexible')
  const [prompt, setPrompt] = useState('')

  // Sync if initialMode changes (e.g. from URL param)
  useEffect(() => {
    if (initialMode) setMode(initialMode)
  }, [initialMode])

  return (
    <div className="border-b border-espresso pb-6">
      <div className="flex items-center gap-4">
        <div className="flex-1">
          <input
            type="text"
            value={prompt}
            onChange={(e) => setPrompt(e.target.value)}
            placeholder={isGenerating ? 'Sifting through your pantry...' : 'describe what you want (e.g. grilled chicken)'}
            readOnly={isGenerating}
            onKeyDown={(e) => {
              if (e.key === 'Enter' && !isGenerating) onGenerate(mode, prompt.trim())
            }}
            className="w-full bg-transparent font-serif text-2xl md:text-3xl text-espresso placeholder:text-espresso/35 placeholder:font-serif focus:outline-none"
          />
        </div>

        {/* mode toggle */}
        <div className="flex-shrink-0 flex items-center border border-espresso/15 overflow-hidden">
          <button
            onClick={() => setMode('pantry_only')}
            className={`px-3 py-1.5 font-sans text-[9px] uppercase tracking-[0.18em] transition-colors duration-150 ${
              mode === 'pantry_only'
                ? 'bg-espresso text-cream'
                : 'text-espresso/40 hover:text-espresso/70'
            }`}
          >
            Pantry Only
          </button>
          <button
            onClick={() => setMode('flexible')}
            className={`px-3 py-1.5 font-sans text-[9px] uppercase tracking-[0.18em] transition-colors duration-150 ${
              mode === 'flexible'
                ? 'bg-sage text-cream'
                : 'text-espresso/40 hover:text-espresso/70'
            }`}
          >
            Flexible
          </button>
          <button
            onClick={() => setMode('spoiling')}
            className={`px-3 py-1.5 font-sans text-[9px] uppercase tracking-[0.18em] transition-colors duration-150 ${
              mode === 'spoiling'
                ? 'bg-red-600 text-cream'
                : 'text-espresso/40 hover:text-espresso/70'
            }`}
          >
            Use It Up
          </button>
        </div>

        {/* generate button */}
        <button
          onClick={() => onGenerate(mode, prompt.trim())}
          disabled={isGenerating}
          className="flex-shrink-0 flex items-center gap-2 border border-espresso/20 px-4 py-2 font-sans text-[10px] uppercase tracking-[0.2em] text-espresso/60 hover:border-espresso hover:text-espresso transition-all duration-150 disabled:opacity-30 disabled:cursor-not-allowed"
        >
          <Image src="/leaf.png" alt="" width={14} height={14} unoptimized />
          {isGenerating ? 'Generating' : 'Generate'}
        </button>
      </div>
    </div>
  )
}

// ─── Tag filter bar ───────────────────────────────────────────────────────────

function TagFilter({
  tags,
  active,
  onToggle,
}: {
  tags: string[]
  active: Set<string>
  onToggle: (tag: string) => void
}) {
  if (tags.length === 0) return null
  return (
    <div className="flex items-center gap-2 flex-wrap">
      <span className="font-sans text-[9px] uppercase tracking-[0.2em] text-espresso/30 mr-1">
        Filter
      </span>
      {tags.map((tag) => (
        <button
          key={tag}
          onClick={() => onToggle(tag)}
          className={`font-sans text-[9px] uppercase tracking-[0.15em] border px-2 py-[3px] transition-colors duration-150 ${
            active.has(tag)
              ? 'bg-espresso text-cream border-espresso'
              : 'text-espresso/40 border-espresso/15 hover:border-espresso/40 hover:text-espresso/60'
          }`}
        >
          {tag}
        </button>
      ))}
      {active.size > 0 && (
        <button
          onClick={() => active.forEach(onToggle)}
          className="font-sans text-[9px] uppercase tracking-[0.15em] text-espresso/30 hover:text-espresso/60 transition-colors ml-1"
        >
          Clear
        </button>
      )}
    </div>
  )
}

// ─── Empty state ──────────────────────────────────────────────────────────────

function EmptyState({ onGenerate }: { onGenerate: () => void }) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      className="py-24 text-center"
    >
      <p className="font-serif text-4xl text-espresso/15 mb-4">No recipes yet</p>
      <p className="font-sans text-[11px] uppercase tracking-[0.2em] text-espresso/30 mb-8">
        Generate your first recipe from your pantry
      </p>
      <button
        onClick={onGenerate}
        className="font-sans text-[10px] uppercase tracking-[0.2em] border border-sage text-sage px-6 py-3 hover:bg-sage hover:text-cream transition-colors duration-150"
      >
        Generate Recipes
      </button>
    </motion.div>
  )
}

// ─── Favorite toggle ──────────────────────────────────────────────────────────

function FavoriteButton({
  recipe,
  onToggle,
}: {
  recipe: SavedRecipe
  onToggle: (id: string) => void
}) {
  return (
    <button
      onClick={(e) => {
        e.stopPropagation()
        onToggle(recipe.id)
      }}
      className={`font-sans text-[9px] uppercase tracking-[0.15em] border px-2 py-[3px] transition-colors duration-150 ${
        recipe.is_favorite
          ? 'border-sage/40 text-sage'
          : 'border-espresso/10 text-espresso/25 hover:text-espresso/50 hover:border-espresso/25'
      }`}
      title={recipe.is_favorite ? 'Remove from favorites' : 'Add to favorites'}
    >
      {recipe.is_favorite ? '♥ Saved' : '♡ Save'}
    </button>
  )
}

// ─── Page ─────────────────────────────────────────────────────────────────────

export default function RecipePage() {
  const { user, isLoading: isUserLoading } = useUser()
  const searchParams = useSearchParams()
  const [savedRecipes, setSavedRecipes] = useState<SavedRecipe[]>([])
  const [generatedResult, setGeneratedResult] = useState<GenerateResult | null>(null)
  const [isGenerating, setIsGenerating] = useState(false)
  const [isLoading, setIsLoading] = useState(true)

  const [selectedRecipe, setSelectedRecipe] = useState<GeneratedRecipe | SavedRecipe | null>(null)
  const [activeTags, setActiveTags] = useState<Set<string>>(new Set())

  // Read mode from URL params (e.g. ?mode=spoiling)
  const urlMode = searchParams.get('mode') as 'pantry_only' | 'flexible' | 'spoiling' | null
  const autoGenerateTriggered = useRef(false)

  // load saved recipes
  useEffect(() => {
    if (!user) return
    listRecipes(user.id)
      .then(setSavedRecipes)
      .catch(() => toast.error('Could not load saved recipes'))
      .finally(() => setIsLoading(false))
  }, [user])

  const handleGenerate = useCallback(
    async (mode: 'pantry_only' | 'flexible' | 'spoiling' = 'flexible', prompt = '') => {
      if (isGenerating || !user) return
      setIsGenerating(true)
      setGeneratedResult(null)

      try {
        const result = await generateRecipes(user.id, mode, 2, prompt)

        // Deduplicate by title (case-insensitive)
        const seen = new Set<string>()
        const unique = result.all_recipes.filter((r) => {
          const key = r.title.toLowerCase()
          if (seen.has(key)) return false
          seen.add(key)
          return true
        })
        result.all_recipes = unique
        result.total_count = unique.length

        setGeneratedResult(result)

        // Only auto-save real AI-generated recipes, not hardcoded mocks
        if (result.is_mock) {
          toast.info('Showing sample recipes (backend unavailable)')
          return
        }

        // auto-save each recipe, then refresh saved list
        const saved = await Promise.all(
          result.all_recipes.map((r) => saveRecipe(user.id, r).catch(() => null))
        )
        const newlySaved = saved.filter(Boolean) as SavedRecipe[]

        if (newlySaved.length > 0) {
          // Merge without duplicating titles already in the list
          setSavedRecipes((prev) => {
            const existingTitles = new Set(prev.map((r) => r.title.toLowerCase()))
            const fresh = newlySaved.filter(
              (r) => !existingTitles.has(r.title.toLowerCase())
            )
            return [...fresh, ...prev]
          })
          toast.success(
            `${newlySaved.length} recipe${newlySaved.length > 1 ? 's' : ''} generated & saved`
          )
        } else if (result.total_count > 0) {
          toast.info('Recipes generated (could not auto-save)')
        } else {
          toast.info('No recipes could be generated — add items to your pantry')
        }
      } catch (err: unknown) {
        const msg = err instanceof Error ? err.message : 'Generation failed'
        toast.error(msg)
      } finally {
        setIsGenerating(false)
      }
    },
    [isGenerating, user]
  )

  // Auto-generate when navigating with ?mode=spoiling (from expiring toast)
  useEffect(() => {
    if (autoGenerateTriggered.current) return
    if (!user || isUserLoading || isLoading) return
    if (urlMode === 'spoiling') {
      autoGenerateTriggered.current = true
      handleGenerate('spoiling', '')
    }
  }, [user, isUserLoading, isLoading, urlMode, handleGenerate])

  const handleToggleFavorite = async (id: string) => {
    if (!user) return
    try {
      const updated = await toggleFavorite(user.id, id)
      setSavedRecipes((prev) => prev.map((r) => (r.id === id ? updated : r)))
    } catch {
      toast.error('Could not update favorite')
    }
  }

  const handleTagToggle = (tag: string) => {
    setActiveTags((prev) => {
      const next = new Set(prev)
      next.has(tag) ? next.delete(tag) : next.add(tag)
      return next
    })
  }

  // collect all unique tags from saved recipes
  const allTags = Array.from(
    new Set(savedRecipes.flatMap((r) => r.tags ?? []))
  ).sort()

  // filter saved recipes by active tags
  const filteredRecipes =
    activeTags.size === 0
      ? savedRecipes
      : savedRecipes.filter((r) =>
          Array.from(activeTags).every((t) => r.tags?.includes(t))
        )

  const showGenerated =
    generatedResult && generatedResult.all_recipes.length > 0 && !isGenerating

  return (
    <>
      {/* page header */}
      <div className="border-b border-espresso/10">
        <div className="px-8 py-6 flex items-baseline justify-between">
          <h1 className="font-serif text-2xl text-espresso">Recipes</h1>
          {savedRecipes.length > 0 && (
            <p className="font-sans text-[10px] uppercase tracking-[0.2em] text-espresso/30">
              {savedRecipes.length} saved
            </p>
          )}
        </div>
      </div>

      <div className="p-8 space-y-10">
        {/* generate bar */}
        <RecipeSearchBar onGenerate={handleGenerate} isGenerating={isGenerating} initialMode={urlMode ?? undefined} />

        {/* generating spinner */}
        <AnimatePresence>
          {isGenerating && (
            <motion.div
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0 }}
              className="py-16 text-center"
            >
              <p className="font-serif text-3xl text-espresso/20 mb-3">
                {urlMode === 'spoiling' ? 'Rescuing expiring ingredients' : 'Finding the right recipe'}
              </p>
              <p className="font-sans text-[10px] uppercase tracking-[0.25em] text-espresso/20">
                {urlMode === 'spoiling' ? 'building recipes to reduce waste' : 'analysing your pantry'}
              </p>
            </motion.div>
          )}
        </AnimatePresence>

        {/* newly generated recipes */}
        <AnimatePresence>
          {showGenerated && (
            <motion.section
              initial={{ opacity: 0, y: 12 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0 }}
            >
              <div className="flex items-center gap-4 mb-2">
                <p className="font-sans text-[10px] uppercase tracking-[0.2em] text-espresso/40">
                  Just generated
                </p>
                <div className="flex-1 h-px bg-espresso/8" />
                <p className="font-sans text-[10px] uppercase tracking-[0.2em] text-espresso/25">
                  {generatedResult!.total_count} new
                </p>
              </div>
              {generatedResult!.all_recipes.map((recipe, i) => (
                <RecipeCard
                  key={`gen-${i}`}
                  recipe={recipe}
                  index={i}
                  onClick={() => setSelectedRecipe(recipe)}
                />
              ))}
            </motion.section>
          )}
        </AnimatePresence>

        {/* saved recipes */}
        {isLoading ? (
          <div className="py-12 text-center">
            <p className="font-sans text-[10px] uppercase tracking-[0.2em] text-espresso/25">
              Loading...
            </p>
          </div>
        ) : savedRecipes.length === 0 && !isGenerating ? (
          <EmptyState onGenerate={() => handleGenerate('flexible')} />
        ) : savedRecipes.length > 0 ? (
          <section>
            <div className="flex items-center gap-4 mb-6">
              <p className="font-sans text-[10px] uppercase tracking-[0.2em] text-espresso/40">
                All recipes
              </p>
              <div className="flex-1 h-px bg-espresso/8" />
            </div>

            {/* tag filter */}
            {allTags.length > 0 && (
              <div className="mb-6">
                <TagFilter
                  tags={allTags}
                  active={activeTags}
                  onToggle={handleTagToggle}
                />
              </div>
            )}

            {filteredRecipes.length === 0 ? (
              <p className="font-sans text-[12px] text-espresso/30 py-8 text-center">
                No recipes match the selected tags
              </p>
            ) : (
              <div>
                {filteredRecipes.map((recipe, i) => (
                  <div key={recipe.id} className="relative group/row">
                    <RecipeCard
                      recipe={recipe}
                      index={i}
                      onClick={() => setSelectedRecipe(recipe)}
                    />
                    {/* favorite button floated top-right */}
                    <div className="absolute top-6 right-0 opacity-0 group-hover/row:opacity-100 transition-opacity duration-150">
                      <FavoriteButton
                        recipe={recipe}
                        onToggle={handleToggleFavorite}
                      />
                    </div>
                  </div>
                ))}
              </div>
            )}
          </section>
        ) : null}
      </div>

      {/* recipe detail modal */}
      <RecipeDetail recipe={selectedRecipe} onClose={() => setSelectedRecipe(null)} />
    </>
  )
}
