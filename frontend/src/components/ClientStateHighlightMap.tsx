'use client';

import React, { useEffect, useState } from 'react';
import { MapContainer, TileLayer, GeoJSON } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import * as topojson from 'topojson-client';
import L from 'leaflet';

// Import all state topojson files
import andamanAndNicobarIslands from '../../topojson/states/andaman-and-nicobar-islands.json';
import andhraPradesh from '../../topojson/states/andhra-pradesh.json';
import arunachalPradesh from '../../topojson/states/arunachal-pradesh.json';
import assam from '../../topojson/states/assam.json';
import bihar from '../../topojson/states/bihar.json';
import chandigarh from '../../topojson/states/chandigarh.json';
import chhattisgarh from '../../topojson/states/chhattisgarh.json';
import delhi from '../../topojson/states/delhi.json';
import dnhAndDd from '../../topojson/states/dnh-and-dd.json';
import goa from '../../topojson/states/goa.json';
import gujarat from '../../topojson/states/gujarat.json';
import haryana from '../../topojson/states/haryana.json';
import himachalPradesh from '../../topojson/states/himachal-pradesh.json';
import jammuAndKashmir from '../../topojson/states/jammu-and-kashmir.json';
import jharkhand from '../../topojson/states/jharkhand.json';
import karnataka from '../../topojson/states/karnataka.json';
import kerala from '../../topojson/states/kerala.json';
import ladakh from '../../topojson/states/ladakh.json';
import lakshadweep from '../../topojson/states/lakshadweep.json';
import madhyaPradesh from '../../topojson/states/madhya-pradesh.json';
import maharashtra from '../../topojson/states/maharashtra.json';
import manipur from '../../topojson/states/manipur.json';
import meghalaya from '../../topojson/states/meghalaya.json';
import mizoram from '../../topojson/states/mizoram.json';
import nagaland from '../../topojson/states/nagaland.json';
import odisha from '../../topojson/states/odisha.json';
import puducherry from '../../topojson/states/puducherry.json';
import punjab from '../../topojson/states/punjab.json';
import rajasthan from '../../topojson/states/rajasthan.json';
import sikkim from '../../topojson/states/sikkim.json';
import tamilNadu from '../../topojson/states/tamilnadu.json';
import telangana from '../../topojson/states/telangana.json';
import tripura from '../../topojson/states/tripura.json';
import uttarPradesh from '../../topojson/states/uttar-pradesh.json';
import uttarakhand from '../../topojson/states/uttarakhand.json';
import westBengal from '../../topojson/states/west-bengal.json';

// Fix for Leaflet marker icons
const iconRetinaUrl =
  'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png';
const iconUrl =
  'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png';
const shadowUrl =
  'https://cdnjs.githubusercontent.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png';

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

interface StateHighlightMapProps {
  stateName: string;
  stateDisplayName?: string;
}

const ClientStateHighlightMap: React.FC<StateHighlightMapProps> = ({ stateName, stateDisplayName }) => {
  const [geoJsonData, setGeoJsonData] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadStateData = () => {
      if (!stateName) return;

      setLoading(true);
      setError(null);
      
      try {
        // Map state names to imported modules
        const stateMap: Record<string, any> = {
          'andaman-and-nicobar-islands': andamanAndNicobarIslands,
          'andhra-pradesh': andhraPradesh,
          'arunachal-pradesh': arunachalPradesh,
          'assam': assam,
          'bihar': bihar,
          'chandigarh': chandigarh,
          'chhattisgarh': chhattisgarh,
          'delhi': delhi,
          'dnh-and-dd': dnhAndDd,
          'goa': goa,
          'gujarat': gujarat,
          'haryana': haryana,
          'himachal-pradesh': himachalPradesh,
          'jammu-and-kashmir': jammuAndKashmir,
          'jharkhand': jharkhand,
          'karnataka': karnataka,
          'kerala': kerala,
          'ladakh': ladakh,
          'lakshadweep': lakshadweep,
          'madhya-pradesh': madhyaPradesh,
          'maharashtra': maharashtra,
          'manipur': manipur,
          'meghalaya': meghalaya,
          'mizoram': mizoram,
          'nagaland': nagaland,
          'odisha': odisha,
          'puducherry': puducherry,
          'punjab': punjab,
          'rajasthan': rajasthan,
          'sikkim': sikkim,
          'tamil-nadu': tamilNadu, // Special mapping for inconsistent filename
          'telangana': telangana,
          'tripura': tripura,
          'uttar-pradesh': uttarPradesh,
          'uttarakhand': uttarakhand,
          'west-bengal': westBengal,
        };
        
        // Get the state topojson from the map
        const stateTopology = stateMap[stateName];
        
        if (!stateTopology) {
          console.warn(`State topology not found for: ${stateName}`);
          setError(`State data not found for: ${stateName}`);
          setLoading(false);
          return;
        }
        
        // Convert topology to GeoJSON
        const objectKey = Object.keys(stateTopology.objects)[0];
        const featureCollection = topojson.feature(stateTopology, stateTopology.objects[objectKey]);
        
        setGeoJsonData(featureCollection);
      } catch (err) {
        console.error('Error loading state topology:', err);
        setError('Failed to load state map data');
      } finally {
        setLoading(false);
      }
    };
    
    loadStateData();
  }, [stateName]);

  const style = (feature: any) => ({
    fillColor: '#3b82f6', // Blue color for highlighted state
    weight: 2,
    opacity: 1,
    color: '#1d4ed8', // Darker blue border
    fillOpacity: 0.7,
  });

  // Calculate bounds for the specific state to center the map
  const calculateBounds = () => {
    if (!geoJsonData) return null;

    try {
      const geoJsonLayer = L.geoJSON(geoJsonData);
      return geoJsonLayer.getBounds();
    } catch (error) {
      console.error('Error calculating state bounds:', error);
      return null;
    }
  };

  if (loading) {
    return (
      <div className="h-96 w-full bg-gray-100 rounded-lg flex items-center justify-center">
        <p>Loading map for {stateDisplayName || stateName}...</p>
      </div>
    );
  }

  if (error || !geoJsonData) {
    return (
      <div className="h-96 w-full bg-gray-100 rounded-lg flex items-center justify-center">
        <p>Map data not available for {stateDisplayName || stateName}: {error || 'No data'}</p>
      </div>
    );
  }

  const bounds = calculateBounds();

  return (
    <div className="h-96 w-full rounded-lg overflow-hidden border border-gray-300">
      <MapContainer
        center={bounds ? [bounds.getCenter().lat, bounds.getCenter().lng] : [20.5937, 78.9629]}
        zoom={bounds ? 7 : 5}
        style={{ height: '100%', width: '100%' }}
        zoomControl={true}
        bounds={bounds || undefined}
        boundsOptions={{ 
          padding: [50, 50],
          maxZoom: 10
        }}
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        <GeoJSON
          data={geoJsonData}
          style={style}
        />
      </MapContainer>
    </div>
  );
};

export default ClientStateHighlightMap;