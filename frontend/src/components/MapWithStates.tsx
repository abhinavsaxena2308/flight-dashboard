'use client';

import React, { useMemo } from 'react';
import { MapContainer, TileLayer, GeoJSON, CircleMarker } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import * as topojson from 'topojson-client';
import indiaTopology from '../../topojson/india.json';
import L from 'leaflet';

const iconRetinaUrl =
  'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png';
const iconUrl =
  'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png';
const shadowUrl =
  'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png';

// Fix for Leaflet marker icons
try {
  delete (L as any).Icon.Default.prototype._getIconUrl;
} catch (e) {
  // Ignore if property doesn't exist
}

L.Icon.Default.mergeOptions({
  iconRetinaUrl,
  iconUrl,
  shadowUrl,
});

interface StateData {
  state: string;
  totalFlights: number;
  incomingFlights: number;
  outgoingFlights: number;
  routes: number;
  airlines: string[];
}

interface TooltipState {
  visible: boolean;
  content: string | StateData;
  x: number;
  y: number;
  loading?: boolean;
  error?: string;
}

interface MapWithStatesProps {
  setTooltip: React.Dispatch<React.SetStateAction<TooltipState>>;
  onStateClick?: (stateName: string) => void;
}

const STATE_COLORS: { [key: string]: string } = {
  'andaman-and-nicobar-islands': '#0ea5e9',
  'andhra-pradesh': '#ef4444',
  'arunachal-pradesh': '#f97316',
  'assam': '#84cc16',
  'bihar': '#eab308',
  'chandigarh': '#14b8a6',
  'chhattisgarh': '#f59e0b',
  'delhi': '#6366f1',
  'dnh-and-dd': '#ec4899',
  'goa': '#06b6d4',
  'gujarat': '#f43f5e',
  'haryana': '#8b5cf6',
  'himachal-pradesh': '#10b981',
  'jammu-and-kashmir': '#3b82f6',
  'jharkhand': '#d946ef',
  'karnataka': '#a855f7',
  'kerala': '#14b8a6',
  'ladakh': '#64748b',
  'lakshadweep': '#fffff',
  'madhya-pradesh': '#f59e0b',
  'maharashtra': '#f97316',
  'manipur': '#ec4899',
  'meghalaya': '#22c55e',
  'mizoram': '#eab308',
  'nagaland': '#ef4444',
  'odisha': '#06b6d4',
  'puducherry': '#6366f1',
  'punjab': '#84cc16',
  'rajasthan': '#eab308',
  'sikkim': '#14b8a6',
  'tamil-nadu': '#f43f5e',
  'telangana': '#8b5cf6',
  'tripura': '#d946ef',
  'uttar-pradesh': '#a855f7',
  'uttarakhand': '#10b981',
  'west-bengal': '#f43f5e',
};

const getColor = (name: string) => {
  if (!name) return '#CCCCCC'; // Default color for unknown states
  
  // Convert name to kebab-case to match keys
  const normalizedName = name
    .toLowerCase()
    .replace(/&/g, 'and')
    .replace(/\s+/g, '-')
    .replace(/[^\w-]/g, '');

  return STATE_COLORS[normalizedName] || '#CCCCCC';
};

const ISLAND_MARKERS = [
  { name: 'Lakshadweep', coords: [10.57, 72.64] as [number, number] },
];

