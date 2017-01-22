
import { Date } from '../common/date';

export class Member {
    id: number;
    email: string;
    nickname: string;
    givenName: string;
    firstNames: string;
    initials: string;
    lastName: string;
    birthdate: Date;
    isAdult: boolean;
    phonenumber: string;
    postalcode: string;
    housenumber: string;
    housenumber_suffix: string;
    streetname: string;
    city: string;
    province: string;
    country: string;
    feeLastPaymentDate: Date;
    feePaid: boolean;
    verifiedEmail: boolean;
    verifiedIdentity: boolean;
    verifiedVotingEntitlement: boolean;

    nameParts(): string {
        return this.givenName + ` ` + this.firstNames + ` ` + this.initials + ` ` + this.lastName;
    }

    officialName(): string {
        return this.initials + ` ` + this.lastName;
    }
}