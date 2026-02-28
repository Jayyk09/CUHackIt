import { SiteHeader } from '@/components/site-header'
import { Hero } from '@/components/hero'
import { HealthFacts } from '@/components/health-facts'
import { FoodPyramid } from '@/components/food-pyramid'
import { SiteFooter } from '@/components/site-footer'

export default function Home() {
  return (
    <main className="min-h-screen">
      <SiteHeader />
      <Hero />
      <div id="facts">
        <HealthFacts />
      </div>
      <div id="pyramid">
        <FoodPyramid />
      </div>
      <SiteFooter />
    </main>
  )
}
