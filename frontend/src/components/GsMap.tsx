import Map from 'ol/Map.js';
import OSM from 'ol/source/OSM.js';
import TileLayer from 'ol/layer/Tile.js';
import View from 'ol/View.js';
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
    let view: View;
    let map: Map;
      useEffect(() => {        
        view = !view ? new View({
            center: [0,0],
            zoom: 1,
          }) : view;
        map = !map ?  new Map({
            layers: [
              new TileLayer({
                source: new OSM(),
              })          
            ],
            target: 'map',
            view: view,
          }) : map;
      }, []);

      

    return (<div className={style.MyStyle}>
        <div id="map" className={style.gsMap}></div>
    </div>);
}