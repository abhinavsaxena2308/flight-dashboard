'use client';

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { useRouter } from 'next/navigation';
import StateHighlightMap from '@/components/StateHighlightMap';

interface StateData {
  state: string;
  totalFlights: number;
  incomingFlights: number;
  outgoingFlights: number;
  routes: number;
  airlines: string[];
}

const StateDashboardPage = () => {
  const router = useRouter();
  const { stateName } = useParams();
  const [stateData, setStateData] = useState<StateData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchStateData = async () => {
      if (!stateName) return;
      
      try {
        setLoading(true);
        // The stateName from params will be in kebab-case, which is what our API expects
        const response = await fetch(`/api/state/${stateName}`);
        
        if (!response.ok) {
          throw new Error(`Failed to fetch data for state: ${stateName}`);
        }
        
        const data: StateData = await response.json();
        setStateData(data);
        setError(null);
      } catch (err) {
        console.error('Error fetching state data:', err);
        setError(err instanceof Error ? err.message : 'Failed to load state data');
      } finally {
        setLoading(false);
      }
    };

    fetchStateData();
  }, [stateName]);

  if (loading) {
    return (
      <div className="p-6">
        <div className="max-w-7xl mx-auto">
          <div className="bg-white rounded-lg shadow-md p-6 mb-6">
            <h1 className="text-2xl font-bold text-gray-800">State Dashboard</h1>
            <p>Loading state data...</p>
          </div>
        </div>
      </div>
    );
  }

  if (error || !stateData) {
    return (
      <div className="p-6">
        <div className="max-w-7xl mx-auto">
          <div className="bg-white rounded-lg shadow-md p-6 mb-6">
            <h1 className="text-2xl font-bold text-gray-800">State Dashboard</h1>
            <p className="text-red-600">Error: {error || 'State data not found'}</p>
            <button 
              onClick={() => router.push('/')}
              className="mt-4 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Go Back
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <div className="flex justify-between items-center">
            <h1 className="text-3xl font-bold text-gray-800 capitalize">
              {stateData.state} Flight Dashboard
            </h1>
            <button 
              onClick={() => router.push('/')}
              className="px-4 py-2 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
            >
              ← Back to Map
            </button>
          </div>
          <p className="text-gray-600 mt-2">Detailed flight statistics for {stateData.state}</p>
        </div>

        {/* Map Visualization */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <h2 className="text-xl font-bold text-gray-800 mb-4">Location Map</h2>
          <StateHighlightMap stateName={stateName as string} stateDisplayName={stateData.state} />
        </div>

        {/* Stats Overview */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-700">Total Flights</h3>
            <p className="text-3xl font-bold text-blue-600">{stateData.totalFlights.toLocaleString()}</p>
          </div>
          
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-700">Incoming Flights</h3>
            <p className="text-3xl font-bold text-green-600">{stateData.incomingFlights.toLocaleString()}</p>
          </div>
          
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-700">Outgoing Flights</h3>
            <p className="text-3xl font-bold text-orange-600">{stateData.outgoingFlights.toLocaleString()}</p>
          </div>
          
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-700">Unique Routes</h3>
            <p className="text-3xl font-bold text-purple-600">{stateData.routes.toLocaleString()}</p>
          </div>
        </div>

        {/* Airlines Section */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <h2 className="text-xl font-bold text-gray-800 mb-4">Airlines Serving {stateData.state}</h2>
          {stateData.airlines.length > 0 ? (
            <div className="flex flex-wrap gap-2">
              {stateData.airlines.map((airline, index) => (
                <span 
                  key={index}
                  className="px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm font-medium"
                >
                  {airline}
                </span>
              ))}
            </div>
          ) : (
            <p className="text-gray-600">No airline data available for this state.</p>
          )}
        </div>

        {/* Additional Info */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-bold text-gray-800 mb-4">State Flight Analysis</h2>
          <div className="space-y-4">
            <div>
              <h3 className="font-semibold text-gray-700">Flight Balance</h3>
              <p className="text-gray-600">
                {stateData.incomingFlights > stateData.outgoingFlights 
                  ? `${stateData.state} receives more flights than it sends out.`
                  : stateData.outgoingFlights > stateData.incomingFlights
                    ? `${stateData.state} sends out more flights than it receives.`
                    : `${stateData.state} has balanced flight traffic.`}
              </p>
            </div>
            
            <div>
              <h3 className="font-semibold text-gray-700">Connectivity</h3>
              <p className="text-gray-600">
                {stateData.state} is connected to {stateData.routes} unique routes, 
                served by {stateData.airlines.length} different airlines.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default StateDashboardPage;