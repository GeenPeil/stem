import { Component } from '@angular/core';

import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Observable } from 'rxjs/Observable';

import { MemberList } from './member-list';
import { Member } from './member';
import { MemberService } from './member.service';

@Component({
	templateUrl: 'app/members/member-list.component.html',
	styles: [`
		.actions {
			border-left: 1px dashed #333;
			text-align: right;
		}`]
})
export class MemberListComponent {

	constructor(private memberService: MemberService) { }

	searchQuery: string;
	private searchQueryFilterStream = new BehaviorSubject<string>(``);
	search() { this.searchQueryFilterStream.next(this.searchQuery); }

	memberList: Observable<MemberList> = this.searchQueryFilterStream
		.filter((val: string) => val != ``)
		.debounceTime(400)
		.distinctUntilChanged()
		.switchMap((searchQuery: string) => this.memberService.getMemberList(searchQuery));
}
