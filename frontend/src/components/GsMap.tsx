import Map from 'ol/Map';
import OSM from 'ol/source/OSM';
import { Tile as TileLayer, Vector as VectorLayer } from 'ol/layer';
import VectorSource from 'ol/source/Vector';
import Point from 'ol/geom/Point';
import Feature from 'ol/Feature';
import View from 'ol/View';
import { MapBrowserEvent, Overlay } from 'ol';
import { fromLonLat } from 'ol/proj';
import myStyle from './gsmap.module.scss';
import { Icon, Style } from 'ol/style.js';
import { useEffect } from 'react';
import { nanoid } from 'nanoid';

export interface GsValue {
  location: string;
  e5: number;
  e10: number;
  diesel: number;
  date: Date;
  longitude: number;
  latitude: number;
}

export interface CenterLocation {
  Longitude: number;
  Latitude: number;
}

interface InputProps {
  center: CenterLocation;
  gsValues: GsValue[];
}

export default function GsMap(inputProps: InputProps) {
  let map: Map;    
  let currentOverlay: Overlay | null = null;
  useEffect(() => {
    if (!map) {
      // eslint-disable-next-line react-hooks/exhaustive-deps
      map = new Map({
        layers: [
          new TileLayer({
            source: new OSM(),
          })
        ],
        target: 'map',
        view: new View({
          center: [0, 0],
          zoom: 1,
        }),
      });
    }
    const myOverlays = inputProps.gsValues.map((gsValue, index) => {
      const element = document.createElement('div');
      element.id = nanoid();
      element.innerHTML = `${gsValue.location}<br/>E5: ${gsValue.e5}<br/>E10: ${gsValue.e10}<br/>Diesel: ${gsValue.diesel}`;
      const overlay = new Overlay({
        element: element,
        offset: [-5, 0],
        positioning: 'bottom-center',
        className: 'ol-tooltip-measure ol-tooltip .ol-tooltip-static'
      });
      overlay.setPosition(fromLonLat([gsValue.longitude, gsValue.latitude]));
      const myStyle = element?.style;
      if (!!myStyle) {
        myStyle.display = 'block';
      }
      //map.addOverlay(overlay);                 
      addPins(gsValue, element, index);
      return overlay;
    });
    addClickListener(myOverlays);
    map.setView(new View({
      center: fromLonLat([inputProps.center.Longitude, inputProps.center.Latitude]),
      zoom: 12,
    }));
  }, []);

  function addPins(gsValue: GsValue, element: HTMLDivElement, index: number) {
    const iconFeature = new Feature({
      geometry: new Point(fromLonLat([gsValue.longitude, gsValue.latitude])),
      ttId: element.id,
      ttIndex: index
    });

    const iconStyle = new Style({
      image: new Icon({
        anchor: [20, 20],
        anchorXUnits: 'pixels',
        anchorYUnits: 'pixels',
        src: '/public/assets/map-pin.png',
      }),
    });
    iconFeature.setStyle(iconStyle);
    const vectorSource = new VectorSource({
      features: [iconFeature],
    });

    const vectorLayer = new VectorLayer({
      source: vectorSource,
    });
    map.addLayer(vectorLayer);
  }

  function addClickListener(myOverlays: Overlay[]) {
    map.on('click', (event: MapBrowserEvent<UIEvent>) => {
      const feature = map.forEachFeatureAtPixel(event.pixel, (feature) => {
        return feature;
      });
      if (!!currentOverlay) {
        map.removeOverlay(currentOverlay);
        // eslint-disable-next-line react-hooks/exhaustive-deps
        currentOverlay = null;
      }
      //console.log(feature);
      //console.log(feature?.get('ttId') + ' ' + feature?.get('ttIndex'));
      if (!!feature?.get('ttIndex')) {
        //console.log(myOverlays[feature?.get('ttIndex')]);
        // eslint-disable-next-line react-hooks/exhaustive-deps
        currentOverlay = myOverlays[feature?.get('ttIndex')];
        map.addOverlay(currentOverlay as Overlay);
      }
    });
  }

  return (<div className={myStyle.MyStyle}>
    <div id="map" className={myStyle.gsMap}></div>
  </div>);

}