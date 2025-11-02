/*
  - Copyright 2022 Sven Loesekann
    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/
import * as React from 'react';
import Map from 'ol/Map';
import OSM from 'ol/source/OSM';
import { Tile as TileLayer, Vector as VectorLayer } from 'ol/layer';
import VectorSource from 'ol/source/Vector';
import Point from 'ol/geom/Point';
import Feature from 'ol/Feature';
import View from 'ol/View';
import { MapBrowserEvent, Overlay } from 'ol';
import { fromLonLat } from 'ol/proj';
import myStyle from './gsmap.module.css';
import { Icon, Style } from 'ol/style.js';
import { useEffect } from 'react';
import { nanoid } from 'nanoid';
import { type CenterLocation } from '../model/location';
import { type GsValue } from '../model/gs-point';
import { useTranslation } from 'node_modules/react-i18next';

interface InputProps {
  center: CenterLocation;
  gsValues: GsValue[];
}

export default function GsMap(inputProps: InputProps) {
  const { i18n } = useTranslation();
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
    const myOverlays = createOverlays();
    addClickListener(myOverlays);
    map.setView(new View({
      center: fromLonLat([inputProps.center.Longitude, inputProps.center.Latitude]),
      zoom: 12,
    }));
  }, []);

  function createOverlays(): Overlay[] {
    return inputProps.gsValues.map((gsValue, index) => {
      const element = document.createElement('div');
      element.id = nanoid();
      element.innerHTML = `${gsValue.location}<br/>${i18n.t('table.e5')}: ${gsValue.e5}<br/>${i18n.t('table.e10')}: ${gsValue.e10}<br/>${i18n.t('table.diesel')}: ${gsValue.diesel}`;
      const overlay = new Overlay({
        element: element,
        offset: [-5, 0],
        positioning: 'bottom-center',
        className: 'ol-tooltip-measure ol-tooltip ol-tooltip-static'
      });
      overlay.setPosition(fromLonLat([gsValue.longitude, gsValue.latitude]));
      const myStyle1 = element?.style;
      if (!!myStyle1) {
        myStyle1.display = 'block';
      }           
      if(element?.classList) {
        element.classList.add(myStyle['gsTooltip']);
      }       
      //map.addOverlay(overlay);                 
      addPins(gsValue, element, index);
      return overlay;
    });
  }

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
        src: '/public/map-pin.png',
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
    map.on('click', (event: MapBrowserEvent<any>) => {
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
