import { JsonObject, JsonMember } from 'typedjson-npm';

@JsonObject()
export class Profile {
	@JsonMember() id: number;
	@JsonMember() nickname: string;
	@JsonMember() email: string;
}
