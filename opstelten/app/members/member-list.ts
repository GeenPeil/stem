import { JsonObject, JsonMember } from 'typedjson-npm';

import { Member } from './member';

@JsonObject()
export class MemberList {
	@JsonMember() error: string;
	@JsonMember({ elements: Member }) members: Member[];
}
