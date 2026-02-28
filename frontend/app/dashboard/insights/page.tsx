const highlights = [
  { label: 'Items Rescued', value: '24', detail: 'past 30 days' },
  { label: 'Waste Avoided', value: '6.3 kg', detail: 'food saved' },
  { label: 'Money Saved', value: '$42', detail: 'estimated' },
  { label: 'CO2 Offset', value: '9.1 kg', detail: 'impact score' },
]

const weeklyPlan = [
  { day: 'Mon', focus: 'Pantry curry', status: 'Cooked' },
  { day: 'Tue', focus: 'Veggie wraps', status: 'Planned' },
  { day: 'Wed', focus: 'Bean chili', status: 'Planned' },
  { day: 'Thu', focus: 'Pasta night', status: 'Planned' },
  { day: 'Fri', focus: 'Leftover bowl', status: 'Flexible' },
]

const expiringSoon = [
  { name: 'Spinach', days: 2 },
  { name: 'Greek yogurt', days: 3 },
  { name: 'Mushrooms', days: 1 },
]

export default function InsightsPage() {
  return (
    <div className="min-h-screen">
      <div className="border-b border-espresso/10">
        <div className="px-8 py-6">
          <p className="font-sans text-xs uppercase tracking-[0.3em] text-espresso/40">
            Pantry Intelligence
          </p>
          <h1 className="font-serif text-3xl text-espresso">Insights and Weekly Rhythm</h1>
        </div>
      </div>

      <div className="px-8 py-10 space-y-10">
        <section className="grid gap-6 lg:grid-cols-4">
          {highlights.map((item) => (
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

        <section className="grid gap-6 lg:grid-cols-[1.1fr_0.9fr]">
          <div className="rounded-3xl border border-espresso/10 bg-cream p-8">
            <div className="flex items-center justify-between">
              <h2 className="font-serif text-2xl text-espresso">Weekly Plan</h2>
              <span className="text-xs uppercase tracking-[0.2em] font-sans text-espresso/40">
                Gemini-assisted
              </span>
            </div>
            <div className="mt-6 space-y-4">
              {weeklyPlan.map((plan) => (
                <div
                  key={plan.day}
                  className="flex items-center justify-between rounded-2xl border border-espresso/10 px-4 py-3"
                >
                  <div className="flex items-center gap-4">
                    <span className="font-serif text-lg text-espresso">{plan.day}</span>
                    <span className="text-espresso/70 font-sans">{plan.focus}</span>
                  </div>
                  <span className="text-[10px] uppercase tracking-[0.2em] font-sans text-espresso/50">
                    {plan.status}
                  </span>
                </div>
              ))}
            </div>
          </div>

          <div className="rounded-3xl border border-espresso/10 bg-gradient-to-br from-cream via-cream to-sage/10 p-8">
            <h2 className="font-serif text-2xl text-espresso">Expiring Soon</h2>
            <p className="text-espresso/70 font-sans mt-3">
              Prioritize these ingredients to reduce waste and boost freshness.
            </p>
            <div className="mt-6 space-y-3">
              {expiringSoon.map((item) => (
                <div
                  key={item.name}
                  className="flex items-center justify-between rounded-2xl border border-espresso/10 px-4 py-3"
                >
                  <span className="font-sans text-espresso">{item.name}</span>
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
        </section>

        <section className="rounded-3xl border border-espresso/10 bg-espresso text-cream p-8">
          <div className="flex flex-wrap items-center justify-between gap-6">
            <div>
              <p className="text-xs uppercase tracking-[0.3em] font-sans text-cream/60">
                Demo Note
              </p>
              <h2 className="font-serif text-2xl text-cream">Personalized insights, coming soon.</h2>
              <p className="text-cream/70 font-sans mt-2 max-w-xl">
                These metrics are pre-filled for the demo. Next, we will connect live pantry activity
                and real-time recipe performance to unlock predictive insights.
              </p>
            </div>
            <button className="rounded-full border border-cream/30 px-6 py-2 text-xs uppercase tracking-[0.25em] font-sans text-cream hover:border-cream/60 transition">
              Share story
            </button>
          </div>
        </section>
      </div>
    </div>
  )
}
