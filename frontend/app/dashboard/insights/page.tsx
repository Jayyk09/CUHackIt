const benefitHighlights = [
  { label: 'Items Rescued', value: '24', detail: 'past 30 days' },
  { label: 'Waste Avoided', value: '6.3 kg', detail: 'food saved' },
  { label: 'Money Saved', value: '$42', detail: 'estimated' },
  { label: 'CO2 Offset', value: '9.1 kg', detail: 'impact score' },
]

const pantrySignals = [
  { label: 'Avg Days to Expiry', value: '6.8', detail: 'across pantry' },
  { label: 'Avg Nutri-Score', value: 'B', detail: 'quality signal' },
  { label: 'Avg Eco Score', value: '71', detail: 'sustainability' },
  { label: 'Freshness Index', value: '88%', detail: 'current mix' },
]

const expiringSoon = [
  { name: 'Spinach', days: 2, action: 'Use in a skillet' },
  { name: 'Greek yogurt', days: 3, action: 'Blend into sauce' },
  { name: 'Mushrooms', days: 1, action: 'Roast tonight' },
]

const pantryMix = [
  { label: 'Produce', value: '32%' },
  { label: 'Proteins', value: '21%' },
  { label: 'Whole grains', value: '18%' },
  { label: 'Dairy', value: '14%' },
  { label: 'Pantry staples', value: '15%' },
]

export default function InsightsPage() {
  return (
    <div className="min-h-screen">
      <div className="border-b border-espresso/10">
        <div className="px-8 py-6">
          <p className="font-sans text-xs uppercase tracking-[0.3em] text-espresso/40">
            Insights
          </p>
          <h1 className="font-serif text-3xl text-espresso">Your Impact and Pantry Health</h1>
        </div>
      </div>

      <div className="px-8 py-10 space-y-10">
        <section className="grid gap-6 lg:grid-cols-4">
          {benefitHighlights.map((item) => (
            <div
              key={item.label}
              className="rounded-3xl border border-espresso/10 bg-cream p-6"
            >
              <p className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/50">
                {item.label}
              </p>
              <p className="font-serif text-3xl text-espresso mt-4">{item.value}</p>
              <p className="text-espresso/60 font-sans text-sm mt-2">{item.detail}</p>
            </div>
          ))}
        </section>

        <section className="rounded-3xl border border-espresso/10 bg-gradient-to-br from-cream via-cream to-sage/10 p-8">
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div>
              <p className="text-xs uppercase tracking-[0.3em] font-sans text-espresso/40">
                Pantry Signals
              </p>
              <h2 className="font-serif text-2xl text-espresso">At-a-Glance Pantry Metrics</h2>
            </div>
            <button className="rounded-full border border-espresso/20 px-5 py-2 text-xs uppercase tracking-[0.25em] font-sans text-espresso hover:border-espresso/40 transition">
              Optimize Pantry
            </button>
          </div>
          <div className="mt-6 grid gap-4 md:grid-cols-2 lg:grid-cols-4">
            {pantrySignals.map((item) => (
              <div
                key={item.label}
                className="rounded-2xl border border-espresso/10 bg-cream/80 p-5"
              >
                <p className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/50">
                  {item.label}
                </p>
                <p className="font-serif text-3xl text-espresso mt-3">{item.value}</p>
                <p className="text-espresso/60 font-sans text-sm mt-2">{item.detail}</p>
              </div>
            ))}
          </div>
        </section>

        <section className="grid gap-6 lg:grid-cols-[1.1fr_0.9fr]">
          <div className="rounded-3xl border border-espresso/10 bg-cream p-8">
            <div className="flex items-center justify-between">
              <h2 className="font-serif text-2xl text-espresso">Expiring Soon</h2>
              <span className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/40">
                next 72 hours
              </span>
            </div>
            <div className="mt-6 space-y-4">
              {expiringSoon.map((item) => (
                <div
                  key={item.name}
                  className="flex items-center justify-between rounded-2xl border border-espresso/10 px-4 py-3"
                >
                  <div>
                    <p className="font-sans text-espresso">{item.name}</p>
                    <p className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/40 mt-1">
                      {item.action}
                    </p>
                  </div>
                  <span className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/50">
                    {item.days} days
                  </span>
                </div>
              ))}
            </div>
            <button className="mt-6 w-full rounded-full border border-espresso/20 px-5 py-2 text-xs uppercase tracking-[0.25em] font-sans text-espresso hover:border-espresso/40 transition">
              Build a rescue recipe
            </button>
          </div>

          <div className="rounded-3xl border border-espresso/10 bg-cream p-8">
            <div className="flex items-center justify-between">
              <h2 className="font-serif text-2xl text-espresso">Pantry Mix</h2>
              <span className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/40">
                balance view
              </span>
            </div>
            <div className="mt-6 space-y-4">
              {pantryMix.map((item) => (
                <div key={item.label} className="space-y-2">
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-sans text-espresso/70">{item.label}</span>
                    <span className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/50">
                      {item.value}
                    </span>
                  </div>
                  <div className="h-2 rounded-full bg-espresso/10">
                    <div
                      className="h-2 rounded-full bg-sage"
                      style={{ width: item.value }}
                    />
                  </div>
                </div>
              ))}
            </div>
          </div>
        </section>

      </div>
    </div>
  )
}
