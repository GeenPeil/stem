import { JsonObject, JsonMember } from 'typedjson-npm';

@JsonObject()
export class Date {
	@JsonMember() year: number;
	@JsonMember() month: number;
	@JsonMember() day: number;

	toString(): string {
		return this.year + `-` + (this.month < 10 ? '0' : '') + this.month + `-` + (this.day < 10 ? '0' : '') + this.day
	}
}
