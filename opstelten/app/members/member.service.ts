import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';

import { Observable } from 'rxjs/Observable';

import { Api } from '../api/api.service';
import { CreateResponse } from '../api/create-response';
import { UpdateResponse } from '../api/update-response';

import { MemberList } from './member-list';
import { Member } from './member';

@Injectable()
export class MemberService {

    constructor(private api: Api) { }

    public getMemberList(searchQuery: string): Observable<MemberList> {
        console.log("getMembers", searchQuery);
        return this.api.get(MemberList, 'member-list?searchQuery=' + encodeURIComponent(searchQuery)).catch(this.handleError);
    }

    public postMember(member: Member): Promise<CreateResponse> {
        return this.api.post('member', JSON.stringify(member)).toPromise();
    }

    public getMember(id: number): Promise<Member> {
        console.log("getMember", id);
        return this.api.get(Member, 'member/' + id).toPromise();
    }

    public putMember(id: number, member: Member): Promise<UpdateResponse> {
        return this.api.put('member/' + id, JSON.stringify(member)).toPromise();
    }

    private handleError(error: any) {
        let errMsg = (error.message) ? error.message : error.status ? `${error.status} - ${error.statusText}` : 'Server error';
        console.error(errMsg);
        return Observable.throw(errMsg);
    }
}
