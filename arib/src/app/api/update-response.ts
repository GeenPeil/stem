import { JsonObject, JsonMember } from "typedjson-npm";

@JsonObject()
export class UpdateResponse {
	@JsonMember({ elements: String }) public errors: string[];

	public hasErrors(): boolean {
		return this.errors.length > 0;
	}
}
