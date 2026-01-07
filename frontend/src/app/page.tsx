
import IndiaMap from "@/components/IndiaMap"

export default function Home() {
  return (
    <div className="min-h-screen bg-dashboard-bg p-6">
      <h1 className="text-3xl font-bold text-center mb-6">India Map</h1>
      <div className="h-[70vh]">
              <IndiaMap />
            </div>
    </div>
    
  )
}