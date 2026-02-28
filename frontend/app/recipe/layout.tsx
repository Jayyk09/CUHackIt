import { DashboardSidebar } from '@/components/dashboard-sidebar'
import { Toaster } from '@/components/ui/sonner'

export default function RecipeLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-cream">
      <DashboardSidebar />
      <main className="ml-72 min-h-screen">{children}</main>
      <Toaster />
    </div>
  )
}
