import { JsonObject, JsonMember } from 'typedjson-npm';

@JsonObject()
export class Country {
	@JsonMember() code: string;
	@JsonMember() name: string;
}

@JsonObject()
export class CountryList {
	@JsonMember({ elements: Country }) countries: Country[];
}