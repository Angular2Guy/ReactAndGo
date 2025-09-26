import { UserDataState } from "../GlobalState";
import { GasPriceAvgs, GasStation, Notification, TimeSlotResponse } from "./dtos";

const fetchGasStations = async function(jwtToken: string, controller: AbortController | null, globalUserDataState: UserDataState): Promise<GasStation[]> {
  const requestOptions2 = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}` },
      body: JSON.stringify({ Longitude: globalUserDataState.Longitude, Latitude: globalUserDataState.Latitude, Radius: globalUserDataState.SearchRadius }),
      signal: controller?.signal
    }      
    const result = await fetch('/gasstation/search/location', requestOptions2);
    const myResult = result.json() as Promise<GasStation[]>;
    return myResult;
};

const fetchPriceAvgs = async function(jwtToken: string, controller: AbortController | null, myPostcode: string): Promise<GasPriceAvgs> {
    const requestOptions3 = {
        method: 'GET',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}` },        
        signal: controller?.signal
      }  
    const result = await fetch(`/gasprice/avgs/${myPostcode}`, requestOptions3);
    return result.json() as Promise<GasPriceAvgs>;
}

const fetchUserNotifications = async function(jwtToken: string, controller: AbortController | null, globalUserUuidState: string): Promise<Notification[]> {
    const requestOptions1 = {
      method: 'GET',
      headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}` },
      signal: controller?.signal
    }
    const result = await fetch(`/usernotification/current/${globalUserUuidState}`, requestOptions1);
    return result.json() as Promise<Notification[]>;
}

const fetchTimeSlots = async function(jwtToken: string, controller: AbortController | null, myPostcode: string): Promise<TimeSlotResponse[]> {
const requestOptions2 = {
        method: 'GET',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}` },
        signal: controller?.signal
      }      
      const result = await fetch(`/postcode/countytimeslots/${myPostcode}`, requestOptions2);
      const myResult = result.json() as Promise<TimeSlotResponse[]>;
      return myResult;
    }

export { fetchGasStations };
export { fetchPriceAvgs };
export { fetchUserNotifications };
export { fetchTimeSlots };