export interface GasPriceAvgs {
	Postcode:        string
	County:          string
	State:           string
	CountyAvgDiesel: number
	CountyAvgE10:    number
	CountyAvgE5:     number
	StateAvgDiesel:  number
	StateAvgE10:     number
	StateAvgE5:      number
}

export interface GasStation {  
  StationName: string;
  Brand: string;
  Street: string;
  Place: string;
  HouseNumber: string;
  PostCode: string;
  Latitude: number;
  Longitude: number;
  PublicHolidayIdentifier: string;
  OtJson: string;
  FirstActive: Date;
  GasPrices: GasPrice[];
}

export interface GasPrice {
  E5: number;
  E10: number;
  Diesel: number;
  Date: string;
  Changed: number;
}

export interface Notification {
  Timestamp: Date;
  UserUuid: string;
  Title: string;
  Message: string;
  DataJson: string;
}

export interface MyDataJson {
  StationName: string;
  Brand: string;
  Street: string;
  Place: string;
  HouseNumber: string;
  PostCode: string;
  Latitude: number;
  Longitude: number;
  E5: number;
  E10: number;
  Diesel: number;
  Timestamp: string;
}

export interface TimeSlotResponse {
  ID:            number;
  CreatedAt:     Date;
  UpdatedAt:     Date;
  DeletedAt:     Date;
  GasStationNum: number;
  AvgE5:         number;
  AvgE10:        number;
  AvgDiesel:     number;
  GsNumE5:       number;
  GsNumE10:      number;
  GsNumDiesel:   number;
  StartDate:     Date;    
  CountyDataID:  number;
}

export interface GsPoint {
  timestamp: string;
  price: number;
}

export interface TimeSlot {
  x: string;
  e5: number;
  e10: number;
  diesel: number;
}

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

export interface GsPoint {
  timestamp: string;
  price: number;
}

export interface TimeSlot {
  x: string;
  e5: number;
  e10: number;
  diesel: number;
}

export interface UserRequest {
  Username: string;
  Password: string;
  Latitude?: number;
  Longitude?: number;
  SearchRadius?: number;
  PostCode?: number;
  TargetDiesel?: string;
  TargetE10?: string;
  TargetE5?: string;
}

export interface UserResponse {
  Token?: string;
  Message?: string;
  PostCode?: number;
  Uuid?: string;
  Longitude?: number;
  Latitude?: number;
  SearchRadius?: number;
  TargetDiesel?: number;
  TargetE5?: number;
  TargetE10?: number;
}

export interface PostCodeLocation {
    Message: string;
  Longitude:  number;
  Latitude:  number;
  Label:      string;
  PostCode:   number
  SquareKM:   number;
    Population: number;
}
