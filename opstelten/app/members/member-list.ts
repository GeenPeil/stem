import { JsonObject, JsonMember } from 'typedjson-npm';

import { Member } from './member';

@JsonObject()
export class MemberList {
    @JsonMember({ elements: Member }) members: Member[];
}