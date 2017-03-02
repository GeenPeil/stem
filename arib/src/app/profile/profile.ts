import { JsonObject, JsonMember } from "typedjson-npm";

@JsonObject()
export class Profile {
	@JsonMember() public id: number;
	@JsonMember() public nickname: string;
	@JsonMember() public email: string;
}
