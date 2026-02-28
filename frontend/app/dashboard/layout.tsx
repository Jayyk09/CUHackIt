import { DashboardSidebar } from '@/components/dashboard-sidebar'

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <div className="min-h-screen bg-cream">
      <DashboardSidebar />
      <main className="ml-72 min-h-screen">
        <div className="border-b border-espresso/10">
          <div className="px-8 py-6">
            <h1 className="font-serif text-2xl text-espresso">Dashboard</h1>
          </div>
        </div>
        <div className="p-8">
          {children}
        </div>
      </main>
    </div>
  )
}
