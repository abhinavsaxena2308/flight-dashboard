'use client';

import dynamic from 'next/dynamic';
import React, { useState } from 'react';

const MapComponent = dynamic(() => import('./MapWithStates'), {
  ssr: false,
  loading: () => <div className="w-full h-full bg-gray-100 flex items-center justify-center">Loading map...</div>
});

const IndiaMap: React.FC = () => {
  const [tooltip, setTooltip] = useState<{ visible: boolean; content: string; x: number; y: number }>({ 
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
            minWidth: '80px',
            textAlign: 'center'
          }}
        >
          {tooltip.content}
        </div>
      )}
    </div>
  );
};

export default IndiaMap;