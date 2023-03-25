import Map from 'ol/Map.js';
import OSM from 'ol/source/OSM.js';
import TileLayer from 'ol/layer/Tile.js';
import View from 'ol/View.js';
import {fromLonLat} from 'ol/proj';
import style from './gsmap.module.scss';
import { useEffect } from 'react';

export interface GsValue {
    name: string;
    e5Price: number;
    e10Price: number;
    dieselPrice: number;
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
      useEffect(() => {        
        if(!!map) {
            map.setView(new View({
                center: fromLonLat([inputProps.center.Longitude,inputProps.center.Latitude]),
                zoom: 12,
            }));            
        }
        map = !map ?  new Map({
            layers: [
              new TileLayer({
                source: new OSM(),
              })          
            ],
            target: 'map',
            view: new View({
                center: [0,0],
                zoom: 1,
              }),
          }) : map;        
      }, []);

    return (<div className={style.MyStyle}>
        <div id="map" className={style.gsMap}></div>
    </div>);
}