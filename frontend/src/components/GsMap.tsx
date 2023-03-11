import style from './gsmap.module.scss';

export interface GsValue {
    name: string;
    e5Price: number;
    e10Price: number;
    dieselPrice: number;
    longitude: number;
    latitude: number;

}

export default function GsMap(gsValues: GsValue[]) {

    return <div className={style.MyStyle}></div>;
}