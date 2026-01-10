'use client';

import React from 'react';
import dynamic from 'next/dynamic';

// define the props interface
interface StateHighlightMapProps {
  stateName: string;
  stateDisplayName?: string;
}

// dynamically importing the client-side map component with no SSR
const ClientStateHighlightMap = dynamic(
  () => import('./ClientStateHighlightMap'),
  { 
    ssr: false,
    loading: () => <div className="h-96 w-full bg-gray-100 rounded-lg flex items-center justify-center"><p>Loading map...</p></div>
  }
);

// wrapper component that handles the dynamic import
const StateHighlightMap: React.FC<StateHighlightMapProps> = (props) => {
  return <ClientStateHighlightMap {...props} />;
};

export default StateHighlightMap;