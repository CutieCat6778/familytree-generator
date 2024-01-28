import { Codes } from "../assets/ISO3166-1.alpha2";
import Forenames from "../assets/common-forenames-by-country.min.json";
import Surenames from "../assets/common-surnames-by-country.min.json";

export function randomEnumKey<T extends Record<string, string | number>>(
	anEnum: T
): keyof T {
	const enumKeys = Object.keys(anEnum) as (keyof T)[];
	const randomIndex = Math.floor(Math.random() * enumKeys.length);
	const randomEnumKey = enumKeys[randomIndex];
	return randomEnumKey;
}

export function generateAFamily(country: string[], config?: {
	childrenThreshold?: number;
	maleThreshold?: number;
	bothParentsSameCountry?: boolean;
	mothersCountry?: string[];
    id?: number;
}) {
	const rndValue = Math.random() < 0.5;
    if (!country) country = [randomEnumKey(Codes)];
	if (!config) {
		config = {
			childrenThreshold: Math.random(),
			maleThreshold: Math.random(),
			bothParentsSameCountry: rndValue,
			mothersCountry: rndValue ? country : [Codes[randomEnumKey(Codes)]],
            id: 0,
		};
	} else {
		config = {
			childrenThreshold: Math.random(),
			maleThreshold: Math.random(),
			bothParentsSameCountry: rndValue,
			mothersCountry: rndValue ? country : [Codes[randomEnumKey(Codes)]],
            id: 0,
			...config,
		};
	}
}

function generateAPerson(country: Codes[], config?: {
    isParent: boolean
}) {

}

function generateAName(country: string[], sex: string) {
    let surname = "";
    let forename = "";
    const randomCountry: string = country[Math.round(Math.random() * country.length)] ?? country[0];
	const SurnameCountry = Surenames[randomCountry as keyof typeof Surenames];
	const ForenamesCountry = Forenames[randomCountry as keyof typeof Forenames];
    if(sex == "Female") {
		const name = SurnameCountry[Math.round(Math.random() * SurnameCountry.length)].localized;
    } else {

    }
}

function getEnumKey(country: Codes) {
    console.log(country);
}