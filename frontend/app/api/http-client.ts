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
import { type UserDataState } from "../GlobalState";
import { type GasPriceAvgs } from "../model/gas-price-avg";
import { type GasStation } from "../model/gas-station";
import { type PostCodeLocation } from "../model/location";
import { type TimeSlotResponse } from "../model/time-slot-response";
import { type UserRequest, type UserResponse } from "../model/user";
import { type Notification } from "../model/notification";
import type { NavigateFunction } from "react-router";

const apiPrefix = '/api';

async function handleResponse<T>(response: Response,navigate: NavigateFunction): Promise<T> {
  if (!response.ok) {
    const error = await response.text();
    if(navigate) navigate('/');
    throw new Error(error || `HTTP error! status: ${response.status}`);
  }
  return response.json();
}

const fetchGasStations = async function (navigate: NavigateFunction, jwtToken: string, controller: AbortController | null, globalUserDataState: UserDataState): Promise<GasStation[]> {
  const requestOptions2 = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}` },
    body: JSON.stringify({ Longitude: globalUserDataState.Longitude, Latitude: globalUserDataState.Latitude, Radius: globalUserDataState.SearchRadius }),
    signal: controller?.signal
  }
  const result = await fetch(`${apiPrefix}/gasstation/search/location`, requestOptions2);
  return handleResponse<GasStation[]>(result,navigate);
};

const fetchPriceAvgs = async function (navigate: NavigateFunction, jwtToken: string, controller: AbortController | null, myPostcode: string): Promise<GasPriceAvgs> {
  const requestOptions3 = {
    method: 'GET',
    headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}` },
    signal: controller?.signal
  }
  const result = await fetch(`${apiPrefix}/gasprice/avgs/${myPostcode}`, requestOptions3);
  return handleResponse<GasPriceAvgs>(result,navigate);
}

const fetchUserNotifications = async function (navigate: NavigateFunction, jwtToken: string, controller: AbortController | null, globalUserUuidState: string): Promise<Notification[]> {
  const requestOptions1 = {
    method: 'GET',
    headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}` },
    signal: controller?.signal
  }
  const result = await fetch(`${apiPrefix}/usernotification/current/${globalUserUuidState}`, requestOptions1);
  return handleResponse<Notification[]>(result,navigate);
}

const fetchTimeSlots = async function (navigate: NavigateFunction, jwtToken: string, controller: AbortController | null, myPostcode: string): Promise<TimeSlotResponse[]> {
  const requestOptions2 = {
    method: 'GET',
    headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}` },
    signal: controller?.signal
  }
  const result = await fetch(`${apiPrefix}/postcode/countytimeslots/${myPostcode}`, requestOptions2);
  return handleResponse<TimeSlotResponse[]>(result,navigate);
}

const postLogin = async function (userName: string, password1: string, controller: AbortController | null): Promise<UserResponse> {
  const requestOptions = loginSigninOptions(userName, password1, controller);
  const result = await fetch(`${apiPrefix}/appuser/login`, requestOptions);
  return handleResponse<UserResponse>(result, null as unknown as NavigateFunction);
}

const postSignin = async function (userName: string, password1: string, controller: AbortController | null): Promise<UserResponse> {
  const requestOptions = loginSigninOptions(userName, password1, controller);
  const result = await fetch(`${apiPrefix}/appuser/signin`, requestOptions);
  return handleResponse<UserResponse>(result, null as unknown as NavigateFunction);
}

const loginSigninOptions = (userName: string, password1: string, controller: AbortController | null) => {
  return {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ Username: userName, Password: password1 } as UserRequest),
    signal: controller?.signal
  };
};

const postLocationRadius = async function (navigate: NavigateFunction, jwtToken: string, controller: AbortController | null, requestString: string): Promise<UserResponse> {
  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}` },
    body: requestString,
    signal: controller?.signal
  };
  const response = await fetch(`${apiPrefix}/appuser/locationradius`, requestOptions);
  return handleResponse<UserResponse>(response,navigate);
}

const fetchLocation = async function (navigate: NavigateFunction, jwtToken: string, controller: AbortController | null, location: string): Promise<PostCodeLocation[]> {
  const requestOptions = {
    method: 'GET',
    headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}` },
    signal: controller?.signal
  };
  const response = await fetch(`${apiPrefix}/appuser/location?location=${location}`, requestOptions);
  return handleResponse<PostCodeLocation[]>(response,navigate);
}

const postTargetPrices = async function (navigate: NavigateFunction, jwtToken: string, controller: AbortController | null, requestString: string): Promise<UserResponse> {
  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${jwtToken}` },
    body: requestString,
    signal: controller?.signal
  };
  const response = await fetch(`${apiPrefix}/appuser/targetprices`, requestOptions);
  return handleResponse<UserResponse>(response,navigate);
}

export { fetchGasStations };
export { fetchPriceAvgs };
export { fetchUserNotifications };
export { fetchTimeSlots };
export { postLogin };
export { postSignin };
export { postLocationRadius };
export { fetchLocation };
export { postTargetPrices };
