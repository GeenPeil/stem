import { JsonObject, JsonMember } from "typedjson-npm";

@JsonObject()
export class Country {
	@JsonMember() public code: string;
	@JsonMember() public name: string;
}

@JsonObject()
export class CountryList {
	@JsonMember({ elements: Country }) public countries: Country[];
}
