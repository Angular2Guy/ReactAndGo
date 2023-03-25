import style from './gsmap.module.scss';

export interface GsValue {
    name: string;
    e5Price: number;
    e10Price: number;
    dieselPrice: number;
    longitude: number;
    latitude: number;

}

interface InputProps {
    gsValues: GsValue[];
}

export default function GsMap(inputProps: InputProps) {

    return <div className={style.MyStyle}>Hello</div>;
}