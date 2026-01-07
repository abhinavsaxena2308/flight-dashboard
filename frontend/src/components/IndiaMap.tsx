'use client';

import dynamic from 'next/dynamic';
import React, { useState } from 'react';

const MapComponent = dynamic(() => import('./MapWithStates'), {
  ssr: false,
  loading: () => <div className="w-full h-full bg-gray-100 flex items-center justify-center">Loading map...</div>
});

interface StateData {
  state: string;
  totalFlights: number;
  incomingFlights: number;
  outgoingFlights: number;
  routes: number;
  airlines: string[];
}

const IndiaMap: React.FC = () => {
  const [tooltip, setTooltip] = useState<{ visible: boolean; content: string | StateData; x: number; y: number; loading?: boolean; error?: string }>({ 
    visible: false, 
    content: '', 
    x: 0, 
    y: 0 
  });

  return (
    <div className="w-full h-full relative">
      <MapComponent setTooltip={setTooltip} />
      {/* Custom tooltip */}
      {tooltip.visible && (
        <div 
          style={{
            position: 'fixed',
            left: `${tooltip.x + 10}px`,
            top: `${tooltip.y + 10}px`,
            backgroundColor: 'white',
            border: '1px solid #ccc',
            borderRadius: '4px',
            padding: '8px',
            fontSize: '14px',
            fontWeight: 'normal',
            zIndex: 1000,
            pointerEvents: 'none',
            boxShadow: '0 2px 4px rgba(0,0,0,0.1)',
            minWidth: '200px',
            maxWidth: '300px',
            textAlign: 'left'
          }}
        >
          {typeof tooltip.content === 'string' ? (
            <div>
              {tooltip.loading ? 'Loading...' : tooltip.content}
              {tooltip.error && <div style={{color: 'red', fontSize: '12px'}}>Error: {tooltip.error}</div>}
            </div>
          ) : (
            <div>
              <div style={{fontWeight: 'bold', marginBottom: '4px'}}>{tooltip.content.state}</div>
              <div>Total Flights: {tooltip.content.totalFlights}</div>
              <div>Incoming: {tooltip.content.incomingFlights}</div>
              <div>Outgoing: {tooltip.content.outgoingFlights}</div>
              <div>Routes: {tooltip.content.routes}</div>
              <div style={{marginTop: '4px'}}>
                Airlines: {tooltip.content.airlines.slice(0, 3).join(', ')}
                {tooltip.content.airlines.length > 3 && ` +${tooltip.content.airlines.length - 3}`}
              </div>
              {tooltip.error && <div style={{color: 'red', fontSize: '12px', marginTop: '4px'}}>Error: {tooltip.error}</div>}
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default IndiaMap;