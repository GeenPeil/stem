import { JsonObject, JsonMember } from 'typedjson-npm';

import { Date } from '../common/date';

@JsonObject()
export class Member {
    @JsonMember() id: number;
    @JsonMember() email: string;
    @JsonMember() nickname: string;
    @JsonMember() givenName: string;
    @JsonMember() firstNames: string;
    @JsonMember() initials: string;
    @JsonMember() lastName: string;
    @JsonMember() birthdate: Date;
    @JsonMember() isAdult: boolean;
    @JsonMember() phonenumber: string;
    @JsonMember() postalcode: string;
    @JsonMember() housenumber: string;
    @JsonMember() housenumberSuffix: string;
    @JsonMember() streetname: string;
    @JsonMember() city: string;
    @JsonMember() province: string;
    @JsonMember() country: string;
    @JsonMember() feeLastPaymentDate: Date;
    @JsonMember() feePaid: boolean;
    @JsonMember() verifiedEmail: boolean;
    @JsonMember() verifiedIdentity: boolean;
    @JsonMember() verifiedVotingEntitlement: boolean;

    nameParts(): string {
        return this.givenName + ` ` + this.firstNames + ` ` + this.initials + ` ` + this.lastName;
    }

    officialName(): string {
        return this.initials + ` ` + this.lastName;
    }
}