const MapWithStates: React.FC<MapWithStatesProps> = ({ setTooltip, onStateClick }) => {
  const geoJsonData = useMemo(() => {
    // @ts-ignore
    const topology = indiaTopology as any;
    // Get the first object key (usually 'india' or 'states')
    const objectKey = Object.keys(topology.objects)[0];
    const geometries = topology.objects[objectKey].geometries;

    // Group geometries by state name
    const groupedByState = geometries.reduce((acc: any, geom: any) => {
      const stateName = geom.properties.st_nm || geom.properties.name || 'Unknown';
      if (!acc[stateName]) {
        acc[stateName] = [];
      }
      acc[stateName].push(geom);
      return acc;
    }, {});

    // Merge geometries for each state
    const features = Object.keys(groupedByState).map(stateName => {
      const stateGeometries = groupedByState[stateName];
      const mergedGeometry = topojson.merge(topology, stateGeometries);
      return {
        type: 'Feature',
        geometry: mergedGeometry,
        properties: {
          st_nm: stateName,
          name: stateName // Ensure compatibility with existing code
        }
      };
    });

    return {
      type: 'FeatureCollection',
      features: features
    };
  }, []);

  const style = (feature: any) => ({
    fillColor: getColor(
      feature.properties?.st_nm || feature.properties?.name || 'Unknown'
    ),
    weight: 0.7,
    opacity: 1,
    color: 'black',
    fillOpacity: 0.7,
  });

  const onEachFeature = (feature: any, layer: L.Layer) => {
    layer.on({
      mouseover: async (e: any) => {
        e.target.setStyle({
          weight: 2,
          color: '#666',
          fillOpacity: 0.9,
        });

        const stateName =
          feature.properties?.st_nm || feature.properties?.name || 'Unknown';
        
        // Set loading state
        setTooltip({
          visible: true,
          content: '',
          x: e.originalEvent.pageX,
          y: e.originalEvent.pageY,
          loading: true,
        });
        
        try {
          // Convert state name to kebab-case to match API expectations
          const formattedStateName = stateName
            .toLowerCase()
            .replace(/&/g, 'and')
            .replace(/\s+/g, '-')
            .replace(/[^\w-]/g, '');
            
          // Fetch state data from backend API (proxied through Next.js)
          const response = await fetch(`/api/state/${formattedStateName}`);
          
          if (response.ok) {
            const stateData: StateData = await response.json();
            
            setTooltip({
              visible: true,
              content: stateData,
              x: e.originalEvent.pageX,
              y: e.originalEvent.pageY,
              loading: false,
            });
          } else {
            // If API call fails, show just the state name
            setTooltip({
              visible: true,
              content: stateName,
              x: e.originalEvent.pageX,
              y: e.originalEvent.pageY,
              loading: false,
              error: `Failed to load data (${response.status})`
            });
          }
        } catch (error) {
          console.error('Error fetching state data:', error);
          setTooltip({
            visible: true,
            content: stateName,
            x: e.originalEvent.pageX,
            y: e.originalEvent.pageY,
            loading: false,
            error: 'Network error'
          });
        }
      },

      mousemove: (e: any) => {
        setTooltip((prev) => ({
          ...prev,
          x: e.originalEvent.pageX,
          y: e.originalEvent.pageY,
        }));
      },

      click: (e: any) => {
        const stateName =
          feature.properties?.st_nm || feature.properties?.name || 'Unknown';
        
        // Convert state name to kebab-case for URL
        const formattedStateName = stateName
          .toLowerCase()
          .replace(/&/g, 'and')
          .replace(/\s+/g, '-')
          .replace(/[^\w-]/g, '');
        
        // Call the click handler if provided
        if (onStateClick) {
          onStateClick(formattedStateName);
        } else {
          // Fallback to direct navigation
          window.location.href = `/state/${formattedStateName}`;
        }
      },
      
      mouseout: (e: any) => {
        e.target.setStyle({
          weight: 1,
          color: 'white',
          fillOpacity: 0.7,
        });

        setTooltip({
          visible: false,
          content: '',
          x: 0,
          y: 0,
        });
      },
    });
  };

  return (
    <MapContainer
  center={[20.5937, 78.9629]}
  zoom={5}
  style={{ height: '100%', width: '100%' }}
  zoomControl={false}
>
      {geoJsonData && (
        <GeoJSON
          data={geoJsonData as any}
          style={style}
          onEachFeature={onEachFeature}
        />
      )}

      
      {ISLAND_MARKERS.map((island) => (
        <CircleMarker
          key={island.name}
          center={island.coords}
          radius={6}
          pathOptions={{
            fillColor: getColor(island.name),
            color: 'white',
            weight: 1,
            fillOpacity: 0.8
          }}
          eventHandlers={{
            mouseover: (e) => {
              e.target.setStyle({
                weight: 2,
                color: '#666',
                fillOpacity: 1
              });
              setTooltip({
                visible: true,
                content: island.name,
                x: e.originalEvent.pageX,
                y: e.originalEvent.pageY
              });
            },
            mousemove: (e) => {
              setTooltip((prev) => ({
                ...prev,
                x: e.originalEvent.pageX,
                y: e.originalEvent.pageY
              }));
            },
            mouseout: (e) => {
              e.target.setStyle({
                weight: 1,
                color: 'white',
                fillOpacity: 0.8
              });
              setTooltip((prev) => ({ ...prev, visible: false }));
            }
          }}
        />
      ))}
    </MapContainer>
  );
};

export default MapWithStates;
