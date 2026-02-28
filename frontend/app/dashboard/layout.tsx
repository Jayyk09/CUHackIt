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
        {children}
      </main>
    </div>
  )
}
