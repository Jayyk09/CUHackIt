const recipes = [
  {
    title: 'Golden Chickpea Skillet',
    time: '25 min',
    servings: 2,
    tags: ['pantry-only', 'high-protein'],
    summary:
      'Crispy chickpeas, wilted greens, and lemony tahini for a fast weeknight win.',
    ingredients: ['chickpeas', 'spinach', 'garlic', 'olive oil', 'tahini', 'lemon'],
  },
  {
    title: 'Harvest Lentil Bowl',
    time: '35 min',
    servings: 3,
    tags: ['fiber-rich', 'meal-prep'],
    summary:
      'Earthy lentils, roasted carrots, and herbed yogurt with a bright vinegar splash.',
    ingredients: ['lentils', 'carrots', 'yogurt', 'parsley', 'red wine vinegar'],
  },
  {
    title: 'Pantry Tomato Rigatoni',
    time: '20 min',
    servings: 2,
    tags: ['family-favorite', '20-min'],
    summary:
      'Slow-simmered tomato sauce with basil and buttery breadcrumbs for crunch.',
    ingredients: ['rigatoni', 'tomato puree', 'garlic', 'basil', 'breadcrumbs'],
  },
]

const pantryHighlights = [
  { label: 'Pantry Match', value: '92%' },
  { label: 'Expiring Soon', value: '3 items' },
  { label: 'Estimated Cost', value: '$8.40' },
  { label: 'CO2 Saved', value: '1.6 kg' },
]

const pantryItems = [
  'chickpeas',
  'rigatoni',
  'tomato puree',
  'spinach',
  'olive oil',
  'garlic',
  'lentils',
  'carrots',
  'basil',
]

export default function RecipesPage() {
  return (
    <div className="min-h-screen">
      <div className="border-b border-espresso/10">
        <div className="px-8 py-6 flex flex-wrap items-center justify-between gap-4">
          <div>
            <p className="font-sans text-xs uppercase tracking-[0.3em] text-espresso/40">
              Gemini Recipes
            </p>
            <h1 className="font-serif text-3xl text-espresso">Recipes Crafted From Your Pantry</h1>
          </div>
          <button className="rounded-full border border-espresso/20 bg-cream px-6 py-2 text-xs uppercase tracking-[0.25em] font-sans text-espresso hover:border-espresso/40 transition">
            Generate New Set
          </button>
        </div>
      </div>

      <div className="px-8 py-10 space-y-10">
        <section className="grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
          <div className="rounded-3xl border border-espresso/10 bg-gradient-to-br from-cream via-cream to-sage/10 p-8">
            <h2 className="font-serif text-2xl text-espresso mb-4">Tonight’s Spotlight</h2>
            <p className="text-espresso/70 font-sans leading-relaxed max-w-xl">
              We used your pantry inventory and expiring items to curate a balanced set of recipes.
              Gemini highlights the best matches so you can cook confidently and reduce waste.
            </p>
            <div className="mt-6 flex flex-wrap gap-3">
              {pantryItems.map((item) => (
                <span
                  key={item}
                  className="rounded-full border border-espresso/20 px-3 py-1 text-xs uppercase tracking-[0.2em] font-sans text-espresso/70"
                >
                  {item}
                </span>
              ))}
            </div>
          </div>

          <div className="rounded-3xl border border-espresso/10 bg-cream p-8">
            <h3 className="font-serif text-xl text-espresso mb-6">Session Snapshot</h3>
            <div className="grid gap-4">
              {pantryHighlights.map((highlight) => (
                <div
                  key={highlight.label}
                  className="flex items-center justify-between border-b border-espresso/10 pb-3"
                >
                  <span className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/50">
                    {highlight.label}
                  </span>
                  <span className="font-serif text-xl text-espresso">{highlight.value}</span>
                </div>
              ))}
            </div>
          </div>
        </section>

        <section className="grid gap-6 lg:grid-cols-3">
          {recipes.map((recipe) => (
            <article
              key={recipe.title}
              className="rounded-3xl border border-espresso/10 bg-cream p-6 flex flex-col"
            >
              <div className="flex items-center justify-between">
                <h3 className="font-serif text-xl text-espresso">{recipe.title}</h3>
                <span className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/60">
                  {recipe.time}
                </span>
              </div>
              <p className="mt-3 text-espresso/70 font-sans leading-relaxed">
                {recipe.summary}
              </p>
              <div className="mt-4 flex flex-wrap gap-2">
                {recipe.tags.map((tag) => (
                  <span
                    key={tag}
                    className="rounded-full bg-sage/20 px-3 py-1 text-[10px] uppercase tracking-[0.2em] font-sans text-espresso"
                  >
                    {tag}
                  </span>
                ))}
              </div>
              <div className="mt-6 space-y-2">
                <p className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/50">
                  Key ingredients
                </p>
                <div className="flex flex-wrap gap-2">
                  {recipe.ingredients.map((ingredient) => (
                    <span
                      key={ingredient}
                      className="rounded-full border border-espresso/15 px-3 py-1 text-xs font-sans text-espresso/70"
                    >
                      {ingredient}
                    </span>
                  ))}
                </div>
              </div>
              <div className="mt-auto pt-6 flex items-center justify-between text-xs uppercase tracking-[0.2em] font-sans text-espresso/50">
                <span>{recipe.servings} servings</span>
                <button className="text-espresso hover:text-espresso/70 transition">
                  View recipe
                </button>
              </div>
            </article>
          ))}
        </section>

        <section className="rounded-3xl border border-espresso/10 bg-espresso text-cream p-8 flex flex-wrap items-center justify-between gap-6">
          <div>
            <p className="text-xs uppercase tracking-[0.3em] font-sans text-cream/60">
              Demo Mode
            </p>
            <h2 className="font-serif text-2xl text-cream">Gemini-powered planning, ready for prime time.</h2>
            <p className="text-cream/70 font-sans mt-2 max-w-xl">
              We are showcasing curated results today. With more time, we would connect live pantry data
              and personalize every recipe in real time.
            </p>
          </div>
          <button className="rounded-full border border-cream/30 px-6 py-2 text-xs uppercase tracking-[0.25em] font-sans text-cream hover:border-cream/60 transition">
            Request full demo
          </button>
        </section>
      </div>
    </div>
  )
}
