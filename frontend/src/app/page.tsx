
import IndiaMap from "@/components/IndiaMap"

export default function Home() {
  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto">
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <div className="flex justify-between items-center">
            <h1 className="text-3xl font-bold text-gray-800">India Flight Dashboard</h1>
            <div className="text-sm text-gray-600">
              Click on a state to view detailed flight statistics
            </div>
          </div>
        </div>
        <div className="h-[70vh] bg-white rounded-lg shadow-md p-4">
          <IndiaMap />
        </div>
      </div>
    </div>
  )
}