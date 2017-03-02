import { JsonObject, JsonMember } from "typedjson-npm";

@JsonObject()
export class CreateResponse {
	@JsonMember public id: number;
	@JsonMember({ elements: String }) public errors: string[];

	public hasErrors(): boolean {
		return this.errors.length > 0;
	}
}
