import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';

import { Observable } from 'rxjs/Observable';

import { APIResponse } from '../common/api-response';
import { AuthHttp } from '../auth/auth-http.service';

import { Member } from './member';

@Injectable()
export class MemberService {

    constructor(private http: AuthHttp) { }

    // public getMembers(searchQuery: string): Observable<Member[]> {
    //     console.log("getMembers", searchQuery);
    //     return this.http.get(Member, 'members')
    //         .map((res: Response) => {
    //             let members = res.json().members.filter((member: Member) => {
    //                 if (searchQuery === `` || searchQuery == null) {
    //                     return true;
    //                 }
    // 
    //                 let memberNameMatches: boolean = true;
    //                 for (var searchQueryPart of searchQuery.split(` `)) {
    //                     if (member.nameParts().indexOf(searchQueryPart) == -1) {
    //                         memberNameMatches = false;
    //                         break;
    //                     }
    //                 }
    //                 if (memberNameMatches) {
    //                     return true;
    //                 }
    // 
    //                 if (member.postalcode === searchQuery) {
    //                     return true;
    //                 }
    //             });
    //             console.dir(members);
    //             return members;
    //         })
    //         .catch(this.handleError)
    // }

    public getMember(id: number): Promise<Member> {
        console.log("getMember", id);
        return this.http.get(Member, 'member/' + id).toPromise();
    }

    public putMember(id: number, member: Member): Promise<APIResponse> {
        return this.http.put('member/' + id, JSON.stringify(member)).toPromise();
    }

    public postMember(member: Member): Promise<APIResponse> {
        return this.http.post('member', JSON.stringify(member)).toPromise();
    }

    private handleError(error: any) {
        // In a real world app, we might use a remote logging infrastructure
        // We'd also dig deeper into the error to get a better message
        let errMsg = (error.message) ? error.message :
            error.status ? `${error.status} - ${error.statusText}` : 'Server error';
        console.error(errMsg); // log to console instead
        return Observable.throw(errMsg);
    }
}
