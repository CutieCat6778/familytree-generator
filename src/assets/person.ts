import type { Codes } from "./ISO3166-1.alpha2";

export type Person = {
	firstName: string;
	lastName: string;
	country: Codes[];
	age: number;
	death: boolean;
	childrens: number[]; // represent id
    isChild: boolean;
    marriage?: {
        person: number,
        isDivorced: boolean
    },
    sex: string,
};
